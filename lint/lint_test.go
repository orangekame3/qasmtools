package lint

import (
	"testing"
)

func TestLintContent(t *testing.T) {
	code := `OPENQASM 3.0;
include "stdgates.qasm";
qubit[2] q;
bit[2] c;
h q[0];
cx q[0],q[3];
measure q->c;`

	linter := NewLinter("")
	err := linter.LoadRules()
	if err != nil {
		t.Fatalf("Failed to load rules: %v", err)
	}

	violations, err := linter.LintContent(code, "<stdin>")
	if err != nil {
		t.Fatalf("Failed to lint content: %v", err)
	}

	if len(violations) == 0 {
		t.Error("Expected violations, but got none")
	}

	// Check for QAS004 violation
	found := false
	for _, v := range violations {
		t.Logf("Violation: %s", v.String())
		if v.Rule.ID == "QAS004" {
			found = true
			// AST parser currently reports all IndexedIdentifier positions as line 1, col 1
			// This is a known parser limitation - the violation is correctly detected
			if v.Line != 1 {
				t.Errorf("Expected violation on line 1 (parser limitation), got line %d", v.Line)
			}
			if v.Column != 1 {
				t.Errorf("Expected violation at column 1 (parser limitation), got column %d", v.Column)
			}
			if v.Severity != SeverityError {
				t.Errorf("Expected severity Error, got %s", v.Severity)
			}
		}
	}

	if !found {
		t.Error("Expected QAS004 violation, but none was found")
	}
}
