package lint

import (
	"os"
	"regexp"
	"strings"

	"github.com/orangekame3/qasmtools/parser"
)

// QubitDeclaredInLocalScopeChecker checks for qubit declarations in local scope (QAS008)
type QubitDeclaredInLocalScopeChecker struct{}

func (c *QubitDeclaredInLocalScopeChecker) Check(node parser.Node, context *CheckContext) []*Violation {
	// This method is required by RuleChecker but not used for program-level analysis
	return nil
}

// CheckProgram implements ProgramChecker interface for program-level analysis
func (c *QubitDeclaredInLocalScopeChecker) CheckProgram(context *CheckContext) []*Violation {
	// Use text-based analysis due to AST parsing issues
	return c.CheckFile(context)
}

// CheckFile performs file-level qubit local scope analysis
func (c *QubitDeclaredInLocalScopeChecker) CheckFile(context *CheckContext) []*Violation {
	var violations []*Violation

	// Read file content for text-based analysis
	content, err := os.ReadFile(context.File)
	if err != nil {
		return violations
	}

	text := string(content)
	lines := strings.Split(text, "\n")

	// Track scope depth and find qubit declarations
	scopeDepth := 0

	for i, line := range lines {
		// Skip comments and empty lines
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "//") {
			continue
		}

		// Remove comments from the line for processing
		codeOnly := c.removeComments(line)

		// For same-line braces and declarations, we need special handling
		// Check if this line has both a brace and a qubit declaration
		openBraces := strings.Count(codeOnly, "{")
		closeBraces := strings.Count(codeOnly, "}")

		// Check for qubit declarations first, before updating scope
		qubitDeclarations := c.findQubitDeclarations(codeOnly)

		// For same-line handling: if we have both an opening brace and a qubit declaration,
		// the declaration is likely in the local scope
		effectiveScopeDepth := scopeDepth
		if openBraces > 0 && len(qubitDeclarations) > 0 {
			// Check the position of the brace vs the qubit declaration
			bracePos := strings.Index(codeOnly, "{")
			for _, decl := range qubitDeclarations {
				qubitPos := strings.Index(codeOnly, "qubit")
				if bracePos >= 0 && qubitPos >= 0 && bracePos < qubitPos {
					// Brace comes before qubit declaration on same line
					effectiveScopeDepth = scopeDepth + 1
				} else {
					effectiveScopeDepth = scopeDepth
				}

				// Check if we're in local scope (not global)
				if effectiveScopeDepth > 0 {
					violation := &Violation{
						Rule:     nil, // Will be set by the runner
						File:     context.File,
						Line:     i + 1,
						Column:   decl.column,
						NodeName: decl.name,
						Message:  "Qubit '" + decl.name + "' can only be declared in global scope.",
						Severity: SeverityError,
					}
					violations = append(violations, violation)
				}
			}
		} else {
			// Normal case: check declarations with current scope
			for _, decl := range qubitDeclarations {
				if scopeDepth > 0 {
					violation := &Violation{
						Rule:     nil, // Will be set by the runner
						File:     context.File,
						Line:     i + 1,
						Column:   decl.column,
						NodeName: decl.name,
						Message:  "Qubit '" + decl.name + "' can only be declared in global scope.",
						Severity: SeverityError,
					}
					violations = append(violations, violation)
				}
			}
		}

		// Update scope depth for next line
		scopeDepth += openBraces - closeBraces
	}

	return violations
}

type qubitDeclaration struct {
	name   string
	column int
}

// findQubitDeclarations finds all qubit declarations in a line
func (c *QubitDeclaredInLocalScopeChecker) findQubitDeclarations(line string) []qubitDeclaration {
	var declarations []qubitDeclaration

	// Patterns for qubit declarations
	patterns := []*regexp.Regexp{
		// qubit name; or qubit[size] name;
		regexp.MustCompile(`\bqubit(?:\[\s*\d+\s*\])?\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*;`),
	}

	for _, pattern := range patterns {
		matches := pattern.FindAllStringSubmatch(line, -1)
		indices := pattern.FindAllStringIndex(line, -1)

		for i, match := range matches {
			if len(match) >= 2 {
				qubitName := match[1]
				// Find the position of the qubit name within the match
				matchStart := indices[i][0]
				qubitPos := strings.Index(line[matchStart:], qubitName)
				if qubitPos != -1 {
					column := matchStart + qubitPos + 1 // Convert to 1-based indexing

					declarations = append(declarations, qubitDeclaration{
						name:   qubitName,
						column: column,
					})
				}
			}
		}
	}

	return declarations
}

// removeComments removes comments from a line
func (c *QubitDeclaredInLocalScopeChecker) removeComments(line string) string {
	if idx := strings.Index(line, "//"); idx != -1 {
		return line[:idx]
	}
	return line
}
