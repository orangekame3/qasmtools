package lint

import (
	"os"
	"testing"
)

func TestQAS005_NamingConventionViolation(t *testing.T) {
	tests := []struct {
		name               string
		input              string
		expectedViolations int
		expectedMessage    string
	}{
		{
			name: "valid lowercase identifier should not trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit myQubit;
bit myBit;`,
			expectedViolations: 0,
		},
		{
			name: "identifier starting with uppercase should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit MyQubit;  // Invalid - starts with uppercase`,
			expectedViolations: 1,
			expectedMessage:    "Identifier 'MyQubit' violates naming conventions. Follow pattern: ^[a-z][a-zA-Z0-9_]*$.",
		},
		{
			name: "identifier starting with number should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit 2qubit;  // Invalid - starts with number`,
			expectedViolations: 1,
			expectedMessage:    "Identifier '2qubit' violates naming conventions. Follow pattern: ^[a-z][a-zA-Z0-9_]*$.",
		},
		{
			name: "identifier starting with underscore should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit _qubit;  // Invalid - starts with underscore`,
			expectedViolations: 1,
			expectedMessage:    "Identifier '_qubit' violates naming conventions. Follow pattern: ^[a-z][a-zA-Z0-9_]*$.",
		},
		{
			name: "valid identifiers with numbers and underscores should not trigger violations",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit myQubit1;
bit myBit_2;
qubit q123_test;`,
			expectedViolations: 0,
		},
		{
			name: "array declarations should also be checked",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit[2] ValidArray;  // Invalid - starts with uppercase
bit[3] invalid_Array;  // Valid - starts with lowercase`,
			expectedViolations: 1,
			expectedMessage:    "Identifier 'ValidArray' violates naming conventions. Follow pattern: ^[a-z][a-zA-Z0-9_]*$.",
		},
		{
			name: "gate declarations should be checked",
			input: `OPENQASM 3.0;

gate MyGate(theta) q {  // Invalid - starts with uppercase
    rx(theta) q;
}

gate validGate(phi) q {  // Valid - starts with lowercase
    ry(phi) q;
}`,
			expectedViolations: 1,
			expectedMessage:    "Identifier 'MyGate' violates naming conventions. Follow pattern: ^[a-z][a-zA-Z0-9_]*$.",
		},
		{
			name: "function declarations should be checked",
			input: `OPENQASM 3.0;

def MyFunction(angle x) -> angle {  // Invalid - starts with uppercase
    return x * 2;
}

def validFunction(angle y) -> angle {  // Valid - starts with lowercase
    return y / 2;
}`,
			expectedViolations: 1,
			expectedMessage:    "Identifier 'MyFunction' violates naming conventions. Follow pattern: ^[a-z][a-zA-Z0-9_]*$.",
		},
		{
			name: "multiple violations should be detected",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit[2] MyQubit;     // Invalid
bit[2] MyBit;         // Invalid
qubit validQubit;     // Valid
bit ValidBit2;        // Invalid`,
			expectedViolations: 3,
		},
		{
			name: "const declarations should be checked",
			input: `OPENQASM 3.0;

const MyConst = 3.14;  // Invalid - starts with uppercase
const validConst = 2.71;  // Valid - starts with lowercase`,
			expectedViolations: 1,
			expectedMessage:    "Identifier 'MyConst' violates naming conventions. Follow pattern: ^[a-z][a-zA-Z0-9_]*$.",
		},
		{
			name: "register type declarations should be checked",
			input: `OPENQASM 3.0;

int MyInt;        // Invalid - starts with uppercase
float validFloat; // Valid - starts with lowercase
angle MyAngle;    // Invalid - starts with uppercase`,
			expectedViolations: 2,
		},
		{
			name: "comments should not affect naming convention checking",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

// This comment has MyQubit but should be ignored
qubit MyQubit;  // This should trigger violation`,
			expectedViolations: 1,
		},
		{
			name: "input/output declarations should be checked",
			input: `OPENQASM 3.0;

input MyInput;     // Invalid - starts with uppercase
output validOutput; // Valid - starts with lowercase`,
			expectedViolations: 1,
			expectedMessage:    "Identifier 'MyInput' violates naming conventions. Follow pattern: ^[a-z][a-zA-Z0-9_]*$.",
		},
		{
			name: "special characters in identifiers should trigger violations",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit my-qubit;  // Invalid - contains hyphen
bit my.bit;      // Invalid - contains dot`,
			expectedViolations: 2,
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
			checker := &NamingConventionViolationChecker{}

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

func TestQAS005_Integration(t *testing.T) {
	// Create temporary file with naming convention violation
	tmpFile, err := os.CreateTemp("", "test_*.qasm")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	testContent := `OPENQASM 3.0;
include "stdgates.qasm";

qubit MyQubit;  // Violates naming convention`

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

	// Should have at least one QAS005 violation
	hasQAS005 := false
	for _, v := range violations {
		if v.Rule.ID == "QAS005" {
			hasQAS005 = true
			expectedMessage := "Identifier 'MyQubit' violates naming conventions. Follow pattern: ^[a-z][a-zA-Z0-9_]*$."
			if v.Message != expectedMessage {
				t.Errorf("Unexpected QAS005 message: %s", v.Message)
			}
			break
		}
	}

	if !hasQAS005 {
		t.Error("Expected QAS005 violation but none found")
	}
}