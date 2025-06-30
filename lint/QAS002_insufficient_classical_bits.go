package lint

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/orangekame3/qasmtools/parser"
)

// InsufficientClassicalBitsChecker checks for insufficient classical bits (QAS002)
type InsufficientClassicalBitsChecker struct{}

func (c *InsufficientClassicalBitsChecker) Check(node parser.Node, context *CheckContext) []*Violation {
	// This method is required by RuleChecker but not used for program-level analysis
	return nil
}

// CheckProgram implements ProgramChecker interface for program-level analysis
func (c *InsufficientClassicalBitsChecker) CheckProgram(context *CheckContext) []*Violation {
	// Always use text-based analysis since AST parsing has issues (see CLAUDE.md)
	return c.checkFileTextBased(context)
}

func (c *InsufficientClassicalBitsChecker) checkFileTextBased(context *CheckContext) []*Violation {
	var violations []*Violation

	// Read file content
	content, err := os.ReadFile(context.File)
	if err != nil {
		return violations
	}

	text := string(content)
	lines := strings.Split(text, "\n")

	// Count classical bits
	classicalBits := 0
	bitPattern := regexp.MustCompile(`^\s*bit\[(\d+)\]`)
	for _, line := range lines {
		if matches := bitPattern.FindStringSubmatch(line); len(matches) > 1 {
			if count, err := strconv.Atoi(matches[1]); err == nil {
				classicalBits += count
			}
		}
	}

	// Count measurements
	measurements := 0
	measurePattern := regexp.MustCompile(`^\s*measure\s+`)
	for _, line := range lines {
		if measurePattern.MatchString(line) {
			measurements++
		}
	}

	// Only report if there's actually a problem
	if measurements > 0 && classicalBits > 0 && measurements > classicalBits {
		// Find first measurement line for reporting
		for i, line := range lines {
			if measurePattern.MatchString(line) {
				violation := &Violation{
					File:     context.File,
					Line:     i + 1,
					Column:   1,
					NodeName: "",
					Message:  fmt.Sprintf("Insufficient classical bits for measurements. Need %d but only %d declared.", measurements, classicalBits),
					Severity: SeverityError,
				}
				violations = append(violations, violation)
				break
			}
		}
	}

	return violations
}