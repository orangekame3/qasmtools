package lint

import (
	"testing"

	"github.com/orangekame3/qasmtools/parser"
)

// TestNewASTBuilder tests the new proper ANTLR visitor implementation
func TestNewASTBuilder(t *testing.T) {
	code := `OPENQASM 3.0;
qubit unused_q;
qubit[2] q;
bit[2] c;
h q[0];
cx q[0], q[1];
measure q -> c;`

	// Parse with new AST builder
	p := parser.NewParser()
	result := p.ParseWithErrors(code)
	
	t.Logf("Parse errors: %d", len(result.Errors))
	for _, err := range result.Errors {
		t.Logf("  Error: %s at line %d:%d", err.Message, err.Position.Line, err.Position.Column)
	}

	if result.Program == nil {
		t.Fatal("Program is nil")
	}

	t.Logf("Program version: %v", result.Program.Version)
	t.Logf("Program statements: %d", len(result.Program.Statements))

	// Check individual statements
	for i, stmt := range result.Program.Statements {
		t.Logf("Statement %d: %T -> %s", i, stmt, stmt.String())
		
		switch s := stmt.(type) {
		case *parser.QuantumDeclaration:
			t.Logf("  Quantum: type=%s, id=%s, size=%v", s.Type, s.Identifier, s.Size)
		case *parser.ClassicalDeclaration:
			t.Logf("  Classical: type=%s, id=%s, size=%v", s.Type, s.Identifier, s.Size)
		case *parser.GateCall:
			t.Logf("  Gate: name=%s, qubits=%d", s.Name, len(s.Qubits))
		case *parser.Measurement:
			t.Logf("  Measurement: qubit=%v, target=%v", s.Qubit, s.Target)
		case *parser.Include:
			t.Logf("  Include: path=%s", s.Path)
		}
	}

	// Should have parsed statements now
	if len(result.Program.Statements) == 0 {
		t.Error("Expected to parse statements, but got none")
	}
}

// TestASTBuilderWithSimpleCode tests with minimal code
func TestASTBuilderWithSimpleCode(t *testing.T) {
	tests := []struct {
		name         string
		code         string
		expectedStmt int
	}{
		{
			name:         "simple qubit declaration",
			code:         "OPENQASM 3.0;\nqubit q;",
			expectedStmt: 1,
		},
		{
			name:         "array qubit declaration",
			code:         "OPENQASM 3.0;\nqubit[2] q;",
			expectedStmt: 1,
		},
		{
			name:         "bit declaration",
			code:         "OPENQASM 3.0;\nbit c;",
			expectedStmt: 1,
		},
		{
			name:         "gate call",
			code:         "OPENQASM 3.0;\nqubit q;\nh q;",
			expectedStmt: 2,
		},
		{
			name:         "include statement",
			code:         "OPENQASM 3.0;\ninclude \"stdgates.qasm\";",
			expectedStmt: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewParser()
			result := p.ParseWithErrors(tt.code)
			
			if result.HasErrors() {
				t.Logf("Parse errors for %s:", tt.name)
				for _, err := range result.Errors {
					t.Logf("  %s", err.Message)
				}
			}

			if result.Program == nil {
				t.Fatalf("Program is nil for %s", tt.name)
			}

			actual := len(result.Program.Statements)
			if actual != tt.expectedStmt {
				t.Errorf("Expected %d statements for %s, got %d", tt.expectedStmt, tt.name, actual)
				for i, stmt := range result.Program.Statements {
					t.Logf("  Statement %d: %T", i, stmt)
				}
			}
		})
	}
}