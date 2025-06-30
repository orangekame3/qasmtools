package formatter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func readFile(t *testing.T, path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file %s: %v", path, err)
	}
	return string(content)
}

func TestFormatQASM(t *testing.T) {
	files, err := os.ReadDir("../testdata/formatter/input")
	if err != nil {
		t.Fatalf("failed to read test data directory: %v", err)
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".qasm") {
			continue
		}

		testName := strings.TrimSuffix(file.Name(), ".qasm")
		t.Run(testName, func(t *testing.T) {
			inputPath := filepath.Join("../testdata/formatter/input", file.Name())
			expectedPath := filepath.Join("../testdata/formatter/expected", file.Name())

			input := readFile(t, inputPath)
			expected := readFile(t, expectedPath)

			formatted, err := FormatQASM(input)
			if err != nil {
				t.Errorf("FormatQASM() error = %v", err)
				return
			}

			// Compare line by line for better error messages
			expectedLines := strings.Split(expected, "\n")
			actualLines := strings.Split(formatted, "\n")

			// Check if number of lines match
			if len(expectedLines) != len(actualLines) {
				t.Errorf("Line count mismatch:\nexpected %d lines:\n%s\ngot %d lines:\n%s",
					len(expectedLines), expected, len(actualLines), formatted)
				return
			}

			// Compare each line
			for i := range expectedLines {
				if expectedLines[i] != actualLines[i] {
					t.Errorf("Line %d mismatch:\nexpected: %q\ngot: %q",
						i+1, expectedLines[i], actualLines[i])
				}
			}

			// Test idempotency: format(format(input)) should equal format(input)
			formattedTwice, err := FormatQASM(formatted)
			if err != nil {
				t.Errorf("Second FormatQASM() error = %v", err)
				return
			}

			if formatted != formattedTwice {
				t.Errorf("Formatter is not idempotent!\nFirst format:\n%s\nSecond format:\n%s",
					formatted, formattedTwice)
			}
		})
	}
}

func TestFormatQASMWithConfig(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		config   *Config
		expected string
	}{
		{
			name: "custom indent size",
			input: `gate h q {
x q;
}`,
			config: &Config{
				Indent:  4,
				Newline: true,
			},
			expected: `gate h q {
    x q;
}
`,
		},
		{
			name: "no trailing newline",
			input: `OPENQASM 3.0;
include "stdgates.qasm";`,
			config: &Config{
				Indent:  2,
				Newline: false,
			},
			expected: `OPENQASM 3.0;
include "stdgates.qasm";`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatted, err := FormatQASMWithConfig(tt.input, tt.config)
			if err != nil {
				t.Errorf("FormatQASMWithConfig() error = %v", err)
				return
			}

			if formatted != tt.expected {
				t.Errorf("FormatQASMWithConfig() = %v, want %v", formatted, tt.expected)
			}

			// Test idempotency for config-based formatting
			formattedTwice, err := FormatQASMWithConfig(formatted, tt.config)
			if err != nil {
				t.Errorf("Second FormatQASMWithConfig() error = %v", err)
				return
			}

			if formatted != formattedTwice {
				t.Errorf("Config-based formatter is not idempotent!\nFirst format:\n%s\nSecond format:\n%s",
					formatted, formattedTwice)
			}
		})
	}
}

func TestValidateQASM(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name: "valid QASM",
			input: `OPENQASM 3.0;
include "stdgates.qasm";
qubit[2] q;
h q[0];`,
			wantErr: false,
		},
		{
			name:    "empty input",
			input:   "",
			wantErr: true,
		},
		{
			name: "invalid syntax",
			input: `OPENQASM 3.0;
invalid command;`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateQASM(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateQASM() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestFormatterIdempotency tests specific edge cases for idempotency
func TestFormatterIdempotency(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "compound statements",
			input: `OPENQASM 3.0;include"stdgates.qasm";qubit[2]q;bit[2]c;hq[0];cxq[0],q[1];measureq->c;`,
		},
		{
			name: "already formatted with empty lines",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

qubit[2] q;
bit[2] c;
h q[0];
cx q[0], q[1];
measure q -> c;`,
		},
		{
			name: "comments and empty lines",
			input: `// Header comment
OPENQASM 3.0;
// Include comment
include "stdgates.qasm";

/* Block comment */
qubit[2] q;
bit[2] c;  // Trailing comment`,
		},
		{
			name: "complex nested structures",
			input: `OPENQASM 3.0;
include "stdgates.qasm";

gate bell_prep q0, q1 {
    h q0;
    cx q0, q1;
}

qubit[2] q;
bit[2] c;

if (c[0] == 1) {
    x q[1];
}`,
		},
		{
			name: "malformed spacing",
			input: `OPENQASM   3.0  ; include   "stdgates.qasm"   ;qubit [ 2 ] q ; bit[ 2]c;h  q[  0 ];`,
		},
		{
			name: "empty and whitespace lines",
			input: `OPENQASM 3.0;


include "stdgates.qasm";   


qubit[2] q;

    
bit[2] c;`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// First format
			formatted1, err := FormatQASM(tt.input)
			if err != nil {
				t.Errorf("First FormatQASM() error = %v", err)
				return
			}

			// Second format (should be identical)
			formatted2, err := FormatQASM(formatted1)
			if err != nil {
				t.Errorf("Second FormatQASM() error = %v", err)
				return
			}

			// Third format (extra verification)
			formatted3, err := FormatQASM(formatted2)
			if err != nil {
				t.Errorf("Third FormatQASM() error = %v", err)
				return
			}

			// Check that all three are identical
			if formatted1 != formatted2 {
				t.Errorf("Formatter is not idempotent (1st vs 2nd)!\nFirst format:\n%s\nSecond format:\n%s",
					formatted1, formatted2)
			}

			if formatted2 != formatted3 {
				t.Errorf("Formatter is not idempotent (2nd vs 3rd)!\nSecond format:\n%s\nThird format:\n%s",
					formatted2, formatted3)
			}

			// Additional check: verify length consistency
			if len(formatted1) != len(formatted2) || len(formatted2) != len(formatted3) {
				t.Errorf("Formatter output length inconsistent!\nLengths: %d, %d, %d",
					len(formatted1), len(formatted2), len(formatted3))
			}
		})
	}
}
