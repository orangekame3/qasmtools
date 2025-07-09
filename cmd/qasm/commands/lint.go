package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/orangekame3/qasmtools/lint"
	"github.com/spf13/cobra"
)

// NewLintCommand creates the lint command
func NewLintCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lint [files...]",
		Short: "Lint QASM files for errors and style issues",
		Long:  `Analyze QASM files for potential errors, style violations, and best practice issues.`,
		RunE:  runLint,
	}

	// Add flags
	cmd.Flags().String("rules", "", "Rules directory")
	cmd.Flags().StringSlice("disable", []string{}, "Disable specific rules (comma-separated)")
	cmd.Flags().StringSlice("enable-only", []string{}, "Enable only specific rules (comma-separated)")
	cmd.Flags().String("format", "text", "Output format (text, json)")
	cmd.Flags().BoolP("quiet", "q", false, "Only show errors, not warnings")
	cmd.Flags().Bool("no-color", false, "Disable colored output")
	cmd.Flags().BoolP("verbose", "v", false, "Verbose output")
	cmd.Flags().Bool("use-ast", true, "Use AST-based analysis (faster and more accurate)")
	cmd.Flags().Bool("parallel", true, "Enable parallel processing for multiple files")
	cmd.Flags().Int("workers", 4, "Number of worker threads for parallel processing")
	cmd.Flags().Bool("performance", false, "Show performance statistics")
	cmd.Flags().Bool("stdin", false, "Read from stdin")

	return cmd
}

func runLint(cmd *cobra.Command, args []string) error {
	stdin, _ := cmd.Flags().GetBool("stdin")

	// Check if we should read from stdin
	if stdin {
		return runLintStdin(cmd)
	}

	if len(args) == 0 {
		return fmt.Errorf("at least one file is required")
	}

	rulesDir, _ := cmd.Flags().GetString("rules")
	disabled, _ := cmd.Flags().GetStringSlice("disable")
	enabledOnly, _ := cmd.Flags().GetStringSlice("enable-only")
	format, _ := cmd.Flags().GetString("format")
	quiet, _ := cmd.Flags().GetBool("quiet")
	noColor, _ := cmd.Flags().GetBool("no-color")
	useAST, _ := cmd.Flags().GetBool("use-ast")
	parallel, _ := cmd.Flags().GetBool("parallel")
	workers, _ := cmd.Flags().GetInt("workers")
	showPerf, _ := cmd.Flags().GetBool("performance")

	// Create optimized linter based on configuration
	var violations []*lint.Violation
	var err error

	if parallel && len(args) > 1 {
		// Use batch linter for multiple files
		batchLinter := lint.NewBatchLinter(rulesDir, workers)
		err = batchLinter.LoadRules()
		if err != nil {
			return fmt.Errorf("failed to load rules: %w", err)
		}
		violations, err = batchLinter.LintFilesParallel(args)
		if showPerf {
			printPerformanceStats(batchLinter.GetStats())
		}
	} else {
		// Use standard linter
		linter := lint.NewLinterWithAST(rulesDir, useAST)
		err = linter.LoadRules()
		if err != nil {
			return fmt.Errorf("failed to load rules: %w", err)
		}
		violations, err = linter.LintFiles(args)
	}

	if err != nil {
		return fmt.Errorf("failed to lint files: %w", err)
	}

	// Filter violations based on flags
	filteredViolations := filterViolations(violations, disabled, enabledOnly, quiet)

	// Output results
	switch format {
	case "json":
		return outputJSON(filteredViolations)
	default:
		return outputTextWithColor(filteredViolations, !noColor)
	}
}

func filterViolations(violations []*lint.Violation, disabled []string, enabledOnly []string, quiet bool) []*lint.Violation {
	var filtered []*lint.Violation

	disabledMap := make(map[string]bool)
	for _, rule := range disabled {
		disabledMap[rule] = true
	}

	enabledMap := make(map[string]bool)
	for _, rule := range enabledOnly {
		enabledMap[rule] = true
	}

	for _, violation := range violations {
		// Skip if rule is disabled
		if disabledMap[violation.Rule.ID] {
			continue
		}

		// Skip if only specific rules are enabled and this isn't one of them
		if len(enabledOnly) > 0 && !enabledMap[violation.Rule.ID] {
			continue
		}

		// Skip warnings if quiet mode is enabled
		if quiet && violation.Severity != lint.SeverityError {
			continue
		}

		filtered = append(filtered, violation)
	}

	return filtered
}

// outputJSON outputs violations in JSON format
func outputJSON(violations []*lint.Violation) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(violations)
}

