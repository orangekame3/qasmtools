package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"time"

	"github.com/mattn/go-isatty"
	"github.com/orangekame3/qasmtools/parser"
	"github.com/spf13/cobra"
)

// NewParseCommand creates the parse command
func NewParseCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "parse [files...]",
		Short: "Parse QASM files and output AST",
		Long:  `Parse QASM files and output the Abstract Syntax Tree (AST) in various formats.`,
		RunE:  runParse,
	}

	// Add flags
	cmd.Flags().Bool("stdin", false, "Read from stdin")
	cmd.Flags().String("format", "json", "Output format (json, tree, summary)")
	cmd.Flags().Bool("pretty", true, "Pretty print JSON output")
	cmd.Flags().Bool("include-positions", false, "Include line and column positions in output")
	cmd.Flags().BoolP("verbose", "v", false, "Verbose output with parsing details")
	cmd.Flags().Bool("benchmark", false, "Show parsing performance metrics")
	cmd.Flags().Int("repeat", 1, "Number of times to repeat parsing (for benchmarking)")

	return cmd
}

func runParse(cmd *cobra.Command, args []string) error {
	stdin, _ := cmd.Flags().GetBool("stdin")
	if !stdin && len(args) == 0 {
		// Check if input is being piped
		if !isatty.IsTerminal(os.Stdin.Fd()) {
			return runParseStdin(cmd)
		}
		return fmt.Errorf("at least one file is required (or use --stdin flag)")
	}

	if stdin {
		return runParseStdin(cmd)
	}

	return runParseFiles(cmd, args)
}

func runParseStdin(cmd *cobra.Command) error {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("failed to read from stdin: %w", err)
	}

	return parseAndOutput(cmd, string(input), "<stdin>")
}

func runParseFiles(cmd *cobra.Command, files []string) error {
	for _, filename := range files {
		content, err := os.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", filename, err)
		}

		if len(files) > 1 {
			fmt.Printf("=== %s ===\n", filename)
		}

		err = parseAndOutput(cmd, string(content), filename)
		if err != nil {
			return err
		}

		if len(files) > 1 {
			fmt.Println()
		}
	}

	return nil
}

func parseAndOutput(cmd *cobra.Command, content, filename string) error {
	format, _ := cmd.Flags().GetString("format")
	pretty, _ := cmd.Flags().GetBool("pretty")
	includePositions, _ := cmd.Flags().GetBool("include-positions")
	verbose, _ := cmd.Flags().GetBool("verbose")
	benchmark, _ := cmd.Flags().GetBool("benchmark")
	repeat, _ := cmd.Flags().GetInt("repeat")

	// Performance measurement
	var program *parser.Program
	var err error
	var totalDuration time.Duration
	var minDuration, maxDuration time.Duration
	var durations []time.Duration

	if benchmark || repeat > 1 {
		// Multiple parsing runs for benchmarking
		for i := 0; i < repeat; i++ {
			p := parser.NewParser()
			start := time.Now()
			prog, parseErr := p.ParseString(content)
			duration := time.Since(start)

			if parseErr != nil {
				return fmt.Errorf("failed to parse QASM (run %d): %w", i+1, parseErr)
			}

			if i == 0 {
				program = prog
				err = parseErr
				minDuration = duration
				maxDuration = duration
			} else {
				if duration < minDuration {
					minDuration = duration
				}
				if duration > maxDuration {
					maxDuration = duration
				}
			}

			totalDuration += duration
			durations = append(durations, duration)
		}

		// Show performance stats
		avgDuration := totalDuration / time.Duration(repeat)

		fmt.Fprintf(os.Stderr, "\n⚡ Parse Performance Results for %s:\n", filename)
		fmt.Fprintf(os.Stderr, "   File size: %d bytes\n", len(content))
		if program != nil {
			fmt.Fprintf(os.Stderr, "   Statements: %d\n", len(program.Statements))
		}
		fmt.Fprintf(os.Stderr, "   Runs: %d\n", repeat)
		fmt.Fprintf(os.Stderr, "   Total time: %v\n", totalDuration)
		fmt.Fprintf(os.Stderr, "   Average time: %v\n", avgDuration)
		fmt.Fprintf(os.Stderr, "   Min time: %v\n", minDuration)
		fmt.Fprintf(os.Stderr, "   Max time: %v\n", maxDuration)

		// Calculate throughput
		throughput := float64(len(content)) / avgDuration.Seconds()
		fmt.Fprintf(os.Stderr, "   Throughput: %.2f bytes/sec\n", throughput)

		if repeat > 1 {
			// Calculate standard deviation
			var variance float64
			for _, d := range durations {
				diff := d.Seconds() - avgDuration.Seconds()
				variance += diff * diff
			}
			variance /= float64(repeat)
			stdDev := time.Duration(math.Sqrt(variance) * float64(time.Second))
			fmt.Fprintf(os.Stderr, "   Std deviation: %v\n", stdDev)
		}
		fmt.Fprintf(os.Stderr, "\n")
	} else {
		// Single parse
		p := parser.NewParser()
		start := time.Now()
		program, err = p.ParseString(content)
		duration := time.Since(start)

		if err != nil {
			return fmt.Errorf("failed to parse QASM: %w", err)
		}

		if verbose {
			fmt.Fprintf(os.Stderr, "Parsing completed successfully for %s in %v\n", filename, duration)
			if program != nil && len(program.Statements) > 0 {
				fmt.Fprintf(os.Stderr, "Found %d statements\n", len(program.Statements))
			}
		}
	}

	switch format {
	case "json":
		return outputJSONForParse(program, pretty, includePositions)
	case "tree":
		return outputTree(program, includePositions)
	case "summary":
		return outputSummary(program, filename)
	default:
		return fmt.Errorf("unsupported format: %s (supported: json, tree, summary)", format)
	}
}

