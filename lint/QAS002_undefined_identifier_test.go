package lint

import (
	"os"
	"testing"
)

func TestQAS002_UndefinedIdentifier(t *testing.T) {
	tests := []struct {
		name               string
		input              string
		expectedViolations int
		expectedMessage    string
	}{
		{
			name: "undefined qubit usage should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

h q[0];  // q is not declared`,
			expectedViolations: 1,
			expectedMessage:    "Identifier 'q' is not declared.",
		},
		{
			name: "declared qubit usage should not trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit q;
h q;`,
			expectedViolations: 0,
		},
		{
			name: "undefined gate usage should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit q;
my_gate q;  // my_gate is not defined`,
			expectedViolations: 1,
			expectedMessage:    "Identifier 'my_gate' is not declared.",
		},
		{
			name: "standard gates should not trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit q;
h q;
x q;
cx q, q;`,
			expectedViolations: 0,
		},
		{
			name: "defined gate usage should not trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

gate my_gate(q) {
    h q;
}

qubit q;
my_gate q;`,
			expectedViolations: 0,
		},
		{
			name: "multiple undefined identifiers should trigger multiple violations",
			input: `OPENQASM 3.0;

undefined_gate1 q1;
undefined_gate2 q2;`,
			expectedViolations: 4, // undefined_gate1, q1, undefined_gate2, q2
		},
		{
			name: "keywords should not trigger violations",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit q;
if (true) {
    h q;
}`,
			expectedViolations: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tmpFile, err := os.CreateTemp("", "test_*.qasm")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			// Write test input
			if _, err := tmpFile.WriteString(tt.input); err != nil {
				t.Fatalf("Failed to write test input: %v", err)
			}
			tmpFile.Close()

			// Create checker
			checker := NewUndefinedIdentifierChecker()

			// Create context
			context := &CheckContext{
				File: tmpFile.Name(),
			}

			// Run check
			violations := checker.CheckProgram(context)

			// Verify results
			if len(violations) != tt.expectedViolations {
				t.Errorf("Expected %d violations, got %d", tt.expectedViolations, len(violations))
				for i, v := range violations {
					t.Errorf("Violation %d: %s", i, v.Message)
				}
				return
			}

			if tt.expectedViolations > 0 && len(violations) > 0 && tt.expectedMessage != "" {
				found := false
				for _, v := range violations {
					if v.Message == tt.expectedMessage {
						found = true
						if v.Severity != SeverityError {
							t.Errorf("Expected severity %v, got %v", SeverityError, v.Severity)
						}
						break
					}
				}
				if !found {
					t.Errorf("Expected message %q not found in violations", tt.expectedMessage)
					for i, v := range violations {
						t.Errorf("Actual violation %d: %s", i, v.Message)
					}
				}
			}
		})
	}
}

func TestQAS002_Integration(t *testing.T) {
	// Create temporary file with undefined identifier
	tmpFile, err := os.CreateTemp("", "test_*.qasm")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	testContent := `OPENQASM 3.0;
include "stdgates.qasm";

h q[0];  // q is not declared`

	if _, err := tmpFile.WriteString(testContent); err != nil {
		t.Fatalf("Failed to write test input: %v", err)
	}
	tmpFile.Close()

	// Create linter with built-in rules
	linter := NewLinter("")
	if err := linter.LoadRules(); err != nil {
		t.Fatalf("Failed to load rules: %v", err)
	}

	// Run linter
	violations, err := linter.LintFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to lint file: %v", err)
	}

	// Should have at least one QAS002 violation
	hasQAS002 := false
	for _, v := range violations {
		if v.Rule.ID == "QAS002" {
			hasQAS002 = true
			if v.Message != "Identifier 'q' is not declared." {
				t.Errorf("Unexpected QAS002 message: %s", v.Message)
			}
			break
		}
	}

	if !hasQAS002 {
		t.Error("Expected QAS002 violation but none found")
	}
}
