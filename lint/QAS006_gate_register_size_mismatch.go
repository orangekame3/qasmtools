package lint

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/orangekame3/qasmtools/parser"
)

// GateRegisterSizeMismatchChecker is the new implementation using BaseChecker framework
type GateRegisterSizeMismatchChecker struct {
	*BaseChecker
	registerSizes map[string]int
	gateParams    map[string]gateDefinition
}

// NewGateRegisterSizeMismatchChecker creates a new GateRegisterSizeMismatchChecker
func NewGateRegisterSizeMismatchChecker() *GateRegisterSizeMismatchChecker {
	return &GateRegisterSizeMismatchChecker{
		BaseChecker: NewBaseChecker("QAS006"),
	}
}

// CheckFile performs file-level gate register size mismatch analysis
func (c *GateRegisterSizeMismatchChecker) CheckFile(context *CheckContext) []*Violation {
	var violations []*Violation

	// Get content using BaseChecker method
	content, err := c.getContent(context)
	if err != nil {
		return violations
	}

	lines := strings.Split(content, "\n")

	// First pass: collect all register declarations using shared utilities
	c.registerSizes = c.findRegisterDeclarations(lines)

	// Second pass: collect all gate definitions
	c.gateParams = c.findGateDefinitions(lines)

	// Third pass: find gate calls and check register size consistency
	for i, line := range lines {
		// Skip comments and empty lines using shared utility
		if SkipCommentAndEmptyLine(line) {
			continue
		}

		// Find gate calls in this line
		gateCalls := c.findGateCalls(line)

		for _, call := range gateCalls {
			// Check register size consistency for this gate call
			if violation := c.checkRegisterSizes(call, context.File, i+1); violation != nil {
				violations = append(violations, violation)
			}
		}
	}

	return violations
}

type gateCall struct {
	gateName  string
	registers []string
	column    int
}

type gateDefinition struct {
	name       string
	paramCount int
}

// findRegisterDeclarations finds all register declarations using shared patterns
func (c *GateRegisterSizeMismatchChecker) findRegisterDeclarations(lines []string) map[string]int {
	registerSizes := make(map[string]int)

	// Use shared patterns for register declarations
	patterns := []*regexp.Regexp{
		QubitDeclarationPattern,      // qubit name;
		ArrayQubitDeclarationPattern, // qubit[size] name;
		BitDeclarationPattern,        // bit name;
		ArrayBitDeclarationPattern,   // bit[size] name;
	}

	for _, line := range lines {
		// Skip comments and empty lines using shared utility
		if SkipCommentAndEmptyLine(line) {
			continue
		}

		for _, pattern := range patterns {
			matches := pattern.FindStringSubmatch(line)
			if len(matches) >= 3 {
				sizeStr := matches[1]
				registerName := matches[2]

				size := 1 // Default for single registers
				if sizeStr != "" {
					if s, err := strconv.Atoi(sizeStr); err == nil {
						size = s
					}
				}

				registerSizes[registerName] = size
			}
		}
	}

	return registerSizes
}

// findGateDefinitions finds all gate definitions and their parameter counts
func (c *GateRegisterSizeMismatchChecker) findGateDefinitions(lines []string) map[string]gateDefinition {
	gateParams := make(map[string]gateDefinition)

	// Pattern for gate definitions: gate name(params...) qubits... { ... }
	gatePattern := regexp.MustCompile(`^\s*gate\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*(?:\([^)]*\))?\s+([^{]+)\{`)

	for _, line := range lines {
		// Skip comments and empty lines using shared utility
		if SkipCommentAndEmptyLine(line) {
			continue
		}

		matches := gatePattern.FindStringSubmatch(line)
		if len(matches) >= 3 {
			gateName := matches[1]
			qubitParams := strings.TrimSpace(matches[2])

			// Count the number of qubit parameters
			if qubitParams != "" {
				// Split by commas and count parameters
				params := strings.Split(qubitParams, ",")
				paramCount := 0
				for _, param := range params {
					param = strings.TrimSpace(param)
					if param != "" {
						paramCount++
					}
				}

				gateParams[gateName] = gateDefinition{
					name:       gateName,
					paramCount: paramCount,
				}
			}
		}
	}

	return gateParams
}

