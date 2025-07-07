package lint

import (
	"os"
	"testing"
)

func TestQAS003_ConstantMeasuredBit(t *testing.T) {
	tests := []struct {
		name               string
		input              string
		expectedViolations int
		expectedMessage    string
	}{
		{
			name: "measuring qubit without gates should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit q;
bit c;

measure q -> c;  // q has no gates applied`,
			expectedViolations: 1,
			expectedMessage:    "Measuring qubit 'q' that has no gates applied. The result will always be |0⟩.",
		},
		{
			name: "measuring qubit with gates should not trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit q;
bit c;

h q;          // Gate applied to q
measure q -> c;`,
			expectedViolations: 0,
		},
		{
			name: "measuring array qubit without gates should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit[2] q;
bit[2] c;

measure q[0] -> c[0];  // q[0] has no gates applied
measure q[1] -> c[1];  // q[1] has no gates applied`,
			expectedViolations: 2,
		},
		{
			name: "measuring array qubit with gates should not trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit[2] q;
bit[2] c;

h q[0];       // Gate applied to q
cx q[0], q[1]; // Gates applied to both q[0] and q[1]
measure q[0] -> c[0];
measure q[1] -> c[1];`,
			expectedViolations: 0,
		},
		{
			name: "mixed case - some qubits with gates, some without",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit q1;
qubit q2;
bit c1;
bit c2;

h q1;         // Gate applied to q1
measure q1 -> c1;  // Should not trigger violation
measure q2 -> c2;  // Should trigger violation`,
			expectedViolations: 1,
			expectedMessage:    "Measuring qubit 'q2' that has no gates applied. The result will always be |0⟩.",
		},
		{
			name: "two-qubit gates should mark both qubits",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit q1;
qubit q2;
bit c1;
bit c2;

cx q1, q2;    // Both q1 and q2 affected
measure q1 -> c1;
measure q2 -> c2;`,
			expectedViolations: 0,
		},
		{
			name: "parameterized gates should be recognized",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit q;
bit c;

rx(pi/2) q;   // Parameterized gate applied
measure q -> c;`,
			expectedViolations: 0,
		},
		{
			name: "comments should not affect detection",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit q;
bit c;

// This is a comment with h q; inside
measure q -> c;  // Should trigger violation`,
			expectedViolations: 1,
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
			checker := NewConstantMeasuredBitChecker()

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
						if v.Severity != SeverityWarning {
							t.Errorf("Expected severity %v, got %v", SeverityWarning, v.Severity)
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

func TestQAS003_Integration(t *testing.T) {
	// Create temporary file with constant measured bit
	tmpFile, err := os.CreateTemp("", "test_*.qasm")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	testContent := `OPENQASM 3.0;
include "stdgates.qasm";

qubit q;
bit c;

measure q -> c;  // q has no gates applied`

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

	// Should have at least one QAS003 violation
	hasQAS003 := false
	for _, v := range violations {
		if v.Rule.ID == "QAS003" {
			hasQAS003 = true
			expectedMessage := "Measuring qubit 'q' that has no gates applied. The result will always be |0⟩."
			if v.Message != expectedMessage {
				t.Errorf("Unexpected QAS003 message: %s", v.Message)
			}
			break
		}
	}

	if !hasQAS003 {
		t.Error("Expected QAS003 violation but none found")
	}
}
