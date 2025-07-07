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
		if v.Rule.ID == "QAS004" {
			found = true
			if v.Line != 6 {
				t.Errorf("Expected violation on line 6, got line %d", v.Line)
			}
			if v.Column != 11 {
				t.Errorf("Expected violation at column 11, got column %d", v.Column)
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
