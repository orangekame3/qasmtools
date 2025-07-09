package commands

import (
	"fmt"
	"os"

	"github.com/mattn/go-isatty"
	"github.com/orangekame3/qasmtools/highlight"
	"github.com/spf13/cobra"
)

// NewHighlightCommand creates the highlight command
func NewHighlightCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "highlight [file]",
		Short: "Highlight QASM files with syntax coloring",
		Long:  `Apply syntax highlighting to QASM files for better readability.`,
		RunE:  runHighlight,
	}

	// Add flags
	cmd.Flags().Bool("stdin", false, "Read from stdin")
	cmd.Flags().Bool("semantic", false, "Use semantic highlighting (slower but more accurate)")

	return cmd
}

func runHighlight(cmd *cobra.Command, args []string) error {
	stdin, _ := cmd.Flags().GetBool("stdin")
	semantic, _ := cmd.Flags().GetBool("semantic")

	if !stdin && len(args) == 0 {
		// Check if input is being piped
		if !isatty.IsTerminal(os.Stdin.Fd()) {
			return runHighlightStdin(semantic)
		}
		return fmt.Errorf("at least one file is required (or use --stdin flag)")
	}

	if stdin {
		return runHighlightStdin(semantic)
	}

	if len(args) != 1 {
		return fmt.Errorf("exactly one file is required")
	}

	return runHighlightFile(args[0], semantic)
}

func runHighlightStdin(semantic bool) error {
	content, err := os.ReadFile("/dev/stdin")
	if err != nil {
		return fmt.Errorf("failed to read from stdin: %w", err)
	}

	var highlighted string
	if semantic {
		highlighter := highlight.NewASTHighlighter()
		highlighted, err = highlighter.HighlightWithAST(string(content))
	} else {
		highlighter := highlight.New()
		highlighted, err = highlighter.Highlight(string(content))
	}

	if err != nil {
		return fmt.Errorf("failed to highlight QASM: %w", err)
	}

	fmt.Print(highlighted)
	return nil
}

func runHighlightFile(filename string, semantic bool) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	var highlighted string
	if semantic {
		highlighter := highlight.NewASTHighlighter()
		highlighted, err = highlighter.HighlightWithAST(string(content))
	} else {
		highlighter := highlight.New()
		highlighted, err = highlighter.Highlight(string(content))
	}

	if err != nil {
		return fmt.Errorf("failed to highlight QASM: %w", err)
	}

	fmt.Print(highlighted)
	return nil
}
