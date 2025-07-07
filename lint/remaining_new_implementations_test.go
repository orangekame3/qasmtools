package lint

import (
	"testing"
)

// TestRemainingNewImplementations tests QAS006-QAS011 new implementations
func TestRemainingNewImplementations(t *testing.T) {
	testCases := []struct {
		name           string
		content        string
		expectedQAS006 int // gate register size mismatch violations
		expectedQAS007 int // gate parameter indexing violations
		expectedQAS008 int // qubit declared in local scope violations
		expectedQAS009 int // illegal break/continue violations
		expectedQAS010 int // invalid instruction in gate violations
		expectedQAS011 int // reserved prefix usage violations
	}{
		{
			name: "no_violations",
			content: `OPENQASM 3.0;
qubit[2] q;
bit[2] c;
gate my_gate a, b { h a; cx a, b; }
my_gate q[0], q[1];
measure q -> c;`,
			expectedQAS006: 0,
			expectedQAS007: 0,
			expectedQAS008: 0,
			expectedQAS009: 0,
			expectedQAS010: 0,
			expectedQAS011: 0,
		},
		{
			name: "gate_size_mismatch_test",
			content: `OPENQASM 3.0;
qubit[2] q2;
qubit[3] q3;
gate my_gate a, b { h a; cx a, b; }
my_gate q2, q3;`,
			expectedQAS006: 1, // size mismatch between q2 and q3
			expectedQAS007: 0,
			expectedQAS008: 0,
			expectedQAS009: 0,
			expectedQAS010: 0,
			expectedQAS011: 0,
		},
		{
			name: "gate_parameter_indexing_test",
			content: `OPENQASM 3.0;
gate my_gate q1, q2 {
    h q1[0];
    cx q1, q2;
}`,
			expectedQAS006: 0,
			expectedQAS007: 1, // q1[0] is invalid parameter indexing
			expectedQAS008: 0,
			expectedQAS009: 0,
			expectedQAS010: 0,
			expectedQAS011: 0,
		},
		{
			name: "local_scope_qubit_test",
			content: `OPENQASM 3.0;
gate my_gate q {
    qubit local_qubit;
    h q;
}`,
			expectedQAS006: 0,
			expectedQAS007: 0,
			expectedQAS008: 1, // local_qubit declared in gate scope
			expectedQAS009: 0,
			expectedQAS010: 0,
			expectedQAS011: 0,
		},
		{
			name: "illegal_break_continue_test",
			content: `OPENQASM 3.0;
qubit q;
break;
continue;`,
			expectedQAS006: 0,
			expectedQAS007: 0,
			expectedQAS008: 0,
			expectedQAS009: 2, // break and continue outside loop
			expectedQAS010: 0,
			expectedQAS011: 0,
		},
		{
			name: "invalid_instruction_in_gate_test",
			content: `OPENQASM 3.0;
gate my_gate q {
    measure q;
    reset q;
}`,
			expectedQAS006: 0,
			expectedQAS007: 0,
			expectedQAS008: 0,
			expectedQAS009: 0,
			expectedQAS010: 2, // measure and reset in gate
			expectedQAS011: 0,
		},
		{
			name: "reserved_prefix_test",
			content: `OPENQASM 3.0;
qubit __reserved_qubit;
bit __reserved_bit;`,
			expectedQAS006: 0,
			expectedQAS007: 0,
			expectedQAS008: 0,
			expectedQAS009: 0,
			expectedQAS010: 0,
			expectedQAS011: 2, // both identifiers use reserved prefix
		},
		{
			name: "comprehensive_violations_test",
			content: `OPENQASM 3.0;
qubit[2] q2;
qubit[3] q3;
qubit __bad_qubit;

gate bad_gate q1, q2 {
    h q1[0];
    measure q1;
    qubit local_q;
}

bad_gate q2, q3;
break;`,
			expectedQAS006: 1, // size mismatch q2, q3
			expectedQAS007: 1, // q1[0] parameter indexing
			expectedQAS008: 1, // local_q in gate scope
			expectedQAS009: 1, // break outside loop
			expectedQAS010: 1, // measure in gate
			expectedQAS011: 1, // __bad_qubit reserved prefix
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			context := &CheckContext{
				File:    "test.qasm",
				Content: tc.content,
			}

			// Test new QAS006 implementation
			qas006Checker := NewGateRegisterSizeMismatchChecker()
			qas006Violations := qas006Checker.CheckProgram(context)
			if len(qas006Violations) != tc.expectedQAS006 {
				t.Errorf("QAS006: Expected %d violations, got %d", tc.expectedQAS006, len(qas006Violations))
				for _, v := range qas006Violations {
					t.Errorf("  - %s", v.Message)
				}
			}

			// Test new QAS007 implementation
			qas007Checker := NewGateParameterIndexingChecker()
			qas007Violations := qas007Checker.CheckProgram(context)
			if len(qas007Violations) != tc.expectedQAS007 {
				t.Errorf("QAS007: Expected %d violations, got %d", tc.expectedQAS007, len(qas007Violations))
				for _, v := range qas007Violations {
					t.Errorf("  - %s", v.Message)
				}
			}

			// Test new QAS008 implementation
			qas008Checker := NewQubitDeclaredInLocalScopeChecker()
			qas008Violations := qas008Checker.CheckProgram(context)
			if len(qas008Violations) != tc.expectedQAS008 {
				t.Errorf("QAS008: Expected %d violations, got %d", tc.expectedQAS008, len(qas008Violations))
				for _, v := range qas008Violations {
					t.Errorf("  - %s", v.Message)
				}
			}

			// Test new QAS009 implementation
			qas009Checker := NewIllegalBreakContinueChecker()
			qas009Violations := qas009Checker.CheckProgram(context)
			if len(qas009Violations) != tc.expectedQAS009 {
				t.Errorf("QAS009: Expected %d violations, got %d", tc.expectedQAS009, len(qas009Violations))
				for _, v := range qas009Violations {
					t.Errorf("  - %s", v.Message)
				}
			}

			// Test new QAS010 implementation
			qas010Checker := NewInvalidInstructionInGateChecker()
			qas010Violations := qas010Checker.CheckProgram(context)
			if len(qas010Violations) != tc.expectedQAS010 {
				t.Errorf("QAS010: Expected %d violations, got %d", tc.expectedQAS010, len(qas010Violations))
				for _, v := range qas010Violations {
					t.Errorf("  - %s", v.Message)
				}
			}

			// Test new QAS011 implementation
			qas011Checker := NewReservedPrefixUsageChecker()
			qas011Violations := qas011Checker.CheckProgram(context)
			if len(qas011Violations) != tc.expectedQAS011 {
				t.Errorf("QAS011: Expected %d violations, got %d", tc.expectedQAS011, len(qas011Violations))
				for _, v := range qas011Violations {
					t.Errorf("  - %s", v.Message)
				}
			}
		})
	}
}
