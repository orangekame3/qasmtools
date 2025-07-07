package lint

import (
	"regexp"
	"strings"

	"github.com/orangekame3/qasmtools/parser"
)

// IllegalBreakContinueChecker is the new implementation using BaseChecker framework
type IllegalBreakContinueChecker struct {
	*BaseChecker
}

// NewIllegalBreakContinueChecker creates a new IllegalBreakContinueChecker
func NewIllegalBreakContinueChecker() *IllegalBreakContinueChecker {
	return &IllegalBreakContinueChecker{
		BaseChecker: NewBaseChecker("QAS009"),
	}
}

// CheckFile performs file-level break/continue analysis
func (c *IllegalBreakContinueChecker) CheckFile(context *CheckContext) []*Violation {
	var violations []*Violation

	// Get content using BaseChecker method
	content, err := c.getContent(context)
	if err != nil {
		return violations
	}

	lines := strings.Split(content, "\n")

	// Track loop nesting depth and find break/continue statements
	loopDepth := 0

	for i, line := range lines {
		// Skip comments and empty lines using shared utility
		if SkipCommentAndEmptyLine(line) {
			continue
		}

		// Remove comments using shared utility
		codeOnly := RemoveComments(line)

		// Track loop nesting changes
		loopDepth += c.countLoopChanges(codeOnly)

		// Find break/continue statements in this line
		breakContinueStatements := c.findBreakContinueStatements(codeOnly)

		for _, stmt := range breakContinueStatements {
			// Check if we're outside of any loop
			if loopDepth <= 0 {
				violation := c.NewViolationBuilder().
					WithMessage("'"+stmt.keyword+"' cannot be used outside of a loop.").
					WithFile(context.File).
					WithPosition(i+1, stmt.column).
					WithNodeName(stmt.keyword).
					AsError().
					Build()
				violations = append(violations, violation)
			}
		}

		// Update loop depth after processing the line (for closing braces)
		loopDepth += c.countLoopEndChanges(codeOnly)
	}

	return violations
}

type breakContinueStatement struct {
	keyword string
	column  int
}

// countLoopChanges counts the increase in loop depth for a line (opening loops)
func (c *IllegalBreakContinueChecker) countLoopChanges(line string) int {
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
func (c *IllegalBreakContinueChecker) countLoopEndChanges(line string) int {
	// Count closing braces that end loops
	// This is a simplified approach - we assume each closing brace ends one scope level
	closeBraces := strings.Count(line, "}")
	return -closeBraces
}

// findBreakContinueStatements finds all break/continue statements in a line
func (c *IllegalBreakContinueChecker) findBreakContinueStatements(line string) []breakContinueStatement {
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

// Check implements RuleChecker interface (required but delegates to CheckProgram)
func (c *IllegalBreakContinueChecker) Check(node parser.Node, context *CheckContext) []*Violation {
	return nil
}

// CheckProgram implements ProgramChecker interface
func (c *IllegalBreakContinueChecker) CheckProgram(context *CheckContext) []*Violation {
	return c.CheckFile(context)
}
