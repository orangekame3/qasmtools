package lint

import (
	"regexp"
	"strings"

	"github.com/orangekame3/qasmtools/parser"
)

// GateParameterIndexingChecker is the new implementation using BaseChecker framework
type GateParameterIndexingChecker struct {
	*BaseChecker
}

// NewGateParameterIndexingChecker creates a new GateParameterIndexingChecker
func NewGateParameterIndexingChecker() *GateParameterIndexingChecker {
	return &GateParameterIndexingChecker{
		BaseChecker: NewBaseChecker("QAS007"),
	}
}

// CheckFile performs file-level gate parameter indexing analysis
func (c *GateParameterIndexingChecker) CheckFile(context *CheckContext) []*Violation {
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
		// Check each line within the gate body for parameter indexing
		for i := gateDef.startLine; i <= gateDef.endLine; i++ {
			if i < len(lines) {
				line := lines[i]

				// Find any index accesses in this line
				indexAccesses := c.findIndexAccesses(line)

				for _, access := range indexAccesses {
					// Check if the accessed register is a gate parameter
					if c.isGateParameter(access.registerName, gateDef.parameters) {
						violation := c.NewViolationBuilder().
							WithMessage("Cannot perform index access on gate argument '"+access.registerName+"'.").
							WithFile(context.File).
							WithPosition(i+1, access.column).
							WithNodeName(access.registerName).
							AsError().
							Build()
						violations = append(violations, violation)
					}
				}
			}
		}
	}

	return violations
}

type gateDefWithBody struct {
	name       string
	parameters []string
	startLine  int
	endLine    int
}

type indexAccess struct {
	registerName string
	index        string
	column       int
}

// findGateDefinitions finds all gate definitions in the file
func (c *GateParameterIndexingChecker) findGateDefinitions(lines []string) []gateDefWithBody {
	var gates []gateDefWithBody

	// Pattern for gate definitions: gate name(params...) qubits... {
	gateStartPattern := regexp.MustCompile(`^\s*gate\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*(?:\([^)]*\))?\s+([^{]+)\s*\{`)

	for i, line := range lines {
		// Skip comments and empty lines using shared utility
		if SkipCommentAndEmptyLine(line) {
			continue
		}

		matches := gateStartPattern.FindStringSubmatch(line)
		if len(matches) >= 3 {
			gateName := matches[1]
			qubitParams := strings.TrimSpace(matches[2])

			// Parse qubit parameters
			parameters := c.parseGateParameters(qubitParams)

			// Find the end of this gate definition
			endLine := c.findGateEnd(lines, i)

			gates = append(gates, gateDefWithBody{
				name:       gateName,
				parameters: parameters,
				startLine:  i,
				endLine:    endLine,
			})
		}
	}

	return gates
}

// parseGateParameters parses the qubit parameters from a gate definition
func (c *GateParameterIndexingChecker) parseGateParameters(paramStr string) []string {
	var parameters []string

	if paramStr == "" {
		return parameters
	}

	// Split by commas and extract parameter names
	parts := strings.Split(paramStr, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			// Extract just the parameter name (without any type information)
			paramName := c.extractParameterName(part)
			if paramName != "" {
				parameters = append(parameters, paramName)
			}
		}
	}

	return parameters
}

// extractParameterName extracts the parameter name from a parameter declaration
func (c *GateParameterIndexingChecker) extractParameterName(param string) string {
	// Parameter can be just a name: "q" or "qubit"
	// For now, assume it's just the identifier
	param = strings.TrimSpace(param)

	// Simple pattern to match identifier
	pattern := regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)`)
	matches := pattern.FindStringSubmatch(param)
	if len(matches) >= 2 {
		return matches[1]
	}

	return ""
}

// findGateEnd finds the closing brace of a gate definition
func (c *GateParameterIndexingChecker) findGateEnd(lines []string, startLine int) int {
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

// findIndexAccesses finds all index accesses in a line
func (c *GateParameterIndexingChecker) findIndexAccesses(line string) []indexAccess {
	var accesses []indexAccess

	// Remove comments using shared utility
	codeOnly := RemoveComments(line)

	// Use shared pattern for array access
	matches := ArrayAccessPattern.FindAllStringSubmatch(codeOnly, -1)
	indices := ArrayAccessPattern.FindAllStringIndex(codeOnly, -1)

	for i, match := range matches {
		if len(match) >= 3 {
			registerName := match[1]
			indexValue := match[2]
			column := indices[i][0] + 1 // Convert to 1-based indexing

			accesses = append(accesses, indexAccess{
				registerName: registerName,
				index:        indexValue,
				column:       column,
			})
		}
	}

	return accesses
}

// isGateParameter checks if a register name is a parameter of the gate
func (c *GateParameterIndexingChecker) isGateParameter(registerName string, parameters []string) bool {
	for _, param := range parameters {
		if param == registerName {
			return true
		}
	}
	return false
}

// Check implements RuleChecker interface (required but delegates to CheckProgram)
func (c *GateParameterIndexingChecker) Check(node parser.Node, context *CheckContext) []*Violation {
	return nil
}

// CheckProgram implements ProgramChecker interface
func (c *GateParameterIndexingChecker) CheckProgram(context *CheckContext) []*Violation {
	return c.CheckFile(context)
}
