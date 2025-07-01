package lint

import (
	"os"
	"regexp"
	"strings"

	"github.com/orangekame3/qasmtools/parser"
)

// GateParameterIndexingChecker checks for prohibited index access on gate parameters (QAS007)
type GateParameterIndexingChecker struct{}

func (c *GateParameterIndexingChecker) Check(node parser.Node, context *CheckContext) []*Violation {
	// This method is required by RuleChecker but not used for program-level analysis
	return nil
}

// CheckProgram implements ProgramChecker interface for program-level analysis
func (c *GateParameterIndexingChecker) CheckProgram(context *CheckContext) []*Violation {
	// Use text-based analysis due to AST parsing issues
	return c.CheckFile(context)
}

// CheckFile performs file-level gate parameter indexing analysis
func (c *GateParameterIndexingChecker) CheckFile(context *CheckContext) []*Violation {
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
		// Check each line within the gate body for parameter indexing
		for i := gateDef.startLine; i <= gateDef.endLine; i++ {
			if i < len(lines) {
				line := lines[i]

				// Find any index accesses in this line
				indexAccesses := c.findIndexAccesses(line)

				for _, access := range indexAccesses {
					// Check if the accessed register is a gate parameter
					if c.isGateParameter(access.registerName, gateDef.parameters) {
						violation := &Violation{
							Rule:     nil, // Will be set by the runner
							File:     context.File,
							Line:     i + 1,
							Column:   access.column,
							NodeName: access.registerName,
							Message:  "Cannot perform index access on gate argument '" + access.registerName + "'.",
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
		// Skip comments and empty lines
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "//") {
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

// findIndexAccesses finds all index accesses in a line
func (c *GateParameterIndexingChecker) findIndexAccesses(line string) []indexAccess {
	var accesses []indexAccess

	// Remove comments from the line
	codeOnly := c.removeComments(line)

	// Pattern for index access: identifier[something]
	indexPattern := regexp.MustCompile(`\b([a-zA-Z_][a-zA-Z0-9_]*)\[([^\]]*)\]`)
	matches := indexPattern.FindAllStringSubmatch(codeOnly, -1)
	indices := indexPattern.FindAllStringIndex(codeOnly, -1)

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

// removeComments removes comments from a line
func (c *GateParameterIndexingChecker) removeComments(line string) string {
	if idx := strings.Index(line, "//"); idx != -1 {
		return line[:idx]
	}
	return line
}
