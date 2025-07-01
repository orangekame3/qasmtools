package lint

import (
	"os"
	"testing"
)

func TestQAS006_GateRegisterSizeMismatch(t *testing.T) {
	tests := []struct {
		name               string
		input              string
		expectedViolations int
		expectedMessage    string
	}{
		{
			name: "matching register sizes should not trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit[2] q1;
qubit[2] q2;

cx q1, q2;  // Both size 2`,
			expectedViolations: 0,
		},
		{
			name: "mismatched register sizes should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit[2] q1;
qubit[3] q2;

cx q1, q2;  // Size mismatch: 2 vs 3`,
			expectedViolations: 1,
			expectedMessage:    "Register lengths passed to gate 'cx' do not match.",
		},
		{
			name: "single register with array register should not trigger violation (broadcasting)",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit q_single;
qubit[3] q_array;

cx q_single, q_array;  // Broadcasting allowed`,
			expectedViolations: 0,
		},
		{
			name: "array register with single register should not trigger violation (broadcasting)",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit[3] q_array;
qubit q_single;

cx q_array, q_single;  // Broadcasting allowed`,
			expectedViolations: 0,
		},
		{
			name: "multiple array registers with different sizes should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit[2] q1;
qubit[3] q2;
qubit[4] q3;

gate mygate q, r, s {
    cx q, r;
    cx r, s;
}

mygate q1, q2, q3;  // All different sizes: 2, 3, 4`,
			expectedViolations: 1,
			expectedMessage:    "Register lengths passed to gate 'mygate' do not match.",
		},
		{
			name: "custom gate with matching sizes should not trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

gate mygate q, r {
    cx q, r;
    h q;
}

qubit[2] q1;
qubit[2] q2;

mygate q1, q2;  // Both size 2`,
			expectedViolations: 0,
		},
		{
			name: "single qubit gates should not trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit[3] q;

h q;  // Single register gate`,
			expectedViolations: 0,
		},
		{
			name: "mixed single and matching array registers should not trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit q_single;
qubit[2] q1;
qubit[2] q2;

gate threequbit a, b, c {
    cx a, b;
    cx b, c;
}

threequbit q_single, q1, q2;  // Single + matching arrays`,
			expectedViolations: 0,
		},
		{
			name: "comments should not affect size checking",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit[2] q1;
qubit[3] q2;

// This comment mentions cx q1, q2 but should be ignored
cx q1, q2;  // This should trigger violation`,
			expectedViolations: 1,
		},
		{
			name: "indexed register access should use base register size",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit[2] q1;
qubit[3] q2;

cx q1[0], q2[0];  // Individual qubit access, should not trigger violation`,
			expectedViolations: 0,
		},
		{
			name: "multiple violations should be detected",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit[2] q1;
qubit[3] q2;
qubit[4] q3;

cx q1, q2;      // Size mismatch
cx q2, q3;      // Size mismatch`,
			expectedViolations: 2,
		},
		{
			name: "parameterized gates should be checked",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit[2] q1;
qubit[3] q2;

rx(pi/2) q1, q2;  // Size mismatch in parameterized gate`,
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
			checker := &GateRegisterSizeMismatchChecker{}

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

func TestQAS006_Integration(t *testing.T) {
	// Create temporary file with register size mismatch
	tmpFile, err := os.CreateTemp("", "test_*.qasm")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	testContent := `OPENQASM 3.0;
include "stdgates.qasm";

qubit[2] q1;
qubit[3] q2;

cx q1, q2;  // Size mismatch`

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

	// Should have at least one QAS006 violation
	hasQAS006 := false
	for _, v := range violations {
		if v.Rule.ID == "QAS006" {
			hasQAS006 = true
			expectedMessage := "Register lengths passed to gate 'cx' do not match."
			if v.Message != expectedMessage {
				t.Errorf("Unexpected QAS006 message: %s", v.Message)
			}
			break
		}
	}

	if !hasQAS006 {
		t.Error("Expected QAS006 violation but none found")
	}
}