// findGateCalls finds all gate calls in a line
func (c *GateRegisterSizeMismatchChecker) findGateCalls(line string) []gateCall {
	var calls []gateCall

	// Remove comments using shared utility
	codeOnly := RemoveComments(line)

	// Pattern for gate calls: gatename(...) register1, register2, ...; (parameterized)
	// or gatename register1, register2, ...; (non-parameterized)
	patterns := []*regexp.Regexp{
		// Parameterized gate calls: gatename(params) registers;
		regexp.MustCompile(`\b([a-zA-Z_][a-zA-Z0-9_]*)\s*\([^)]*\)\s+((?:[a-zA-Z_][a-zA-Z0-9_]*(?:\[[^\]]*\])?\s*,?\s*)+)\s*;`),
		// Non-parameterized gate calls: gatename registers;
		regexp.MustCompile(`\b([a-zA-Z_][a-zA-Z0-9_]*)\s+((?:[a-zA-Z_][a-zA-Z0-9_]*(?:\[[^\]]*\])?\s*,?\s*)+)\s*;`),
	}

	for _, pattern := range patterns {
		matches := pattern.FindAllStringSubmatch(codeOnly, -1)
		indices := pattern.FindAllStringIndex(codeOnly, -1)

		for i, match := range matches {
			if len(match) >= 3 {
				gateName := match[1]
				registerList := strings.TrimSpace(match[2])

				// Skip keywords using shared utility
				if IsKeyword(gateName) {
					continue
				}

				// Parse the register list
				registers := c.parseRegisterList(registerList)

				if len(registers) > 1 { // Only check multi-register gate calls
					column := indices[i][0] + 1 // Convert to 1-based indexing

					calls = append(calls, gateCall{
						gateName:  gateName,
						registers: registers,
						column:    column,
					})
				}
			}
		}
	}

	return calls
}

// parseRegisterList parses a comma-separated list of registers
func (c *GateRegisterSizeMismatchChecker) parseRegisterList(registerList string) []string {
	var registers []string

	// Split by commas and clean up
	parts := strings.Split(registerList, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			registers = append(registers, part)
		}
	}

	return registers
}

// hasIndexAccess checks if a register string contains index access
func (c *GateRegisterSizeMismatchChecker) hasIndexAccess(register string) bool {
	return strings.Contains(register, "[") && strings.Contains(register, "]")
}

// checkRegisterSizes checks if registers in a gate call have consistent sizes
func (c *GateRegisterSizeMismatchChecker) checkRegisterSizes(call gateCall, file string, line int) *Violation {
	if len(call.registers) < 2 {
		return nil // Single register calls don't need size checking
	}

	// Check if any register is accessed with an index (individual qubit access)
	// If so, skip size checking as these are individual qubit operations
	for _, registerName := range call.registers {
		if c.hasIndexAccess(registerName) {
			return nil // Individual qubit access, skip size checking
		}
	}

	var sizes []int

	// Get sizes for all registers in the call
	for _, registerName := range call.registers {
		baseRegisterName := ExtractRegisterName(registerName) // Use shared utility
		if size, exists := c.registerSizes[baseRegisterName]; exists {
			sizes = append(sizes, size)
		}
	}

	// Check if all sizes are the same (or if some are single registers which can broadcast)
	if len(sizes) >= 2 {
		firstSize := sizes[0]
		hasMismatch := false

		for i := 1; i < len(sizes); i++ {
			// Allow broadcasting: single registers (size 1) can work with any size
			// But array registers must match exactly
			if sizes[i] != firstSize && sizes[i] != 1 && firstSize != 1 {
				hasMismatch = true
				break
			}
		}

		if hasMismatch {
			return c.NewViolationBuilder().
				WithMessage("Register lengths passed to gate '"+call.gateName+"' do not match.").
				WithFile(file).
				WithPosition(line, call.column).
				WithNodeName(call.gateName).
				AsError().
				Build()
		}
	}

	return nil
}

// Check implements RuleChecker interface (required but delegates to CheckProgram)
func (c *GateRegisterSizeMismatchChecker) Check(node parser.Node, context *CheckContext) []*Violation {
	return nil
}

// CheckProgram implements ProgramChecker interface
func (c *GateRegisterSizeMismatchChecker) CheckProgram(context *CheckContext) []*Violation {
	return c.CheckFile(context)
}
