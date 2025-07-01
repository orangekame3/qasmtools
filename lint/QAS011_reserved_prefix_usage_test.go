package lint

import (
	"os"
	"testing"
)

func TestQAS011_ReservedPrefixUsage(t *testing.T) {
	tests := []struct {
		name               string
		input              string
		expectedViolations int
		expectedMessage    string
	}{
		{
			name: "normal identifiers should not trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit myQubit;
bit myBit;
int myInt;`,
			expectedViolations: 0,
		},
		{
			name: "qubit with reserved prefix should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit __reserved_qubit;  // Invalid - uses reserved prefix`,
			expectedViolations: 1,
			expectedMessage:    "Identifier '__reserved_qubit' uses reserved prefix '__'.",
		},
		{
			name: "bit with reserved prefix should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

bit __reserved_bit;  // Invalid - uses reserved prefix`,
			expectedViolations: 1,
			expectedMessage:    "Identifier '__reserved_bit' uses reserved prefix '__'.",
		},
		{
			name: "gate with reserved prefix should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

gate __reserved_gate q {  // Invalid - uses reserved prefix
    h q;
}`,
			expectedViolations: 1,
			expectedMessage:    "Identifier '__reserved_gate' uses reserved prefix '__'.",
		},
		{
			name: "function with reserved prefix should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

def __reserved_function() {  // Invalid - uses reserved prefix
    return 0;
}`,
			expectedViolations: 1,
			expectedMessage:    "Identifier '__reserved_function' uses reserved prefix '__'.",
		},
		{
			name: "const with reserved prefix should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

const __reserved_const = 3.14;  // Invalid - uses reserved prefix`,
			expectedViolations: 1,
			expectedMessage:    "Identifier '__reserved_const' uses reserved prefix '__'.",
		},
		{
			name: "register types with reserved prefix should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

int __reserved_int;      // Invalid - uses reserved prefix
float __reserved_float;  // Invalid - uses reserved prefix`,
			expectedViolations: 2,
		},
		{
			name: "array declarations with reserved prefix should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit[2] __reserved_array;  // Invalid - uses reserved prefix
bit[3] __another_reserved;  // Invalid - uses reserved prefix`,
			expectedViolations: 2,
		},
		{
			name: "input/output with reserved prefix should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

input __reserved_input;   // Invalid - uses reserved prefix
output __reserved_output; // Invalid - uses reserved prefix`,
			expectedViolations: 2,
		},
		{
			name: "single underscore prefix should not trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit _single_underscore;  // Valid - only single underscore
bit _another_one;          // Valid - only single underscore`,
			expectedViolations: 0,
		},
		{
			name: "triple underscore prefix should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit ___triple_underscore;  // Invalid - starts with __ (even though it's triple)`,
			expectedViolations: 1,
			expectedMessage:    "Identifier '___triple_underscore' uses reserved prefix '__'.",
		},
		{
			name: "mixed valid and invalid identifiers",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit validQubit;          // Valid
bit __invalid_bit;         // Invalid - reserved prefix
int anotherValid;          // Valid
float __invalid_float;     // Invalid - reserved prefix`,
			expectedViolations: 2,
		},
		{
			name: "comments should not affect reserved prefix detection",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

// This comment mentions __reserved but should be ignored
qubit __actual_reserved;  // This should trigger violation`,
			expectedViolations: 1,
		},
		{
			name: "reserved prefix in middle of identifier should not trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit valid__middle;  // Valid - __ not at start
bit end__valid;       // Valid - __ not at start`,
			expectedViolations: 0,
		},
		{
			name: "circuit with reserved prefix should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

circuit __reserved_circuit(q) {  // Invalid - uses reserved prefix
    h q;
}`,
			expectedViolations: 1,
			expectedMessage:    "Identifier '__reserved_circuit' uses reserved prefix '__'.",
		},
		{
			name: "parameterized declarations with reserved prefix should trigger violations",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

gate __reserved_gate(theta) q {  // Invalid - uses reserved prefix
    rx(theta) q;
}

def __reserved_def(angle x) -> angle {  // Invalid - uses reserved prefix
    return x * 2;
}`,
			expectedViolations: 2,
		},
		{
			name: "empty declarations should not cause issues",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit validQubit;  // Valid declaration`,
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
			checker := &ReservedPrefixUsageChecker{}

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

func TestQAS011_Integration(t *testing.T) {
	// Create temporary file with reserved prefix usage
	tmpFile, err := os.CreateTemp("", "test_*.qasm")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	testContent := `OPENQASM 3.0;
include "stdgates.qasm";

qubit __reserved_qubit;  // Should trigger violation`

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

	// Should have at least one QAS011 violation
	hasQAS011 := false
	for _, v := range violations {
		if v.Rule.ID == "QAS011" {
			hasQAS011 = true
			expectedMessage := "Identifier '__reserved_qubit' uses reserved prefix '__'."
			if v.Message != expectedMessage {
				t.Errorf("Unexpected QAS011 message: %s", v.Message)
			}
			break
		}
	}

	if !hasQAS011 {
		t.Error("Expected QAS011 violation but none found")
	}
}