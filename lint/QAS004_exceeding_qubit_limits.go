package lint

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/orangekame3/qasmtools/parser"
)

// ExceedingQubitLimitsChecker checks for array index out of bounds violations (QAS004)
type ExceedingQubitLimitsChecker struct{}

func NewExceedingQubitLimitsChecker(maxQubits int) *ExceedingQubitLimitsChecker {
	// maxQubits parameter is ignored since this checker now focuses on array bounds
	return &ExceedingQubitLimitsChecker{}
}

func (c *ExceedingQubitLimitsChecker) Check(node parser.Node, context *CheckContext) []*Violation {
	// This method is required by RuleChecker but not used for program-level analysis
	return nil
}

// CheckProgram implements ProgramChecker interface for program-level analysis
func (c *ExceedingQubitLimitsChecker) CheckProgram(context *CheckContext) []*Violation {
	// Always use text-based analysis since AST parsing has issues (see CLAUDE.md)
	return c.checkFileTextBased(context)
}

func (c *ExceedingQubitLimitsChecker) checkFileTextBased(context *CheckContext) []*Violation {
	var violations []*Violation

	// Use content from context (works for both file-based and LSP content)
	text := context.Content
	if text == "" {
		// Fallback to reading file if content not provided
		content, err := os.ReadFile(context.File)
		if err != nil {
			return violations
		}
		text = string(content)
	}

	lines := strings.Split(text, "\n")

	// Build map of declared arrays and their sizes
	arrayDeclarations := c.findArrayDeclarations(lines)

	// Check all array accesses for bounds violations
	for i, line := range lines {
		violations = append(violations, c.checkLineForBoundsViolations(line, i+1, arrayDeclarations, context)...)
	}

	return violations
}

// arrayDeclaration represents a declared array and its size
type arrayDeclaration struct {
	name string
	size int
	line int
}

func (c *ExceedingQubitLimitsChecker) findArrayDeclarations(lines []string) map[string]arrayDeclaration {
	declarations := make(map[string]arrayDeclaration)

	// Patterns for array declarations
	qubitArrayPattern := regexp.MustCompile(`^\s*qubit\[(\d+)\]\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*;`)
	bitArrayPattern := regexp.MustCompile(`^\s*bit\[(\d+)\]\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*;`)

	for i, line := range lines {
		// Check qubit array declarations
		if matches := qubitArrayPattern.FindStringSubmatch(line); len(matches) > 2 {
			if size, err := strconv.Atoi(matches[1]); err == nil {
				declarations[matches[2]] = arrayDeclaration{
					name: matches[2],
					size: size,
					line: i + 1,
				}
			}
		}

		// Check bit array declarations
		if matches := bitArrayPattern.FindStringSubmatch(line); len(matches) > 2 {
			if size, err := strconv.Atoi(matches[1]); err == nil {
				declarations[matches[2]] = arrayDeclaration{
					name: matches[2],
					size: size,
					line: i + 1,
				}
			}
		}
	}

	return declarations
}

func (c *ExceedingQubitLimitsChecker) checkLineForBoundsViolations(line string, lineNum int, declarations map[string]arrayDeclaration, context *CheckContext) []*Violation {
	var violations []*Violation

	// Pattern for array access: identifier[index]
	arrayAccessPattern := regexp.MustCompile(`\b([a-zA-Z_][a-zA-Z0-9_]*)\[(\d+)\]`)

	matches := arrayAccessPattern.FindAllStringSubmatch(line, -1)
	for _, match := range matches {
		if len(match) > 2 {
			arrayName := match[1]
			indexStr := match[2]

			if index, err := strconv.Atoi(indexStr); err == nil {
				if decl, exists := declarations[arrayName]; exists {
					// Check if index is out of bounds (arrays are 0-indexed)
					if index >= decl.size {
						violation := &Violation{
							File:     context.File,
							Line:     lineNum,
							Column:   strings.Index(line, match[0]) + 1,
							NodeName: arrayName,
							Message:  fmt.Sprintf("Array index out of bounds: %s[%d] exceeds declared size %d (valid indices: 0-%d)", arrayName, index, decl.size, decl.size-1),
							Severity: SeverityError,
						}
						violations = append(violations, violation)
					}
				}
			}
		}
	}

	return violations
}
