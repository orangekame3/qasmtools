package lint

import (
	"os"
	"testing"
)

func TestQAS008_QubitDeclaredInLocalScope(t *testing.T) {
	tests := []struct {
		name               string
		input              string
		expectedViolations int
		expectedMessage    string
	}{
		{
			name: "global qubit declaration should not trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit q;       // Global scope - valid
qubit[2] q_array;  // Global scope - valid`,
			expectedViolations: 0,
		},
		{
			name: "qubit declaration in function should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

def myfunction() {
    qubit local_q;  // Local scope - invalid
    h local_q;
}`,
			expectedViolations: 1,
			expectedMessage:    "Qubit 'local_q' can only be declared in global scope.",
		},
		{
			name: "qubit declaration in gate should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

gate mygate q, r {
    qubit temp;  // Local scope - invalid
    h temp;
    cx q, r;
}`,
			expectedViolations: 1,
			expectedMessage:    "Qubit 'temp' can only be declared in global scope.",
		},
		{
			name: "qubit declaration in if block should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

bit c;

if (c == 0) {
    qubit conditional_q;  // Local scope - invalid
    h conditional_q;
}`,
			expectedViolations: 1,
			expectedMessage:    "Qubit 'conditional_q' can only be declared in global scope.",
		},
		{
			name: "qubit declaration in for loop should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

for i in [0:5] {
    qubit loop_q;  // Local scope - invalid
    h loop_q;
}`,
			expectedViolations: 1,
			expectedMessage:    "Qubit 'loop_q' can only be declared in global scope.",
		},
		{
			name: "multiple qubit declarations in nested scopes should trigger multiple violations",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

def outer_function() {
    qubit outer_q;  // Local scope - invalid
    
    if (true) {
        qubit inner_q;  // Nested local scope - invalid
    }
}`,
			expectedViolations: 2,
		},
		{
			name: "array qubit declaration in local scope should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

def myfunction() {
    qubit[3] local_array;  // Local scope - invalid
}`,
			expectedViolations: 1,
			expectedMessage:    "Qubit 'local_array' can only be declared in global scope.",
		},
		{
			name: "mixed global and local qubit declarations",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit global_q;      // Global scope - valid

def myfunction() {
    qubit local_q;   // Local scope - invalid
    h local_q;
}

qubit another_global;  // Global scope - valid`,
			expectedViolations: 1,
			expectedMessage:    "Qubit 'local_q' can only be declared in global scope.",
		},
		{
			name: "comments should not affect scope detection",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

def myfunction() {
    // This comment mentions qubit but should be ignored
    qubit local_q;  // This should trigger violation
}`,
			expectedViolations: 1,
		},
		{
			name: "deeply nested scope should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

def outer() {
    if (true) {
        for i in [0:2] {
            qubit deeply_nested;  // Deep local scope - invalid
        }
    }
}`,
			expectedViolations: 1,
			expectedMessage:    "Qubit 'deeply_nested' can only be declared in global scope.",
		},
		{
			name: "empty scopes should not cause issues",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit global_q;  // Global scope - valid

def empty_function() {
    // No qubit declarations here
}

qubit another_global;  // Global scope - valid`,
			expectedViolations: 0,
		},
		{
			name: "qubit declarations with same-line braces",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

def compact_function() { qubit compact_q; }  // Local scope - invalid`,
			expectedViolations: 1,
			expectedMessage:    "Qubit 'compact_q' can only be declared in global scope.",
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
			checker := &QubitDeclaredInLocalScopeChecker{}

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

func TestQAS008_Integration(t *testing.T) {
	// Create temporary file with qubit declared in local scope
	tmpFile, err := os.CreateTemp("", "test_*.qasm")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	testContent := `OPENQASM 3.0;
include "stdgates.qasm";

def myfunction() {
    qubit local_q;  // Should trigger violation
    h local_q;
}`

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

	// Should have at least one QAS008 violation
	hasQAS008 := false
	for _, v := range violations {
		if v.Rule.ID == "QAS008" {
			hasQAS008 = true
			expectedMessage := "Qubit 'local_q' can only be declared in global scope."
			if v.Message != expectedMessage {
				t.Errorf("Unexpected QAS008 message: %s", v.Message)
			}
			break
		}
	}

	if !hasQAS008 {
		t.Error("Expected QAS008 violation but none found")
	}
}
