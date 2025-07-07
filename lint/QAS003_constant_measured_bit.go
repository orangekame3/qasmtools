package lint

import (
	"regexp"
	"strings"

	"github.com/orangekame3/qasmtools/parser"
)

// ConstantMeasuredBitChecker is the new implementation using BaseChecker framework
type ConstantMeasuredBitChecker struct {
	*BaseChecker
	qubitGateMap map[string]bool // tracks which qubits have had gates applied
}

// NewConstantMeasuredBitChecker creates a new ConstantMeasuredBitChecker
func NewConstantMeasuredBitChecker() *ConstantMeasuredBitChecker {
	return &ConstantMeasuredBitChecker{
		BaseChecker: NewBaseChecker("QAS003"),
	}
}

// CheckFile performs file-level constant measured bit analysis
func (c *ConstantMeasuredBitChecker) CheckFile(context *CheckContext) []*Violation {
	var violations []*Violation

	// Get content using BaseChecker method
	content, err := c.getContent(context)
	if err != nil {
		return violations
	}

	lines := strings.Split(content, "\n")

	// First pass: collect all qubits and their gate applications
	c.qubitGateMap = c.analyzeQubitGateUsage(lines)

	// Second pass: find measurements and check if qubits have gates applied
	for i, line := range lines {
		// Skip comments and empty lines using shared utility
		if SkipCommentAndEmptyLine(line) {
			continue
		}

		// Find measurements in this line
		measurements := c.findMeasurements(line)

		for _, measurement := range measurements {
			qubitName := ExtractRegisterName(measurement.qubitName) // Use shared utility
			if !c.qubitGateMap[qubitName] {
				violation := c.NewViolationBuilder().
					WithMessage("Measuring qubit '"+qubitName+"' that has no gates applied. The result will always be |0⟩.").
					WithFile(context.File).
					WithPosition(i+1, measurement.column).
					WithNodeName(qubitName).
					AsWarning().
					Build()
				violations = append(violations, violation)
			}
		}
	}

	return violations
}

// measurementInfo represents a measurement with position
type measurementInfo struct {
	qubitName string
	column    int
}

// analyzeQubitGateUsage tracks which qubits have had gates applied
func (c *ConstantMeasuredBitChecker) analyzeQubitGateUsage(lines []string) map[string]bool {
	qubitGateMap := make(map[string]bool)

	for _, line := range lines {
		// Skip comments and empty lines using shared utility
		if SkipCommentAndEmptyLine(line) {
			continue
		}

		// Find gate applications using pattern matching
		c.analyzeGateApplicationsInLine(line, qubitGateMap)
	}

	return qubitGateMap
}

// analyzeGateApplicationsInLine finds gate applications in a single line
func (c *ConstantMeasuredBitChecker) analyzeGateApplicationsInLine(line string, qubitGateMap map[string]bool) {
	// Remove comments using shared utility
	cleanLine := RemoveComments(line)

	// Pattern for gates: gate_name(optional_params) qubits;
	// This handles both parameterized gates like rx(pi/2) q; and regular gates like h q;
	gatePattern := regexp.MustCompile(`\b([a-zA-Z_][a-zA-Z0-9_]*)(?:\([^)]*\))?\s+([a-zA-Z_][a-zA-Z0-9_]*(?:\[[^\]]*\])?(?:\s*,\s*[a-zA-Z_][a-zA-Z0-9_]*(?:\[[^\]]*\])?)*)\s*;`)
	matches := gatePattern.FindAllStringSubmatch(cleanLine, -1)

	for _, match := range matches {
		if len(match) >= 3 {
			gateName := match[1]
			qubitsStr := match[2]

			// Skip if this is not actually a gate (e.g., declarations, measurements)
			if IsKeyword(gateName) || gateName == "measure" {
				continue
			}

			// Extract individual qubit names from the qubit string
			qubitNames := strings.Split(qubitsStr, ",")
			for _, qubitName := range qubitNames {
				qubitName = strings.TrimSpace(qubitName)
				baseQubitName := ExtractRegisterName(qubitName) // Use shared utility
				if baseQubitName != "" {
					qubitGateMap[baseQubitName] = true
				}
			}
		}
	}
}

// findMeasurements finds all measurements in a line
func (c *ConstantMeasuredBitChecker) findMeasurements(line string) []measurementInfo {
	var measurements []measurementInfo

	// Remove comments using shared utility
	cleanLine := RemoveComments(line)

	// Pattern for measurements: measure qubit -> bit;
	measurePattern := regexp.MustCompile(`\bmeasure\s+([a-zA-Z_][a-zA-Z0-9_]*(?:\[[^\]]*\])?)\s*(?:->\s*[^;]+)?\s*;`)
	matches := measurePattern.FindAllStringSubmatch(cleanLine, -1)
	indices := measurePattern.FindAllStringIndex(cleanLine, -1)

	for i, match := range matches {
		if len(match) >= 2 {
			qubitName := match[1]
			column := indices[i][0] + 1 // Convert to 1-based indexing

			measurements = append(measurements, measurementInfo{
				qubitName: qubitName,
				column:    column,
			})
		}
	}

	return measurements
}

// Check implements RuleChecker interface (required but delegates to CheckProgram)
func (c *ConstantMeasuredBitChecker) Check(node parser.Node, context *CheckContext) []*Violation {
	return nil
}

// CheckProgram implements ProgramChecker interface
func (c *ConstantMeasuredBitChecker) CheckProgram(context *CheckContext) []*Violation {
	return c.CheckFile(context)
}
