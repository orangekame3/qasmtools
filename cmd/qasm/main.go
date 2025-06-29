package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/orangekame3/qasmtools/formatter"
	"github.com/orangekame3/qasmtools/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: qasm <command> [arguments]")
		fmt.Println("Commands:")
		fmt.Println("  fmt <input_file> [-o <output_file>] [-i <indent_size>] [-n <newline>]")
		fmt.Println("  parse <input_file>")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "fmt":
		handleFormatCommand(os.Args[2:])
	case "parse":
		handleParseCommand(os.Args[2:])
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func handleFormatCommand(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: qasm format <input_file> [-o <output_file>] [-i <indent_size>] [-n <newline>]")
		os.Exit(1)
	}
	inputFile := args[0]
	outputFile := ""
	config := &formatter.Config{
		Indent:  2,
		Newline: true,
	}

	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "-o", "--output_file":
			if i+1 < len(args) {
				outputFile = args[i+1]
				i++
			}
		case "-i", "--indent":
			if i+1 < len(args) {
				indent, err := strconv.ParseUint(args[i+1], 10, 32)
				if err != nil {
					fmt.Printf("Invalid indent size: %v\n", err)
					os.Exit(1)
				}
				config.Indent = uint(indent)
				i++
			}
		case "-n", "--newline":
			config.Newline = true
		}
	}

	content, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		os.Exit(1)
	}

	formattedContent, err := formatter.FormatQASMWithConfig(string(content), config)
	if err != nil {
		fmt.Printf("Error formatting QASM: %v\n", err)
		os.Exit(1)
	}

	if outputFile != "" {
		err = ioutil.WriteFile(outputFile, []byte(formattedContent), 0644)
		if err != nil {
			fmt.Printf("Error writing output file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Formatted content written to %s\n", outputFile)
	} else {
		fmt.Print(formattedContent)
	}
}

func handleParseCommand(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: qasm parse <input_file>")
		os.Exit(1)
	}
	inputFile := args[0]

	content, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		os.Exit(1)
	}

	p := parser.NewParser()
	program, err := p.ParseString(string(content))
	if err != nil {
		fmt.Printf("Error parsing QASM: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully parsed %s. Program: %+v\n", inputFile, program)
}