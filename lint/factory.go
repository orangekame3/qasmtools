package lint

import "github.com/orangekame3/qasmtools/parser"

// CreateChecker creates appropriate checker based on rule ID using new BaseChecker framework
func CreateChecker(rule *Rule) RuleChecker {
	switch rule.ID {
	case "QAS001":
		return NewUnusedQubitChecker()
	case "QAS002":
		return NewUndefinedIdentifierChecker()
	case "QAS003":
		return NewConstantMeasuredBitChecker()
	case "QAS004":
		return NewOutOfBoundsIndexChecker()
	case "QAS005":
		return NewNamingConventionViolationChecker()
	case "QAS006":
		return NewGateRegisterSizeMismatchChecker()
	case "QAS007":
		return NewGateParameterIndexingChecker()
	case "QAS008":
		return NewQubitDeclaredInLocalScopeChecker()
	case "QAS009":
		return NewIllegalBreakContinueChecker()
	case "QAS010":
		return NewInvalidInstructionInGateChecker()
	case "QAS011":
		return NewReservedPrefixUsageChecker()
	case "QAS012":
		return NewSnakeCaseRequiredChecker()
	default:
		return NewNoOpChecker()
	}
}

// NoOpChecker is a checker that does nothing
type NoOpChecker struct{}

// NewNoOpChecker creates a new NoOpChecker
func NewNoOpChecker() *NoOpChecker {
	return &NoOpChecker{}
}

func (c *NoOpChecker) Check(node parser.Node, context *CheckContext) []*Violation {
	return nil
}
