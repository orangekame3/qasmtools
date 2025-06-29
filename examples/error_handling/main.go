package main

import (
	"fmt"
	"os"

	"github.com/orangekame3/qasmparser/parser"
)

func main() {
	fmt.Println("=== Error Handling Example ===")
	fmt.Println()

	// Example 1: Valid QASM program
	validQasm := `
OPENQASM 3.0;
include "stdgates.qasm";

qubit[2] q;
h q[0];
cx q[0], q[1];
`

	// Example 2: QASM with syntax errors
	invalidQasm := `
OPENQASM 3.0;
include "stdgates.qasm"  // Missing semicolon

qubit[2] q;
invalid_gate q[0];       // Unknown gate
h q[999];                // Index out of bounds (semantic error)
measure q ->             // Incomplete measure statement
`

	// Example 3: QASM with lexer errors
	lexerErrorQasm := `
OPENQASM 3.0;
qubit[2] q;
h q[0];
// Unclosed string literal
gate my_gate(theta) q {
    rx("incomplete string
}
`

	// Create parser with different configurations
	fmt.Println("1. Parsing valid QASM:")
	basicParser := parser.NewParser()

	// Note: This will panic until ANTLR files are generated
	// Demonstrating the error handling interface
	result := basicParser.ParseWithErrors(validQasm)
	handleResult("Valid QASM", result)

	fmt.Println("\n2. Parsing QASM with syntax errors:")
	result = basicParser.ParseWithErrors(invalidQasm)
	handleResult("Invalid QASM", result)

	fmt.Println("\n3. Parsing QASM with lexer errors:")
	result = basicParser.ParseWithErrors(lexerErrorQasm)
	handleResult("Lexer Error QASM", result)

	// Example 4: Parser with custom options
	fmt.Println("\n4. Parser with custom error handling options:")
	strictParser := parser.NewParserWithOptions(&parser.ParseOptions{
		StrictMode:      true,
		ErrorRecovery:   false,
		MaxErrors:       3,
		IncludeComments: false,
	})

	result = strictParser.ParseWithErrors(invalidQasm)
	handleResult("Strict Mode", result)

	// Example 5: Simple validation
	fmt.Println("\n5. Simple validation (returns first error only):")
	if err := basicParser.Validate(validQasm); err != nil {
		fmt.Printf("Validation failed: %s\n", err)
	} else {
		fmt.Println("✓ Valid QASM program")
	}

	if err := basicParser.Validate(invalidQasm); err != nil {
		fmt.Printf("✗ Validation failed: %s\n", err)
	} else {
		fmt.Println("✓ Valid QASM program")
	}

	// Example 6: Error types and custom error creation
	fmt.Println("\n6. Different error types:")

	syntaxErr := parser.NewSyntaxError("unexpected token ';'", parser.Position{Line: 5, Column: 12})
	fmt.Printf("Syntax Error: %s\n", syntaxErr.Error())

	semanticErr := parser.NewSemanticError("undefined gate 'my_gate'", parser.Position{Line: 8, Column: 1})
	fmt.Printf("Semantic Error: %s\n", semanticErr.Error())

	lexerErr := parser.NewLexerError("unterminated string literal", parser.Position{Line: 3, Column: 15})
	fmt.Printf("Lexer Error: %s\n", lexerErr.Error())

	// Example 7: File parsing with error handling
	fmt.Println("\n7. File parsing example:")
	demonstrateFileParsing(basicParser)

	fmt.Println("\nNote: The actual parsing will work once ANTLR files are generated with 'task generate'")
}

func handleResult(title string, result *parser.ParseResult) {
	fmt.Printf("=== %s ===\n", title)

	if result.HasErrors() {
		fmt.Printf("Found %d errors:\n", len(result.Errors))
		for i, err := range result.Errors {
			fmt.Printf("  %d. %s\n", i+1, err.Error())
		}

		// Show error summary by type
		syntaxCount := 0
		semanticCount := 0
		lexerCount := 0

		for _, err := range result.Errors {
			switch err.Type {
			case "syntax":
				syntaxCount++
			case "semantic":
				semanticCount++
			case "lexer":
				lexerCount++
			}
		}

		fmt.Printf("Error summary: %d syntax, %d semantic, %d lexer\n",
			syntaxCount, semanticCount, lexerCount)
	} else {
		fmt.Println("✓ No errors found")
	}

	if result.Program != nil {
		fmt.Printf("Program parsed with %d statements\n", len(result.Program.Statements))
		if result.Program.Version != nil {
			fmt.Printf("OpenQASM version: %s\n", result.Program.Version.Number)
		}
	} else {
		fmt.Println("No AST generated due to errors")
	}
}

func demonstrateFileParsing(p *parser.Parser) {
	// Create a temporary file with invalid QASM
	tempFile := "/tmp/test.qasm"
	content := `
OPENQASM 3.0;
qubit q;
invalid_statement;  // This will cause an error
h q;
`

	err := os.WriteFile(tempFile, []byte(content), 0644)
	if err != nil {
		fmt.Printf("Failed to create temp file: %s\n", err)
		return
	}
	defer os.Remove(tempFile)

	// Parse file
	program, err := p.ParseFile(tempFile)
	if err != nil {
		if parseErr, ok := err.(*parser.ParseError); ok {
			fmt.Printf("Parse error in file: %s\n", parseErr.Error())
		} else {
			fmt.Printf("File error: %s\n", err)
		}
	} else {
		fmt.Printf("File parsed successfully: %d statements\n", len(program.Statements))
	}

	// Alternative: use ParseWithErrors for more detailed error information
	fileContent, err := os.ReadFile(tempFile)
	if err != nil {
		fmt.Printf("Failed to read file: %s\n", err)
		return
	}

	result := p.ParseWithErrors(string(fileContent))
	handleResult("File Parsing", result)
}
