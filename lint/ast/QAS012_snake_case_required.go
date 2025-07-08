package ast

import (
	"fmt"

	"github.com/orangekame3/qasmtools/lint/astutil"
	"github.com/orangekame3/qasmtools/parser"
)

// QAS012SnakeCaseRequiredRule detects identifiers that don't follow snake_case convention
type QAS012SnakeCaseRequiredRule struct {
	*ASTRuleBase
}

// NewQAS012SnakeCaseRequiredRule creates a new QAS012 rule instance
func NewQAS012SnakeCaseRequiredRule() *QAS012SnakeCaseRequiredRule {
	return &QAS012SnakeCaseRequiredRule{
		ASTRuleBase: NewASTRuleBase("QAS012"),
	}
}

// ID returns the rule identifier
func (r *QAS012SnakeCaseRequiredRule) ID() string {
	return "QAS012"
}

// CheckAST performs AST-based analysis for snake_case naming convention
func (r *QAS012SnakeCaseRequiredRule) CheckAST(program *parser.Program, ctx *CheckContext) []*Violation {
	if program == nil {
		return nil
	}

	violations := make([]*Violation, 0)
	checkedIdentifiers := make(map[string]bool) // To avoid duplicate violations
	
	// Check all declarations for snake_case convention
	declarations := astutil.FindDeclarations(program)
	
	// Check quantum declarations
	for _, decl := range declarations.Quantum {
		if !astutil.IsSnakeCase(decl.Identifier) && !checkedIdentifiers[decl.Identifier] {
			message := fmt.Sprintf("Identifier '%s' should be written in snake_case.", decl.Identifier)
			
			violation := r.NewViolationBuilder().
				WithMessage(message).
				WithFile(ctx.File).
				WithNode(decl).
				WithNodeName(decl.Identifier).
				AsWarning().
				Build()
			
			violations = append(violations, violation)
			checkedIdentifiers[decl.Identifier] = true
		}
	}
	
	// Check classical declarations
	for _, decl := range declarations.Classical {
		if !astutil.IsSnakeCase(decl.Identifier) && !checkedIdentifiers[decl.Identifier] {
			message := fmt.Sprintf("Identifier '%s' should be written in snake_case.", decl.Identifier)
			
			violation := r.NewViolationBuilder().
				WithMessage(message).
				WithFile(ctx.File).
				WithNode(decl).
				WithNodeName(decl.Identifier).
				AsWarning().
				Build()
			
			violations = append(violations, violation)
			checkedIdentifiers[decl.Identifier] = true
		}
	}
	
	// Check gate declarations
	for _, decl := range declarations.Gates {
		if !astutil.IsSnakeCase(decl.Name) && !checkedIdentifiers[decl.Name] {
			message := fmt.Sprintf("Gate identifier '%s' should be written in snake_case.", decl.Name)
			
			violation := r.NewViolationBuilder().
				WithMessage(message).
				WithFile(ctx.File).
				WithNode(decl).
				WithNodeName(decl.Name).
				AsWarning().
				Build()
			
			violations = append(violations, violation)
			checkedIdentifiers[decl.Name] = true
		}
		
		// Check gate parameters for snake_case
		if decl.Parameters != nil {
			for _, param := range decl.Parameters {
				if !astutil.IsSnakeCase(param.Name) && !checkedIdentifiers[param.Name] {
					message := fmt.Sprintf("Gate parameter '%s' should be written in snake_case.", param.Name)
					
					violation := r.NewViolationBuilder().
						WithMessage(message).
						WithFile(ctx.File).
						WithNode(&param).
						WithNodeName(param.Name).
						AsWarning().
						Build()
					
					violations = append(violations, violation)
					checkedIdentifiers[param.Name] = true
				}
			}
		}
		
		// Check gate qubits for snake_case
		if decl.Qubits != nil {
			for _, qubit := range decl.Qubits {
				if !astutil.IsSnakeCase(qubit.Name) && !checkedIdentifiers[qubit.Name] {
					message := fmt.Sprintf("Gate qubit parameter '%s' should be written in snake_case.", qubit.Name)
					
					violation := r.NewViolationBuilder().
						WithMessage(message).
						WithFile(ctx.File).
						WithNode(&qubit).
						WithNodeName(qubit.Name).
						AsWarning().
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