// outputTextWithColor outputs violations with colored text
func outputTextWithColor(violations []*lint.Violation, useColor bool) error {
	if len(violations) == 0 {
		if useColor {
			style := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
			fmt.Println(style.Render("✅ No issues found"))
		} else {
			fmt.Println("✅ No issues found")
		}
		return nil
	}

	// Define styles for colored output
	var fileStyle, errorStyle, warningStyle, infoStyle, ruleStyle, urlStyle lipgloss.Style
	if useColor {
		fileStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true)    // Blue
		errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)    // Red
		warningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(true) // Yellow
		infoStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Bold(true)    // Cyan
		ruleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Bold(true)    // Magenta
		urlStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Italic(true)    // Gray
	} else {
		fileStyle = lipgloss.NewStyle()
		errorStyle = lipgloss.NewStyle()
		warningStyle = lipgloss.NewStyle()
		infoStyle = lipgloss.NewStyle()
		ruleStyle = lipgloss.NewStyle()
		urlStyle = lipgloss.NewStyle()
	}

	for _, violation := range violations {
		// Format colored output
		filePart := fileStyle.Render(fmt.Sprintf("%s:%d:%d:", violation.File, violation.Line, violation.Column))

		var severityPart string
		switch violation.Severity {
		case lint.SeverityError:
			severityPart = errorStyle.Render(string(violation.Severity))
		case lint.SeverityWarning:
			severityPart = warningStyle.Render(string(violation.Severity))
		case lint.SeverityInfo:
			severityPart = infoStyle.Render(string(violation.Severity))
		}

		rulePart := ruleStyle.Render(fmt.Sprintf("[%s]", violation.Rule.ID))

		var result string
		if violation.Rule.DocumentationURL != "" {
			urlPart := urlStyle.Render(fmt.Sprintf("(%s)", violation.Rule.DocumentationURL))
			result = fmt.Sprintf("%s %s %s %s %s", filePart, severityPart, rulePart, violation.Message, urlPart)
		} else {
			result = fmt.Sprintf("%s %s %s %s", filePart, severityPart, rulePart, violation.Message)
		}

		fmt.Println(result)
	}

	// Summary
	errorCount := 0
	warningCount := 0
	infoCount := 0

	for _, violation := range violations {
		switch violation.Severity {
		case lint.SeverityError:
			errorCount++
		case lint.SeverityWarning:
			warningCount++
		case lint.SeverityInfo:
			infoCount++
		}
	}

	fmt.Printf("\n📊 Found %d issues: %d errors, %d warnings, %d info\n",
		len(violations), errorCount, warningCount, infoCount)

	return nil
}

// printPerformanceStats outputs performance statistics
func printPerformanceStats(stats lint.PerformanceStats) {
	fmt.Printf("\n⚡ Performance Stats:\n")
	fmt.Printf("   Files processed: %d\n", stats.TotalFiles)
	fmt.Printf("   Total time: %v\n", stats.TotalTime)
	fmt.Printf("   Parse time: %v\n", stats.ParseTime)
	fmt.Printf("   Analysis time: %v\n", stats.AnalysisTime)
	fmt.Printf("   AST rules used: %d\n", stats.ASTRulesUsed)
	fmt.Printf("   Text rules used: %d\n", stats.TextRulesUsed)
	if stats.TotalFiles > 0 {
		avgTime := stats.TotalTime / time.Duration(stats.TotalFiles)
		fmt.Printf("   Average time per file: %v\n", avgTime)
	}
}

// runLintStdin handles linting from stdin
func runLintStdin(cmd *cobra.Command) error {
	// Read from stdin
	content, err := io.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("failed to read from stdin: %w", err)
	}

	// Get flags
	rulesDir, _ := cmd.Flags().GetString("rules")
	disabled, _ := cmd.Flags().GetStringSlice("disable")
	enabledOnly, _ := cmd.Flags().GetStringSlice("enable-only")
	format, _ := cmd.Flags().GetString("format")
	quiet, _ := cmd.Flags().GetBool("quiet")
	noColor, _ := cmd.Flags().GetBool("no-color")
	useAST, _ := cmd.Flags().GetBool("use-ast")

	// Create linter
	linter := lint.NewLinterWithAST(rulesDir, useAST)
	err = linter.LoadRules()
	if err != nil {
		return fmt.Errorf("failed to load rules: %w", err)
	}

	// Lint content
	violations, err := linter.LintContent(string(content), "<stdin>")
	if err != nil {
		return fmt.Errorf("failed to lint content: %w", err)
	}

	// Filter violations
	filteredViolations := filterViolations(violations, disabled, enabledOnly, quiet)

	// Output results
	switch format {
	case "json":
		return outputJSON(filteredViolations)
	default:
		return outputTextWithColor(filteredViolations, !noColor)
	}
}
