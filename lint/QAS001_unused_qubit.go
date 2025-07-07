package lint

import (
	"regexp"
	"strings"

	"github.com/orangekame3/qasmtools/parser"
)

// UnusedQubitChecker is the new implementation using the BaseChecker framework
type UnusedQubitChecker struct {
	*BaseChecker
	qubitDeclarations []qubitDecl
}

// NewUnusedQubitChecker creates a new UnusedQubitChecker
func NewUnusedQubitChecker() *UnusedQubitChecker {
	return &UnusedQubitChecker{
		BaseChecker: NewBaseChecker("QAS001"),
	}
}

type qubitDecl struct {
	name   string
	line   int
	column int
}

// CheckFile performs file-level unused qubit analysis
func (c *UnusedQubitChecker) CheckFile(context *CheckContext) []*Violation {
	var violations []*Violation

	// Get content for text-based analysis using BaseChecker method
	content, err := c.getContent(context)
	if err != nil {
		return violations
	}

	// First pass: collect all qubit declarations
	c.qubitDeclarations = c.findQubitDeclarations(content)

	// Second pass: check for usage
	for _, decl := range c.qubitDeclarations {
		if !c.isQubitUsed(decl.name, content) {
			violation := c.NewViolationBuilder().
				WithMessage("Qubit '"+decl.name+"' is declared but never used.").
				WithFile(context.File).
				WithPosition(decl.line, decl.column).
				WithNodeName(decl.name).
				AsWarning().
				Build()
			violations = append(violations, violation)
		}
	}

	return violations
}

// findQubitDeclarations finds all qubit declarations using shared regex patterns
func (c *UnusedQubitChecker) findQubitDeclarations(content string) []qubitDecl {
	var declarations []qubitDecl
	lines := strings.Split(content, "\n")

	for i, line := range lines {
		// Skip comments and empty lines using shared utility
		if SkipCommentAndEmptyLine(line) {
			continue
		}

		// Check for array qubit declarations first: qubit[size] name;
		if matches := ArrayQubitDeclarationPattern.FindStringSubmatch(line); len(matches) > 2 {
			varName := matches[2] // Variable name is in matches[2] for array declarations
			declarations = append(declarations, qubitDecl{
				name:   varName,
				line:   i + 1,
				column: strings.Index(line, varName) + 1,
			})
		} else if matches := QubitDeclarationPattern.FindStringSubmatch(line); len(matches) > 1 {
			// Handle single qubit declarations only if it's not an array declaration: qubit name;
			varName := matches[1]
			declarations = append(declarations, qubitDecl{
				name:   varName,
				line:   i + 1,
				column: strings.Index(line, varName) + 1,
			})
		}
	}

	return declarations
}

// isQubitUsed checks if a qubit is used anywhere in the content
func (c *UnusedQubitChecker) isQubitUsed(qubitName string, content string) bool {
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		// Skip declaration lines
		declPattern := `^\s*(qubit(\[\d+\])?\s+` + regexp.QuoteMeta(qubitName) + `\s*;)`
		if matched, _ := regexp.MatchString(declPattern, line); matched {
			continue
		}

		// Remove comments using shared utility
		cleanLine := RemoveComments(line)

		// Look for usage patterns using shared patterns where applicable
		usagePatterns := []string{
			`\b` + regexp.QuoteMeta(qubitName) + `\[\d+\]`,      // Array access
			`\b[a-z]+\s+` + regexp.QuoteMeta(qubitName) + `\b`,  // Gate application
			`\b` + regexp.QuoteMeta(qubitName) + `\s*,`,         // Usage in gate parameters (first param)
			`,\s*` + regexp.QuoteMeta(qubitName) + `\b`,         // Usage in gate parameters (second param)
			`\bmeasure\s+` + regexp.QuoteMeta(qubitName) + `\b`, // Measurement
		}

		for _, pattern := range usagePatterns {
			if matched, _ := regexp.MatchString(pattern, cleanLine); matched {
				return true
			}
		}
	}

	return false
}

// Check implements RuleChecker interface (required but delegates to CheckProgram)
func (c *UnusedQubitChecker) Check(node parser.Node, context *CheckContext) []*Violation {
	return nil
}

// CheckProgram implements ProgramChecker interface
func (c *UnusedQubitChecker) CheckProgram(context *CheckContext) []*Violation {
	return c.CheckFile(context)
}
