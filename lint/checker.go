package lint

import (
	"fmt"
	"regexp"

	"github.com/orangekame3/qasmtools/parser"
)

// UnusedQubitChecker checks for unused qubit declarations
type UnusedQubitChecker struct{}

func (c *UnusedQubitChecker) Check(node parser.Node, context *CheckContext) []*Violation {
	var violations []*Violation

	// Only check qubit declarations
	if decl, ok := node.(*parser.QuantumDeclaration); ok && decl.Type == "qubit" {
		qubitName := decl.Identifier
		
		// Check if the qubit is used anywhere
		if usages, exists := context.UsageMap[qubitName]; !exists || len(usages) == 0 {
			violation := &Violation{
				File:     context.File,
				Line:     decl.Position.Line,
				Column:   decl.Position.Column,
				NodeName: qubitName,
				Message:  "Qubit '" + qubitName + "' is declared but never used.",
				Severity: SeverityWarning,
			}
			violations = append(violations, violation)
		}
	}

	return violations
}

// GateNamingChecker checks gate naming conventions
type GateNamingChecker struct {
	Pattern *regexp.Regexp
}

func NewGateNamingChecker(pattern string) *GateNamingChecker {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		// Fallback to default pattern
		regex = regexp.MustCompile("^[a-z][a-z0-9_]*$")
	}
	
	return &GateNamingChecker{
		Pattern: regex,
	}
}

func (c *GateNamingChecker) Check(node parser.Node, context *CheckContext) []*Violation {
	var violations []*Violation

	// Only check gate declarations
	if gate, ok := node.(*parser.GateDefinition); ok {
		gateName := gate.Name
		
		if !c.Pattern.MatchString(gateName) {
			violation := &Violation{
				File:     context.File,
				Line:     gate.Position.Line,
				Column:   gate.Position.Column,
				NodeName: gateName,
				Message:  "Gate '" + gateName + "' does not follow naming convention.",
				Severity: SeverityWarning,
			}
			violations = append(violations, violation)
		}
	}

	return violations
}

// MaxQubitsChecker checks for maximum qubit count
type MaxQubitsChecker struct {
	MaxCount int
}

func NewMaxQubitsChecker(maxCount int) *MaxQubitsChecker {
	return &MaxQubitsChecker{
		MaxCount: maxCount,
	}
}

func (c *MaxQubitsChecker) Check(node parser.Node, context *CheckContext) []*Violation {
	var violations []*Violation

	// Count qubits in the entire program
	qubitCount := 0
	for _, stmt := range context.Program.Statements {
		if decl, ok := stmt.(*parser.QuantumDeclaration); ok && decl.Type == "qubit" {
			// TODO: Handle array sizes properly
			qubitCount += 1 // For now, count each declaration as 1
		}
	}

	if qubitCount > c.MaxCount {
		// Only report once per program
		if decl, ok := node.(*parser.QuantumDeclaration); ok && decl.Type == "qubit" {
			violation := &Violation{
				File:     context.File,
				Line:     decl.Position.Line,
				Column:   decl.Position.Column,
				NodeName: "program",
				Message:  fmt.Sprintf("Program uses %d qubits, exceeding maximum of %d", qubitCount, c.MaxCount),
				Severity: SeverityWarning,
			}
			violations = append(violations, violation)
		}
	}

	return violations
}

// CreateChecker creates a checker based on rule configuration
func CreateChecker(rule *Rule) RuleChecker {
	switch rule.Match.Type {
	case "declaration":
		switch rule.Match.Kind {
		case "qubit":
			for _, check := range rule.Check {
				switch check.Type {
				case "usage":
					if check.NotFound {
						return &UnusedQubitChecker{}
					}
				case "count":
					if check.Max > 0 {
						return NewMaxQubitsChecker(check.Max)
					}
				}
			}
		case "gate":
			for _, check := range rule.Check {
				if check.Type == "naming" && check.Pattern != "" {
					return NewGateNamingChecker(check.Pattern)
				}
			}
		}
	}

	// Return a no-op checker if no specific checker is found
	return &NoOpChecker{}
}

// NoOpChecker is a checker that does nothing
type NoOpChecker struct{}

func (c *NoOpChecker) Check(node parser.Node, context *CheckContext) []*Violation {
	return nil
}