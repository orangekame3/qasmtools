package ioutils

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss/v2"
)

// ShowDiff displays a formatted diff between original and formatted content
func ShowDiff(filename, original, formatted string) error {
	// Simple diff implementation
	originalLines := strings.Split(original, "\n")
	formattedLines := strings.Split(formatted, "\n")

	fmt.Fprintf(os.Stderr, "--- %s (original)\n", filename)
	fmt.Fprintf(os.Stderr, "+++ %s (formatted)\n", filename)

	// Find differences line by line
	maxLines := len(originalLines)
	if len(formattedLines) > maxLines {
		maxLines = len(formattedLines)
	}

	for i := 0; i < maxLines; i++ {
		var origLine, formLine string
		if i < len(originalLines) {
			origLine = originalLines[i]
		}
		if i < len(formattedLines) {
			formLine = formattedLines[i]
		}

		if origLine != formLine {
			if origLine != "" {
				fmt.Fprintf(os.Stderr, "-%s\n", origLine)
			}
			if formLine != "" {
				fmt.Fprintf(os.Stderr, "+%s\n", formLine)
			}
		}
	}

	return nil
}

// ReadFileOrStdin reads content from a file or stdin
func ReadFileOrStdin(filename string) ([]byte, error) {
	if filename == "" || filename == "-" {
		return os.ReadFile("/dev/stdin")
	}
	return os.ReadFile(filename)
}

// WriteFileOrStdout writes content to a file or stdout
func WriteFileOrStdout(filename, content string) error {
	if filename == "" || filename == "-" {
		fmt.Print(content)
		return nil
	}
	return os.WriteFile(filename, []byte(content), 0644)
}

// PrintSuccess prints a success message with styling
func PrintSuccess(message string) {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	fmt.Fprintln(os.Stderr, style.Render("✅ "+message))
}

// PrintError prints an error message with styling
func PrintError(message string) {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
	fmt.Fprintln(os.Stderr, style.Render("❌ "+message))
}

// PrintWarning prints a warning message with styling
func PrintWarning(message string) {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(true)
	fmt.Fprintln(os.Stderr, style.Render("⚠️  "+message))
}

// PrintInfo prints an info message with styling
func PrintInfo(message string) {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true)
	fmt.Fprintln(os.Stderr, style.Render("ℹ️  "+message))
}

// IsStdinPiped checks if stdin is being piped to
func IsStdinPiped() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) == 0
}

// FileExists checks if a file exists
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// EnsureFileExists ensures a file exists, creating it if necessary
func EnsureFileExists(filename string) error {
	if !FileExists(filename) {
		file, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", filename, err)
		}
		defer file.Close()
	}
	return nil
}
