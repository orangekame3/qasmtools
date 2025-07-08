package ast

import (
	"fmt"
	"strings"

	"github.com/orangekame3/qasmtools/lint/astutil"
	"github.com/orangekame3/qasmtools/parser"
)

// QAS011ReservedPrefixUsageRule detects usage of reserved prefix __ in identifiers
type QAS011ReservedPrefixUsageRule struct {
	*ASTRuleBase
}

// NewQAS011ReservedPrefixUsageRule creates a new QAS011 rule instance
func NewQAS011ReservedPrefixUsageRule() *QAS011ReservedPrefixUsageRule {
	return &QAS011ReservedPrefixUsageRule{
		ASTRuleBase: NewASTRuleBase("QAS011"),
	}
}

// ID returns the rule identifier
func (r *QAS011ReservedPrefixUsageRule) ID() string {
	return "QAS011"
}

// CheckAST performs AST-based analysis for reserved prefix usage
func (r *QAS011ReservedPrefixUsageRule) CheckAST(program *parser.Program, ctx *CheckContext) []*Violation {
	if program == nil {
		return nil
	}

	violations := make([]*Violation, 0)
	checkedIdentifiers := make(map[string]bool) // To avoid duplicate violations
	
	// Check all declarations for reserved prefix usage
	declarations := astutil.FindDeclarations(program)
	
	// Check quantum declarations
	for _, decl := range declarations.Quantum {
		if r.hasReservedPrefix(decl.Identifier) && !checkedIdentifiers[decl.Identifier] {
			message := fmt.Sprintf("Identifier '%s' uses reserved prefix '__'.", decl.Identifier)
			
			violation := r.NewViolationBuilder().
				WithMessage(message).
				WithFile(ctx.File).
				WithNode(decl).
				WithNodeName(decl.Identifier).
				AsError().
				Build()
			
			violations = append(violations, violation)
			checkedIdentifiers[decl.Identifier] = true
		}
	}
	
	// Check classical declarations
	for _, decl := range declarations.Classical {
		if r.hasReservedPrefix(decl.Identifier) && !checkedIdentifiers[decl.Identifier] {
			message := fmt.Sprintf("Identifier '%s' uses reserved prefix '__'.", decl.Identifier)
			
			violation := r.NewViolationBuilder().
				WithMessage(message).
				WithFile(ctx.File).
				WithNode(decl).
				WithNodeName(decl.Identifier).
				AsError().
				Build()
			
			violations = append(violations, violation)
			checkedIdentifiers[decl.Identifier] = true
		}
	}
	
	// Check gate declarations
	for _, decl := range declarations.Gates {
		if r.hasReservedPrefix(decl.Name) && !checkedIdentifiers[decl.Name] {
			message := fmt.Sprintf("Gate identifier '%s' uses reserved prefix '__'.", decl.Name)
			
			violation := r.NewViolationBuilder().
				WithMessage(message).
				WithFile(ctx.File).
				WithNode(decl).
				WithNodeName(decl.Name).
				AsError().
				Build()
			
			violations = append(violations, violation)
			checkedIdentifiers[decl.Name] = true
		}
		
		// Check gate parameters for reserved prefix
		if decl.Parameters != nil {
			for _, param := range decl.Parameters {
				if r.hasReservedPrefix(param.Name) && !checkedIdentifiers[param.Name] {
					message := fmt.Sprintf("Gate parameter '%s' uses reserved prefix '__'.", param.Name)
					
					violation := r.NewViolationBuilder().
						WithMessage(message).
						WithFile(ctx.File).
						WithNode(&param).
						WithNodeName(param.Name).
						AsError().
						Build()
					
					violations = append(violations, violation)
					checkedIdentifiers[param.Name] = true
				}
			}
		}
		
		// Check gate qubits for reserved prefix
		if decl.Qubits != nil {
			for _, qubit := range decl.Qubits {
				if r.hasReservedPrefix(qubit.Name) && !checkedIdentifiers[qubit.Name] {
					message := fmt.Sprintf("Gate qubit parameter '%s' uses reserved prefix '__'.", qubit.Name)
					
					violation := r.NewViolationBuilder().
						WithMessage(message).
						WithFile(ctx.File).
						WithNode(&qubit).
						WithNodeName(qubit.Name).
						AsError().
						Build()
					
					violations = append(violations, violation)
					checkedIdentifiers[qubit.Name] = true
				}
			}
		}
	}
	
	// Note: Function definitions are not currently supported in the AST
	// This check can be added when FunctionDefinition is implemented

	return violations
}

// hasReservedPrefix checks if an identifier has the reserved prefix __
func (r *QAS011ReservedPrefixUsageRule) hasReservedPrefix(identifier string) bool {
	return strings.HasPrefix(identifier, "__")
}