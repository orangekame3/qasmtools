package highlight

import (
	"strings"
	"testing"
)

func TestHighlighter_ColoredOutput(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name: "basic keywords and operators",
			input: `OPENQASM 3.0;
include "stdgates.inc";
gate h q { }`,
		},
		{
			name: "measurements and comments",
			input: `// This is a comment
measure q -> c; // Measurement`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := New()
			output, err := h.Highlight(tt.input)
			if err != nil {
				t.Errorf("Highlight() error = %v", err)
				return
			}

			// Verify that the output contains ANSI color codes
			if !strings.Contains(output, "\033[") {
				t.Error("Output does not contain ANSI color codes")
			}

			// Verify that each line ends with reset code
			lines := strings.Split(output, "\n")
			for _, line := range lines {
				if line != "" && !strings.HasSuffix(line, string(Reset)) {
					t.Errorf("Line does not end with reset code: %q", line)
				}
			}
		})
	}
}

func TestHighlighter_TokenInfo(t *testing.T) {
	input := `OPENQASM 3.0;
include "stdgates.inc";

// Initialize qubits
qubit[2] q;
bit[2] c;

// Apply gates
h q[0];
cx q[0], q[1];

// Measure
measure q -> c;`

	h := New()
	_, err := h.Highlight(input)
	if err != nil {
		t.Errorf("Highlight() error = %v", err)
		return
	}

	tokens := h.GetTokens()
	if len(tokens) == 0 {
		t.Error("GetTokens() returned no tokens")
		return
	}

	// Check specific token positions
	expectedChecks := []struct {
		pos      int
		wantType string
		content  string
	}{
		{0, "keyword", "OPENQASM"},
		{1, "number", "3.0"},
		{2, "punctuation", ";"},
		{3, "keyword", "include"},
		{4, "string", "\"stdgates.inc\""},
	}

	for _, check := range expectedChecks {
		if check.pos >= len(tokens) {
			t.Errorf("Token position %d out of range (got %d tokens)", check.pos, len(tokens))
			continue
		}

		token := tokens[check.pos]
		if token.TypeName != check.wantType {
			t.Errorf("Token[%d] type = %v, want %v", check.pos, token.TypeName, check.wantType)
		}
		if token.Content != check.content {
			t.Errorf("Token[%d] content = %v, want %v", check.pos, token.Content, check.content)
		}
	}
}

func TestHighlighter_CustomColorScheme(t *testing.T) {
	scheme := &ColorScheme{
		Keyword:     Blue,
		Operator:    Green,
		Identifier:  Yellow,
		Number:      Red,
		String:      Magenta,
		Comment:     Cyan,
		Gate:        Red,
		Measurement: Blue,
		Register:    Green,
		Punctuation: White,
	}

	h := NewWithColorScheme(scheme)
	input := `OPENQASM 3.0;
// Test comment
measure q -> c;`

	output, err := h.Highlight(input)
	if err != nil {
		t.Errorf("Highlight() error = %v", err)
		return
	}

	// Verify that the output contains the custom colors
	if !strings.Contains(output, string(Blue)) ||
		!strings.Contains(output, string(Cyan)) ||
		!strings.Contains(output, string(Yellow)) {
		t.Errorf("Output does not contain expected custom colors. Output: %q", output)
		t.Errorf("Looking for Blue: %q, Cyan: %q, Yellow: %q", string(Blue), string(Cyan), string(Yellow))
	}
}
