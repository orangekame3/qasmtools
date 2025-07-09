package commands

import (
	"fmt"
	"io"
	"os"

	"github.com/mattn/go-isatty"
	"github.com/orangekame3/qasmtools/cmd/qasm/config"
	"github.com/orangekame3/qasmtools/cmd/qasm/ioutils"
	"github.com/orangekame3/qasmtools/formatter"
	"github.com/spf13/cobra"
)

// NewFormatCommand creates the format command
func NewFormatCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fmt [files...]",
		Short: "Format QASM files",
		Long:  `Format one or more QASM files according to the standard style.`,
		RunE:  runFormat,
	}

	// Add flags
	cmd.Flags().BoolP("write", "w", false, "Write result to file instead of stdout")
	cmd.Flags().BoolP("check", "c", false, "Check if files are formatted")
	cmd.Flags().Bool("stdin", false, "Read from stdin")
	cmd.Flags().UintP("indent", "i", 2, "Number of spaces for indentation")
	cmd.Flags().Bool("newline", true, "Add newline at end of file")
	cmd.Flags().BoolP("verbose", "v", false, "Verbose output")

	return cmd
}

func runFormat(cmd *cobra.Command, args []string) error {
	stdin, _ := cmd.Flags().GetBool("stdin")
	if !stdin && len(args) == 0 {
		// Check if input is being piped
		if !isatty.IsTerminal(os.Stdin.Fd()) {
			config, err := config.GetConfigFromFlags(cmd)
			if err != nil {
				return err
			}
			return RunFormatStdin(cmd, config)
		}
		return fmt.Errorf("at least one file is required (or use --stdin flag)")
	}

	config, err := config.GetConfigFromFlags(cmd)
	if err != nil {
		return err
	}

	if stdin {
		return RunFormatStdin(cmd, config)
	}

	return runFormatWithConfig(cmd, args, config)
}

func RunFormatStdin(cmd *cobra.Command, config *formatter.Config) error {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("failed to read from stdin: %w", err)
	}

	formatted, err := formatter.FormatQASMWithConfig(string(input), config)
	if err != nil {
		return fmt.Errorf("failed to format QASM: %w", err)
	}

	if config.Check {
		if string(input) != formatted {
			fmt.Fprintln(os.Stderr, "stdin: not formatted")
			os.Exit(1)
		}
		return nil
	}

	fmt.Print(formatted)
	return nil
}

func runFormatWithConfig(cmd *cobra.Command, args []string, config *formatter.Config) error {
	hasError := false
	hasChanges := false

	for _, filename := range args {
		changed, err := FormatFileWithConfig(filename, config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error formatting %s: %v\n", filename, err)
			hasError = true
			continue
		}
		if changed {
			hasChanges = true
		}
	}

	if hasError {
		return fmt.Errorf("formatting failed for one or more files")
	}

	if config.Check && hasChanges {
		os.Exit(1)
	}

	return nil
}

func FormatFileWithConfig(filename string, config *formatter.Config) (bool, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return false, fmt.Errorf("failed to read file: %w", err)
	}

	formatted, err := formatter.FormatQASMWithConfig(string(content), config)
	if err != nil {
		return false, fmt.Errorf("failed to format QASM: %w", err)
	}

	original := string(content)
	changed := original != formatted

	if config.Check {
		if changed {
			fmt.Fprintf(os.Stderr, "%s: not formatted\n", filename)
		} else if config.Verbose {
			fmt.Fprintf(os.Stderr, "%s: already formatted\n", filename)
		}
		return changed, nil
	}

	if config.Write {
		if changed {
			err = os.WriteFile(filename, []byte(formatted), 0644)
			if err != nil {
				return false, fmt.Errorf("failed to write file: %w", err)
			}
			if config.Verbose {
				fmt.Fprintf(os.Stderr, "Formatted %s\n", filename)
			}
		} else if config.Verbose {
			fmt.Fprintf(os.Stderr, "%s: no changes\n", filename)
		}
	} else {
		if changed && config.Verbose {
			err = ioutils.ShowDiff(filename, original, formatted)
			if err != nil {
				return false, err
			}
		}
		fmt.Print(formatted)
	}

	return changed, nil
}

// CheckFileWithConfig checks if a file is properly formatted
func CheckFileWithConfig(filename string, config *formatter.Config) (bool, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return false, fmt.Errorf("failed to read file: %w", err)
	}

	formatted, err := formatter.FormatQASMWithConfig(string(content), config)
	if err != nil {
		return false, fmt.Errorf("failed to format QASM: %w", err)
	}

	// Returns true if already formatted (no changes needed)
	return string(content) == formatted, nil
}
