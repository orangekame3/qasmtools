package main

import (
	"fmt"
	"os"

	"github.com/orangekame3/qasmtools/formatter"
	"github.com/orangekame3/qasmtools/parser"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "qasm",
		Short: "QASM tools for formatting and parsing QASM files",
		Long:  `A collection of tools for working with QASM files, including formatting and parsing capabilities.`,
	}

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

	// Add parse subcommand
	parseCmd := &cobra.Command{
		Use:   "parse [file]",
		Short: "Parse QASM files",
		Long:  `Parse a QASM file and validate its syntax.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("a file is required")
			}
			return handleParseCommand(args[0])
		},
	}

	rootCmd.AddCommand(fmtCmd, parseCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
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

func handleParseCommand(filename string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading input file: %w", err)
	}

	p := parser.NewParser()
	program, err := p.ParseString(string(content))
	if err != nil {
		return fmt.Errorf("error parsing QASM: %w", err)
	}

	fmt.Printf("Successfully parsed %s. Program: %+v\n", filename, program)
	return nil
}

func showDiff(filename, original, formatted string) error {
	// TODO: Implement diff functionality
	fmt.Printf("Diff functionality not implemented yet for %s\n", filename)
	return nil
}
