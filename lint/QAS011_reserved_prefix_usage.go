package lint

import (
	"regexp"
	"strings"

	"github.com/orangekame3/qasmtools/parser"
)

// ReservedPrefixUsageChecker is the new implementation using BaseChecker framework
type ReservedPrefixUsageChecker struct {
	*BaseChecker
}

// NewReservedPrefixUsageChecker creates a new ReservedPrefixUsageChecker
func NewReservedPrefixUsageChecker() *ReservedPrefixUsageChecker {
	return &ReservedPrefixUsageChecker{
		BaseChecker: NewBaseChecker("QAS011"),
	}
}

// CheckFile performs file-level reserved prefix usage analysis
func (c *ReservedPrefixUsageChecker) CheckFile(context *CheckContext) []*Violation {
	var violations []*Violation

	// Get content using BaseChecker method
	content, err := c.getContent(context)
	if err != nil {
		return violations
	}

	lines := strings.Split(content, "\n")

	// Check each line for identifier declarations with reserved prefixes
	for i, line := range lines {
		// Skip comments and empty lines using shared utility
		if SkipCommentAndEmptyLine(line) {
			continue
		}

		// Find identifier declarations using shared utility and check for reserved prefixes
		identifiers := FindIdentifierDeclarations(line)

		for _, identifier := range identifiers {
			if c.hasReservedPrefix(identifier.Name) {
				violation := c.NewViolationBuilder().
					WithMessage("Identifier '"+identifier.Name+"' uses reserved prefix '__'.").
					WithFile(context.File).
					WithPosition(i+1, identifier.Column).
					WithNodeName(identifier.Name).
					AsError().
					Build()
				violations = append(violations, violation)
			}
		}

		// Also check for additional declaration patterns not covered by shared utility
		additionalIdentifiers := c.findAdditionalDeclarations(line)
		for _, identifier := range additionalIdentifiers {
			if c.hasReservedPrefix(identifier.name) {
				violation := c.NewViolationBuilder().
					WithMessage("Identifier '"+identifier.name+"' uses reserved prefix '__'.").
					WithFile(context.File).
					WithPosition(i+1, identifier.column).
					WithNodeName(identifier.name).
					AsError().
					Build()
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

// findAdditionalDeclarations finds additional identifier declarations not covered by shared utility
func (c *ReservedPrefixUsageChecker) findAdditionalDeclarations(line string) []identifierDeclarationWithPrefix {
	var declarations []identifierDeclarationWithPrefix

	// Remove comments using shared utility
	codeOnly := RemoveComments(line)

	// Additional patterns for declarations not covered by shared utility
	patterns := []*regexp.Regexp{
		// function declarations: def name(...) {...}
		regexp.MustCompile(`\bdef\s+([^\s\(]+)\s*\(`),
		// circuit declarations: circuit name(...) {...}
		regexp.MustCompile(`\bcircuit\s+([^\s\(]+)\s*\(`),
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

// Check implements RuleChecker interface (required but delegates to CheckProgram)
func (c *ReservedPrefixUsageChecker) Check(node parser.Node, context *CheckContext) []*Violation {
	return nil
}

// CheckProgram implements ProgramChecker interface
func (c *ReservedPrefixUsageChecker) CheckProgram(context *CheckContext) []*Violation {
	return c.CheckFile(context)
}
