package lint

import (
	"regexp"
	"strings"

	"github.com/orangekame3/qasmtools/parser"
)

// UnusedQubitChecker checks for unused qubit declarations (QAS001)
type UnusedQubitChecker struct{}

func (c *UnusedQubitChecker) Check(node parser.Node, context *CheckContext) []*Violation {
	// This method is required by RuleChecker but not used for program-level analysis
	return nil
}

// CheckProgram implements ProgramChecker interface for program-level analysis
func (c *UnusedQubitChecker) CheckProgram(context *CheckContext) []*Violation {
	// Always use text-based analysis since AST parsing has issues (see CLAUDE.md)
	return c.CheckFile(context)
}

// extractQubitName extracts the base name from qubit identifier (removes array brackets)
func (c *UnusedQubitChecker) extractQubitName(identifier string) string {
	if idx := strings.Index(identifier, "["); idx != -1 {
		return identifier[:idx]
	}
	return identifier
}

// CheckFile performs file-level unused qubit analysis as fallback when AST parsing fails
func (c *UnusedQubitChecker) CheckFile(context *CheckContext) []*Violation {
	var violations []*Violation

	// Get content for text-based analysis
	text, err := context.GetContent()
	if err != nil {
		return violations
	}
	lines := strings.Split(text, "\n")

	// Find qubit declarations
	qubitDeclarations := c.findQubitDeclarations(lines)

	// Check each declared qubit for usage
	for _, decl := range qubitDeclarations {
		if !c.isQubitUsed(decl.name, text) {
			violation := &Violation{
				File:     context.File,
				Line:     decl.line,
				Column:   decl.column,
				NodeName: decl.name,
				Message:  "Qubit '" + decl.name + "' is declared but never used.",
				Severity: SeverityWarning,
			}
			violations = append(violations, violation)
		}
	}

	return violations
}

type qubitDecl struct {
	name   string
	line   int
	column int
}

func (c *UnusedQubitChecker) findQubitDeclarations(lines []string) []qubitDecl {
	var declarations []qubitDecl

	// Regex patterns for qubit declarations
	singleQubit := regexp.MustCompile(`^\s*qubit\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*;`)
	arrayQubit := regexp.MustCompile(`^\s*qubit\[\s*\d+\s*\]\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*;`)

	for i, line := range lines {
		// Skip comments
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "//") {
			continue
		}

		// Check for single qubit declaration
		if matches := singleQubit.FindStringSubmatch(line); len(matches) > 1 {
			declarations = append(declarations, qubitDecl{
				name:   matches[1],
				line:   i + 1,
				column: strings.Index(line, matches[1]) + 1,
			})
		}

		// Check for array qubit declaration
		if matches := arrayQubit.FindStringSubmatch(line); len(matches) > 1 {
			declarations = append(declarations, qubitDecl{
				name:   matches[1],
				line:   i + 1,
				column: strings.Index(line, matches[1]) + 1,
			})
		}
	}

	return declarations
}

func (c *UnusedQubitChecker) isQubitUsed(qubitName string, content string) bool {
	// Look for usage patterns, but exclude declarations
	// Split content into lines to analyze individually
	lines := strings.Split(content, "\n")
	
	for _, line := range lines {
		// Skip declaration lines
		declPattern := `^\s*(qubit(\[\d+\])?\s+` + regexp.QuoteMeta(qubitName) + `\s*;)`
		if matched, _ := regexp.MatchString(declPattern, line); matched {
			continue
		}
		
		// Look for usage patterns in this line:
		// 1. Array access: qubit_name[0]
		// 2. Gate calls: h qubit_name;
		// 3. Gate parameters: cx qubit1, qubit_name;
		// 4. Measurements: measure qubit_name
		
		patterns := []string{
			`\b` + regexp.QuoteMeta(qubitName) + `\[\d+\]`,           // Array access
			`\b[a-z]+\s+` + regexp.QuoteMeta(qubitName) + `\b`,       // Gate application
			`\b` + regexp.QuoteMeta(qubitName) + `\s*,`,              // Usage in gate parameters (first param)
			`,\s*` + regexp.QuoteMeta(qubitName) + `\b`,              // Usage in gate parameters (second param)
			`\bmeasure\s+` + regexp.QuoteMeta(qubitName) + `\b`,      // Measurement
		}

		for _, pattern := range patterns {
			matched, _ := regexp.MatchString(pattern, line)
			if matched {
				return true
			}
		}
	}

	return false
}