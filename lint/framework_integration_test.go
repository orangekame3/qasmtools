package lint

import (
	"testing"
)

// TestNewImplementationCompatibility tests that new implementations work correctly
func TestNewImplementationCompatibility(t *testing.T) {
	testCases := []struct {
		name           string
		content        string
		expectedQAS001 int // unused qubit violations
		expectedQAS005 int // naming violations
		expectedQAS012 int // snake_case violations
	}{
		{
			name: "no_violations",
			content: `OPENQASM 3.0;
qubit valid_qubit;
gate valid_gate q { h q; }
h valid_qubit;`,
			expectedQAS001: 0,
			expectedQAS005: 0,
			expectedQAS012: 0,
		},
		{
			name: "unused_qubit_only",
			content: `OPENQASM 3.0;
qubit unused_qubit;
qubit used_qubit;
h used_qubit;`,
			expectedQAS001: 1,
			expectedQAS005: 0,
			expectedQAS012: 0,
		},
		{
			name: "naming_violations",
			content: `OPENQASM 3.0;
qubit CamelCase;
qubit snake_case_good;
gate PascalCaseGate q { h q; }
h snake_case_good;`,
			expectedQAS001: 1, // CamelCase is unused
			expectedQAS005: 2, // CamelCase and PascalCaseGate violate naming
			expectedQAS012: 2, // CamelCase and PascalCaseGate not snake_case
		},
		{
			name: "comprehensive_test",
			content: `OPENQASM 3.0;
qubit test_array;
qubit unused_qubit;
qubit CamelCase;
bit result_bits;
gate valid_gate q { h q; }
h test_array;
measure test_array -> result_bits;`,
			expectedQAS001: 2, // unused_qubit and CamelCase are unused
			expectedQAS005: 1, // CamelCase violates naming
			expectedQAS012: 1, // CamelCase not snake_case
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			context := &CheckContext{
				File:    "test.qasm",
				Content: tc.content,
			}

			// Test new QAS001 implementation
			qas001Checker := NewUnusedQubitChecker()
			qas001Violations := qas001Checker.CheckProgram(context)
			if len(qas001Violations) != tc.expectedQAS001 {
				t.Errorf("QAS001: Expected %d violations, got %d", tc.expectedQAS001, len(qas001Violations))
				for _, v := range qas001Violations {
					t.Errorf("  - %s", v.Message)
				}
			}

			// Test new QAS005 implementation
			qas005Checker := NewNamingConventionViolationChecker()
			qas005Violations := qas005Checker.CheckProgram(context)
			if len(qas005Violations) != tc.expectedQAS005 {
				t.Errorf("QAS005: Expected %d violations, got %d", tc.expectedQAS005, len(qas005Violations))
				for _, v := range qas005Violations {
					t.Errorf("  - %s", v.Message)
				}
			}

			// Test new QAS012 implementation
			qas012Checker := NewSnakeCaseRequiredChecker()
			qas012Violations := qas012Checker.CheckProgram(context)
			if len(qas012Violations) != tc.expectedQAS012 {
				t.Errorf("QAS012: Expected %d violations, got %d", tc.expectedQAS012, len(qas012Violations))
				for _, v := range qas012Violations {
					t.Errorf("  - %s", v.Message)
				}
			}
		})
	}
}

// TestSharedUtilitiesBehavior tests specific shared utility edge cases
func TestSharedUtilitiesBehavior(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "identifier_naming_edge_cases",
			input:    "gate validGateName q { h q; }",
			expected: true, // should find one identifier
		},
		{
			name:     "comment_handling",
			input:    "qubit q; // this is a comment",
			expected: true, // should find qubit q
		},
		{
			name:     "complex_declaration",
			input:    "qubit[5] complex_array;",
			expected: true, // should find complex_array
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			declarations := FindIdentifierDeclarations(tt.input)
			hasDeclarations := len(declarations) > 0
			if hasDeclarations != tt.expected {
				t.Errorf("Input '%s': expected declarations=%v, got %d declarations",
					tt.input, tt.expected, len(declarations))
			}
		})
	}
}
