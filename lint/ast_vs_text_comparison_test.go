package lint

import (
	"testing"

	"github.com/orangekame3/qasmtools/lint/ast"
	"github.com/orangekame3/qasmtools/parser"
)

// TestASTvsTextComparison compares AST-based and text-based rule results
func TestASTvsTextComparison(t *testing.T) {
	testCases := []struct {
		name        string
		code        string
		ruleID      string
		description string
	}{
		{
			name:        "QAS001_unused_qubit",
			code:        "OPENQASM 3.0;\nqubit unused_q;\nqubit[2] q;\nh q[0];",
			ruleID:      "QAS001",
			description: "Should detect unused_q as violation",
		},
		{
			name:        "QAS002_undefined_identifier",
			code:        "OPENQASM 3.0;\nqubit q;\nh undefined_gate q;",
			ruleID:      "QAS002", 
			description: "Should detect undefined identifier",
		},
		{
			name:        "QAS003_constant_measurement",
			code:        "OPENQASM 3.0;\nqubit q;\nbit c;\nmeasure q -> c;",
			ruleID:      "QAS003",
			description: "Should detect measurement of unaffected qubit",
		},
		{
			name:        "QAS004_out_of_bounds",
			code:        "OPENQASM 3.0;\nqubit[2] q;\nh q[3];",
			ruleID:      "QAS004",
			description: "Should detect out of bounds access",
		},
		{
			name:        "QAS005_naming_convention",
			code:        "OPENQASM 3.0;\nqubit BadName;\nh BadName;",
			ruleID:      "QAS005",
			description: "Should detect naming convention violation",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Parse the code
			p := parser.NewParser()
			result := p.ParseWithErrors(tc.code)
			if result.HasErrors() {
				t.Logf("Parse errors: %v", result.Errors)
			}

			if result.Program == nil {
				t.Fatal("Program is nil")
			}

			// Create contexts
			astCtx := &ast.CheckContext{
				File:     "test.qasm",
				Content:  tc.code,
				Program:  result.Program,
				UsageMap: make(map[string][]parser.Node),
			}

			textCtx := &CheckContext{
				File:     "test.qasm",
				Content:  tc.code,
				Program:  result.Program,
				UsageMap: make(map[string][]parser.Node),
			}

			// Test AST rule if available
			var astViolations []*ast.Violation
			if astRule := CreateASTRule(tc.ruleID); astRule != nil {
				astViolations = astRule.CheckAST(result.Program, astCtx)
				t.Logf("AST rule %s found %d violations", tc.ruleID, len(astViolations))
				for _, v := range astViolations {
					t.Logf("  AST: %s", v.Message)
				}
			} else {
				t.Logf("No AST rule available for %s", tc.ruleID)
			}

			// Test text-based rule
			var textViolations []*Violation
			if checker := CreateChecker(&Rule{ID: tc.ruleID}); checker != nil {
				if programChecker, ok := checker.(ProgramChecker); ok {
					textViolations = programChecker.CheckProgram(textCtx)
				}
				t.Logf("Text rule %s found %d violations", tc.ruleID, len(textViolations))
				for _, v := range textViolations {
					t.Logf("  Text: %s", v.Message)
				}
			}

			// Compare results
			if len(astViolations) != len(textViolations) {
				t.Logf("WARNING: AST and text rules produced different counts for %s", tc.ruleID)
				t.Logf("  AST: %d violations, Text: %d violations", len(astViolations), len(textViolations))
				
				// This is not necessarily a failure - AST rules might be more accurate
				// But we should understand why
			}

			// Both should detect violations for these test cases
			if len(astViolations) == 0 && len(textViolations) == 0 {
				t.Errorf("Both AST and text rules failed to detect expected violation for %s: %s", tc.ruleID, tc.description)
			}
		})
	}
}

// TestASTRulesDirectIntegration tests AST rules through the linter
func TestASTRulesDirectIntegration(t *testing.T) {
	code := `OPENQASM 3.0;
qubit unused_q;  // Should trigger QAS001
qubit[2] q;
bit c;           // Should trigger QAS002 (insufficient bits)
h q[0];
measure q[0] -> c;
measure q[1] -> c;  // Two measurements, one bit`

	// Create AST-enabled linter
	linter := NewLinterWithAST("", true)
	err := linter.LoadRules()
	if err != nil {
		t.Fatalf("Failed to load rules: %v", err)
	}

	t.Logf("Linter useAST: %v", linter.useAST)
	t.Logf("Available AST rules: %d", len(linter.astRules))
	for ruleID := range linter.astRules {
		t.Logf("  AST rule: %s", ruleID)
	}

	// Lint the content
	violations, err := linter.LintContent(code, "test.qasm")
	if err != nil {
		t.Fatalf("Failed to lint content: %v", err)
	}

	t.Logf("Found %d violations using AST rules", len(violations))
	for _, v := range violations {
		t.Logf("  %s [%s]: %s", v.Rule.ID, v.Severity, v.Message)
	}

	// Should find some violations
	if len(violations) == 0 {
		t.Error("Expected AST rules to find violations, but none were found")
	}

	// Create text-only linter for comparison
	textLinter := NewLinterWithAST("", false)
	err = textLinter.LoadRules()
	if err != nil {
		t.Fatalf("Failed to load text rules: %v", err)
	}

	textViolations, err := textLinter.LintContent(code, "test.qasm")
	if err != nil {
		t.Fatalf("Failed to lint with text rules: %v", err)
	}

	t.Logf("Found %d violations using text rules", len(textViolations))
	for _, v := range textViolations {
		t.Logf("  %s [%s]: %s", v.Rule.ID, v.Severity, v.Message)
	}

	// Compare the results
	t.Logf("AST vs Text comparison:")
	t.Logf("  AST rules: %d violations", len(violations))
	t.Logf("  Text rules: %d violations", len(textViolations))
}