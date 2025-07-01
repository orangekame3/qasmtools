package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/orangekame3/qasmtools/lint"
)

func main() {
	var outputDir string
	flag.StringVar(&outputDir, "output", "docs/rules", "output directory for generated documentation")
	flag.Parse()

	// Create linter to load rules
	linter := lint.NewLinter("")

	// Load rules
	err := linter.LoadRules()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to load rules: %v\n", err)
		os.Exit(1)
	}

	// Get rules from linter
	rules := linter.GetRules()

	// Create documentation generator
	generator := lint.NewDocumentationGenerator(outputDir)

	// Generate documentation
	err = generator.GenerateAllDocumentation(rules)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to generate documentation: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Generated documentation for %d rules in %s\n", len(rules), outputDir)
}
