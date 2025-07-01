package lint

import (
	"os"
	"testing"
)

func TestQAS009_IllegalBreakContinue(t *testing.T) {
	tests := []struct {
		name               string
		input              string
		expectedViolations int
		expectedMessage    string
	}{
		{
			name: "break inside for loop should not trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

for i in [0:1:10] {
    if (i == 5) {
        break;  // Valid - inside loop
    }
}`,
			expectedViolations: 0,
		},
		{
			name: "continue inside for loop should not trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

for i in [0:1:10] {
    if (i == 5) {
        continue;  // Valid - inside loop
    }
}`,
			expectedViolations: 0,
		},
		{
			name: "break outside loop should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit q;
break;  // Invalid - outside loop`,
			expectedViolations: 1,
			expectedMessage:    "'break' cannot be used outside of a loop.",
		},
		{
			name: "continue outside loop should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit q;
continue;  // Invalid - outside loop`,
			expectedViolations: 1,
			expectedMessage:    "'continue' cannot be used outside of a loop.",
		},
		{
			name: "break inside while loop should not trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

bit c;
while (c == 0) {
    break;  // Valid - inside while loop
}`,
			expectedViolations: 0,
		},
		{
			name: "continue inside while loop should not trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

bit c;
while (c == 0) {
    continue;  // Valid - inside while loop
}`,
			expectedViolations: 0,
		},
		{
			name: "break/continue inside nested loops should not trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

for i in [0:1:5] {
    for j in [0:1:3] {
        if (i == j) {
            break;     // Valid - inside nested loop
        }
        continue;      // Valid - inside nested loop
    }
}`,
			expectedViolations: 0,
		},
		{
			name: "multiple break/continue outside loops should trigger multiple violations",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit q;
break;      // Invalid - outside loop
continue;   // Invalid - outside loop`,
			expectedViolations: 2,
		},
		{
			name: "break/continue after loop ends should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

for i in [0:1:5] {
    // Valid loop content
}
break;      // Invalid - loop has ended`,
			expectedViolations: 1,
			expectedMessage:    "'break' cannot be used outside of a loop.",
		},
		{
			name: "break/continue inside function but outside loop should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

def myfunction() {
    break;      // Invalid - inside function but outside loop
    continue;   // Invalid - inside function but outside loop
}`,
			expectedViolations: 2,
		},
		{
			name: "break/continue inside gate but outside loop should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

gate mygate q, r {
    break;      // Invalid - inside gate but outside loop
    continue;   // Invalid - inside gate but outside loop
}`,
			expectedViolations: 2,
		},
		{
			name: "comments should not affect break/continue detection",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

// This comment mentions break but should be ignored
break;  // This should trigger violation`,
			expectedViolations: 1,
		},
		{
			name: "break/continue in complex nested structure",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

def myfunction() {
    for i in [0:1:5] {
        if (i == 2) {
            break;      // Valid - inside loop within function
        }
    }
    break;              // Invalid - outside loop but inside function
}`,
			expectedViolations: 1,
			expectedMessage:    "'break' cannot be used outside of a loop.",
		},
		{
			name: "empty loops should not cause issues",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

for i in [0:1:5] {
    // Empty loop body
}`,
			expectedViolations: 0,
		},
		{
			name: "loop with only break/continue statements",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

for i in [0:1:5] {
    break;
    continue;  // Unreachable but syntactically valid
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
			checker := &IllegalBreakContinueChecker{}

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

func TestQAS009_Integration(t *testing.T) {
	// Create temporary file with break/continue outside loop
	tmpFile, err := os.CreateTemp("", "test_*.qasm")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	testContent := `OPENQASM 3.0;
include "stdgates.qasm";

qubit q;
break;  // Should trigger violation`

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

	// Should have at least one QAS009 violation
	hasQAS009 := false
	for _, v := range violations {
		if v.Rule.ID == "QAS009" {
			hasQAS009 = true
			expectedMessage := "'break' cannot be used outside of a loop."
			if v.Message != expectedMessage {
				t.Errorf("Unexpected QAS009 message: %s", v.Message)
			}
			break
		}
	}

	if !hasQAS009 {
		t.Error("Expected QAS009 violation but none found")
	}
}
