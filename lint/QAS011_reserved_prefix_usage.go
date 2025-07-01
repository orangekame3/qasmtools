package lint

import (
	"os"
	"regexp"
	"strings"

	"github.com/orangekame3/qasmtools/parser"
)

// ReservedPrefixUsageChecker checks for use of reserved prefixes in identifiers (QAS011)
type ReservedPrefixUsageChecker struct{}

func (c *ReservedPrefixUsageChecker) Check(node parser.Node, context *CheckContext) []*Violation {
	// This method is required by RuleChecker but not used for program-level analysis
	return nil
}

// CheckProgram implements ProgramChecker interface for program-level analysis
func (c *ReservedPrefixUsageChecker) CheckProgram(context *CheckContext) []*Violation {
	// Use text-based analysis due to AST parsing issues
	return c.CheckFile(context)
}

// CheckFile performs file-level reserved prefix usage analysis
func (c *ReservedPrefixUsageChecker) CheckFile(context *CheckContext) []*Violation {
	var violations []*Violation

	// Read file content for text-based analysis
	content, err := os.ReadFile(context.File)
	if err != nil {
		return violations
	}

	text := string(content)
	lines := strings.Split(text, "\n")

	// Check each line for identifier declarations with reserved prefixes
	for i, line := range lines {
		// Skip comments and empty lines
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "//") {
			continue
		}

		// Find identifier declarations and check for reserved prefixes
		identifiers := c.findIdentifierDeclarations(line)

		for _, identifier := range identifiers {
			if c.hasReservedPrefix(identifier.name) {
				violation := &Violation{
					Rule:     nil, // Will be set by the runner
					File:     context.File,
					Line:     i + 1,
					Column:   identifier.column,
					NodeName: identifier.name,
					Message:  "Identifier '" + identifier.name + "' uses reserved prefix '__'.",
					Severity: SeverityError,
				}
				violations = append(violations, violation)
			}
		}
	}

	return violations
}

type identifierDeclarationWithPrefix struct {
	name   string
	column int
}

// findIdentifierDeclarations finds all identifier declarations in a line
func (c *ReservedPrefixUsageChecker) findIdentifierDeclarations(line string) []identifierDeclarationWithPrefix {
	var declarations []identifierDeclarationWithPrefix

	// Remove comments from the line
	codeOnly := c.removeComments(line)

	// Patterns for different types of declarations
	patterns := []*regexp.Regexp{
		// qubit declarations: qubit name; or qubit[size] name;
		regexp.MustCompile(`\bqubit(?:\[\s*\d+\s*\])?\s+([^\s;]+)\s*;`),
		// bit declarations: bit name; or bit[size] name;
		regexp.MustCompile(`\bbit(?:\[\s*\d+\s*\])?\s+([^\s;]+)\s*;`),
		// gate declarations: gate name(...) {...} or gate name params {...}
		regexp.MustCompile(`\bgate\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*(?:\(|[a-zA-Z_])`),
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

					declarations = append(declarations, identifierDeclarationWithPrefix{
						name:   identifierName,
						column: column,
					})
				}
			}
		}
	}

	return declarations
}

// hasReservedPrefix checks if an identifier starts with a reserved prefix
func (c *ReservedPrefixUsageChecker) hasReservedPrefix(name string) bool {
	// Check for double underscore prefix (reserved)
	return strings.HasPrefix(name, "__")
}

// removeComments removes comments from a line
func (c *ReservedPrefixUsageChecker) removeComments(line string) string {
	if idx := strings.Index(line, "//"); idx != -1 {
		return line[:idx]
	}
	return line
}
