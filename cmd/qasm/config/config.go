package config

import (
	"github.com/orangekame3/qasmtools/formatter"
	"github.com/spf13/cobra"
)

// GetConfigFromFlags extracts formatter configuration from command flags
func GetConfigFromFlags(cmd *cobra.Command) (*formatter.Config, error) {
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

	return &formatter.Config{
		Write:   write,
		Check:   check,
		Indent:  indent,
		Newline: newline,
		Verbose: verbose,
	}, nil
}

// AddCommonFormatFlags adds common formatting flags to a command
func AddCommonFormatFlags(cmd *cobra.Command) {
	cmd.Flags().BoolP("write", "w", false, "Write result to file instead of stdout")
	cmd.Flags().BoolP("check", "c", false, "Check if files are formatted")
	cmd.Flags().UintP("indent", "i", 2, "Number of spaces for indentation")
	cmd.Flags().Bool("newline", true, "Add newline at end of file")
	cmd.Flags().BoolP("verbose", "v", false, "Verbose output")
}

// AddCommonLintFlags adds common linting flags to a command
func AddCommonLintFlags(cmd *cobra.Command) {
	cmd.Flags().StringSlice("disable", []string{}, "Disable specific rules (comma-separated)")
	cmd.Flags().StringSlice("enable-only", []string{}, "Enable only specific rules (comma-separated)")
	cmd.Flags().BoolP("quiet", "q", false, "Only show errors, not warnings")
	cmd.Flags().BoolP("verbose", "v", false, "Verbose output")
}

// AddCommonIOFlags adds common I/O flags to a command
func AddCommonIOFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("stdin", false, "Read from stdin")
}