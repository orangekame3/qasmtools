package ast

import (
	"github.com/orangekame3/qasmtools/lint/astutil"
	"github.com/orangekame3/qasmtools/parser"
)

// ConstantMeasuredBitRule implements QAS003 using AST-based analysis
type ConstantMeasuredBitRule struct {
	*ASTRuleBase
}

// NewConstantMeasuredBitRule creates a new AST-based constant measured bit rule
func NewConstantMeasuredBitRule() ASTRule {
	return &ConstantMeasuredBitRule{
		ASTRuleBase: NewASTRuleBase("QAS003"),
	}
}

// CheckAST performs AST-based constant measured bit analysis
func (r *ConstantMeasuredBitRule) CheckAST(program *parser.Program, ctx *CheckContext) []*Violation {
	var violations []*Violation

	// Find all measurements
	measurements := astutil.FindNodesByType(program, (*parser.Measurement)(nil))
	
	// For each measurement, check if the qubit has been affected by any gates
	for _, measurement := range measurements {
		if qubitExpr := measurement.Qubit; qubitExpr != nil {
			qubitName := r.extractQubitName(qubitExpr)
			if qubitName != "" && !r.isQubitAffectedByGates(qubitName, program) {
				violation := r.NewViolationBuilder().
					WithMessage("Measuring qubit '"+qubitName+"' that has no gates applied. The result will always be |0⟩.").
					WithFile(ctx.File).
					WithNode(measurement).
					WithNodeName(qubitName).
					AsWarning().
					Build()
				violations = append(violations, violation)
			}
		}
	}

	return violations
}

// extractQubitName extracts the qubit name from a measurement expression
func (r *ConstantMeasuredBitRule) extractQubitName(expr parser.Expression) string {
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

// isQubitAffectedByGates checks if a qubit has been affected by any gate operations
func (r *ConstantMeasuredBitRule) isQubitAffectedByGates(qubitName string, program *parser.Program) bool {
	var affected bool

	astutil.VisitAllNodes(program, func(node parser.Node) {
		if affected {
			return
		}

		if gateCall, ok := node.(*parser.GateCall); ok {
			// Check if this gate call affects the qubit
			for _, qubit := range gateCall.Qubits {
				targetQubitName := r.extractQubitName(qubit)
				if targetQubitName == qubitName {
					affected = true
					return
				}
			}
		}
	})

	return affected
}