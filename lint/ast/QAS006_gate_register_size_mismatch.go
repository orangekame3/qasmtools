package ast

import (
	"fmt"

	"github.com/orangekame3/qasmtools/lint/astutil"
	"github.com/orangekame3/qasmtools/parser"
)

// QAS006GateRegisterSizeMismatchRule detects gate register size mismatches
type QAS006GateRegisterSizeMismatchRule struct {
	*ASTRuleBase
}

// NewQAS006GateRegisterSizeMismatchRule creates a new QAS006 rule instance
func NewQAS006GateRegisterSizeMismatchRule() *QAS006GateRegisterSizeMismatchRule {
	return &QAS006GateRegisterSizeMismatchRule{
		ASTRuleBase: NewASTRuleBase("QAS006"),
	}
}

// ID returns the rule identifier
func (r *QAS006GateRegisterSizeMismatchRule) ID() string {
	return "QAS006"
}

// CheckAST performs AST-based analysis for gate register size mismatches
func (r *QAS006GateRegisterSizeMismatchRule) CheckAST(program *parser.Program, ctx *CheckContext) []*Violation {
	if program == nil {
		return nil
	}

	violations := make([]*Violation, 0)
	
	// Build register size map from declarations
	registerSizes := make(map[string]int)
	
	// Find all quantum and classical declarations
	declarations := astutil.FindDeclarations(program)
	for _, decl := range declarations.Quantum {
		if decl.Size != nil {
			if intLit, ok := decl.Size.(*parser.IntegerLiteral); ok {
				registerSizes[decl.Identifier] = int(intLit.Value)
			}
		} else {
			// Single qubit has size 1
			registerSizes[decl.Identifier] = 1
		}
	}
	
	for _, decl := range declarations.Classical {
		if decl.Size != nil {
			if intLit, ok := decl.Size.(*parser.IntegerLiteral); ok {
				registerSizes[decl.Identifier] = int(intLit.Value)
			}
		} else {
			// Single bit has size 1
			registerSizes[decl.Identifier] = 1
		}
	}

	// Find all gate calls and check register size consistency
	gateCalls := astutil.FindNodesByType(program, (*parser.GateCall)(nil))
	for _, gateCall := range gateCalls {
		if gateCall.Qubits == nil || len(gateCall.Qubits) < 2 {
			continue // Need at least 2 qubits to check size mismatch
		}

		// Get sizes for all qubits in the gate call
		qubitSizes := make([]int, 0)
		qubitNames := make([]string, 0)
		
		for _, qubit := range gateCall.Qubits {
			var size int
			var name string
			
			if id, ok := qubit.(*parser.Identifier); ok {
				name = id.Name
				if declSize, exists := registerSizes[name]; exists {
					size = declSize
				} else {
					size = 1 // Default to single qubit
				}
			} else if idxId, ok := qubit.(*parser.IndexedIdentifier); ok {
				name = idxId.Name
				// For indexed access, size is 1 (single element)
				size = 1
			} else {
				continue // Skip unknown qubit expressions
			}
			
			qubitSizes = append(qubitSizes, size)
			qubitNames = append(qubitNames, name)
		}

		// Check if all sizes are the same (or compatible)
		if len(qubitSizes) > 1 {
			baseSize := qubitSizes[0]
			hasMismatch := false
			
			for i := 1; i < len(qubitSizes); i++ {
				if qubitSizes[i] != baseSize {
					// Size mismatch found
					hasMismatch = true
					break
				}
			}
			
			if hasMismatch {
				message := fmt.Sprintf("Register lengths passed to gate '%s' do not match.", gateCall.Name)
				
				violation := r.NewViolationBuilder().
					WithMessage(message).
					WithFile(ctx.File).
					WithNode(gateCall).
					WithNodeName(gateCall.Name).
					AsError().
					Build()
				
				violations = append(violations, violation)
			}
		}
	}

	return violations
}