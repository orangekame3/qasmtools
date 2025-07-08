package ast

import (
	"fmt"
	"strings"

	"github.com/orangekame3/qasmtools/lint/astutil"
	"github.com/orangekame3/qasmtools/parser"
)

// UndefinedIdentifierRule implements QAS002 using AST-based analysis
type UndefinedIdentifierRule struct {
	*ASTRuleBase
}

// NewUndefinedIdentifierRule creates a new AST-based undefined identifier rule
func NewUndefinedIdentifierRule() ASTRule {
	return &UndefinedIdentifierRule{
		ASTRuleBase: NewASTRuleBase("QAS002"),
	}
}

// CheckAST performs AST-based undefined identifier analysis
func (r *UndefinedIdentifierRule) CheckAST(program *parser.Program, ctx *CheckContext) []*Violation {
	var violations []*Violation

	// Find all declarations to build a symbol table
	declarations := astutil.FindDeclarations(program)
	declaredIdentifiers := make(map[string]bool)
	
	// Add quantum identifiers
	for _, decl := range declarations.Quantum {
		declaredIdentifiers[decl.Identifier] = true
	}
	
	// Add classical identifiers
	for _, decl := range declarations.Classical {
		declaredIdentifiers[decl.Identifier] = true
	}
	
	// Add gate identifiers
	for _, decl := range declarations.Gates {
		declaredIdentifiers[decl.Name] = true
	}
	
	// Fallback: extract gate definitions from text if AST parsing missed them
	if len(declarations.Gates) == 0 {
		r.extractGateDefinitionsFromText(ctx.Content, declaredIdentifiers)
	}
	
	// Add built-in identifiers (gates from stdgates.qasm, etc.)
	builtins := r.getBuiltinIdentifiers()
	for _, builtin := range builtins {
		declaredIdentifiers[builtin] = true
	}

	// Find all identifier usages and check if they're declared
	r.checkIdentifierUsages(program, declaredIdentifiers, ctx, &violations)

	return violations
}

// checkIdentifierUsages traverses the AST and checks all identifier usages
func (r *UndefinedIdentifierRule) checkIdentifierUsages(program *parser.Program, declaredIdentifiers map[string]bool, ctx *CheckContext, violations *[]*Violation) {
	astutil.VisitAllNodes(program, func(node parser.Node) {
		switch n := node.(type) {
		case *parser.GateCall:
			// Check gate name
			if n.Name != "" && !declaredIdentifiers[n.Name] && !r.isKeyword(n.Name) {
				violation := r.NewViolationBuilder().
					WithMessage(fmt.Sprintf("Identifier '%s' is not declared.", n.Name)).
					WithFile(ctx.File).
					WithNode(n).
					WithNodeName(n.Name).
					AsError().
					Build()
				*violations = append(*violations, violation)
			}
			
			// Check qubit arguments
			for _, qubit := range n.Qubits {
				if identifier := r.extractIdentifierName(qubit); identifier != "" {
					if !declaredIdentifiers[identifier] && !r.isKeyword(identifier) {
						violation := r.NewViolationBuilder().
							WithMessage(fmt.Sprintf("Identifier '%s' is not declared.", identifier)).
							WithFile(ctx.File).
							WithNode(qubit).
							WithNodeName(identifier).
							AsError().
							Build()
						*violations = append(*violations, violation)
					}
				}
			}
			
		case *parser.Measurement:
			// Check measurement source (qubit)
			if n.Qubit != nil {
				if identifier := r.extractIdentifierName(n.Qubit); identifier != "" {
					if !declaredIdentifiers[identifier] && !r.isKeyword(identifier) {
						violation := r.NewViolationBuilder().
							WithMessage(fmt.Sprintf("Identifier '%s' is not declared.", identifier)).
							WithFile(ctx.File).
							WithNode(n.Qubit).
							WithNodeName(identifier).
							AsError().
							Build()
						*violations = append(*violations, violation)
					}
				}
			}
			
			// Check measurement target (classical bit)
			if n.Target != nil {
				if identifier := r.extractIdentifierName(n.Target); identifier != "" {
					if !declaredIdentifiers[identifier] && !r.isKeyword(identifier) {
						violation := r.NewViolationBuilder().
							WithMessage(fmt.Sprintf("Identifier '%s' is not declared.", identifier)).
							WithFile(ctx.File).
							WithNode(n.Target).
							WithNodeName(identifier).
							AsError().
							Build()
						*violations = append(*violations, violation)
					}
				}
			}
		}
	})
}

// extractIdentifierName extracts the base identifier name from various expression types
func (r *UndefinedIdentifierRule) extractIdentifierName(expr parser.Expression) string {
	switch e := expr.(type) {
	case *parser.Identifier:
		return e.Name
	case *parser.IndexedIdentifier:
		return e.Name
	case *parser.RangedIdentifier:
		return e.Name
	default:
		return ""
	}
}

// getBuiltinIdentifiers returns a list of built-in identifiers
func (r *UndefinedIdentifierRule) getBuiltinIdentifiers() []string {
	return []string{
		// Standard gates
		"h", "x", "y", "z", "s", "sdg", "t", "tdg",
		"cx", "cnot", "cz", "ccx", "toffoli",
		"rx", "ry", "rz", "u1", "u2", "u3",
		"phase", "cphase", "crx", "cry", "crz",
		"swap", "cswap", "fredkin",
		"p", "cp", "u", "cu",
		// Built-in functions
		"sin", "cos", "tan", "exp", "ln", "sqrt",
		"pi", "euler", "tau",
		// OpenQASM keywords
		"OPENQASM", "include", "qreg", "creg", "gate", "measure",
		"reset", "barrier", "if", "else", "for", "while",
		"def", "return", "input", "output",
	}
}

// isKeyword checks if an identifier is a reserved keyword
func (r *UndefinedIdentifierRule) isKeyword(identifier string) bool {
	keywords := []string{
		"OPENQASM", "include", "qreg", "creg", "gate", "measure",
		"reset", "barrier", "if", "else", "for", "while",
		"def", "return", "input", "output", "qubit", "bit",
		"int", "float", "complex", "bool", "duration", "stretch",
		"const", "mutable", "readonly", "let", "in", "out", "inout",
	}
	
	for _, keyword := range keywords {
		if strings.EqualFold(identifier, keyword) {
			return true
		}
	}
	return false
}

// extractGateDefinitionsFromText extracts gate definitions from text when AST parsing fails
func (r *UndefinedIdentifierRule) extractGateDefinitionsFromText(content string, declaredIdentifiers map[string]bool) {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Look for gate definitions: "gate gatename ..."
		if strings.HasPrefix(line, "gate ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				gateName := parts[1]
				// Remove any parameter parentheses or qubit parameters
				if idx := strings.IndexAny(gateName, "( {"); idx != -1 {
					gateName = gateName[:idx]
				}
				if gateName != "" {
					declaredIdentifiers[gateName] = true
				}
			}
		}
	}
}