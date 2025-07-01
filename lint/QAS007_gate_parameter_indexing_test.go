package lint

import (
	"os"
	"testing"
)

func TestQAS007_GateParameterIndexing(t *testing.T) {
	tests := []struct {
		name               string
		input              string
		expectedViolations int
		expectedMessage    string
	}{
		{
			name: "valid gate without index access should not trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

gate mygate q, r {
    h q;
    cx q, r;
}`,
			expectedViolations: 0,
		},
		{
			name: "index access on gate parameter should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

gate mygate q, r {
    h q[0];  // Invalid - index access on parameter
    cx q, r;
}`,
			expectedViolations: 1,
			expectedMessage:    "Cannot perform index access on gate argument 'q'.",
		},
		{
			name: "index access on multiple gate parameters should trigger multiple violations",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

gate mygate q, r {
    h q[0];     // Invalid - index access on parameter q
    x r[1];     // Invalid - index access on parameter r
    cx q, r;    // Valid - no index access
}`,
			expectedViolations: 2,
		},
		{
			name: "index access on local variables (not parameters) should not trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit[2] localQubit;

gate mygate q, r {
    h q;           // Valid - parameter without index
    cx q, r;       // Valid - parameters without index
}

h localQubit[0];   // Valid - not inside gate, can use index`,
			expectedViolations: 0,
		},
		{
			name: "mixed parameter and non-parameter index access",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

gate mygate q, r {
    qubit[2] temp;  // Local declaration inside gate (hypothetical)
    h q[0];         // Invalid - index access on parameter
    // temp[0] would be valid if it were possible
    cx q, r;
}`,
			expectedViolations: 1,
			expectedMessage:    "Cannot perform index access on gate argument 'q'.",
		},
		{
			name: "parameterized gate with index access should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

gate mygate(theta) q, r {
    rx(theta) q[0];  // Invalid - index access on parameter
    ry(theta) r;     // Valid - no index access
}`,
			expectedViolations: 1,
			expectedMessage:    "Cannot perform index access on gate argument 'q'.",
		},
		{
			name: "nested gate calls with parameter index access should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

gate innergate q {
    h q[0];  // Invalid - index access on parameter
}

gate outergate a, b {
    innergate a;     // Valid - no index access here
    cx a[1], b;      // Invalid - index access on parameter a
}`,
			expectedViolations: 2,
		},
		{
			name: "comments should not affect parameter indexing detection",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

gate mygate q, r {
    // This comment has q[0] but should be ignored
    h q[0];  // This should trigger violation
    cx q, r;
}`,
			expectedViolations: 1,
		},
		{
			name: "empty gate body should not trigger violations",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

gate emptygate q, r {
    // No operations
}`,
			expectedViolations: 0,
		},
		{
			name: "single parameter gate with index access should trigger violation",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

gate singlegate q {
    h q[0];  // Invalid - index access on parameter
}`,
			expectedViolations: 1,
			expectedMessage:    "Cannot perform index access on gate argument 'q'.",
		},
		{
			name: "complex gate with multiple operations and index violations",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

gate complexgate a, b, c {
    h a;        // Valid
    cx a, b;    // Valid
    h b[0];     // Invalid - index access on parameter b
    rx(pi/2) c[1];  // Invalid - index access on parameter c
    cx b, c;    // Valid
}`,
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
			checker := &GateParameterIndexingChecker{}

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

func TestQAS007_Integration(t *testing.T) {
	// Create temporary file with gate parameter indexing violation
	tmpFile, err := os.CreateTemp("", "test_*.qasm")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	testContent := `OPENQASM 3.0;
include "stdgates.qasm";

gate mygate q, r {
    h q[0];  // Index access on parameter
    cx q, r;
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

	// Should have at least one QAS007 violation
	hasQAS007 := false
	for _, v := range violations {
		if v.Rule.ID == "QAS007" {
			hasQAS007 = true
			expectedMessage := "Cannot perform index access on gate argument 'q'."
			if v.Message != expectedMessage {
				t.Errorf("Unexpected QAS007 message: %s", v.Message)
			}
			break
		}
	}

	if !hasQAS007 {
		t.Error("Expected QAS007 violation but none found")
	}
}
