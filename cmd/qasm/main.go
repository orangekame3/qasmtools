package main

import (
	"fmt"
	"os"

	"github.com/orangekame3/qasmtools/formatter"
	"github.com/orangekame3/qasmtools/highlight"
	"github.com/spf13/cobra"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "qasm",
	Short: "QASM tools",
	Long:  `A collection of tools for working with QASM files.`,
}

func init() {
	// Add fmt subcommand
	fmtCmd := &cobra.Command{
		Use:   "fmt [files...]",
		Short: "Format QASM files",
		Long:  `Format one or more QASM files according to the standard style.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("at least one file is required")
			}
			config, err := getConfigFromFlags(cmd)
			if err != nil {
				return err
			}

			if config.Check {
				return runCheckWithConfig(cmd, args, config)
			}
			return runFormatWithConfig(cmd, args, config)
		},
	}

	// Add flags to fmt command
	fmtCmd.Flags().BoolP("write", "w", false, "write result to (source) file instead of stdout")
	fmtCmd.Flags().Bool("check", false, "check if files are formatted without modifying them")
	fmtCmd.Flags().UintP("indent", "i", 2, "indentation size")
	fmtCmd.Flags().BoolP("newline", "n", true, "ensure files end with a newline")
	fmtCmd.Flags().BoolP("verbose", "v", false, "verbose output")
	fmtCmd.Flags().Bool("diff", false, "display diffs instead of rewriting files")

	// Add highlight subcommand
	highlightCmd := &cobra.Command{
		Use:   "highlight [file]",
		Short: "Highlight QASM file",
		Long:  `Display QASM file with syntax highlighting.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			file, err := cmd.Flags().GetString("file")
			if err != nil {
				return err
			}

			if file == "" && len(args) == 0 {
				return fmt.Errorf("a file is required (use -f flag or provide as argument)")
			}

			if file != "" {
				return runHighlight(file)
			}
			return runHighlight(args[0])
		},
	}

	// Add flags to highlight command
	highlightCmd.Flags().StringP("file", "f", "", "file to highlight")

	rootCmd.AddCommand(fmtCmd, highlightCmd)
}

func getConfigFromFlags(cmd *cobra.Command) (*formatter.Config, error) {
	write, err := cmd.Flags().GetBool("write")
	if err != nil {
		return nil, err
	}

	check, err := cmd.Flags().GetBool("check")
	if err != nil {
		return nil, err
	}

	indent, err := cmd.Flags().GetUint("indent")
	if err != nil {
		return nil, err
	}

	newline, err := cmd.Flags().GetBool("newline")
	if err != nil {
		return nil, err
	}

	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		return nil, err
	}

	diff, err := cmd.Flags().GetBool("diff")
	if err != nil {
		return nil, err
	}

	return &formatter.Config{
		Write:   write,
		Check:   check,
		Indent:  indent,
		Newline: newline,
		Verbose: verbose,
		Diff:    diff,
	}, nil
}

func runFormatWithConfig(cmd *cobra.Command, args []string, config *formatter.Config) error {
	var processedFiles int
	var modifiedFiles int

	for _, file := range args {
		modified, err := formatFileWithConfig(file, config)
		if err != nil {
			if config.Verbose {
				fmt.Printf("❌ %s: %v\n", file, err)
			}
			return fmt.Errorf("failed to format %s: %w", file, err)
		}

		processedFiles++
		if modified {
			modifiedFiles++
			if config.Verbose && config.Write {
				fmt.Printf("✅ %s formatted and saved\n", file)
			}
		} else if config.Verbose {
			fmt.Printf("ℹ️  %s already formatted\n", file)
		}
	}

	if config.Verbose && config.Write {
		fmt.Printf("\n📊 Processed %d files, modified %d files\n", processedFiles, modifiedFiles)
	}

	return nil
}

func runCheckWithConfig(cmd *cobra.Command, args []string, config *formatter.Config) error {
	var hasUnformatted bool
	var checkedFiles int

	for _, file := range args {
		formatted, err := checkFileWithConfig(file, config)
		if err != nil {
			fmt.Printf("❌ %s: %v\n", file, err)
			hasUnformatted = true
			continue
		}

		checkedFiles++
		if !formatted {
			hasUnformatted = true
			fmt.Printf("❌ %s is not formatted\n", file)
		} else {
			fmt.Printf("✅ %s is formatted correctly\n", file)
		}
	}

	if checkedFiles > 0 {
		if hasUnformatted {
			fmt.Printf("\n❌ Some files are not properly formatted\n")
			os.Exit(1)
		} else {
			fmt.Printf("\n✅ All %d files are formatted correctly\n", checkedFiles)
		}
	}

	return nil
}

func formatFileWithConfig(filename string, config *formatter.Config) (bool, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return false, fmt.Errorf("failed to read file: %w", err)
	}

	formatted, err := formatter.FormatQASMWithConfig(string(content), config)
	if err != nil {
		return false, fmt.Errorf("failed to format QASM: %w", err)
	}

	modified := string(content) != formatted

	if config.Write {
		if modified {
			err := os.WriteFile(filename, []byte(formatted), 0600)
			return modified, err
		}
		return false, nil
	}

	if config.Diff {
		err := showDiff(filename, string(content), formatted)
		return modified, err
	}

	fmt.Print(formatted)
	return modified, nil
}

func checkFileWithConfig(filename string, config *formatter.Config) (bool, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return false, fmt.Errorf("failed to read file: %w", err)
	}

	formatted, err := formatter.FormatQASMWithConfig(string(content), config)
	if err != nil {
		return false, fmt.Errorf("failed to format QASM: %w", err)
	}

	return string(content) == formatted, nil
}

func runHighlight(filename string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	h := highlight.New()
	colored, err := h.Highlight(string(content))
	if err != nil {
		return fmt.Errorf("failed to highlight: %v", err)
	}

	fmt.Print(colored)
	return nil
}

func showDiff(filename, original, formatted string) error {
	// TODO: Implement diff functionality
	fmt.Printf("Diff functionality not implemented yet for %s\n", filename)
	return nil
}
