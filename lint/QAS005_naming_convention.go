package lint

import (
	"os"
	"regexp"
	"strings"

	"github.com/orangekame3/qasmtools/parser"
)

// NamingConventionChecker checks naming conventions (QAS005)
type NamingConventionChecker struct {
	Pattern *regexp.Regexp
}

func NewNamingConventionChecker(pattern string) *NamingConventionChecker {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		// Fallback to default pattern
		regex = regexp.MustCompile("^[a-z][a-zA-Z0-9_]*$")
	}
	
	return &NamingConventionChecker{
		Pattern: regex,
	}
}

func (c *NamingConventionChecker) Check(node parser.Node, context *CheckContext) []*Violation {
	// This method is required by RuleChecker but not used for program-level analysis
	return nil
}

// CheckProgram implements ProgramChecker interface for program-level analysis
func (c *NamingConventionChecker) CheckProgram(context *CheckContext) []*Violation {
	// Always use text-based analysis since AST parsing has issues (see CLAUDE.md)
	return c.checkFileTextBased(context)
}

func (c *NamingConventionChecker) checkFileTextBased(context *CheckContext) []*Violation {
	var violations []*Violation

	// Read file content
	content, err := os.ReadFile(context.File)
	if err != nil {
		return violations
	}

	text := string(content)
	lines := strings.Split(text, "\n")

	// Patterns for different declarations
	qubitPattern := regexp.MustCompile(`^\s*qubit(?:\[\d+\])?\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*;`)
	bitPattern := regexp.MustCompile(`^\s*bit(?:\[\d+\])?\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*;`)
	gatePattern := regexp.MustCompile(`^\s*gate\s+([a-zA-Z_][a-zA-Z0-9_]*)\s+`)

	for i, line := range lines {
		var identifierName string
		var found bool

		// Check qubit declarations
		if matches := qubitPattern.FindStringSubmatch(line); len(matches) > 1 {
			identifierName = matches[1]
			found = true
		}
		// Check bit declarations
		if !found {
			if matches := bitPattern.FindStringSubmatch(line); len(matches) > 1 {
				identifierName = matches[1]
				found = true
			}
		}
		// Check gate definitions
		if !found {
			if matches := gatePattern.FindStringSubmatch(line); len(matches) > 1 {
				identifierName = matches[1]
				found = true
			}
		}

		if found && !c.Pattern.MatchString(identifierName) {
			violation := &Violation{
				File:     context.File,
				Line:     i + 1,
				Column:   strings.Index(line, identifierName) + 1,
				NodeName: identifierName,
				Message:  "Identifier '" + identifierName + "' violates naming convention. Should match pattern: ^[a-z][a-zA-Z0-9_]*$",
				Severity: SeverityWarning,
			}
			violations = append(violations, violation)
		}
	}

	return violations
}