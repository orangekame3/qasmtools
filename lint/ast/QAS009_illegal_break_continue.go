package ast

import (
	"regexp"
	"strings"

	"github.com/orangekame3/qasmtools/parser"
)

// QAS009IllegalBreakContinueRule implements QAS009 using text-based analysis within AST package
// This rule uses text-based analysis since BreakStatement/ContinueStatement are not in current AST
type QAS009IllegalBreakContinueRule struct {
	*ASTRuleBase
}

// NewQAS009IllegalBreakContinueRule creates a new QAS009 rule instance
func NewQAS009IllegalBreakContinueRule() *QAS009IllegalBreakContinueRule {
	return &QAS009IllegalBreakContinueRule{
		ASTRuleBase: NewASTRuleBase("QAS009"),
	}
}

// CheckAST performs text-based analysis within AST rule framework
func (r *QAS009IllegalBreakContinueRule) CheckAST(program *parser.Program, ctx *CheckContext) []*Violation {
	var violations []*Violation

	// Use text-based analysis on the content since break/continue are not in AST
	content := ctx.Content
	lines := strings.Split(content, "\n")

	// Track loop nesting depth and find break/continue statements
	loopDepth := 0

	for i, line := range lines {
		// Skip comments and empty lines
		if r.skipCommentAndEmptyLine(line) {
			continue
		}

		// Remove comments
		codeOnly := r.removeComments(line)

		// Track loop nesting changes
		loopDepth += r.countLoopChanges(codeOnly)

		// Find break/continue statements in this line
		breakContinueStatements := r.findBreakContinueStatements(codeOnly)

		for _, stmt := range breakContinueStatements {
			// Check if we're outside of any loop
			if loopDepth <= 0 {
				violation := r.NewViolationBuilder().
					WithMessage("'"+stmt.keyword+"' cannot be used outside of a loop.").
					WithFile(ctx.File).
					WithLine(i+1).
					WithColumn(stmt.column).
					WithNodeName(stmt.keyword).
					AsError().
					Build()
				violations = append(violations, violation)
			}
		}

		// Update loop depth after processing the line (for closing braces)
		loopDepth += r.countLoopEndChanges(codeOnly)
	}

	return violations
}

type breakContinueStatement struct {
	keyword string
	column  int
}

// skipCommentAndEmptyLine checks if a line should be skipped
func (r *QAS009IllegalBreakContinueRule) skipCommentAndEmptyLine(line string) bool {
	trimmed := strings.TrimSpace(line)
	return trimmed == "" || strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "/*")
}

// removeComments removes comments from a line
func (r *QAS009IllegalBreakContinueRule) removeComments(line string) string {
	// Remove single-line comments
	if idx := strings.Index(line, "//"); idx != -1 {
		line = line[:idx]
	}
	return strings.TrimSpace(line)
}

// countLoopChanges counts the increase in loop depth for a line (opening loops)
func (r *QAS009IllegalBreakContinueRule) countLoopChanges(line string) int {
	loopCount := 0

	// Patterns for loop starts
	patterns := []*regexp.Regexp{
		// for loops: for variable in range
		regexp.MustCompile(`\bfor\s+[a-zA-Z_][a-zA-Z0-9_]*\s+in\s+`),
		// while loops: while condition
		regexp.MustCompile(`\bwhile\s*\(`),
	}

	for _, pattern := range patterns {
		matches := pattern.FindAllStringSubmatch(line, -1)
		loopCount += len(matches)
	}

	return loopCount
}

// countLoopEndChanges counts the decrease in loop depth for a line (closing braces)
func (r *QAS009IllegalBreakContinueRule) countLoopEndChanges(line string) int {
	// Count closing braces that end loops
	closeBraces := strings.Count(line, "}")
	return -closeBraces
}

// findBreakContinueStatements finds all break/continue statements in a line
func (r *QAS009IllegalBreakContinueRule) findBreakContinueStatements(line string) []breakContinueStatement {
	var statements []breakContinueStatement

	// Patterns for break and continue statements
	patterns := []*regexp.Regexp{
		// break; statement
		regexp.MustCompile(`\b(break)\s*;`),
		// continue; statement
		regexp.MustCompile(`\b(continue)\s*;`),
	}

	for _, pattern := range patterns {
		matches := pattern.FindAllStringSubmatch(line, -1)
		indices := pattern.FindAllStringIndex(line, -1)

		for i, match := range matches {
			if len(match) >= 2 {
				keyword := match[1]
				column := indices[i][0] + 1 // Convert to 1-based indexing

				statements = append(statements, breakContinueStatement{
					keyword: keyword,
					column:  column,
				})
			}
		}
	}

	return statements
}