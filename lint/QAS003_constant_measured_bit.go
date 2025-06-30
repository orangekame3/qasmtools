package lint

import (
	"os"
	"regexp"
	"strings"

	"github.com/orangekame3/qasmtools/parser"
)

// ConstantMeasuredBitChecker checks for measurements of unaffected qubits (QAS003)
type ConstantMeasuredBitChecker struct{}

func (c *ConstantMeasuredBitChecker) Check(node parser.Node, context *CheckContext) []*Violation {
	// This method is required by RuleChecker but not used for program-level analysis
	return nil
}

// CheckProgram implements ProgramChecker interface for program-level analysis
func (c *ConstantMeasuredBitChecker) CheckProgram(context *CheckContext) []*Violation {
	// Always use text-based analysis since AST parsing has issues (see CLAUDE.md)
	return c.checkFileTextBased(context)
}

func (c *ConstantMeasuredBitChecker) checkFileTextBased(context *CheckContext) []*Violation {
	var violations []*Violation

	// Read file content
	content, err := os.ReadFile(context.File)
	if err != nil {
		return violations
	}

	text := string(content)
	lines := strings.Split(text, "\n")

	// Find all measurements and check if qubits were affected by gates
	measurePattern := regexp.MustCompile(`^\s*measure\s+([a-zA-Z_][a-zA-Z0-9_]*(?:\[\d+\])?)`)
	for i, line := range lines {
		if matches := measurePattern.FindStringSubmatch(line); len(matches) > 1 {
			qubitRef := matches[1]
			// Extract base qubit name (remove array access)
			qubitName := qubitRef
			if idx := strings.Index(qubitRef, "["); idx != -1 {
				qubitName = qubitRef[:idx]
			}
			
			// Check if this qubit was affected by any gates
			if !c.isQubitAffectedByGatesTextBased(qubitName, text) {
				violation := &Violation{
					File:     context.File,
					Line:     i + 1,
					Column:   1,
					NodeName: qubitName,
					Message:  "Measuring qubit '" + qubitName + "' that has not been affected by any gates will always yield |0⟩.",
					Severity: SeverityWarning,
				}
				violations = append(violations, violation)
			}
		}
	}

	return violations
}

func (c *ConstantMeasuredBitChecker) isQubitAffectedByGatesTextBased(qubitName string, content string) bool {
	lines := strings.Split(content, "\n")
	
	for _, line := range lines {
		// Skip measurement lines and declarations
		if strings.Contains(line, "measure") || strings.Contains(line, "qubit") {
			continue
		}
		
		// Look for gate applications to this qubit
		// Pattern: gate_name qubit_name; or gate_name qubit_name[index];
		gatePatterns := []string{
			`\b[a-z]+\s+` + regexp.QuoteMeta(qubitName) + `\[\d+\]`,  // Array access in gate
			`\b[a-z]+\s+` + regexp.QuoteMeta(qubitName) + `\b`,       // Direct gate application
			`\b` + regexp.QuoteMeta(qubitName) + `\[\d+\]\s*,`,       // First parameter in multi-qubit gate
			`\b` + regexp.QuoteMeta(qubitName) + `\s*,`,              // First parameter 
			`,\s*` + regexp.QuoteMeta(qubitName) + `\[\d+\]`,        // Second parameter with array
			`,\s*` + regexp.QuoteMeta(qubitName) + `\b`,             // Second parameter
		}
		
		for _, pattern := range gatePatterns {
			matched, _ := regexp.MatchString(pattern, line)
			if matched {
				return true
			}
		}
	}
	
	return false
}