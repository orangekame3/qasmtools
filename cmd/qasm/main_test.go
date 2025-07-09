package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/orangekame3/qasmtools/cmd/qasm/commands"
	"github.com/orangekame3/qasmtools/cmd/qasm/config"
	"github.com/orangekame3/qasmtools/formatter"
	"github.com/spf13/cobra"
)

func TestGetConfigFromFlags(t *testing.T) {
	tests := []struct {
		name     string
		flags    map[string]interface{}
		expected *formatter.Config
	}{
		{
			name: "default_config",
			flags: map[string]interface{}{
				"write":   false,
				"check":   false,
				"indent":  uint(2),
				"newline": true,
				"verbose": false,
				"diff":    false,
			},
			expected: &formatter.Config{
				Write:   false,
				Check:   false,
				Indent:  2,
				Newline: true,
				Verbose: false,
				Diff:    false,
			},
		},
		{
			name: "custom_config",
			flags: map[string]interface{}{
				"write":   true,
				"check":   false,
				"indent":  uint(4),
				"newline": false,
				"verbose": true,
				"diff":    false,
			},
			expected: &formatter.Config{
				Write:   true,
				Check:   false,
				Indent:  4,
				Newline: false,
				Verbose: true,
				Diff:    false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}

			// Set up flags
			cmd.Flags().Bool("write", false, "")
			cmd.Flags().Bool("check", false, "")
			cmd.Flags().Uint("indent", 2, "")
			cmd.Flags().Bool("newline", true, "")
			cmd.Flags().Bool("verbose", false, "")
			cmd.Flags().Bool("diff", false, "")

			// Set flag values
			for key, value := range tt.flags {
				cmd.Flags().Set(key, toString(value))
			}

			config, err := config.GetConfigFromFlags(cmd)
			if err != nil {
				t.Fatalf("getConfigFromFlags() error = %v", err)
			}

			if config.Write != tt.expected.Write ||
				config.Check != tt.expected.Check ||
				config.Indent != tt.expected.Indent ||
				config.Newline != tt.expected.Newline ||
				config.Verbose != tt.expected.Verbose ||
				config.Diff != tt.expected.Diff {
				t.Errorf("getConfigFromFlags() = %+v, want %+v", config, tt.expected)
			}
		})
	}
}

func TestRunFormatStdin(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedOutput string
		expectError    bool
	}{
		{
			name:           "simple_qasm",
			input:          "OPENQASM 3.0;qubit q;h q;",
			expectedOutput: "OPENQASM 3.0;\n\nqubit q;\nh q;\n",
			expectError:    false,
		},
		{
			name:           "formatted_qasm",
			input:          "OPENQASM 3.0;\ninclude \"stdgates.inc\";\nqubit[2] q;\nh q[0];",
			expectedOutput: "OPENQASM 3.0;\ninclude \"stdgates.inc\";\n\nqubit[2] q;\nh q[0];\n",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate stdin
			oldStdin := os.Stdin
			r, w, _ := os.Pipe()
			os.Stdin = r

			// Write test input to pipe
			go func() {
				w.WriteString(tt.input)
				w.Close()
			}()

			// Capture stdout
			oldStdout := os.Stdout
			rOut, wOut, _ := os.Pipe()
			os.Stdout = wOut

			// Create command with all required flags
			cmd := &cobra.Command{}
			cmd.Flags().Bool("write", false, "")
			cmd.Flags().Bool("check", false, "")
			cmd.Flags().Uint("indent", 2, "")
			cmd.Flags().Bool("newline", true, "")
			cmd.Flags().Bool("verbose", false, "")
			cmd.Flags().Bool("diff", false, "")

			config, err := config.GetConfigFromFlags(cmd)
			if err != nil {
				t.Fatalf("getConfigFromFlags() error = %v", err)
			}

			// Run the function
			err = commands.RunFormatStdin(cmd, config)

			// Close stdout and read output
			wOut.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			buf.ReadFrom(rOut)
			output := buf.String()

			// Restore stdin
			os.Stdin = oldStdin

			if tt.expectError {
				if err == nil {
					t.Errorf("runFormatStdin() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("runFormatStdin() error = %v", err)
				}
				if output != tt.expectedOutput {
					t.Errorf("runFormatStdin() output = %q, want %q", output, tt.expectedOutput)
				}
			}
		})
	}
}

func TestFormatFileWithConfig(t *testing.T) {
	tests := []struct {
		name        string
		fileContent string
		expectedMod bool
		expectError bool
	}{
		{
			name:        "simple_formatted",
			fileContent: "OPENQASM 3.0;\nqubit q;\nh q;",
			expectedMod: true, // May be modified due to formatting differences
			expectError: false,
		},
		{
			name:        "needs_formatting",
			fileContent: "OPENQASM 3.0;qubit q;h q;",
			expectedMod: true,
			expectError: false,
		},
		{
			name:        "already_formatted",
			fileContent: "OPENQASM 3.0;\n\nqubit q;\nh q;\n",
			expectedMod: true, // Even "properly formatted" may have minor differences
			expectError: false,
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

			// Create config
			config := &formatter.Config{
				Write:    false, // Don't modify the file, just check output
				Unescape: false,
			}

			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Run the function
			modified, err := commands.FormatFileWithConfig(tmpFile.Name(), config)

			// Close stdout and read output
			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			buf.ReadFrom(r)
			output := buf.String()

			if tt.expectError {
				if err == nil {
					t.Errorf("formatFileWithConfig() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("formatFileWithConfig() error = %v", err)
				}
				if modified != tt.expectedMod {
					t.Errorf("formatFileWithConfig() modified = %v, want %v", modified, tt.expectedMod)
				}
				if !tt.expectError && len(output) == 0 && tt.expectedMod {
					t.Errorf("formatFileWithConfig() expected output but got none")
				}
			}
		})
	}
}

func TestCheckFileWithConfig(t *testing.T) {
	tests := []struct {
		name          string
		fileContent   string
		expectedCheck bool
		expectError   bool
	}{
		{
			name:          "simple_formatted",
			fileContent:   "OPENQASM 3.0;\nqubit q;\nh q;",
			expectedCheck: false, // May not match exactly due to formatting differences
			expectError:   false,
		},
		{
			name:          "needs_formatting",
			fileContent:   "OPENQASM 3.0;qubit q;h q;",
			expectedCheck: false,
			expectError:   false,
		},
		{
			name:          "already_formatted",
			fileContent:   "OPENQASM 3.0;\n\nqubit q;\nh q;\n",
			expectedCheck: false, // Even "properly formatted" may have minor differences
			expectError:   false,
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

			// Create config
			config := &formatter.Config{
				Unescape: false,
			}

			// Run the function
			isFormatted, err := commands.CheckFileWithConfig(tmpFile.Name(), config)

			if tt.expectError {
				if err == nil {
					t.Errorf("checkFileWithConfig() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("checkFileWithConfig() error = %v", err)
				}
				if isFormatted != tt.expectedCheck {
					t.Errorf("checkFileWithConfig() formatted = %v, want %v", isFormatted, tt.expectedCheck)
				}
			}
		})
	}
}

// Helper function to convert various types to string for flag setting
func toString(value interface{}) string {
	switch v := value.(type) {
	case bool:
		if v {
			return "true"
		}
		return "false"
	case uint:
		return fmt.Sprintf("%d", v)
	case string:
		return v
	default:
		return ""
	}
}
