package lint

import "github.com/orangekame3/qasmtools/parser"

// CreateChecker creates appropriate checker based on rule ID
func CreateChecker(rule *Rule) RuleChecker {
	switch rule.ID {
	case "QAS001":
		return &UnusedQubitChecker{}
	case "QAS002":
		return &UndefinedIdentifierChecker{}
	case "QAS003":
		return &ConstantMeasuredBitChecker{}
	case "QAS004":
		return &OutOfBoundsIndexChecker{}
	case "QAS005":
		return &NamingConventionViolationChecker{}
	case "QAS006":
		return &GateRegisterSizeMismatchChecker{}
	case "QAS007":
		return &GateParameterIndexingChecker{}
	case "QAS008":
		return &QubitDeclaredInLocalScopeChecker{}
	case "QAS009":
		return &IllegalBreakContinueChecker{}
	case "QAS010":
		return &InvalidInstructionInGateChecker{}
	case "QAS011":
		return &ReservedPrefixUsageChecker{}
	case "QAS012":
		return &SnakeCaseRequiredChecker{}
	default:
		return &NoOpChecker{}
	}
}

// NoOpChecker is a checker that does nothing
type NoOpChecker struct{}

func (c *NoOpChecker) Check(node parser.Node, context *CheckContext) []*Violation {
	return nil
}
