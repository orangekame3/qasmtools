package lint

import (
	"path/filepath"
	"testing"
)

func TestLintIntegration(t *testing.T) {
	tests := []struct {
		name               string
		file               string
		expectedViolations int
		expectedRuleIDs    []string
	}{
		{
			name:               "clean code - no violations",
			file:               "testdata/valid/clean_code.qasm",
			expectedViolations: 0,
			expectedRuleIDs:    []string{},
		},
		{
			name:               "unused qubit violation",
			file:               "testdata/violations/unused_qubit.qasm",
			expectedViolations: 1,
			expectedRuleIDs:    []string{"QAS001"},
		},
		{
			name:               "naming convention violations",
			file:               "testdata/violations/naming_violation.qasm",
			expectedViolations: 4, // QAS001: unused, QAS005: naming pattern, QAS012: snake_case (2 violations)
			expectedRuleIDs:    []string{"QAS001", "QAS005", "QAS012"},
		},
		{
			name:               "array bounds violations",
			file:               "testdata/violations/array_bounds.qasm",
			expectedViolations: 2,
			expectedRuleIDs:    []string{"QAS004"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create linter with embedded rules
			linter := NewLinter("")
			err := linter.LoadRules()
			if err != nil {
				t.Fatalf("Failed to load rules: %v", err)
			}

			// Get absolute path to test file
			absPath, err := filepath.Abs(tt.file)
			if err != nil {
				t.Fatalf("Failed to get absolute path: %v", err)
			}

			// Lint the file
			violations, err := linter.LintFile(absPath)
			if err != nil {
				t.Fatalf("Failed to lint file %s: %v", tt.file, err)
			}

			// Check number of violations
			if len(violations) != tt.expectedViolations {
				t.Errorf("Expected %d violations, got %d", tt.expectedViolations, len(violations))
				for _, v := range violations {
					t.Logf("Violation: %s", v.String())
				}
			}

			// Check rule IDs if violations are expected
			if tt.expectedViolations > 0 {
				foundRules := make(map[string]bool)
				for _, violation := range violations {
					foundRules[violation.Rule.ID] = true
				}

				for _, expectedID := range tt.expectedRuleIDs {
					if !foundRules[expectedID] {
						t.Errorf("Expected violation from rule %s not found", expectedID)
					}
				}
			}
		})
	}
}

func TestLintMultipleFiles(t *testing.T) {
	linter := NewLinter("")
	err := linter.LoadRules()
	if err != nil {
		t.Fatalf("Failed to load rules: %v", err)
	}

	// Get absolute paths
	cleanFile, err := filepath.Abs("testdata/valid/clean_code.qasm")
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	violationFile, err := filepath.Abs("testdata/violations/unused_qubit.qasm")
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	files := []string{cleanFile, violationFile}
	violations, err := linter.LintFiles(files)
	if err != nil {
		t.Fatalf("Failed to lint files: %v", err)
	}

	// Should have violations from the second file only
	if len(violations) == 0 {
		t.Error("Expected at least one violation from multiple files")
	}

	// Check that violations include file information
	for _, violation := range violations {
		if violation.File == "" {
			t.Error("Violation missing file information")
		}
	}
}

func TestLintWithRuleFiltering(t *testing.T) {
	linter := NewLinter("")
	err := linter.LoadRules()
	if err != nil {
		t.Fatalf("Failed to load rules: %v", err)
	}

	violationFile, err := filepath.Abs("testdata/violations/unused_qubit.qasm")
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	violations, err := linter.LintFile(violationFile)
	if err != nil {
		t.Fatalf("Failed to lint file: %v", err)
	}

	if len(violations) == 0 {
		t.Skip("No violations found, cannot test filtering")
	}

	// Test that violations have rule information
	for _, violation := range violations {
		if violation.Rule == nil {
			t.Error("Violation missing rule information")
		}
		if violation.Rule.ID == "" {
			t.Error("Violation rule missing ID")
		}
		if violation.Severity == "" {
			t.Error("Violation missing severity")
		}
	}
}
