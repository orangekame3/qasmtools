package lint

import (
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/orangekame3/qasmtools/parser"
)

// GateRegisterSizeMismatchChecker checks for register size mismatches in gate calls (QAS006)
type GateRegisterSizeMismatchChecker struct{}

func (c *GateRegisterSizeMismatchChecker) Check(node parser.Node, context *CheckContext) []*Violation {
	// This method is required by RuleChecker but not used for program-level analysis
	return nil
}

// CheckProgram implements ProgramChecker interface for program-level analysis
func (c *GateRegisterSizeMismatchChecker) CheckProgram(context *CheckContext) []*Violation {
	// Use text-based analysis due to AST parsing issues
	return c.CheckFile(context)
}

// CheckFile performs file-level gate register size mismatch analysis
func (c *GateRegisterSizeMismatchChecker) CheckFile(context *CheckContext) []*Violation {
	var violations []*Violation

	// Read file content for text-based analysis
	content, err := os.ReadFile(context.File)
	if err != nil {
		return violations
	}

	text := string(content)
	lines := strings.Split(text, "\n")

	// First pass: collect all register declarations and their sizes
	registerSizes := c.findRegisterDeclarations(lines)
	
	// Second pass: collect all gate definitions and their parameter counts
	gateParams := c.findGateDefinitions(lines)

	// Third pass: find gate calls and check register size consistency
	for i, line := range lines {
		// Skip comments and empty lines
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "//") {
			continue
		}

		// Find gate calls in this line
		gateCalls := c.findGateCalls(line)
		
		for _, call := range gateCalls {
			// Check register size consistency for this gate call
			if violation := c.checkRegisterSizes(call, registerSizes, gateParams, context.File, i+1); violation != nil {
				violations = append(violations, violation)
			}
		}
	}

	return violations
}

type registerInfo struct {
	name string
	size int // -1 for single registers (no array)
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

// findRegisterDeclarations finds all register declarations and their sizes
func (c *GateRegisterSizeMismatchChecker) findRegisterDeclarations(lines []string) map[string]int {
	registerSizes := make(map[string]int)

	// Patterns for register declarations
	patterns := []*regexp.Regexp{
		// qubit[size] name; or qubit name;
		regexp.MustCompile(`^\s*qubit(?:\[(\d+)\])?\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*;`),
		// bit[size] name; or bit name;
		regexp.MustCompile(`^\s*bit(?:\[(\d+)\])?\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*;`),
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
		// Skip comments and empty lines
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "//") {
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

	// Remove comments from the line
	codeOnly := c.removeComments(line)

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
				
				// Skip certain keywords that are not gate calls
				if c.isKeyword(gateName) {
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
			// Keep the original register string (with indices if present)
			// We'll extract the base name later when needed
			registers = append(registers, part)
		}
	}
	
	return registers
}

// extractRegisterName extracts the base register name from register[index] format
func (c *GateRegisterSizeMismatchChecker) extractRegisterName(register string) string {
	// Pattern to match register[index] or just register
	pattern := regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)(?:\[[^\]]*\])?`)
	matches := pattern.FindStringSubmatch(register)
	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}

// isKeyword checks if a name is a reserved keyword or statement type
func (c *GateRegisterSizeMismatchChecker) isKeyword(name string) bool {
	keywords := map[string]bool{
		"OPENQASM": true,
		"include":  true,
		"qubit":    true,
		"bit":      true,
		"int":      true,
		"uint":     true,
		"float":    true,
		"angle":    true,
		"complex":  true,
		"bool":     true,
		"const":    true,
		"def":      true,
		"gate":     true,
		"circuit":  true,
		"measure":  true,
		"reset":    true,
		"barrier":  true,
		"if":       true,
		"else":     true,
		"for":      true,
		"while":    true,
		"break":    true,
		"continue": true,
		"return":   true,
		"input":    true,
		"output":   true,
	}
	
	return keywords[name]
}

// checkRegisterSizes checks if registers in a gate call have consistent sizes
func (c *GateRegisterSizeMismatchChecker) checkRegisterSizes(call gateCall, registerSizes map[string]int, gateParams map[string]gateDefinition, file string, line int) *Violation {
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
	var sizeNames []string
	
	// Get sizes for all registers in the call
	for _, registerName := range call.registers {
		baseRegisterName := c.extractRegisterName(registerName)
		if size, exists := registerSizes[baseRegisterName]; exists {
			sizes = append(sizes, size)
			sizeNames = append(sizeNames, baseRegisterName)
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
			return &Violation{
				Rule:     nil, // Will be set by the runner
				File:     file,
				Line:     line,
				Column:   call.column,
				NodeName: call.gateName,
				Message:  "Register lengths passed to gate '" + call.gateName + "' do not match.",
				Severity: SeverityError,
			}
		}
	}

	return nil
}

// hasIndexAccess checks if a register string contains index access
func (c *GateRegisterSizeMismatchChecker) hasIndexAccess(register string) bool {
	return strings.Contains(register, "[") && strings.Contains(register, "]")
}

// removeComments removes comments from a line
func (c *GateRegisterSizeMismatchChecker) removeComments(line string) string {
	if idx := strings.Index(line, "//"); idx != -1 {
		return line[:idx]
	}
	return line
}