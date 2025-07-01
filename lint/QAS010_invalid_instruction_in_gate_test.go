package lint

import (
	"os"
	"testing"
)

func TestQAS010_InvalidInstructionInGate(t *testing.T) {
	tests := []struct {
		name               string
		input              string
		expectedViolations int
		expectedMessage    string
	}{
		{
			name: "valid gate with only unitary operations should not trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

gate mygate q, r {
    h q;
    cx q, r;
    x r;
}`,
			expectedViolations: 0,
		},
		{
			name: "measure instruction in gate should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

gate invalidgate q, r {
    h q;
    measure q -> c;  // Invalid - measurement in gate
}`,
			expectedViolations: 1,
			expectedMessage:    "Invalid instruction 'measure' used within gate definition.",
		},
		{
			name: "reset instruction in gate should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

gate invalidgate q {
    reset q;  // Invalid - reset in gate
    h q;
}`,
			expectedViolations: 1,
			expectedMessage:    "Invalid instruction 'reset' used within gate definition.",
		},
		{
			name: "barrier instruction in gate should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

gate invalidgate q, r {
    h q;
    barrier q, r;  // Invalid - barrier in gate
    cx q, r;
}`,
			expectedViolations: 1,
			expectedMessage:    "Invalid instruction 'barrier' used within gate definition.",
		},
		{
			name: "if statement in gate should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

gate invalidgate q {
    if (true) {  // Invalid - classical control in gate
        h q;
    }
}`,
			expectedViolations: 1,
			expectedMessage:    "Invalid instruction 'if' used within gate definition.",
		},
		{
			name: "while loop in gate should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

gate invalidgate q {
    while (true) {  // Invalid - classical control in gate
        h q;
    }
}`,
			expectedViolations: 1,
			expectedMessage:    "Invalid instruction 'while' used within gate definition.",
		},
		{
			name: "for loop in gate should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

gate invalidgate q {
    for i in [0:1:5] {  // Invalid - classical control in gate
        h q;
    }
}`,
			expectedViolations: 1,
			expectedMessage:    "Invalid instruction 'for' used within gate definition.",
		},
		{
			name: "multiple invalid instructions should trigger multiple violations",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

gate invalidgate q, r {
    measure q -> c;  // Invalid - measurement
    reset r;         // Invalid - reset
    barrier q, r;    // Invalid - barrier
}`,
			expectedViolations: 3,
		},
		{
			name: "valid instructions outside gate should not trigger violations",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit q;
bit c;

gate validgate q {
    h q;  // Valid unitary operation
    x q;
}

measure q -> c;  // Valid - outside gate
reset q;         // Valid - outside gate`,
			expectedViolations: 0,
		},
		{
			name: "multiple gates with invalid instructions should trigger violations",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

gate firstgate q {
    h q;
    measure q -> c;  // Invalid - measurement in gate
}

gate secondgate r {
    reset r;         // Invalid - reset in gate
}`,
			expectedViolations: 2,
		},
		{
			name: "comments should not affect invalid instruction detection",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

gate invalidgate q {
    // This comment mentions measure but should be ignored
    measure q -> c;  // This should trigger violation
}`,
			expectedViolations: 1,
		},
		{
			name: "empty gate should not trigger violations",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

gate emptygate q {
    // No instructions
}`,
			expectedViolations: 0,
		},
		{
			name: "classical assignment in gate should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

gate invalidgate q {
    bit temp = 0;  // Invalid - classical assignment in gate
    h q;
}`,
			expectedViolations: 1,
		},
		{
			name: "parameterized gate with invalid instructions should trigger violations",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

gate invalidgate(theta) q {
    rx(theta) q;     // Valid - parameterized unitary
    measure q -> c;  // Invalid - measurement
}`,
			expectedViolations: 1,
			expectedMessage:    "Invalid instruction 'measure' used within gate definition.",
		},
		{
			name: "gate with only custom gate calls should not trigger violations",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

gate subgate q {
    h q;
    x q;
}

gate maingate q, r {
    subgate q;    // Valid - calling another gate
    cx q, r;      // Valid - standard gate
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
			checker := &InvalidInstructionInGateChecker{}

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

func TestQAS010_Integration(t *testing.T) {
	// Create temporary file with invalid instruction in gate
	tmpFile, err := os.CreateTemp("", "test_*.qasm")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	testContent := `OPENQASM 3.0;
include "stdgates.qasm";

gate invalidgate q {
    h q;
    measure q -> c;  // Invalid instruction in gate
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

	// Should have at least one QAS010 violation
	hasQAS010 := false
	for _, v := range violations {
		if v.Rule.ID == "QAS010" {
			hasQAS010 = true
			expectedMessage := "Invalid instruction 'measure' used within gate definition."
			if v.Message != expectedMessage {
				t.Errorf("Unexpected QAS010 message: %s", v.Message)
			}
			break
		}
	}

	if !hasQAS010 {
		t.Error("Expected QAS010 violation but none found")
	}
}
