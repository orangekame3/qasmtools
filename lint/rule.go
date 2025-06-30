package lint

import (
	"fmt"

	"github.com/orangekame3/qasmtools/parser"
)

// Severity represents the severity level of a lint rule
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
	SeverityInfo    Severity = "info"
)

// Rule represents a lint rule loaded from YAML
type Rule struct {
	ID          string   `yaml:"id"`
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Level       Severity `yaml:"level"`
	Enabled     bool     `yaml:"enabled"`
	Match       Match    `yaml:"match"`
	Check       []Check  `yaml:"check"`
	Message     string   `yaml:"message"`
	Tags        []string `yaml:"tags"`
	Fixable     bool     `yaml:"fixable"`
}

// Match defines what AST nodes to match
type Match struct {
	Type string `yaml:"type"` // declaration, statement, expression, etc.
	Kind string `yaml:"kind"` // qubit, gate, measure, etc.
}

// Check defines what to check on matched nodes
type Check struct {
	Type     string `yaml:"type"`      // usage, naming, count, etc.
	NotFound bool   `yaml:"not_found"` // for usage checks
	Pattern  string `yaml:"pattern"`   // for naming checks
	Max      int    `yaml:"max"`       // for count checks
	Min      int    `yaml:"min"`       // for count checks
}

// Violation represents a rule violation
type Violation struct {
	Rule     *Rule
	Message  string
	File     string
	Line     int
	Column   int
	Severity Severity
	NodeName string
}

// String returns a formatted string representation of the violation
func (v *Violation) String() string {
	return fmt.Sprintf("%s:%d:%d: %s [%s] %s",
		v.File, v.Line, v.Column, v.Severity, v.Rule.ID, v.Message)
}

// RuleChecker checks a specific rule against AST nodes
type RuleChecker interface {
	Check(node parser.Node, context *CheckContext) []*Violation
}

// CheckContext provides context for rule checking
type CheckContext struct {
	File     string
	Program  *parser.Program
	UsageMap map[string][]parser.Node // For tracking symbol usage
}
