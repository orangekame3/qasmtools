package lint

import (
	"os"
	"regexp"
	"strings"

	"github.com/orangekame3/qasmtools/parser"
)

// InvalidInstructionInGateChecker checks for non-unitary instructions in gate definitions (QAS010)
type InvalidInstructionInGateChecker struct{}

func (c *InvalidInstructionInGateChecker) Check(node parser.Node, context *CheckContext) []*Violation {
	// This method is required by RuleChecker but not used for program-level analysis
	return nil
}

// CheckProgram implements ProgramChecker interface for program-level analysis
func (c *InvalidInstructionInGateChecker) CheckProgram(context *CheckContext) []*Violation {
	// Use text-based analysis due to AST parsing issues
	return c.CheckFile(context)
}

// CheckFile performs file-level invalid instruction in gate analysis
func (c *InvalidInstructionInGateChecker) CheckFile(context *CheckContext) []*Violation {
	var violations []*Violation

	// Read file content for text-based analysis
	content, err := os.ReadFile(context.File)
	if err != nil {
		return violations
	}

	text := string(content)
	lines := strings.Split(text, "\n")

	// Find all gate definitions and check their bodies
	gateDefinitions := c.findGateDefinitions(lines)

	for _, gateDef := range gateDefinitions {
		// Check each line within the gate body for invalid instructions
		for i := gateDef.startLine; i <= gateDef.endLine; i++ {
			if i < len(lines) {
				line := lines[i]

				// Find any invalid instructions in this line
				invalidInstructions := c.findInvalidInstructions(line)

				for _, instruction := range invalidInstructions {
					violation := &Violation{
						Rule:     nil, // Will be set by the runner
						File:     context.File,
						Line:     i + 1,
						Column:   instruction.column,
						NodeName: instruction.instruction,
						Message:  "Invalid instruction '" + instruction.instruction + "' used within gate definition.",
						Severity: SeverityError,
					}
					violations = append(violations, violation)
				}
			}
		}
	}

	return violations
}

type gateDefWithInstructions struct {
	name      string
	startLine int
	endLine   int
}

type invalidInstruction struct {
	instruction string
	column      int
}

// findGateDefinitions finds all gate definitions in the file
func (c *InvalidInstructionInGateChecker) findGateDefinitions(lines []string) []gateDefWithInstructions {
	var gates []gateDefWithInstructions

	// Pattern for gate definitions: gate name(params...) qubits... {
	gateStartPattern := regexp.MustCompile(`^\s*gate\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*(?:\([^)]*\))?\s+[^{]*\s*\{`)

	for i, line := range lines {
		// Skip comments and empty lines
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "//") {
			continue
		}

		matches := gateStartPattern.FindStringSubmatch(line)
		if len(matches) >= 2 {
			gateName := matches[1]

			// Find the end of this gate definition
			endLine := c.findGateEnd(lines, i)

			gates = append(gates, gateDefWithInstructions{
				name:      gateName,
				startLine: i,
				endLine:   endLine,
			})
		}
	}

	return gates
}

// findGateEnd finds the closing brace of a gate definition
func (c *InvalidInstructionInGateChecker) findGateEnd(lines []string, startLine int) int {
	braceCount := 0
	foundOpenBrace := false

	for i := startLine; i < len(lines); i++ {
		line := lines[i]

		// Remove comments
		if idx := strings.Index(line, "//"); idx != -1 {
			line = line[:idx]
		}

		// Count braces
		for _, char := range line {
			if char == '{' {
				braceCount++
				foundOpenBrace = true
			} else if char == '}' {
				braceCount--
				if foundOpenBrace && braceCount == 0 {
					return i
				}
			}
		}
	}

	// If no closing brace found, return the last line
	return len(lines) - 1
}

// findInvalidInstructions finds all invalid (non-unitary) instructions in a line
func (c *InvalidInstructionInGateChecker) findInvalidInstructions(line string) []invalidInstruction {
	var instructions []invalidInstruction

	// Remove comments from the line
	codeOnly := c.removeComments(line)

	// List of non-unitary instructions that are invalid in gate definitions
	invalidInstructionPatterns := []*regexp.Regexp{
		// measure statement: measure qubit -> bit;
		regexp.MustCompile(`\b(measure)\s+`),
		// reset statement: reset qubit;
		regexp.MustCompile(`\b(reset)\s+`),
		// barrier statement: barrier qubits;
		regexp.MustCompile(`\b(barrier)\s+`),
		// Classical control flow (not allowed in gates)
		regexp.MustCompile(`\b(if)\s*\(`),
		regexp.MustCompile(`\b(while)\s*\(`),
		regexp.MustCompile(`\b(for)\s+`),
		// Classical assignments (not allowed in gates)
		regexp.MustCompile(`\b([a-zA-Z_][a-zA-Z0-9_]*)\s*=\s*`),
		// Function calls (generally not allowed in gates unless they're unitary gates)
		// Note: This is more complex and may need refinement based on OpenQASM 3.0 spec
	}

	for _, pattern := range invalidInstructionPatterns {
		matches := pattern.FindAllStringSubmatch(codeOnly, -1)
		indices := pattern.FindAllStringIndex(codeOnly, -1)

		for i, match := range matches {
			if len(match) >= 2 {
				instructionName := match[1]
				column := indices[i][0] + 1 // Convert to 1-based indexing

				instructions = append(instructions, invalidInstruction{
					instruction: instructionName,
					column:      column,
				})
			}
		}
	}

	return instructions
}

// removeComments removes comments from a line
func (c *InvalidInstructionInGateChecker) removeComments(line string) string {
	if idx := strings.Index(line, "//"); idx != -1 {
		return line[:idx]
	}
	return line
}
