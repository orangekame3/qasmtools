package lint

import (
	"os"
	"testing"
)

func TestQAS012_SnakeCaseRequired(t *testing.T) {
	tests := []struct {
		name               string
		input              string
		expectedViolations int
		expectedMessage    string
	}{
		{
			name: "valid snake_case identifiers should not trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit my_qubit;
bit my_bit;
qubit single;`,
			expectedViolations: 0,
		},
		{
			name: "camelCase qubit should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit myQubit;  // Invalid - camelCase`,
			expectedViolations: 1,
			expectedMessage:    "Identifier 'myQubit' should be written in snake_case.",
		},
		{
			name: "PascalCase bit should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

bit MyBit;  // Invalid - PascalCase`,
			expectedViolations: 1,
			expectedMessage:    "Identifier 'MyBit' should be written in snake_case.",
		},
		{
			name: "camelCase gate should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

gate myGate q {  // Invalid - camelCase
    h q;
}`,
			expectedViolations: 1,
			expectedMessage:    "Identifier 'myGate' should be written in snake_case.",
		},
		{
			name: "PascalCase gate should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

gate MyGate q {  // Invalid - PascalCase
    h q;
}`,
			expectedViolations: 1,
			expectedMessage:    "Identifier 'MyGate' should be written in snake_case.",
		},
		{
			name: "camelCase circuit should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

circuit myCircuit(q) {  // Invalid - camelCase
    h q;
}`,
			expectedViolations: 1,
			expectedMessage:    "Identifier 'myCircuit' should be written in snake_case.",
		},
		{
			name: "mixed valid and invalid identifiers",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit valid_qubit;     // Valid - snake_case
bit invalidBit;        // Invalid - camelCase
gate validGate q {     // Invalid - camelCase
    h q;
}
qubit another_valid;   // Valid - snake_case`,
			expectedViolations: 2,
		},
		{
			name: "identifiers with numbers should follow snake_case",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit qubit1;          // Valid - ends with number
bit my_bit2;           // Valid - snake_case with number
qubit myQubit3;        // Invalid - camelCase with number`,
			expectedViolations: 1,
			expectedMessage:    "Identifier 'myQubit3' should be written in snake_case.",
		},
		{
			name: "array declarations should follow snake_case",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit[2] valid_array;     // Valid - snake_case
bit[3] invalidArray;      // Invalid - camelCase`,
			expectedViolations: 1,
			expectedMessage:    "Identifier 'invalidArray' should be written in snake_case.",
		},
		{
			name: "identifiers starting with underscore should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit _invalid_start;  // Invalid - starts with underscore`,
			expectedViolations: 1,
			expectedMessage:    "Identifier '_invalid_start' should be written in snake_case.",
		},
		{
			name: "identifiers ending with underscore should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit invalid_end_;  // Invalid - ends with underscore`,
			expectedViolations: 1,
			expectedMessage:    "Identifier 'invalid_end_' should be written in snake_case.",
		},
		{
			name: "identifiers with consecutive underscores should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit invalid__double;  // Invalid - consecutive underscores`,
			expectedViolations: 1,
			expectedMessage:    "Identifier 'invalid__double' should be written in snake_case.",
		},
		{
			name: "single character identifiers should not trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit q;  // Valid - single character
bit c;    // Valid - single character`,
			expectedViolations: 0,
		},
		{
			name: "comments should not affect snake_case detection",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

// This comment mentions myQubit but should be ignored
qubit myQubit;  // This should trigger violation`,
			expectedViolations: 1,
		},
		{
			name: "all uppercase should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit UPPERCASE;  // Invalid - all uppercase`,
			expectedViolations: 1,
			expectedMessage:    "Identifier 'UPPERCASE' should be written in snake_case.",
		},
		{
			name: "parameterized gates should follow snake_case",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

gate valid_gate(theta) q {    // Valid - snake_case
    rx(theta) q;
}

gate invalidGate(phi) q {     // Invalid - camelCase
    ry(phi) q;
}`,
			expectedViolations: 1,
			expectedMessage:    "Identifier 'invalidGate' should be written in snake_case.",
		},
		{
			name: "valid complex snake_case identifiers",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit quantum_register_1;
bit classical_bit_array;
gate controlled_rotation q, r {
    cx q, r;
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
			checker := &SnakeCaseRequiredChecker{}

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

func TestQAS012_Integration(t *testing.T) {
	// Create temporary file with non-snake_case identifier
	tmpFile, err := os.CreateTemp("", "test_*.qasm")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	testContent := `OPENQASM 3.0;
include "stdgates.qasm";

qubit myQubit;  // Should trigger violation`

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

	// Should have at least one QAS012 violation
	hasQAS012 := false
	for _, v := range violations {
		if v.Rule.ID == "QAS012" {
			hasQAS012 = true
			expectedMessage := "Identifier 'myQubit' should be written in snake_case."
			if v.Message != expectedMessage {
				t.Errorf("Unexpected QAS012 message: %s", v.Message)
			}
			break
		}
	}

	if !hasQAS012 {
		t.Error("Expected QAS012 violation but none found")
	}
}
