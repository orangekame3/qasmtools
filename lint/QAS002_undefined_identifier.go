package lint

import (
	"regexp"
	"strings"

	"github.com/orangekame3/qasmtools/parser"
)

// UndefinedIdentifierChecker is the new implementation using BaseChecker framework
type UndefinedIdentifierChecker struct {
	*BaseChecker
	declaredIdentifiers map[string]bool
	builtinIdentifiers  map[string]bool
}

// NewUndefinedIdentifierChecker creates a new UndefinedIdentifierChecker
func NewUndefinedIdentifierChecker() *UndefinedIdentifierChecker {
	return &UndefinedIdentifierChecker{
		BaseChecker:        NewBaseChecker("QAS002"),
		builtinIdentifiers: GetBuiltinIdentifiers(),
	}
}

// CheckFile performs file-level undefined identifier analysis
func (c *UndefinedIdentifierChecker) CheckFile(context *CheckContext) []*Violation {
	var violations []*Violation

	// Get content using BaseChecker method
	content, err := c.getContent(context)
	if err != nil {
		return violations
	}

	// First pass: collect all declared identifiers using shared utilities
	c.declaredIdentifiers = c.findDeclaredIdentifiers(content)

	// Second pass: find usage of undefined identifiers
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		// Skip comments and empty lines using shared utility
		if SkipCommentAndEmptyLine(line) {
			continue
		}

		// Find all identifier usages in this line
		usages := c.findIdentifierUsages(line)

		for _, usage := range usages {
			// Check if the identifier is declared, built-in, or a keyword
			if !c.isIdentifierDefined(usage.name) {
				violation := c.NewViolationBuilder().
					WithMessage("Identifier '"+usage.name+"' is not declared.").
					WithFile(context.File).
					WithPosition(i+1, usage.column).
					WithNodeName(usage.name).
					AsError().
					Build()
				violations = append(violations, violation)
			}
		}
	}

	return violations
}

// identifierUsage represents an identifier usage with position
type identifierUsage struct {
	name   string
	column int
}

// findDeclaredIdentifiers finds all declared identifiers using shared utilities
func (c *UndefinedIdentifierChecker) findDeclaredIdentifiers(content string) map[string]bool {
	declared := make(map[string]bool)
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		// Use shared utility to find identifier declarations
		declarations := FindIdentifierDeclarations(line)
		for _, decl := range declarations {
			declared[decl.Name] = true
		}
	}

	return declared
}

// findIdentifierUsages finds all identifier usages in a line
func (c *UndefinedIdentifierChecker) findIdentifierUsages(line string) []identifierUsage {
	var usages []identifierUsage

	// Remove comments using shared utility
	cleanLine := RemoveComments(line)

	// Skip include statements to avoid false positives on filenames
	if strings.Contains(cleanLine, "include") {
		return usages
	}

	// Skip OPENQASM version declarations
	if strings.Contains(cleanLine, "OPENQASM") {
		return usages
	}

	// Find all identifier-like tokens
	pattern := regexp.MustCompile(`\b([a-zA-Z_][a-zA-Z0-9_]*)\b`)
	matches := pattern.FindAllStringSubmatch(cleanLine, -1)
	indices := pattern.FindAllStringIndex(cleanLine, -1)

	for i, match := range matches {
		if len(match) >= 2 {
			identifier := match[1]
			// Skip keywords using shared utility
			if !IsKeyword(identifier) {
				column := indices[i][0] + 1 // Convert to 1-based indexing
				usages = append(usages, identifierUsage{
					name:   identifier,
					column: column,
				})
			}
		}
	}

	return usages
}

// isIdentifierDefined checks if an identifier is defined (declared, built-in, or keyword)
func (c *UndefinedIdentifierChecker) isIdentifierDefined(name string) bool {
	// Check if it's a declared identifier
	if c.declaredIdentifiers[name] {
		return true
	}

	// Check if it's a built-in identifier using shared utility
	if c.builtinIdentifiers[name] {
		return true
	}

	// Check if it's a keyword using shared utility
	if IsKeyword(name) {
		return true
	}

	return false
}

// Check implements RuleChecker interface (required but delegates to CheckProgram)
func (c *UndefinedIdentifierChecker) Check(node parser.Node, context *CheckContext) []*Violation {
	return nil
}

// CheckProgram implements ProgramChecker interface
func (c *UndefinedIdentifierChecker) CheckProgram(context *CheckContext) []*Violation {
	return c.CheckFile(context)
}
