package main

import (
	"fmt"

	"github.com/orangekame3/qasmparser/parser"
)

func main() {
	// Simple OpenQASM 3.0 program
	qasm := `
OPENQASM 3.0;
include "stdgates.qasm";

qubit[2] q;
bit[2] c;

h q[0];
cx q[0], q[1];
measure q -> c;
`

	// Create parser with default options
	p := parser.NewParser()

	// Parse the QASM code
	fmt.Println("Parsing OpenQASM 3.0 program...")

	// Note: This will panic until ANTLR files are generated
	// Run 'task generate' first to create the ANTLR parser files
	result := p.ParseWithErrors(qasm)

	if result.HasErrors() {
		fmt.Printf("Found %d errors:\n", len(result.Errors))
		for _, err := range result.Errors {
			fmt.Printf("  %s\n", err.Error())
		}
	} else {
		fmt.Println("Parsing successful!")
	}

	if result.Program != nil {
		fmt.Printf("Parsed program with %d statements\n", len(result.Program.Statements))
		if result.Program.Version != nil {
			fmt.Printf("OpenQASM version: %s\n", result.Program.Version.Number)
		}
	}

	// Example of basic validation
	fmt.Println("\nValidating syntax...")
	if err := p.Validate(qasm); err != nil {
		fmt.Printf("Validation failed: %s\n", err)
	} else {
		fmt.Println("Validation successful!")
	}
}
