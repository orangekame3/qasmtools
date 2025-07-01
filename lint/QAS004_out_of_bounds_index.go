package lint

import (
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/orangekame3/qasmtools/parser"
)

// OutOfBoundsIndexChecker checks for array index access that exceeds defined bounds (QAS004)
type OutOfBoundsIndexChecker struct{}

func (c *OutOfBoundsIndexChecker) Check(node parser.Node, context *CheckContext) []*Violation {
	// This method is required by RuleChecker but not used for program-level analysis
	return nil
}

// CheckProgram implements ProgramChecker interface for program-level analysis
func (c *OutOfBoundsIndexChecker) CheckProgram(context *CheckContext) []*Violation {
	// Use text-based analysis due to AST parsing issues
	return c.CheckFile(context)
}

// CheckFile performs file-level out-of-bounds index analysis
func (c *OutOfBoundsIndexChecker) CheckFile(context *CheckContext) []*Violation {
	var violations []*Violation

	// Read file content for text-based analysis
	content, err := os.ReadFile(context.File)
	if err != nil {
		return violations
	}

	text := string(content)
	lines := strings.Split(text, "\n")

	// First pass: collect all array declarations and their sizes
	arraySizes := c.findArrayDeclarations(lines)

	// Second pass: find array accesses and check bounds
	for i, line := range lines {
		// Skip comments and empty lines
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "//") {
			continue
		}

		// Find array accesses in this line
		accesses := c.findArrayAccesses(line)
		
		for _, access := range accesses {
			// Check if the array exists and index is within bounds
			if size, exists := arraySizes[access.arrayName]; exists {
				if access.index >= size {
					violation := &Violation{
						Rule:     nil, // Will be set by the runner
						File:     context.File,
						Line:     i + 1,
						Column:   access.column,
						NodeName: access.arrayName,
						Message:  "Index out of bounds: accessing '" + strconv.Itoa(access.index) + "' on '" + access.arrayName + "' of length " + strconv.Itoa(size) + ".",
						Severity: SeverityError,
					}
					violations = append(violations, violation)
				}
			}
		}
	}

	return violations
}

type arrayAccess struct {
	arrayName string
	index     int
	column    int
}

// findArrayDeclarations finds all array declarations and their sizes
func (c *OutOfBoundsIndexChecker) findArrayDeclarations(lines []string) map[string]int {
	arraySizes := make(map[string]int)

	// Patterns for array declarations
	patterns := []*regexp.Regexp{
		// qubit[size] name;
		regexp.MustCompile(`^\s*qubit\[\s*(\d+)\s*\]\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*;`),
		// bit[size] name;
		regexp.MustCompile(`^\s*bit\[\s*(\d+)\s*\]\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*;`),
		// int[size] name; (for future OpenQASM extensions)
		regexp.MustCompile(`^\s*int\[\s*(\d+)\s*\]\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*;`),
		// float[size] name; (for future OpenQASM extensions)
		regexp.MustCompile(`^\s*float\[\s*(\d+)\s*\]\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*;`),
	}

	for _, line := range lines {
		// Skip comments and empty lines
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "//") {
			continue
		}

		for _, pattern := range patterns {
			matches := pattern.FindStringSubmatch(line)
			if len(matches) >= 3 {
				sizeStr := matches[1]
				arrayName := matches[2]
				
				if size, err := strconv.Atoi(sizeStr); err == nil {
					arraySizes[arrayName] = size
				}
			}
		}
	}

	return arraySizes
}

// findArrayAccesses finds all array accesses in a line
func (c *OutOfBoundsIndexChecker) findArrayAccesses(line string) []arrayAccess {
	var accesses []arrayAccess

	// Remove comments from the line
	codeOnly := c.removeComments(line)

	// Pattern for array access: identifier[number]
	accessPattern := regexp.MustCompile(`\b([a-zA-Z_][a-zA-Z0-9_]*)\[\s*(\d+)\s*\]`)
	matches := accessPattern.FindAllStringSubmatch(codeOnly, -1)
	indices := accessPattern.FindAllStringIndex(codeOnly, -1)

	for i, match := range matches {
		if len(match) >= 3 {
			arrayName := match[1]
			indexStr := match[2]
			
			if index, err := strconv.Atoi(indexStr); err == nil {
				column := indices[i][0] + 1 // Convert to 1-based indexing
				
				accesses = append(accesses, arrayAccess{
					arrayName: arrayName,
					index:     index,
					column:    column,
				})
			}
		}
	}

	return accesses
}

// removeComments removes comments from a line
func (c *OutOfBoundsIndexChecker) removeComments(line string) string {
	if idx := strings.Index(line, "//"); idx != -1 {
		return line[:idx]
	}
	return line
}