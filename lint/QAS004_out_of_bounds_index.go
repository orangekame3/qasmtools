package lint

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/orangekame3/qasmtools/parser"
)

// OutOfBoundsIndexChecker is the new implementation using BaseChecker framework
type OutOfBoundsIndexChecker struct {
	*BaseChecker
	arraySizes map[string]int // maps array names to their sizes
}

// NewOutOfBoundsIndexChecker creates a new OutOfBoundsIndexChecker
func NewOutOfBoundsIndexChecker() *OutOfBoundsIndexChecker {
	return &OutOfBoundsIndexChecker{
		BaseChecker: NewBaseChecker("QAS004"),
	}
}

// CheckFile performs file-level out-of-bounds index analysis
func (c *OutOfBoundsIndexChecker) CheckFile(context *CheckContext) []*Violation {
	var violations []*Violation

	// Get content using BaseChecker method
	content, err := c.getContent(context)
	if err != nil {
		return violations
	}

	lines := strings.Split(content, "\n")

	// First pass: collect all array declarations and their sizes
	c.arraySizes = c.findArrayDeclarations(lines)

	// Second pass: find array accesses and check bounds
	for i, line := range lines {
		// Skip comments and empty lines using shared utility
		if SkipCommentAndEmptyLine(line) {
			continue
		}

		// Find array accesses in this line
		accesses := c.findArrayAccesses(line)

		for _, access := range accesses {
			arrayName := access.arrayName
			index := access.index

			// Check if array is declared and if index is within bounds
			if size, exists := c.arraySizes[arrayName]; exists {
				if index >= size {
					violation := c.NewViolationBuilder().
						WithMessage("Index out of bounds: accessing '"+strconv.Itoa(index)+"' on '"+arrayName+"' of length "+strconv.Itoa(size)+".").
						WithFile(context.File).
						WithPosition(i+1, access.column).
						WithNodeName(arrayName).
						AsError().
						Build()
					violations = append(violations, violation)
				}
			}
		}
	}

	return violations
}

// arrayAccess represents an array access with position
type arrayAccess struct {
	arrayName string
	index     int
	column    int
}

// findArrayDeclarations finds all array declarations and their sizes
func (c *OutOfBoundsIndexChecker) findArrayDeclarations(lines []string) map[string]int {
	arraySizes := make(map[string]int)

	for _, line := range lines {
		// Skip comments and empty lines using shared utility
		if SkipCommentAndEmptyLine(line) {
			continue
		}

		// Use shared patterns for array declarations
		patterns := []*regexp.Regexp{
			ArrayQubitDeclarationPattern, // qubit[size] name;
			ArrayBitDeclarationPattern,   // bit[size] name;
			// Additional patterns for other array types
			regexp.MustCompile(`\b(?:int|uint|float|angle|complex|bool)\[\s*(\d+)\s*\]\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*;`),
		}

		for _, pattern := range patterns {
			matches := pattern.FindAllStringSubmatch(line, -1)
			for _, match := range matches {
				if len(match) >= 3 {
					sizeStr := match[1]
					arrayName := match[2]

					if size, err := strconv.Atoi(sizeStr); err == nil {
						arraySizes[arrayName] = size
					}
				}
			}
		}
	}

	return arraySizes
}

// findArrayAccesses finds all array accesses in a line
func (c *OutOfBoundsIndexChecker) findArrayAccesses(line string) []arrayAccess {
	var accesses []arrayAccess

	// Remove comments using shared utility
	cleanLine := RemoveComments(line)

	// Use shared pattern for array access
	matches := ArrayAccessPattern.FindAllStringSubmatch(cleanLine, -1)
	indices := ArrayAccessPattern.FindAllStringIndex(cleanLine, -1)

	for i, match := range matches {
		if len(match) >= 3 {
			arrayName := match[1]
			indexStr := match[2]

			if index, err := strconv.Atoi(indexStr); err == nil {
				// Find the position of the index number directly in the line
				matchStart := indices[i][0]
				fullMatch := match[0] // e.g., "q[3]"
				// Find the bracket position and then look for the index
				bracketPos := strings.Index(fullMatch, "[")
				if bracketPos != -1 {
					// Index starts after the opening bracket
					indexStartInMatch := bracketPos + 1
					// Skip any whitespace
					for indexStartInMatch < len(fullMatch) && fullMatch[indexStartInMatch] == ' ' {
						indexStartInMatch++
					}
					column := matchStart + indexStartInMatch + 1 // Convert to 1-based indexing

					accesses = append(accesses, arrayAccess{
						arrayName: arrayName,
						index:     index,
						column:    column,
					})
				}
			}
		}
	}

	return accesses
}

// Check implements RuleChecker interface (required but delegates to CheckProgram)
func (c *OutOfBoundsIndexChecker) Check(node parser.Node, context *CheckContext) []*Violation {
	return nil
}

// CheckProgram implements ProgramChecker interface
func (c *OutOfBoundsIndexChecker) CheckProgram(context *CheckContext) []*Violation {
	return c.CheckFile(context)
}
