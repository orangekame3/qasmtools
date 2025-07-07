package lint

import (
	"os"
	"strings"

	"github.com/orangekame3/qasmtools/parser"
)

// BaseChecker provides common functionality for all lint checkers
type BaseChecker struct {
	ruleID string
}

// NewBaseChecker creates a new BaseChecker with the given rule ID
func NewBaseChecker(ruleID string) *BaseChecker {
	return &BaseChecker{
		ruleID: ruleID,
	}
}

// GetRuleID returns the rule ID for this checker
func (b *BaseChecker) GetRuleID() string {
	return b.ruleID
}

// Check implements the RuleChecker interface for AST-based checking
// This is the standard implementation that delegates to CheckProgram
func (b *BaseChecker) Check(node parser.Node, context *CheckContext) []*Violation {
	// This method is required by RuleChecker but not used for program-level analysis
	// All checkers currently use text-based analysis due to AST parsing issues
	return nil
}

// CheckProgram implements the ProgramChecker interface for program-level checking
// This delegates to CheckFile for text-based analysis
func (b *BaseChecker) CheckProgram(context *CheckContext) []*Violation {
	// Use text-based analysis due to AST parsing issues
	return b.CheckFile(context)
}

// CheckFile provides common file processing logic for text-based checkers
func (b *BaseChecker) CheckFile(context *CheckContext) []*Violation {
	var violations []*Violation

	// Get content for text-based analysis
	content, err := b.getContent(context)
	if err != nil {
		return violations
	}

	// Split into lines for processing
	lines := strings.Split(content, "\n")

	// Process each line using the template method pattern
	for i, line := range lines {
		// Skip comments and empty lines using shared utility
		if SkipCommentAndEmptyLine(line) {
			continue
		}

		// Call the specific checker's line processing logic
		// This will be overridden by concrete checkers
		lineViolations := b.ProcessLine(line, i+1, context)
		violations = append(violations, lineViolations...)
	}

	return violations
}

// ProcessLine is a template method that should be overridden by concrete checkers
// Default implementation returns no violations
func (b *BaseChecker) ProcessLine(line string, lineNum int, context *CheckContext) []*Violation {
	// Override this method in concrete checkers
	return nil
}

// getContent returns the content for analysis, preferring provided content over file reading
func (b *BaseChecker) getContent(context *CheckContext) (string, error) {
	if context.Content != "" {
		return context.Content, nil
	}

	content, err := os.ReadFile(context.File)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// CreateViolation is a helper method for creating violations with consistent formatting
func (b *BaseChecker) CreateViolation(message, file string, line, column int, nodeName string, severity Severity) *Violation {
	// Create a minimal rule for the violation to avoid nil pointer issues
	rule := &Rule{
		ID: b.ruleID,
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

// CreateErrorViolation creates an error-level violation
func (b *BaseChecker) CreateErrorViolation(message, file string, line, column int, nodeName string) *Violation {
	return b.CreateViolation(message, file, line, column, nodeName, SeverityError)
}

// CreateWarningViolation creates a warning-level violation
func (b *BaseChecker) CreateWarningViolation(message, file string, line, column int, nodeName string) *Violation {
	return b.CreateViolation(message, file, line, column, nodeName, SeverityWarning)
}

// CreateInfoViolation creates an info-level violation
func (b *BaseChecker) CreateInfoViolation(message, file string, line, column int, nodeName string) *Violation {
	return b.CreateViolation(message, file, line, column, nodeName, SeverityInfo)
}

// ViolationBuilder provides a fluent interface for building violations
type ViolationBuilder struct {
	checker  *BaseChecker
	message  string
	file     string
	line     int
	column   int
	nodeName string
	severity Severity
}

// NewViolationBuilder creates a new violation builder
func (b *BaseChecker) NewViolationBuilder() *ViolationBuilder {
	return &ViolationBuilder{
		checker:  b,
		severity: SeverityWarning, // Default severity
	}
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

// WithPosition sets the line and column
func (vb *ViolationBuilder) WithPosition(line, column int) *ViolationBuilder {
	vb.line = line
	vb.column = column
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
	return vb.checker.CreateViolation(vb.message, vb.file, vb.line, vb.column, vb.nodeName, vb.severity)
}

// LineProcessor provides an interface for line-by-line processing
type LineProcessor interface {
	ProcessLine(line string, lineNum int, context *CheckContext) []*Violation
}

// ProcessFileLines processes file lines using a LineProcessor
func ProcessFileLines(context *CheckContext, processor LineProcessor) []*Violation {
	var violations []*Violation

	// Get content for analysis
	content, err := context.GetContent()
	if err != nil {
		return violations
	}

	// Split into lines for processing
	lines := strings.Split(content, "\n")

	// Process each line
	for i, line := range lines {
		// Skip comments and empty lines
		if SkipCommentAndEmptyLine(line) {
			continue
		}

		// Process the line
		lineViolations := processor.ProcessLine(line, i+1, context)
		violations = append(violations, lineViolations...)
	}

	return violations
}
