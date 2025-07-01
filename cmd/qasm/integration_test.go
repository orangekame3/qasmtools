package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestUnescapeIntegration(t *testing.T) {
	// Build the binary first
	binPath := filepath.Join(os.TempDir(), "qasm_test")
	cmd := exec.Command("go", "build", "-o", binPath, ".")
	if err := cmd.Run(); err != nil {
		t.Skipf("Failed to build qasm binary: %v", err)
	}
	defer os.Remove(binPath)

	tests := []struct {
		name           string
		input          string
		args           []string
		expectedOutput string
		expectError    bool
	}{
		{
			name:           "stdin_normal",
			input:          "OPENQASM 3.0;qubit q;h q;",
			args:           []string{"fmt", "--stdin"},
			expectedOutput: "OPENQASM 3.0;\nqubit q;\nh q;\n",
			expectError:    false,
		},
		{
			name:           "stdin_unescape",
			input:          `"OPENQASM 3.0;\nqubit q;\nh q;"`,
			args:           []string{"fmt", "--stdin", "--unescape"},
			expectedOutput: "OPENQASM 3.0;\nqubit q;\nh q;\n",
			expectError:    false,
		},
		{
			name:           "stdin_complex_unescape",
			input:          `"OPENQASM 3.0;\ninclude \"stdgates.inc\";\nqubit[2] q;\nh q[0];"`,
			args:           []string{"fmt", "--stdin", "--unescape"},
			expectedOutput: "OPENQASM 3.0;\ninclude \"stdgates.inc\";\n\nqubit[2] q;\nh q[0];\n",
			expectError:    false,
		},
		{
			name:        "stdin_invalid_unescape",
			input:       `"OPENQASM 3.0;\nqubit q;`,
			args:        []string{"fmt", "--stdin", "--unescape"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binPath, tt.args...)
			cmd.Stdin = bytes.NewBufferString(tt.input)
			
			output, err := cmd.Output()
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but command succeeded")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if string(output) != tt.expectedOutput {
					t.Errorf("Output = %q, want %q", string(output), tt.expectedOutput)
				}
			}
		})
	}
}

func TestFileUnescapeIntegration(t *testing.T) {
	// Build the binary first
	binPath := filepath.Join(os.TempDir(), "qasm_test")
	cmd := exec.Command("go", "build", "-o", binPath, ".")
	if err := cmd.Run(); err != nil {
		t.Skipf("Failed to build qasm binary: %v", err)
	}
	defer os.Remove(binPath)

	tests := []struct {
		name           string
		fileContent    string
		useUnescape    bool
		expectedOutput string
		expectError    bool
	}{
		{
			name:           "file_normal",
			fileContent:    "OPENQASM 3.0;qubit q;h q;",
			useUnescape:    false,
			expectedOutput: "OPENQASM 3.0;\nqubit q;\nh q;\n",
			expectError:    false,
		},
		{
			name:           "file_escaped",
			fileContent:    `"OPENQASM 3.0;\nqubit q;\nh q;"`,
			useUnescape:    true,
			expectedOutput: "OPENQASM 3.0;\nqubit q;\nh q;\n",
			expectError:    false,
		},
		{
			name:        "file_invalid_escaped",
			fileContent: `"OPENQASM 3.0;\nqubit q;`,
			useUnescape: true,
			expectError: true,
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

			// Write test content
			if _, err := tmpFile.WriteString(tt.fileContent); err != nil {
				t.Fatalf("Failed to write to temp file: %v", err)
			}
			tmpFile.Close()

			// Prepare command arguments
			args := []string{"fmt"}
			if tt.useUnescape {
				args = append(args, "--unescape")
			}
			args = append(args, tmpFile.Name())

			cmd := exec.Command(binPath, args...)
			output, err := cmd.Output()

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but command succeeded")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if string(output) != tt.expectedOutput {
					t.Errorf("Output = %q, want %q", string(output), tt.expectedOutput)
				}
			}
		})
	}
}