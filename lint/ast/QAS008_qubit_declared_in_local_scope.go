package ast

import (
	"fmt"

	"github.com/orangekame3/qasmtools/parser"
)

// QAS008QubitDeclaredInLocalScopeRule detects qubits declared in non-global scope
type QAS008QubitDeclaredInLocalScopeRule struct {
	*ASTRuleBase
}

// NewQAS008QubitDeclaredInLocalScopeRule creates a new QAS008 rule instance
func NewQAS008QubitDeclaredInLocalScopeRule() *QAS008QubitDeclaredInLocalScopeRule {
	return &QAS008QubitDeclaredInLocalScopeRule{
		ASTRuleBase: NewASTRuleBase("QAS008"),
	}
}

// ID returns the rule identifier
func (r *QAS008QubitDeclaredInLocalScopeRule) ID() string {
	return "QAS008"
}

// CheckAST performs AST-based analysis for qubits declared in local scope
func (r *QAS008QubitDeclaredInLocalScopeRule) CheckAST(program *parser.Program, ctx *CheckContext) []*Violation {
	if program == nil {
		return nil
	}

	violations := make([]*Violation, 0)
	
	// Check all statements in the program for local scope violations
	r.checkStatementsForLocalQubitDeclarations(program.Statements, false, ctx, &violations)

	return violations
}

// checkStatementsForLocalQubitDeclarations recursively checks statements for qubit declarations in local scopes
func (r *QAS008QubitDeclaredInLocalScopeRule) checkStatementsForLocalQubitDeclarations(statements []parser.Statement, inLocalScope bool, ctx *CheckContext, violations *[]*Violation) {
	for _, stmt := range statements {
		r.checkStatementForLocalQubitDeclaration(stmt, inLocalScope, ctx, violations)
	}
}

// checkStatementForLocalQubitDeclaration checks a single statement for local qubit declarations
func (r *QAS008QubitDeclaredInLocalScopeRule) checkStatementForLocalQubitDeclaration(stmt parser.Statement, inLocalScope bool, ctx *CheckContext, violations *[]*Violation) {
	if stmt == nil {
		return
	}

	switch s := stmt.(type) {
	case *parser.QuantumDeclaration:
		// If we're in a local scope and this is a qubit declaration, it's a violation
		if inLocalScope {
			message := fmt.Sprintf("Qubit '%s' can only be declared in global scope", s.Identifier)
			
			violation := r.NewViolationBuilder().
				WithMessage(message).
				WithFile(ctx.File).
				WithNode(s).
				WithNodeName(s.Identifier).
				AsError().
				Build()
			
			*violations = append(*violations, violation)
		}
	
	case *parser.GateDefinition:
		// Gate definitions create a local scope
		if s.Body != nil {
			r.checkStatementsForLocalQubitDeclarations(s.Body, true, ctx, violations)
		}
	
	case *parser.IfStatement:
		// If statements create local scopes
		if s.ThenBody != nil {
			r.checkStatementsForLocalQubitDeclarations(s.ThenBody, true, ctx, violations)
		}
		if s.ElseBody != nil {
			r.checkStatementsForLocalQubitDeclarations(s.ElseBody, true, ctx, violations)
		}
	
	case *parser.ForStatement:
		// For loops create local scopes
		if s.Body != nil {
			r.checkStatementsForLocalQubitDeclarations(s.Body, true, ctx, violations)
		}
	
	case *parser.WhileStatement:
		// While loops create local scopes
		if s.Body != nil {
			r.checkStatementsForLocalQubitDeclarations(s.Body, true, ctx, violations)
		}
	}
}