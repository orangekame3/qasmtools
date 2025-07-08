package ast

import (
	"github.com/orangekame3/qasmtools/lint/astutil"
	"github.com/orangekame3/qasmtools/parser"
)

// NamingConventionViolationRule implements QAS005 using AST-based analysis
type NamingConventionViolationRule struct {
	*ASTRuleBase
}

// NewNamingConventionViolationRule creates a new AST-based naming convention rule
func NewNamingConventionViolationRule() ASTRule {
	return &NamingConventionViolationRule{
		ASTRuleBase: NewASTRuleBase("QAS005"),
	}
}

// CheckAST performs AST-based naming convention analysis
func (r *NamingConventionViolationRule) CheckAST(program *parser.Program, ctx *CheckContext) []*Violation {
	var violations []*Violation

	// Find all declarations
	declarations := astutil.FindDeclarations(program)

	// Check quantum declarations
	for _, qubitDecl := range declarations.Quantum {
		if !astutil.IsSnakeCase(qubitDecl.Identifier) {
			violation := r.NewViolationBuilder().
				WithMessage("Identifier '"+qubitDecl.Identifier+"' violates naming conventions. Follow pattern: ^[a-z][a-zA-Z0-9_]*$.").
				WithFile(ctx.File).
				WithNode(qubitDecl).
				WithNodeName(qubitDecl.Identifier).
				AsWarning().
				Build()
			violations = append(violations, violation)
		}
	}

	// Check classical declarations
	for _, classicalDecl := range declarations.Classical {
		if !astutil.IsSnakeCase(classicalDecl.Identifier) {
			violation := r.NewViolationBuilder().
				WithMessage("Identifier '"+classicalDecl.Identifier+"' violates naming conventions. Follow pattern: ^[a-z][a-zA-Z0-9_]*$.").
				WithFile(ctx.File).
				WithNode(classicalDecl).
				WithNodeName(classicalDecl.Identifier).
				AsWarning().
				Build()
			violations = append(violations, violation)
		}
	}

	// Check gate definitions
	for _, gateDecl := range declarations.Gates {
		if !astutil.IsSnakeCase(gateDecl.Name) {
			violation := r.NewViolationBuilder().
				WithMessage("Identifier '"+gateDecl.Name+"' violates naming conventions. Follow pattern: ^[a-z][a-zA-Z0-9_]*$.").
				WithFile(ctx.File).
				WithNode(gateDecl).
				WithNodeName(gateDecl.Name).
				AsWarning().
				Build()
			violations = append(violations, violation)
		}

		// Check gate parameters
		for _, param := range gateDecl.Parameters {
			if !astutil.IsSnakeCase(param.Name) {
				violation := r.NewViolationBuilder().
					WithMessage("Parameter '"+param.Name+"' violates naming conventions. Follow pattern: ^[a-z][a-zA-Z0-9_]*$.").
					WithFile(ctx.File).
					WithNode(&param).
					WithNodeName(param.Name).
					AsWarning().
					Build()
				violations = append(violations, violation)
			}
		}

		// Check gate qubits
		for _, qubit := range gateDecl.Qubits {
			if !astutil.IsSnakeCase(qubit.Name) {
				violation := r.NewViolationBuilder().
					WithMessage("Qubit parameter '"+qubit.Name+"' violates naming conventions. Follow pattern: ^[a-z][a-zA-Z0-9_]*$.").
					WithFile(ctx.File).
					WithNode(&qubit).
					WithNodeName(qubit.Name).
					AsWarning().
					Build()
				violations = append(violations, violation)
			}
		}
	}

	return violations
}