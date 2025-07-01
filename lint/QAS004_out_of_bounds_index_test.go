package lint

import (
	"os"
	"testing"
)

func TestQAS004_OutOfBoundsIndex(t *testing.T) {
	tests := []struct {
		name               string
		input              string
		expectedViolations int
		expectedMessage    string
	}{
		{
			name: "accessing index within bounds should not trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit[2] q;

h q[0];  // Valid index 0
h q[1];  // Valid index 1`,
			expectedViolations: 0,
		},
		{
			name: "accessing index out of bounds should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit[2] q;

h q[2];  // Invalid index 2 (valid range is 0-1)`,
			expectedViolations: 1,
			expectedMessage:    "Index out of bounds: accessing '2' on 'q' of length 2.",
		},
		{
			name: "multiple out of bounds accesses should trigger multiple violations",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit[2] q;
bit[3] c;

h q[2];      // Invalid index 2 for q[2]
h q[5];      // Invalid index 5 for q[2] 
measure q[0] -> c[3];  // Invalid index 3 for c[3]`,
			expectedViolations: 3,
		},
		{
			name: "bit arrays should also be checked",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit[2] q;
bit[2] c;

h q[0];
measure q[0] -> c[2];  // Invalid index 2 for c[2] (valid range is 0-1)`,
			expectedViolations: 1,
			expectedMessage:    "Index out of bounds: accessing '2' on 'c' of length 2.",
		},
		{
			name: "accessing non-existent array should be handled gracefully",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

h undefined_array[0];  // Should not crash, but also not trigger QAS004`,
			expectedViolations: 0, // This would be caught by QAS002 (undefined identifier)
		},
		{
			name: "single qubits and bits should not trigger violations",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit q;
bit c;

h q;
measure q -> c;`,
			expectedViolations: 0,
		},
		{
			name: "comments should not affect bounds checking",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit[2] q;

// This comment has q[5] but should be ignored
h q[2];  // This should trigger violation`,
			expectedViolations: 1,
		},
		{
			name: "zero-length arrays should handle all accesses as out of bounds",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit[0] q;

h q[0];  // Even index 0 is out of bounds for zero-length array`,
			expectedViolations: 1,
			expectedMessage:    "Index out of bounds: accessing '0' on 'q' of length 0.",
		},
		{
			name: "large array with valid and invalid accesses",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit[10] q;

h q[0];   // Valid
h q[5];   // Valid
h q[9];   // Valid (last valid index)
h q[10];  // Invalid`,
			expectedViolations: 1,
			expectedMessage:    "Index out of bounds: accessing '10' on 'q' of length 10.",
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
			checker := &OutOfBoundsIndexChecker{}

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

func TestQAS004_Integration(t *testing.T) {
	// Create temporary file with out of bounds access
	tmpFile, err := os.CreateTemp("", "test_*.qasm")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	testContent := `OPENQASM 3.0;
include "stdgates.qasm";

qubit[2] q;

h q[2];  // Index 2 is out of bounds`

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

	// Should have at least one QAS004 violation
	hasQAS004 := false
	for _, v := range violations {
		if v.Rule.ID == "QAS004" {
			hasQAS004 = true
			expectedMessage := "Index out of bounds: accessing '2' on 'q' of length 2."
			if v.Message != expectedMessage {
				t.Errorf("Unexpected QAS004 message: %s", v.Message)
			}
			break
		}
	}

	if !hasQAS004 {
		t.Error("Expected QAS004 violation but none found")
	}
}
