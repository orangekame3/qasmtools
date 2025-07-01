package lint

import (
	"os"
	"regexp"
	"strings"

	"github.com/orangekame3/qasmtools/parser"
)

// ConstantMeasuredBitChecker checks for measurements of qubits with no gates applied (QAS003)
type ConstantMeasuredBitChecker struct{}

func (c *ConstantMeasuredBitChecker) Check(node parser.Node, context *CheckContext) []*Violation {
	// This method is required by RuleChecker but not used for program-level analysis
	return nil
}

// CheckProgram implements ProgramChecker interface for program-level analysis
func (c *ConstantMeasuredBitChecker) CheckProgram(context *CheckContext) []*Violation {
	// Use text-based analysis due to AST parsing issues
	return c.CheckFile(context)
}

// CheckFile performs file-level constant measured bit analysis
func (c *ConstantMeasuredBitChecker) CheckFile(context *CheckContext) []*Violation {
	var violations []*Violation

	// Read file content for text-based analysis
	content, err := os.ReadFile(context.File)
	if err != nil {
		return violations
	}

	text := string(content)
	lines := strings.Split(text, "\n")

	// First pass: collect all qubits and their gate applications
	qubitGateMap := c.analyzeQubitGateUsage(lines)

	// Second pass: find measurements and check if qubits have gates applied
	for i, line := range lines {
		// Skip comments and empty lines
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "//") {
			continue
		}

		// Find measurements in this line
		measurements := c.findMeasurements(line)
		
		for _, measurement := range measurements {
			// Check if the measured qubit has any gates applied
			if !qubitGateMap[measurement.qubitName] {
				violation := &Violation{
					Rule:     nil, // Will be set by the runner
					File:     context.File,
					Line:     i + 1,
					Column:   measurement.column,
					NodeName: measurement.qubitName,
					Message:  "Measuring qubit '" + measurement.qubitName + "' that has no gates applied. The result will always be |0⟩.",
					Severity: SeverityWarning,
				}
				violations = append(violations, violation)
			}
		}
	}

	return violations
}

type measurementInfo struct {
	qubitName string
	column    int
}

// analyzeQubitGateUsage analyzes which qubits have gates applied to them
func (c *ConstantMeasuredBitChecker) analyzeQubitGateUsage(lines []string) map[string]bool {
	qubitGateMap := make(map[string]bool)

	// First, find all qubit declarations to initialize the map
	for _, line := range lines {
		qubits := c.findQubitDeclarations(line)
		for _, qubit := range qubits {
			qubitGateMap[qubit] = false // Initially, no gates applied
		}
	}

	// Then, find all gate applications
	gatePatterns := []*regexp.Regexp{
		// Single qubit gates: h q; x q; etc.
		regexp.MustCompile(`\b([a-z]+)\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*(?:\[\s*\d+\s*\])?\s*;`),
		// Two qubit gates: cx q1, q2; etc.
		regexp.MustCompile(`\b([a-z]+)\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*(?:\[\s*\d+\s*\])?\s*,\s*([a-zA-Z_][a-zA-Z0-9_]*)\s*(?:\[\s*\d+\s*\])?\s*;`),
		// Parameterized gates: rx(pi/2) q; etc.
		regexp.MustCompile(`\b([a-z]+)\s*\([^)]+\)\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*(?:\[\s*\d+\s*\])?\s*;`),
	}

	for _, line := range lines {
		// Skip comments and empty lines
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "//") {
			continue
		}

		// Remove comments from the line
		codeOnly := c.removeComments(line)

		for _, pattern := range gatePatterns {
			matches := pattern.FindAllStringSubmatch(codeOnly, -1)
			for _, match := range matches {
				if len(match) >= 3 {
					gateName := match[1]
					
					// Skip non-gate keywords
					if c.isNonGateKeyword(gateName) {
						continue
					}

					// Extract qubit name (remove array indices)
					qubitName := c.extractQubitName(match[2])
					if qubitName != "" {
						qubitGateMap[qubitName] = true
					}

					// For two-qubit gates, mark both qubits
					if len(match) >= 4 && match[3] != "" {
						qubitName2 := c.extractQubitName(match[3])
						if qubitName2 != "" {
							qubitGateMap[qubitName2] = true
						}
					}
				}
			}
		}
	}

	return qubitGateMap
}

// findQubitDeclarations finds all qubit declarations in a line
func (c *ConstantMeasuredBitChecker) findQubitDeclarations(line string) []string {
	var qubits []string

	// Patterns for qubit declarations
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`^\s*qubit\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*;`),        // single qubit
		regexp.MustCompile(`^\s*qubit\[\s*\d+\s*\]\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*;`), // array qubit
	}

	for _, pattern := range patterns {
		matches := pattern.FindStringSubmatch(line)
		if len(matches) > 1 {
			qubits = append(qubits, matches[1])
		}
	}

	return qubits
}

// findMeasurements finds all measurements in a line
func (c *ConstantMeasuredBitChecker) findMeasurements(line string) []measurementInfo {
	var measurements []measurementInfo

	// Remove comments from the line
	codeOnly := c.removeComments(line)

	// Pattern for measure statements: measure qubit -> bit;
	measurePattern := regexp.MustCompile(`\bmeasure\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*(?:\[\s*\d+\s*\])?\s*->\s*`)
	matches := measurePattern.FindAllStringSubmatch(codeOnly, -1)
	indices := measurePattern.FindAllStringIndex(codeOnly, -1)

	for i, match := range matches {
		if len(match) > 1 {
			qubitName := c.extractQubitName(match[1])
			column := indices[i][0] + 1 // Convert to 1-based indexing

			measurements = append(measurements, measurementInfo{
				qubitName: qubitName,
				column:    column,
			})
		}
	}

	return measurements
}

// extractQubitName extracts the base name from qubit identifier (removes array brackets)
func (c *ConstantMeasuredBitChecker) extractQubitName(identifier string) string {
	if idx := strings.Index(identifier, "["); idx != -1 {
		return identifier[:idx]
	}
	return identifier
}

// isNonGateKeyword checks if an identifier is a keyword that's not a gate
func (c *ConstantMeasuredBitChecker) isNonGateKeyword(identifier string) bool {
	nonGateKeywords := map[string]bool{
		"measure":  true,
		"reset":    true,
		"barrier":  true,
		"delay":    true,
		"if":       true,
		"else":     true,
		"for":      true,
		"while":    true,
		"break":    true,
		"continue": true,
		"include":  true,
		"qubit":    true,
		"bit":      true,
		"gate":     true,
		"def":      true,
		"let":      true,
		"const":    true,
	}

	return nonGateKeywords[identifier]
}

// removeComments removes comments from a line
func (c *ConstantMeasuredBitChecker) removeComments(line string) string {
	if idx := strings.Index(line, "//"); idx != -1 {
		return line[:idx]
	}
	return line
}