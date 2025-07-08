// Package ast provides AST-based lint rule implementations
package ast

import (
	"github.com/orangekame3/qasmtools/parser"
)

// Severity represents the severity level of a lint rule violation
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
	SeverityInfo    Severity = "info"
)

// Rule represents minimal rule information for violations
type Rule struct {
	ID string
}

// Violation represents a rule violation from AST analysis
type Violation struct {
	Rule     *Rule
	Message  string
	File     string
	Line     int
	Column   int
	Severity Severity
	NodeName string
}

// CheckContext provides context for AST rule checking
type CheckContext struct {
	File     string
	Content  string
	Program  *parser.Program
	UsageMap map[string][]parser.Node
}

// ASTRule interface for AST-based lint rules
type ASTRule interface {
	// ID returns the unique identifier for this rule (e.g., "QAS001")
	ID() string
	
	// CheckAST performs the rule check on the AST and returns any violations
	CheckAST(program *parser.Program, ctx *CheckContext) []*Violation
}

// ASTRuleBase provides common functionality for AST-based rules
type ASTRuleBase struct {
	ruleID string
}

// NewASTRuleBase creates a new base for AST rules
func NewASTRuleBase(ruleID string) *ASTRuleBase {
	return &ASTRuleBase{
		ruleID: ruleID,
	}
}

// ID returns the rule ID
func (r *ASTRuleBase) ID() string {
	return r.ruleID
}

// CreateViolation is a helper method for creating violations with consistent formatting
func (r *ASTRuleBase) CreateViolation(message, file string, line, column int, nodeName string, severity Severity) *Violation {
	// Create a minimal rule for the violation to avoid nil pointer issues
	rule := &Rule{
		ID: r.ruleID,
	}
	return &Violation{
		Rule:     rule,
		Message:  message,
		File:     file,
		Line:     line,
		Column:   column,
		NodeName: nodeName,
		Severity: severity,
	}
}

// NewViolationBuilder creates a new violation builder for this rule
func (r *ASTRuleBase) NewViolationBuilder() *ViolationBuilder {
	return &ViolationBuilder{
		rule:     r,
		severity: SeverityWarning, // Default severity
	}
}

// ViolationBuilder provides a fluent interface for building violations from AST rules
type ViolationBuilder struct {
	rule     *ASTRuleBase
	message  string
	file     string
	line     int
	column   int
	nodeName string
	severity Severity
}

// WithMessage sets the violation message
func (vb *ViolationBuilder) WithMessage(message string) *ViolationBuilder {
	vb.message = message
	return vb
}

// WithFile sets the file path
func (vb *ViolationBuilder) WithFile(file string) *ViolationBuilder {
	vb.file = file
	return vb
}

// WithPosition sets the line and column from a Position
func (vb *ViolationBuilder) WithPosition(pos parser.Position) *ViolationBuilder {
	vb.line = pos.Line
	vb.column = pos.Column
	return vb
}

// WithNode sets position information from an AST node
func (vb *ViolationBuilder) WithNode(node parser.Node) *ViolationBuilder {
	pos := node.Pos()
	vb.line = pos.Line
	vb.column = pos.Column
	return vb
}

// WithLine sets the line number
func (vb *ViolationBuilder) WithLine(line int) *ViolationBuilder {
	vb.line = line
	return vb
}

// WithColumn sets the column number
func (vb *ViolationBuilder) WithColumn(column int) *ViolationBuilder {
	vb.column = column
	return vb
}

// WithNodeName sets the node name
func (vb *ViolationBuilder) WithNodeName(nodeName string) *ViolationBuilder {
	vb.nodeName = nodeName
	return vb
}

// WithSeverity sets the severity level
func (vb *ViolationBuilder) WithSeverity(severity Severity) *ViolationBuilder {
	vb.severity = severity
	return vb
}

// AsError sets severity to error
func (vb *ViolationBuilder) AsError() *ViolationBuilder {
	vb.severity = SeverityError
	return vb
}

// AsWarning sets severity to warning
func (vb *ViolationBuilder) AsWarning() *ViolationBuilder {
	vb.severity = SeverityWarning
	return vb
}

// AsInfo sets severity to info
func (vb *ViolationBuilder) AsInfo() *ViolationBuilder {
	vb.severity = SeverityInfo
	return vb
}

// Build creates the violation
func (vb *ViolationBuilder) Build() *Violation {
	return vb.rule.CreateViolation(vb.message, vb.file, vb.line, vb.column, vb.nodeName, vb.severity)
}