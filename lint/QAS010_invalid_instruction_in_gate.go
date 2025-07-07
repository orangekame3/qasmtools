package lint

import (
	"regexp"
	"strings"

	"github.com/orangekame3/qasmtools/parser"
)

// InvalidInstructionInGateChecker is the new implementation using BaseChecker framework
type InvalidInstructionInGateChecker struct {
	*BaseChecker
}

// NewInvalidInstructionInGateChecker creates a new InvalidInstructionInGateChecker
func NewInvalidInstructionInGateChecker() *InvalidInstructionInGateChecker {
	return &InvalidInstructionInGateChecker{
		BaseChecker: NewBaseChecker("QAS010"),
	}
}

// CheckFile performs file-level invalid instruction in gate analysis
func (c *InvalidInstructionInGateChecker) CheckFile(context *CheckContext) []*Violation {
	var violations []*Violation

	// Get content using BaseChecker method
	content, err := c.getContent(context)
	if err != nil {
		return violations
	}

	lines := strings.Split(content, "\n")

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
					violation := c.NewViolationBuilder().
						WithMessage("Invalid instruction '"+instruction.instruction+"' used within gate definition.").
						WithFile(context.File).
						WithPosition(i+1, instruction.column).
						WithNodeName(instruction.instruction).
						AsError().
						Build()
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
		// Skip comments and empty lines using shared utility
		if SkipCommentAndEmptyLine(line) {
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

		// Remove comments using shared utility
		line = RemoveComments(line)

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

	// Remove comments using shared utility
	codeOnly := RemoveComments(line)

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

// Check implements RuleChecker interface (required but delegates to CheckProgram)
func (c *InvalidInstructionInGateChecker) Check(node parser.Node, context *CheckContext) []*Violation {
	return nil
}

// CheckProgram implements ProgramChecker interface
func (c *InvalidInstructionInGateChecker) CheckProgram(context *CheckContext) []*Violation {
	return c.CheckFile(context)
}
