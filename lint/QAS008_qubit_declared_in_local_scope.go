package lint

import (
	"regexp"
	"strings"

	"github.com/orangekame3/qasmtools/parser"
)

// QubitDeclaredInLocalScopeChecker is the new implementation using BaseChecker framework
type QubitDeclaredInLocalScopeChecker struct {
	*BaseChecker
}

// NewQubitDeclaredInLocalScopeChecker creates a new QubitDeclaredInLocalScopeChecker
func NewQubitDeclaredInLocalScopeChecker() *QubitDeclaredInLocalScopeChecker {
	return &QubitDeclaredInLocalScopeChecker{
		BaseChecker: NewBaseChecker("QAS008"),
	}
}

// CheckFile performs file-level qubit local scope analysis
func (c *QubitDeclaredInLocalScopeChecker) CheckFile(context *CheckContext) []*Violation {
	var violations []*Violation

	// Get content using BaseChecker method
	content, err := c.getContent(context)
	if err != nil {
		return violations
	}

	lines := strings.Split(content, "\n")

	// Track scope depth and find qubit declarations
	scopeDepth := 0

	for i, line := range lines {
		// Skip comments and empty lines using shared utility
		if SkipCommentAndEmptyLine(line) {
			continue
		}

		// Remove comments using shared utility
		codeOnly := RemoveComments(line)

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
					violation := c.NewViolationBuilder().
						WithMessage("Qubit '"+decl.name+"' can only be declared in global scope.").
						WithFile(context.File).
						WithPosition(i+1, decl.column).
						WithNodeName(decl.name).
						AsError().
						Build()
					violations = append(violations, violation)
				}
			}
		} else {
			// Normal case: check declarations with current scope
			for _, decl := range qubitDeclarations {
				if scopeDepth > 0 {
					violation := c.NewViolationBuilder().
						WithMessage("Qubit '"+decl.name+"' can only be declared in global scope.").
						WithFile(context.File).
						WithPosition(i+1, decl.column).
						WithNodeName(decl.name).
						AsError().
						Build()
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

// findQubitDeclarations finds all qubit declarations in a line using custom patterns (not anchored)
func (c *QubitDeclaredInLocalScopeChecker) findQubitDeclarations(line string) []qubitDeclaration {
	var declarations []qubitDeclaration

	// Custom patterns that are not anchored to line start (for scope detection)
	patterns := []*regexp.Regexp{
		// qubit name; or qubit[size] name;
		regexp.MustCompile(`\bqubit(?:\[\s*\d+\s*\])?\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*;`),
	}

	for _, pattern := range patterns {
		matches := pattern.FindAllStringSubmatch(line, -1)
		indices := pattern.FindAllStringIndex(line, -1)

		for i, match := range matches {
			if len(match) >= 2 {
				qubitName := match[1] // Name is at index 1 for the custom pattern
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

// Check implements RuleChecker interface (required but delegates to CheckProgram)
func (c *QubitDeclaredInLocalScopeChecker) Check(node parser.Node, context *CheckContext) []*Violation {
	return nil
}

// CheckProgram implements ProgramChecker interface
func (c *QubitDeclaredInLocalScopeChecker) CheckProgram(context *CheckContext) []*Violation {
	return c.CheckFile(context)
}
