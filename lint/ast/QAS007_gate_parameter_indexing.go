package ast

import (
	"fmt"

	"github.com/orangekame3/qasmtools/lint/astutil"
	"github.com/orangekame3/qasmtools/parser"
)

// QAS007GateParameterIndexingRule detects illegal index access on gate parameters
type QAS007GateParameterIndexingRule struct {
	*ASTRuleBase
}

// NewQAS007GateParameterIndexingRule creates a new QAS007 rule instance
func NewQAS007GateParameterIndexingRule() *QAS007GateParameterIndexingRule {
	return &QAS007GateParameterIndexingRule{
		ASTRuleBase: NewASTRuleBase("QAS007"),
	}
}

// ID returns the rule identifier
func (r *QAS007GateParameterIndexingRule) ID() string {
	return "QAS007"
}

// CheckAST performs AST-based analysis for gate parameter indexing violations
func (r *QAS007GateParameterIndexingRule) CheckAST(program *parser.Program, ctx *CheckContext) []*Violation {
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
		
		// Collect gate parameters
		gateParams := make(map[string]bool)
		if gateDef.Qubits != nil {
			for _, param := range gateDef.Qubits {
				gateParams[param.Name] = true
			}
		}
		
		// Check all statements in gate body for illegal index access
		for _, stmt := range gateDef.Body {
			r.checkStatementForIndexAccess(stmt, gateParams, ctx, &violations)
		}
	}

	return violations
}

// checkStatementForIndexAccess recursively checks statements for illegal index access
func (r *QAS007GateParameterIndexingRule) checkStatementForIndexAccess(stmt parser.Statement, gateParams map[string]bool, ctx *CheckContext, violations *[]*Violation) {
	if stmt == nil {
		return
	}

	switch s := stmt.(type) {
	case *parser.GateCall:
		// Check qubits in gate call
		if s.Qubits != nil {
			for _, qubit := range s.Qubits {
				if idxId, ok := qubit.(*parser.IndexedIdentifier); ok {
					// Check if the indexed identifier is a gate parameter
					if gateParams[idxId.Name] {
						message := fmt.Sprintf("Cannot perform index access on gate parameter '%s'", idxId.Name)
						
						violation := r.NewViolationBuilder().
							WithMessage(message).
							WithFile(ctx.File).
							WithNode(idxId).
							WithNodeName(idxId.Name).
							AsError().
							Build()
						
						*violations = append(*violations, violation)
					}
				}
			}
		}
	
	case *parser.Measurement:
		// Check qubit being measured
		if idxId, ok := s.Qubit.(*parser.IndexedIdentifier); ok {
			if gateParams[idxId.Name] {
				message := fmt.Sprintf("Cannot perform index access on gate parameter '%s'", idxId.Name)
				
				violation := r.NewViolationBuilder().
					WithMessage(message).
					WithFile(ctx.File).
					WithNode(idxId).
					WithNodeName(idxId.Name).
					AsError().
					Build()
				
				*violations = append(*violations, violation)
			}
		}
		
		// Check target being assigned to
		if idxId, ok := s.Target.(*parser.IndexedIdentifier); ok {
			if gateParams[idxId.Name] {
				message := fmt.Sprintf("Cannot perform index access on gate parameter '%s'", idxId.Name)
				
				violation := r.NewViolationBuilder().
					WithMessage(message).
					WithFile(ctx.File).
					WithNode(idxId).
					WithNodeName(idxId.Name).
					AsError().
					Build()
				
				*violations = append(*violations, violation)
			}
		}
	
	case *parser.IfStatement:
		// Check condition for index access
		r.checkExpressionForIndexAccess(s.Condition, gateParams, ctx, violations)
		
		// Check then body statements
		if s.ThenBody != nil {
			for _, thenStmt := range s.ThenBody {
				r.checkStatementForIndexAccess(thenStmt, gateParams, ctx, violations)
			}
		}
		
		// Check else body statements
		if s.ElseBody != nil {
			for _, elseStmt := range s.ElseBody {
				r.checkStatementForIndexAccess(elseStmt, gateParams, ctx, violations)
			}
		}
	}
}

// checkExpressionForIndexAccess recursively checks expressions for illegal index access
func (r *QAS007GateParameterIndexingRule) checkExpressionForIndexAccess(expr parser.Expression, gateParams map[string]bool, ctx *CheckContext, violations *[]*Violation) {
	if expr == nil {
		return
	}

	switch e := expr.(type) {
	case *parser.IndexedIdentifier:
		if gateParams[e.Name] {
			message := fmt.Sprintf("Cannot perform index access on gate parameter '%s'", e.Name)
			
			violation := r.NewViolationBuilder().
				WithMessage(message).
				WithFile(ctx.File).
				WithNode(e).
				WithNodeName(e.Name).
				AsError().
				Build()
			
			*violations = append(*violations, violation)
		}
	
	case *parser.BinaryExpression:
		r.checkExpressionForIndexAccess(e.Left, gateParams, ctx, violations)
		r.checkExpressionForIndexAccess(e.Right, gateParams, ctx, violations)
	
	case *parser.UnaryExpression:
		r.checkExpressionForIndexAccess(e.Operand, gateParams, ctx, violations)
	}
}