package ast

import (
	"github.com/orangekame3/qasmtools/lint/astutil"
	"github.com/orangekame3/qasmtools/parser"
)

// UnusedQubitRule implements QAS001 using AST-based analysis
type UnusedQubitRule struct {
	*ASTRuleBase
}

// NewUnusedQubitRule creates a new AST-based unused qubit rule
func NewUnusedQubitRule() ASTRule {
	return &UnusedQubitRule{
		ASTRuleBase: NewASTRuleBase("QAS001"),
	}
}

// CheckAST performs AST-based unused qubit analysis
func (r *UnusedQubitRule) CheckAST(program *parser.Program, ctx *CheckContext) []*Violation {
	var violations []*Violation

	// Find all qubit declarations
	declarations := astutil.FindDeclarations(program)
	
	// Get all identifier usages in the program
	usageMap := astutil.GetIdentifierUsages(program)

	// Check each qubit declaration for usage
	for _, qubitDecl := range declarations.Quantum {
		if qubitDecl.Type != "qubit" {
			continue
		}

		// Check if this qubit is used anywhere
		if !r.isQubitUsed(qubitDecl.Identifier, usageMap, declarations) {
			violation := r.NewViolationBuilder().
				WithMessage("Qubit '"+qubitDecl.Identifier+"' is declared but never used.").
				WithFile(ctx.File).
				WithNode(qubitDecl).
				WithNodeName(qubitDecl.Identifier).
				AsWarning().
				Build()
			violations = append(violations, violation)
		}
	}

	return violations
}

// isQubitUsed checks if a qubit is used anywhere in the program
func (r *UnusedQubitRule) isQubitUsed(qubitName string, usageMap map[string][]parser.Node, declarations *astutil.Declarations) bool {
	usages, exists := usageMap[qubitName]
	if !exists || len(usages) == 0 {
		return false
	}

	// Filter out the declaration itself - we only care about usage, not declaration
	var nonDeclarationUsages []parser.Node
	for _, usage := range usages {
		if !r.isDeclarationNode(usage, qubitName, declarations) {
			nonDeclarationUsages = append(nonDeclarationUsages, usage)
		}
	}

	return len(nonDeclarationUsages) > 0
}

// isDeclarationNode checks if a node represents the declaration of the given identifier
func (r *UnusedQubitRule) isDeclarationNode(node parser.Node, identifierName string, declarations *astutil.Declarations) bool {
	// Check if this node is the actual declaration
	for _, qubitDecl := range declarations.Quantum {
		if qubitDecl.Identifier == identifierName && node == qubitDecl {
			return true
		}
	}
	return false
}