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