func outputJSONForParse(program *parser.Program, pretty, includePositions bool) error {
	encoder := json.NewEncoder(os.Stdout)
	if pretty {
		encoder.SetIndent("", "  ")
	}

	// If includePositions is false, we might want to create a simplified version
	// For now, we'll output the full program structure
	return encoder.Encode(program)
}

func outputTree(program *parser.Program, includePositions bool) error {
	if program == nil {
		fmt.Println("(empty program)")
		return nil
	}

	fmt.Printf("Program\n")
	if program.Version != nil {
		if includePositions {
			fmt.Printf("├── Version: %s (line: %d, col: %d)\n",
				program.Version.Number, program.Version.Position.Line, program.Version.Position.Column)
		} else {
			fmt.Printf("├── Version: %s\n", program.Version.Number)
		}
	}

	fmt.Printf("└── Statements (%d)\n", len(program.Statements))
	for i, stmt := range program.Statements {
		isLast := i == len(program.Statements)-1
		prefix := "    ├── "
		if isLast {
			prefix = "    └── "
		}

		stmtType := fmt.Sprintf("%T", stmt)
		if includePositions {
			fmt.Printf("%s%s (line: %d, col: %d)\n",
				prefix, stmtType, stmt.Pos().Line, stmt.Pos().Column)
		} else {
			fmt.Printf("%s%s\n", prefix, stmtType)
		}

		// Add more detailed output for specific statement types
		switch s := stmt.(type) {
		case *parser.QuantumDeclaration:
			subPrefix := "        "
			if isLast {
				subPrefix = "        "
			}
			fmt.Printf("%s├── Type: %s\n", subPrefix, s.Type)
			fmt.Printf("%s├── Identifier: %s\n", subPrefix, s.Identifier)
			if s.Size != nil {
				fmt.Printf("%s└── Size: %v\n", subPrefix, s.Size)
			} else {
				fmt.Printf("%s└── Size: (single)\n", subPrefix)
			}
		case *parser.ClassicalDeclaration:
			subPrefix := "        "
			if isLast {
				subPrefix = "        "
			}
			fmt.Printf("%s├── Type: %s\n", subPrefix, s.Type)
			fmt.Printf("%s└── Identifier: %s\n", subPrefix, s.Identifier)
		case *parser.GateCall:
			subPrefix := "        "
			if isLast {
				subPrefix = "        "
			}
			fmt.Printf("%s├── Gate: %s\n", subPrefix, s.Name)
			fmt.Printf("%s└── Qubits: %d\n", subPrefix, len(s.Qubits))
		case *parser.Measurement:
			subPrefix := "        "
			if isLast {
				subPrefix = "        "
			}
			fmt.Printf("%s├── Qubit: %v\n", subPrefix, s.Qubit)
			fmt.Printf("%s└── Target: %v\n", subPrefix, s.Target)
		}
	}

	return nil
}

func outputSummary(program *parser.Program, filename string) error {
	if program == nil {
		fmt.Printf("File: %s - No valid program found\n", filename)
		return nil
	}

	fmt.Printf("=== QASM Parse Summary ===\n")
	fmt.Printf("File: %s\n", filename)

	if program.Version != nil {
		fmt.Printf("Version: %s\n", program.Version.Number)
	}

	fmt.Printf("Total Statements: %d\n", len(program.Statements))

	// Count different statement types
	counts := make(map[string]int)
	for _, stmt := range program.Statements {
		// Simplify type names
		switch stmt.(type) {
		case *parser.QuantumDeclaration:
			counts["Qubit Declarations"]++
		case *parser.ClassicalDeclaration:
			counts["Classical Declarations"]++
		case *parser.GateCall:
			counts["Gate Calls"]++
		case *parser.Measurement:
			counts["Measurements"]++
		case *parser.Include:
			counts["Include Statements"]++
		default:
			counts["Other"]++
		}
	}

	fmt.Println("\nStatement Breakdown:")
	for stmtType, count := range counts {
		fmt.Printf("  %s: %d\n", stmtType, count)
	}

	return nil
}
