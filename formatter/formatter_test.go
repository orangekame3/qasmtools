package formatter

import (
	"strings"
	"testing"
)

func TestFormatQASM(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "version and include statements",
			input: `OPENQASM 3.0;
include"stdgates.qasm";`,
			expected: `OPENQASM 3.0;
include "stdgates.qasm";
`,
		},
		{
			name: "indentation in blocks",
			input: `OPENQASM 3.0;
gate h q {
x q;
}`,
			expected: `OPENQASM 3.0;

gate h q {
  x q;
}
`,
		},
		{
			name: "spacing around binary operators",
			input: `bit[2]c=0;
int[32] x=5+3*2;`,
			expected: `bit[2] c = 0;
int[32] x = 5 + 3 * 2;
`,
		},
		{
			name: "no spaces in brackets",
			input: `qubit[ 2 ] q;
h q[ 0 ];`,
			expected: `qubit[2] q;
h q[0];
`,
		},
		{
			name: "comments preservation",
			input: `// This is a header comment
OPENQASM 3.0;
// Include standard gates
include "stdgates.qasm"; // Standard gates
/* Multi-line
   comment */
qubit[2] q;  // Qubit declaration`,
			expected: `// This is a header comment
OPENQASM 3.0;
// Include standard gates
include "stdgates.qasm"; // Standard gates

/* Multi-line
   comment */
qubit[2] q; // Qubit declaration
`,
		},
		{
			name: "timing expressions",
			input: `delay[100ns] q;
delay[  50  ns  ] q;`,
			expected: `delay[100ns] q;
delay[50ns] q;
`,
		},
		{
			name: "empty lines between major blocks",
			input: `OPENQASM 3.0;
include "stdgates.qasm";
qubit[2] q;
bit c;
gate custom a, b {
  cx a, b;
}
h q[0];
measure q[0] -> c;`,
			expected: `OPENQASM 3.0;
include "stdgates.qasm";

qubit[2] q;
bit c;

gate custom a, b {
  cx a, b;
}

h q[0];
measure q[0] -> c;
`,
		},
		{
			name: "if statement formatting",
			input: `if(c==1){
h q;
}else{
x q;
}`,
			expected: `if (c == 1) {
  h q;
} else {
  x q;
}
`,
		},
		{
			name: "gate call with parameters",
			input: `rz(pi/2)q[0];
cphase(pi/4)q[0],q[1];`,
			expected: `rz(pi/2) q[0];
cphase(pi/4) q[0], q[1];
`,
		},
		{
			name: "measurement formatting",
			input: `measureq[0]->c[0];
measure q->c;`,
			expected: `measure q[0] -> c[0];
measure q -> c;
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatted, err := FormatQASM(tt.input)
			if err != nil {
				t.Errorf("FormatQASM() error = %v", err)
				return
			}

			// Compare line by line for better error messages
			expectedLines := strings.Split(tt.expected, "\n")
			actualLines := strings.Split(formatted, "\n")

			// Check if number of lines match
			if len(expectedLines) != len(actualLines) {
				t.Errorf("Line count mismatch:\nexpected %d lines:\n%s\ngot %d lines:\n%s",
					len(expectedLines), tt.expected, len(actualLines), formatted)
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
