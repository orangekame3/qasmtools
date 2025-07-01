package lint

import (
	"os"
	"regexp"
	"strings"

	"github.com/orangekame3/qasmtools/parser"
)

// IllegalBreakContinueChecker checks for break/continue statements outside loops (QAS009)
type IllegalBreakContinueChecker struct{}

func (c *IllegalBreakContinueChecker) Check(node parser.Node, context *CheckContext) []*Violation {
	// This method is required by RuleChecker but not used for program-level analysis
	return nil
}

// CheckProgram implements ProgramChecker interface for program-level analysis
func (c *IllegalBreakContinueChecker) CheckProgram(context *CheckContext) []*Violation {
	// Use text-based analysis due to AST parsing issues
	return c.CheckFile(context)
}

// CheckFile performs file-level break/continue analysis
func (c *IllegalBreakContinueChecker) CheckFile(context *CheckContext) []*Violation {
	var violations []*Violation

	// Read file content for text-based analysis
	content, err := os.ReadFile(context.File)
	if err != nil {
		return violations
	}

	text := string(content)
	lines := strings.Split(text, "\n")

	// Track loop nesting depth and find break/continue statements
	loopDepth := 0

	for i, line := range lines {
		// Skip comments and empty lines
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "//") {
			continue
		}

		// Remove comments from the line for processing
		codeOnly := c.removeComments(line)

		// Track loop nesting changes
		loopDepth += c.countLoopChanges(codeOnly)

		// Find break/continue statements in this line
		breakContinueStatements := c.findBreakContinueStatements(codeOnly)

		for _, stmt := range breakContinueStatements {
			// Check if we're outside of any loop
			if loopDepth <= 0 {
				violation := &Violation{
					Rule:     nil, // Will be set by the runner
					File:     context.File,
					Line:     i + 1,
					Column:   stmt.column,
					NodeName: stmt.keyword,
					Message:  "'" + stmt.keyword + "' cannot be used outside of a loop.",
					Severity: SeverityError,
				}
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

// removeComments removes comments from a line
func (c *IllegalBreakContinueChecker) removeComments(line string) string {
	if idx := strings.Index(line, "//"); idx != -1 {
		return line[:idx]
	}
	return line
}
