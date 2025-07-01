package lint

import (
	"os"
	"regexp"
	"strings"

	"github.com/orangekame3/qasmtools/parser"
)

// SnakeCaseRequiredChecker checks for snake_case naming in identifiers (QAS012)
type SnakeCaseRequiredChecker struct{}

func (c *SnakeCaseRequiredChecker) Check(node parser.Node, context *CheckContext) []*Violation {
	// This method is required by RuleChecker but not used for program-level analysis
	return nil
}

// CheckProgram implements ProgramChecker interface for program-level analysis
func (c *SnakeCaseRequiredChecker) CheckProgram(context *CheckContext) []*Violation {
	// Use text-based analysis due to AST parsing issues
	return c.CheckFile(context)
}

// CheckFile performs file-level snake_case naming analysis
func (c *SnakeCaseRequiredChecker) CheckFile(context *CheckContext) []*Violation {
	var violations []*Violation

	// Read file content for text-based analysis
	content, err := os.ReadFile(context.File)
	if err != nil {
		return violations
	}

	text := string(content)
	lines := strings.Split(text, "\n")

	// Check each line for identifier declarations with non-snake_case naming
	for i, line := range lines {
		// Skip comments and empty lines
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "//") {
			continue
		}

		// Find identifier declarations and check for snake_case
		identifiers := c.findIdentifierDeclarations(line)

		for _, identifier := range identifiers {
			if !c.isSnakeCase(identifier.name) {
				violation := &Violation{
					Rule:     nil, // Will be set by the runner
					File:     context.File,
					Line:     i + 1,
					Column:   identifier.column,
					NodeName: identifier.name,
					Message:  "Identifier '" + identifier.name + "' should be written in snake_case.",
					Severity: SeverityWarning,
				}
				violations = append(violations, violation)
			}
		}
	}

	return violations
}

type identifierDeclarationSnakeCase struct {
	name   string
	column int
}

// findIdentifierDeclarations finds all identifier declarations in a line
func (c *SnakeCaseRequiredChecker) findIdentifierDeclarations(line string) []identifierDeclarationSnakeCase {
	var declarations []identifierDeclarationSnakeCase

	// Remove comments from the line
	codeOnly := c.removeComments(line)

	// Patterns for different types of declarations that should follow snake_case
	patterns := []*regexp.Regexp{
		// qubit declarations: qubit name; or qubit[size] name;
		regexp.MustCompile(`\bqubit(?:\[\s*\d+\s*\])?\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*;`),
		// bit declarations: bit name; or bit[size] name;
		regexp.MustCompile(`\bbit(?:\[\s*\d+\s*\])?\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*;`),
		// gate declarations: gate name(...) {...} or gate name params {...}
		regexp.MustCompile(`\bgate\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*(?:\(|[a-zA-Z_])`),
		// circuit declarations: circuit name(...) {...}
		regexp.MustCompile(`\bcircuit\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\(`),
	}

	for _, pattern := range patterns {
		matches := pattern.FindAllStringSubmatch(codeOnly, -1)
		indices := pattern.FindAllStringIndex(codeOnly, -1)

		for i, match := range matches {
			if len(match) >= 2 {
				identifierName := match[1]
				// Find the position of the identifier within the match
				matchStart := indices[i][0]
				identifierPos := strings.Index(codeOnly[matchStart:], identifierName)
				if identifierPos != -1 {
					column := matchStart + identifierPos + 1 // Convert to 1-based indexing

					declarations = append(declarations, identifierDeclarationSnakeCase{
						name:   identifierName,
						column: column,
					})
				}
			}
		}
	}

	return declarations
}

// isSnakeCase checks if an identifier follows snake_case naming convention
func (c *SnakeCaseRequiredChecker) isSnakeCase(name string) bool {
	// Pattern: must start with lowercase letter, followed by lowercase letters, digits, or underscores
	// No consecutive underscores, no ending underscore
	snakeCasePattern := regexp.MustCompile(`^[a-z][a-z0-9_]*[a-z0-9]$|^[a-z]$`)

	// Check basic pattern
	if !snakeCasePattern.MatchString(name) {
		return false
	}

	// Additional checks for proper snake_case:
	// 1. No consecutive underscores
	if strings.Contains(name, "__") {
		return false
	}

	// 2. No ending with underscore (unless single character)
	if len(name) > 1 && strings.HasSuffix(name, "_") {
		return false
	}

	// 3. No starting with underscore (handled by pattern but double-check)
	if strings.HasPrefix(name, "_") {
		return false
	}

	return true
}

// removeComments removes comments from a line
func (c *SnakeCaseRequiredChecker) removeComments(line string) string {
	if idx := strings.Index(line, "//"); idx != -1 {
		return line[:idx]
	}
	return line
}
