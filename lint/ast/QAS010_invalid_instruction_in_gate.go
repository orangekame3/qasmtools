package ast

import (
	"fmt"

	"github.com/orangekame3/qasmtools/lint/astutil"
	"github.com/orangekame3/qasmtools/parser"
)

// QAS010InvalidInstructionInGateRule detects non-unitary instructions in gate definitions
type QAS010InvalidInstructionInGateRule struct {
	*ASTRuleBase
}

// NewQAS010InvalidInstructionInGateRule creates a new QAS010 rule instance
func NewQAS010InvalidInstructionInGateRule() *QAS010InvalidInstructionInGateRule {
	return &QAS010InvalidInstructionInGateRule{
		ASTRuleBase: NewASTRuleBase("QAS010"),
	}
}

// ID returns the rule identifier
func (r *QAS010InvalidInstructionInGateRule) ID() string {
	return "QAS010"
}

// CheckAST performs AST-based analysis for invalid instructions in gate definitions
func (r *QAS010InvalidInstructionInGateRule) CheckAST(program *parser.Program, ctx *CheckContext) []*Violation {
	if program == nil {
		return nil
	}

	violations := make([]*Violation, 0)
	
	// Find all gate definitions
	gateDefinitions := astutil.FindNodesByType(program, (*parser.GateDefinition)(nil))
	
	for _, gateDef := range gateDefinitions {
		if gateDef.Body == nil {
			continue
		}
		
		// Check all statements in gate body for non-unitary instructions
		r.checkGateBodyForInvalidInstructions(gateDef.Body, gateDef.Name, ctx, &violations)
	}

	return violations
}

// checkGateBodyForInvalidInstructions recursively checks gate body statements for invalid instructions
func (r *QAS010InvalidInstructionInGateRule) checkGateBodyForInvalidInstructions(statements []parser.Statement, gateName string, ctx *CheckContext, violations *[]*Violation) {
	for _, stmt := range statements {
		r.checkStatementForInvalidInstructions(stmt, gateName, ctx, violations)
	}
}

// checkStatementForInvalidInstructions checks a single statement for invalid instructions in gate context
func (r *QAS010InvalidInstructionInGateRule) checkStatementForInvalidInstructions(stmt parser.Statement, gateName string, ctx *CheckContext, violations *[]*Violation) {
	if stmt == nil {
		return
	}

	switch s := stmt.(type) {
	case *parser.Measurement:
		// Measurements are non-unitary and not allowed in gate definitions
		message := fmt.Sprintf("Measurement instruction is not allowed within gate definition '%s' (non-unitary operation)", gateName)
		
		violation := r.NewViolationBuilder().
			WithMessage(message).
			WithFile(ctx.File).
			WithNode(s).
			WithNodeName("measure").
			AsError().
			Build()
		
		*violations = append(*violations, violation)
	
	case *parser.ClassicalDeclaration:
		// Classical declarations are not allowed in gate definitions
		message := fmt.Sprintf("Classical declaration '%s' is not allowed within gate definition '%s'", s.Identifier, gateName)
		
		violation := r.NewViolationBuilder().
			WithMessage(message).
			WithFile(ctx.File).
			WithNode(s).
			WithNodeName(s.Identifier).
			AsError().
			Build()
		
		*violations = append(*violations, violation)
	
	case *parser.QuantumDeclaration:
		// Quantum declarations are not allowed in gate definitions
		message := fmt.Sprintf("Quantum declaration '%s' is not allowed within gate definition '%s'", s.Identifier, gateName)
		
		violation := r.NewViolationBuilder().
			WithMessage(message).
			WithFile(ctx.File).
			WithNode(s).
			WithNodeName(s.Identifier).
			AsError().
			Build()
		
		*violations = append(*violations, violation)
	
	case *parser.IfStatement:
		// If statements can contain invalid instructions - check recursively
		if s.ThenBody != nil {
			r.checkGateBodyForInvalidInstructions(s.ThenBody, gateName, ctx, violations)
		}
		if s.ElseBody != nil {
			r.checkGateBodyForInvalidInstructions(s.ElseBody, gateName, ctx, violations)
		}
	
	case *parser.ForStatement:
		// For statements are not typically allowed in gates, but check body if present
		message := fmt.Sprintf("For loop is not typically allowed within gate definition '%s'", gateName)
		
		violation := r.NewViolationBuilder().
			WithMessage(message).
			WithFile(ctx.File).
			WithNode(s).
			WithNodeName("for").
			AsError().
			Build()
		
		*violations = append(*violations, violation)
		
		if s.Body != nil {
			r.checkGateBodyForInvalidInstructions(s.Body, gateName, ctx, violations)
		}
	
	case *parser.WhileStatement:
		// While statements are not typically allowed in gates
		message := fmt.Sprintf("While loop is not typically allowed within gate definition '%s'", gateName)
		
		violation := r.NewViolationBuilder().
			WithMessage(message).
			WithFile(ctx.File).
			WithNode(s).
			WithNodeName("while").
			AsError().
			Build()
		
		*violations = append(*violations, violation)
		
		if s.Body != nil {
			r.checkGateBodyForInvalidInstructions(s.Body, gateName, ctx, violations)
		}
	
	// These are typically allowed in gate definitions:
	case *parser.GateCall:
		// Gate calls are allowed (unitary operations)
	}
}