package lint

import (
	"os"
	"regexp"
	"strings"

	"github.com/orangekame3/qasmtools/parser"
)

// NamingConventionViolationChecker checks for naming convention violations (QAS005)
type NamingConventionViolationChecker struct{}

func (c *NamingConventionViolationChecker) Check(node parser.Node, context *CheckContext) []*Violation {
	// This method is required by RuleChecker but not used for program-level analysis
	return nil
}

// CheckProgram implements ProgramChecker interface for program-level analysis
func (c *NamingConventionViolationChecker) CheckProgram(context *CheckContext) []*Violation {
	// Use text-based analysis due to AST parsing issues
	return c.CheckFile(context)
}

// CheckFile performs file-level naming convention analysis
func (c *NamingConventionViolationChecker) CheckFile(context *CheckContext) []*Violation {
	var violations []*Violation

	// Read file content for text-based analysis
	content, err := os.ReadFile(context.File)
	if err != nil {
		return violations
	}

	text := string(content)
	lines := strings.Split(text, "\n")

	// Check each line for identifier declarations
	for i, line := range lines {
		// Skip comments and empty lines
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "//") {
			continue
		}

		// Find identifier declarations and check naming conventions
		identifiers := c.findIdentifierDeclarations(line)
		
		for _, identifier := range identifiers {
			if !c.isValidIdentifierName(identifier.name) {
				violation := &Violation{
					Rule:     nil, // Will be set by the runner
					File:     context.File,
					Line:     i + 1,
					Column:   identifier.column,
					NodeName: identifier.name,
					Message:  "Identifier '" + identifier.name + "' violates naming conventions. Follow pattern: ^[a-z][a-zA-Z0-9_]*$.",
					Severity: SeverityWarning,
				}
				violations = append(violations, violation)
			}
		}
	}

	return violations
}

type identifierDeclaration struct {
	name   string
	column int
}

// findIdentifierDeclarations finds all identifier declarations in a line
func (c *NamingConventionViolationChecker) findIdentifierDeclarations(line string) []identifierDeclaration {
	var declarations []identifierDeclaration

	// Remove comments from the line
	codeOnly := c.removeComments(line)

	// Patterns for different types of declarations
	// Use more permissive pattern to capture any identifier-like token
	patterns := []*regexp.Regexp{
		// qubit declarations: qubit name; or qubit[size] name;
		regexp.MustCompile(`\bqubit(?:\[\s*\d+\s*\])?\s+([^\s;]+)\s*;`),
		// bit declarations: bit name; or bit[size] name;
		regexp.MustCompile(`\bbit(?:\[\s*\d+\s*\])?\s+([^\s;]+)\s*;`),
		// gate declarations: gate name(...) {...}
		regexp.MustCompile(`\bgate\s+([^\s\(]+)\s*\(`),
		// function declarations: def name(...) {...}
		regexp.MustCompile(`\bdef\s+([^\s\(]+)\s*\(`),
		// circuit declarations: circuit name(...) {...}
		regexp.MustCompile(`\bcircuit\s+([^\s\(]+)\s*\(`),
		// register declarations: int name; float name; etc.
		regexp.MustCompile(`\b(?:int|uint|float|angle|complex|bool)(?:\[\s*\d+\s*\])?\s+([^\s;]+)\s*;`),
		// const declarations: const name = value;
		regexp.MustCompile(`\bconst\s+([^\s=]+)\s*=`),
		// input/output declarations: input name; output name;
		regexp.MustCompile(`\b(?:input|output)\s+([^\s;]+)\s*;`),
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
					
					declarations = append(declarations, identifierDeclaration{
						name:   identifierName,
						column: column,
					})
				}
			}
		}
	}

	return declarations
}

// isValidIdentifierName checks if an identifier follows the naming convention
func (c *NamingConventionViolationChecker) isValidIdentifierName(name string) bool {
	// Pattern: must start with lowercase letter, followed by letters, digits, or underscores
	pattern := regexp.MustCompile(`^[a-z][a-zA-Z0-9_]*$`)
	return pattern.MatchString(name)
}

// removeComments removes comments from a line
func (c *NamingConventionViolationChecker) removeComments(line string) string {
	if idx := strings.Index(line, "//"); idx != -1 {
		return line[:idx]
	}
	return line
}