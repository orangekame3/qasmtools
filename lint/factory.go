package lint

import (
	"github.com/orangekame3/qasmtools/lint/ast"
	"github.com/orangekame3/qasmtools/parser"
)

// CreateChecker creates appropriate checker based on rule ID
// Most rules use AST-based implementation via CreateASTRule()
// Only text-based rules that cannot use AST due to parser limitations are listed here
func CreateChecker(rule *Rule) RuleChecker {
	switch rule.ID {
	// All rules now use AST-based implementation via CreateASTRule()
	default:
		return NewNoOpChecker()
	}
}

// CreateASTRule creates AST-based rules for improved analysis
func CreateASTRule(ruleID string) ast.ASTRule {
	switch ruleID {
	case "QAS001":
		return ast.NewUnusedQubitRule()
	case "QAS002":
		return ast.NewUndefinedIdentifierRule()
	case "QAS003":
		return ast.NewConstantMeasuredBitRule()
	case "QAS004":
		return ast.NewOutOfBoundsIndexRule()
	case "QAS005":
		return ast.NewNamingConventionViolationRule()
	case "QAS006":
		return ast.NewQAS006GateRegisterSizeMismatchRule()
	case "QAS007":
		return ast.NewQAS007GateParameterIndexingRule()
	case "QAS008":
		return ast.NewQAS008QubitDeclaredInLocalScopeRule()
	case "QAS010":
		return ast.NewQAS010InvalidInstructionInGateRule()
	case "QAS011":
		return ast.NewQAS011ReservedPrefixUsageRule()
	case "QAS009":
		return ast.NewQAS009IllegalBreakContinueRule()
	case "QAS012":
		return ast.NewQAS012SnakeCaseRequiredRule()
	// All 12 rules now have AST implementations
	default:
		return nil
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
