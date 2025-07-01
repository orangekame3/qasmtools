package lint

import (
	"os"
	"path/filepath"
	"testing"
)

func TestQAS001_UnusedQubit(t *testing.T) {
	tests := []struct {
		name               string
		input              string
		expectedViolations int
		expectedMessage    string
	}{
		{
			name: "unused qubit should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit q;
qubit unused_qubit;

h q;`,
			expectedViolations: 1,
			expectedMessage:    "Qubit 'unused_qubit' is declared but never used.",
		},
		{
			name: "all qubits used should not trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit q;
qubit q2;

h q;
cx q, q2;`,
			expectedViolations: 0,
		},
		{
			name: "unused array qubit should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit[2] q;
qubit[3] unused_array;

h q[0];
cx q[0], q[1];`,
			expectedViolations: 1,
			expectedMessage:    "Qubit 'unused_array' is declared but never used.",
		},
		{
			name: "array qubit with usage should not trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit[2] q;

h q[0];
cx q[0], q[1];`,
			expectedViolations: 0,
		},
		{
			name: "qubit used in measurement should not trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit q;
bit c;

measure q -> c;`,
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
			checker := &UnusedQubitChecker{}

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

			if tt.expectedViolations > 0 && len(violations) > 0 {
				if violations[0].Message != tt.expectedMessage {
					t.Errorf("Expected message %q, got %q", tt.expectedMessage, violations[0].Message)
				}
				if violations[0].Severity != SeverityWarning {
					t.Errorf("Expected severity %v, got %v", SeverityWarning, violations[0].Severity)
				}
			}
		})
	}
}

func TestQAS001_Integration(t *testing.T) {
	// Test with actual lint runner
	tmpDir, err := os.MkdirTemp("", "qas001_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Write test file
	testFile := filepath.Join(tmpDir, "test.qasm")
	testContent := `OPENQASM 3.0;
include "stdgates.qasm";

qubit q;
qubit unused_qubit;

h q;`

	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Create linter with built-in rules
	linter := NewLinter("")
	if err := linter.LoadRules(); err != nil {
		t.Fatalf("Failed to load rules: %v", err)
	}

	// Run linter
	violations, err := linter.LintFile(testFile)
	if err != nil {
		t.Fatalf("Failed to lint file: %v", err)
	}

	// Should have at least one QAS001 violation
	hasQAS001 := false
	for _, v := range violations {
		if v.Rule.ID == "QAS001" {
			hasQAS001 = true
			if v.Message != "Qubit 'unused_qubit' is declared but never used." {
				t.Errorf("Unexpected QAS001 message: %s", v.Message)
			}
			break
		}
	}

	if !hasQAS001 {
		t.Error("Expected QAS001 violation but none found")
	}
}
