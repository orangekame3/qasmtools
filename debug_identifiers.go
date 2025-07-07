package main

import (
	"fmt"
	"github.com/orangekame3/qasmtools/lint"
)

func main() {
	// Debug FindIdentifierDeclarations
	lines := []string{
		"qubit valid_qubit;",
		"bit[2] valid_bit;", 
		"gate valid_gate q { h q; }",
	}
	
	for _, line := range lines {
		declarations := lint.FindIdentifierDeclarations(line)
		fmt.Printf("Line `%s` found %d declarations:
", line, len(declarations))
		for _, d := range declarations {
			fmt.Printf("  - `%s` at column %d
", d.Name, d.Column)
		}
	}
}
