package lint

import (
	"testing"
)

// TestNewImplementations tests the new framework implementations
func TestNewImplementations(t *testing.T) {
	testCases := []struct {
		name           string
		content        string
		expectedQAS002 int // undefined identifier violations
		expectedQAS003 int // constant measured bit violations
		expectedQAS004 int // out of bounds violations
	}{
		{
			name: "undefined_identifier_test",
			content: `OPENQASM 3.0;
qubit defined_qubit;
bit result;
h defined_qubit;
cx defined_qubit, undefined_qubit;
measure defined_qubit -> result;`,
			expectedQAS002: 1, // undefined_qubit is undefined
			expectedQAS003: 0, // defined_qubit has gate applied before measurement
			expectedQAS004: 0, // no array access
		},
		{
			name: "constant_measurement_test",
			content: `OPENQASM 3.0;
qubit constant_qubit;
qubit gate_applied_qubit;
bit[2] result;
h gate_applied_qubit;
measure constant_qubit -> result[0];
measure gate_applied_qubit -> result[1];`,
			expectedQAS002: 0, // all identifiers defined
			expectedQAS003: 1, // constant_qubit measured without gates
			expectedQAS004: 0, // array access within bounds
		},
		{
			name: "out_of_bounds_test",
			content: `OPENQASM 3.0;
qubit[3] test_array;
bit[2] small_array;
h test_array[0];
h test_array[5];
measure test_array[0] -> small_array[3];`,
			expectedQAS002: 0, // all identifiers defined
			expectedQAS003: 0, // qubits have gates applied
			expectedQAS004: 2, // test_array[5] and small_array[3] out of bounds
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			context := &CheckContext{
				File:    "test.qasm",
				Content: tc.content,
			}

			// Test new QAS002 implementation
			qas002Checker := NewUndefinedIdentifierChecker()
			qas002Violations := qas002Checker.CheckProgram(context)
			if len(qas002Violations) != tc.expectedQAS002 {
				t.Errorf("QAS002: Expected %d violations, got %d", tc.expectedQAS002, len(qas002Violations))
				for _, v := range qas002Violations {
					t.Errorf("  - %s", v.Message)
				}
			}

			// Test new QAS003 implementation
			qas003Checker := NewConstantMeasuredBitChecker()
			qas003Violations := qas003Checker.CheckProgram(context)
			if len(qas003Violations) != tc.expectedQAS003 {
				t.Errorf("QAS003: Expected %d violations, got %d", tc.expectedQAS003, len(qas003Violations))
				for _, v := range qas003Violations {
					t.Errorf("  - %s", v.Message)
				}
			}

			// Test new QAS004 implementation
			qas004Checker := NewOutOfBoundsIndexChecker()
			qas004Violations := qas004Checker.CheckProgram(context)
			if len(qas004Violations) != tc.expectedQAS004 {
				t.Errorf("QAS004: Expected %d violations, got %d", tc.expectedQAS004, len(qas004Violations))
				for _, v := range qas004Violations {
					t.Errorf("  - %s", v.Message)
				}
			}
		})
	}
}
