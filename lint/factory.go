package lint

import "github.com/orangekame3/qasmtools/parser"

// CreateChecker creates appropriate checker based on rule ID
func CreateChecker(rule *Rule) RuleChecker {
	switch rule.ID {
	case "QAS001":
		return &UnusedQubitChecker{}
	case "QAS002":
		return &InsufficientClassicalBitsChecker{}
	case "QAS003":
		return &ConstantMeasuredBitChecker{}
	case "QAS004":
		return NewExceedingQubitLimitsChecker(100) // Default limit
	case "QAS005":
		return NewNamingConventionChecker("^[a-z][a-zA-Z0-9_]*$")
	default:
		return &NoOpChecker{}
	}
}

// NoOpChecker is a checker that does nothing
type NoOpChecker struct{}

func (c *NoOpChecker) Check(node parser.Node, context *CheckContext) []*Violation {
	return nil
}