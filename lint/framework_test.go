package lint

import (
	"testing"
)

// TestFrameworkMigration tests that new implementations work correctly
func TestFrameworkMigration(t *testing.T) {
	// This test focuses on validating the new framework works correctly
	// rather than comparing with old implementations that may have bugs
	t.Skip("Replaced by TestNewImplementationCompatibility in framework_integration_test.go")
}

// TestSharedUtilities tests the shared utility functions
func TestSharedUtilities(t *testing.T) {
	tests := []struct {
		name     string
		function func() bool
	}{
		{
			name: "RemoveComments",
			function: func() bool {
				line := "qubit q; // this is a comment"
				result := RemoveComments(line)
				return result == "qubit q; "
			},
		},
		{
			name: "ExtractRegisterName",
			function: func() bool {
				reg := "myQubit[0]"
				result := ExtractRegisterName(reg)
				return result == "myQubit"
			},
		},
		{
			name: "SkipCommentAndEmptyLine",
			function: func() bool {
				return SkipCommentAndEmptyLine("// comment") &&
					SkipCommentAndEmptyLine("   ") &&
					!SkipCommentAndEmptyLine("qubit q;")
			},
		},
		{
			name: "IsKeyword",
			function: func() bool {
				return IsKeyword("qubit") &&
					IsKeyword("gate") &&
					!IsKeyword("myVariable")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.function() {
				t.Errorf("Utility function %s failed", tt.name)
			}
		})
	}
}

// TestBaseCheckerFunctionality tests BaseChecker framework
func TestBaseCheckerFunctionality(t *testing.T) {
	checker := NewBaseChecker("TEST")

	if checker.GetRuleID() != "TEST" {
		t.Errorf("Expected rule ID 'TEST', got '%s'", checker.GetRuleID())
	}

	// Test violation builder
	violation := checker.NewViolationBuilder().
		WithMessage("Test message").
		WithFile("test.qasm").
		WithPosition(1, 1).
		WithNodeName("testNode").
		AsError().
		Build()

	if violation.Message != "Test message" {
		t.Errorf("Expected message 'Test message', got '%s'", violation.Message)
	}

	if violation.Severity != SeverityError {
		t.Errorf("Expected error severity, got %s", violation.Severity)
	}
}

// TestFindIdentifierDeclarations tests the shared identifier finding function
func TestFindIdentifierDeclarations(t *testing.T) {
	testCases := []struct {
		line     string
		expected int
	}{
		{"qubit myQubit;", 1},
		{"bit[2] myBit;", 1},
		{"gate myGate q { h q; }", 1},
		{"// comment line", 0},
		{"", 0},
		{"qubit q1; bit b1;", 2}, // Multiple declarations
	}

	for _, tc := range testCases {
		t.Run(tc.line, func(t *testing.T) {
			declarations := FindIdentifierDeclarations(tc.line)
			if len(declarations) != tc.expected {
				t.Errorf("Line '%s': expected %d declarations, got %d",
					tc.line, tc.expected, len(declarations))
			}
		})
	}
}
