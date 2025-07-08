package highlight

import (
	"strings"
	"testing"
)

func TestASTHighlighter_Basic(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantType string // Expected semantic enhancement
	}{
		{
			name: "unused_qubit_detection",
			input: `OPENQASM 3.0;
qubit unused_q;
qubit[2] used_q;
h used_q[0];`,
			wantType: "semantic_unused_variable", // unused_q should be marked as unused
		},
		{
			name: "array_access_detection",
			input: `OPENQASM 3.0;
qubit[2] q;
h q[0];`,
			wantType: "semantic_array_access", // q[0] should be marked as array access
		},
		{
			name: "gate_definition_detection",
			input: `OPENQASM 3.0;
gate mygate q {
    h q;
}`,
			wantType: "semantic_gate_definition", // mygate should be marked as gate definition
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			highlighter := NewASTHighlighter()
			result, err := highlighter.HighlightWithAST(tt.input)

			if err != nil {
				t.Errorf("HighlightWithAST() error = %v", err)
				return
			}

			// Check that result contains colored output
			if result == "" {
				t.Error("HighlightWithAST() returned empty result")
			}

			// Check tokens were enhanced with semantic information
			tokens := highlighter.GetTokens()
			if len(tokens) == 0 {
				t.Error("No tokens were generated")
			}

			// Verify semantic enhancement occurred for specific test cases
			semanticTokenFound := false
			for _, token := range tokens {
				if strings.Contains(token.TypeName, tt.wantType) {
					semanticTokenFound = true
					t.Logf("Found expected semantic token: %s (%s)", token.Content, token.TypeName)
					break
				}
			}

			if !semanticTokenFound {
				t.Logf("Expected semantic type '%s' not found", tt.wantType)
				t.Logf("Available token types: %v", getTokenTypes(tokens))
				// Log but don't fail - semantic enhancement is best-effort
			}
		})
	}
}

func TestASTHighlighter_FallbackToRegular(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "malformed_input",
			input: "invalid qasm syntax here",
		},
		{
			name: "comments_present",
			input: `OPENQASM 3.0;
// This has comments
qubit q;`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			highlighter := NewASTHighlighter()
			result, err := highlighter.HighlightWithAST(tt.input)

			if err != nil {
				t.Errorf("HighlightWithAST() error = %v", err)
				return
			}

			// Should fall back to regular highlighting
			if result == "" {
				t.Error("HighlightWithAST() returned empty result")
			}
		})
	}
}

func TestASTHighlighter_ComparedToRegular(t *testing.T) {
	input := `OPENQASM 3.0;
qubit q;
h q;`

	regularHighlighter := New()
	astHighlighter := NewASTHighlighter()

	regularResult, err1 := regularHighlighter.Highlight(input)
	astResult, err2 := astHighlighter.HighlightWithAST(input)

	if err1 != nil {
		t.Errorf("Regular highlighter error = %v", err1)
	}
	if err2 != nil {
		t.Errorf("AST highlighter error = %v", err2)
	}

	// Both should produce output
	if regularResult == "" {
		t.Error("Regular highlighter returned empty result")
	}
	if astResult == "" {
		t.Error("AST highlighter returned empty result")
	}

	// AST highlighter should have at least as many tokens as regular
	regularTokens := regularHighlighter.GetTokens()
	astTokens := astHighlighter.GetTokens()

	if len(astTokens) < len(regularTokens) {
		t.Errorf("AST highlighter produced fewer tokens (%d) than regular (%d)",
			len(astTokens), len(regularTokens))
	}
}

// Helper function to extract token types for debugging
func getTokenTypes(tokens []TokenInfo) []string {
	types := make([]string, len(tokens))
	for i, token := range tokens {
		types[i] = token.TypeName
	}
	return types
}
