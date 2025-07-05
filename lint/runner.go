package lint

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/orangekame3/qasmtools/parser"
)

// Linter is the main linter engine
type Linter struct {
	rules    []*Rule
	loader   *RuleLoader
	checkers map[string]RuleChecker
}

// NewLinter creates a new linter instance
func NewLinter(rulesDir string) *Linter {
	return &Linter{
		loader:   NewRuleLoader(rulesDir),
		checkers: make(map[string]RuleChecker),
	}
}

// LoadRules loads all rules and creates corresponding checkers
func (l *Linter) LoadRules() error {
	rules, err := l.loader.LoadRules()
	if err != nil {
		return err
	}

	l.rules = rules

	// Create checkers for each rule
	for _, rule := range rules {
		checker := CreateChecker(rule)
		l.checkers[rule.ID] = checker
	}

	return nil
}

// GetRules returns the loaded rules
func (l *Linter) GetRules() []*Rule {
	return l.rules
}

// LintContent lints QASM content directly from a string
func (l *Linter) LintContent(content string, filename string) ([]*Violation, error) {
	// Parse the content
	p := parser.NewParser()
	result := p.ParseWithErrors(content)
	if result.HasErrors() {
		return nil, fmt.Errorf("failed to parse content: %v", result.Errors)
	}
	if result.Program == nil {
		return nil, fmt.Errorf("failed to parse content: program is nil")
	}

	// Build usage map for symbol tracking
	usageMap := l.buildUsageMap(result.Program)

	// Create check context
	context := &CheckContext{
		File:     filename,
		Content:  content,
		Program:  result.Program,
		UsageMap: usageMap,
	}

	var allViolations []*Violation

	// Run each rule against the AST (same logic as LintFile)
	for _, rule := range l.rules {
		if !rule.Enabled {
			continue
		}

		checker := l.checkers[rule.ID]
		if checker == nil {
			continue
		}

		violations := l.runRuleOnProgram(rule, checker, result.Program, context)

		// Set rule reference for each violation
		for _, violation := range violations {
			violation.Rule = rule
			violation.Severity = rule.Level
		}

		allViolations = append(allViolations, violations...)
	}

	return allViolations, nil
}

// LintFile lints a single QASM file
func (l *Linter) LintFile(filename string) ([]*Violation, error) {
	// Parse the file
	p := parser.NewParser()
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	result := p.ParseWithErrors(string(content))
	if result.HasErrors() {
		return nil, fmt.Errorf("failed to parse file: %v", result.Errors)
	}
	if result.Program == nil {
		return nil, fmt.Errorf("failed to parse file: program is nil")
	}

	// Build usage map for symbol tracking
	usageMap := l.buildUsageMap(result.Program)

	// Create check context
	context := &CheckContext{
		File:     filename,
		Content:  string(content),
		Program:  result.Program,
		UsageMap: usageMap,
	}

	var allViolations []*Violation

	// Run each rule against the AST
	for _, rule := range l.rules {
		checker := l.checkers[rule.ID]
		violations := l.runRuleOnProgram(rule, checker, result.Program, context)

		// Set rule reference for each violation
		for _, violation := range violations {
			violation.Rule = rule
			violation.Severity = rule.Level
		}

		allViolations = append(allViolations, violations...)
	}

	return allViolations, nil
}

// runRuleOnProgram runs a single rule against the entire program
func (l *Linter) runRuleOnProgram(rule *Rule, checker RuleChecker, program *parser.Program, context *CheckContext) []*Violation {
	var violations []*Violation

	// For some rules that need program-level analysis, check if they implement ProgramChecker
	if programChecker, ok := checker.(ProgramChecker); ok {
		programViolations := programChecker.CheckProgram(context)
		return programViolations
	}

	// Check version statement
	if program.Version != nil {
		if l.matchesRule(rule, program.Version) {
			nodeViolations := checker.Check(program.Version, context)
			violations = append(violations, nodeViolations...)
		}
	}

	// Check all statements (includes Include statements)
	for _, stmt := range program.Statements {
		if l.matchesRule(rule, stmt) {
			nodeViolations := checker.Check(stmt, context)
			violations = append(violations, nodeViolations...)
		}
	}

	return violations
}

// matchesRule checks if an AST node matches the rule's match criteria
func (l *Linter) matchesRule(rule *Rule, node parser.Node) bool {
	switch rule.Match.Type {
	case "declaration":
		switch rule.Match.Kind {
		case "qubit":
			if decl, ok := node.(*parser.QuantumDeclaration); ok {
				return decl.Type == "qubit"
			}
			return false
		case "gate":
			_, ok := node.(*parser.GateDefinition)
			return ok
		case "classical":
			_, ok := node.(*parser.ClassicalDeclaration)
			return ok
		}
	case "statement":
		// All statements match
		return true
	case "expression":
		// Handle expression matching if needed
		return false
	}

	return false
}

// buildUsageMap builds a map of symbol names to their usage locations
func (l *Linter) buildUsageMap(program *parser.Program) map[string][]parser.Node {
	usageMap := make(map[string][]parser.Node)

	// Walk through all statements and collect symbol usage
	for _, stmt := range program.Statements {
		l.collectSymbolUsage(stmt, usageMap)
	}

	return usageMap
}

// collectSymbolUsage recursively collects symbol usage from AST nodes
func (l *Linter) collectSymbolUsage(node parser.Node, usageMap map[string][]parser.Node) {
	switch n := node.(type) {
	case *parser.GateCall:
		// Record usage of qubits in gate calls
		for _, qubit := range n.Qubits {
			if id, ok := qubit.(*parser.Identifier); ok {
				// Extract base name for array access (e.g., "q[0]" -> "q")
				baseName := l.extractBaseName(id.Name)
				usageMap[baseName] = append(usageMap[baseName], node)
			}
		}
	case *parser.Measurement:
		// Record usage of qubits and classical bits in measurements
		if id, ok := n.Qubit.(*parser.Identifier); ok {
			baseName := l.extractBaseName(id.Name)
			usageMap[baseName] = append(usageMap[baseName], node)
		}
		if id, ok := n.Target.(*parser.Identifier); ok {
			baseName := l.extractBaseName(id.Name)
			usageMap[baseName] = append(usageMap[baseName], node)
		}
		// Add more cases as needed for other statement types
	}
}

// extractBaseName extracts the base identifier name, removing array access
func (l *Linter) extractBaseName(name string) string {
	if idx := strings.Index(name, "["); idx != -1 {
		return name[:idx]
	}
	return name
}

// LintFiles lints multiple files
func (l *Linter) LintFiles(filenames []string) ([]*Violation, error) {
	var allViolations []*Violation

	for _, filename := range filenames {
		violations, err := l.LintFile(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to lint %s: %w", filename, err)
		}
		allViolations = append(allViolations, violations...)
	}

	return allViolations, nil
}

// LintDirectory lints all QASM files in a directory
func (l *Linter) LintDirectory(dir string) ([]*Violation, error) {
	pattern := filepath.Join(dir, "*.qasm")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	return l.LintFiles(files)
}
