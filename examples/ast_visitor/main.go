package main

import (
	"fmt"
	"strings"

	"github.com/orangekame3/qasmparser/parser"
)

// PrintVisitor demonstrates the visitor pattern by printing AST node information
type PrintVisitor struct {
	parser.BaseVisitor
	indent int
	output strings.Builder
}

func (v *PrintVisitor) print(msg string) {
	v.output.WriteString(strings.Repeat("  ", v.indent))
	v.output.WriteString(msg)
	v.output.WriteString("\n")
}

func (v *PrintVisitor) VisitProgram(node *parser.Program) interface{} {
	v.print("Program:")
	v.indent++

	if node.Version != nil {
		parser.Walk(v, node.Version)
	}

	v.print(fmt.Sprintf("Statements: %d", len(node.Statements)))
	v.indent++
	for _, stmt := range node.Statements {
		parser.Walk(v, stmt)
	}
	v.indent--

	if len(node.Comments) > 0 {
		v.print(fmt.Sprintf("Comments: %d", len(node.Comments)))
		v.indent++
		for _, comment := range node.Comments {
			parser.Walk(v, &comment)
		}
		v.indent--
	}

	v.indent--
	return nil
}

func (v *PrintVisitor) VisitVersion(node *parser.Version) interface{} {
	v.print(fmt.Sprintf("Version: %s", node.Number))
	return nil
}

func (v *PrintVisitor) VisitQuantumDeclaration(node *parser.QuantumDeclaration) interface{} {
	v.print(fmt.Sprintf("Quantum Declaration: %s %s", node.Type, node.Identifier))
	if node.Size != nil {
		v.indent++
		v.print("Size:")
		v.indent++
		parser.Walk(v, node.Size)
		v.indent--
		v.indent--
	}
	return nil
}

func (v *PrintVisitor) VisitClassicalDeclaration(node *parser.ClassicalDeclaration) interface{} {
	v.print(fmt.Sprintf("Classical Declaration: %s %s", node.Type, node.Identifier))
	if node.Size != nil {
		v.indent++
		v.print("Size:")
		v.indent++
		parser.Walk(v, node.Size)
		v.indent--
		v.indent--
	}
	if node.Initializer != nil {
		v.indent++
		v.print("Initializer:")
		v.indent++
		parser.Walk(v, node.Initializer)
		v.indent--
		v.indent--
	}
	return nil
}

func (v *PrintVisitor) VisitGateCall(node *parser.GateCall) interface{} {
	v.print(fmt.Sprintf("Gate Call: %s", node.Name))

	if len(node.Parameters) > 0 {
		v.indent++
		v.print("Parameters:")
		v.indent++
		for _, param := range node.Parameters {
			parser.Walk(v, param)
		}
		v.indent--
		v.indent--
	}

	if len(node.Qubits) > 0 {
		v.indent++
		v.print("Qubits:")
		v.indent++
		for _, qubit := range node.Qubits {
			parser.Walk(v, qubit)
		}
		v.indent--
		v.indent--
	}

	if len(node.Modifiers) > 0 {
		v.indent++
		v.print("Modifiers:")
		v.indent++
		for _, modifier := range node.Modifiers {
			parser.Walk(v, &modifier)
		}
		v.indent--
		v.indent--
	}

	return nil
}

func (v *PrintVisitor) VisitMeasurement(node *parser.Measurement) interface{} {
	v.print("Measurement:")
	v.indent++
	v.print("Qubit:")
	v.indent++
	parser.Walk(v, node.Qubit)
	v.indent--

	if node.Target != nil {
		v.print("Target:")
		v.indent++
		parser.Walk(v, node.Target)
		v.indent--
	}
	v.indent--
	return nil
}

func (v *PrintVisitor) VisitInclude(node *parser.Include) interface{} {
	v.print(fmt.Sprintf("Include: %s", node.Path))
	return nil
}

func (v *PrintVisitor) VisitIdentifier(node *parser.Identifier) interface{} {
	v.print(fmt.Sprintf("Identifier: %s", node.Name))
	return nil
}

func (v *PrintVisitor) VisitIndexedIdentifier(node *parser.IndexedIdentifier) interface{} {
	v.print(fmt.Sprintf("Indexed Identifier: %s", node.Name))
	v.indent++
	v.print("Index:")
	v.indent++
	parser.Walk(v, node.Index)
	v.indent--
	v.indent--
	return nil
}

func (v *PrintVisitor) VisitIntegerLiteral(node *parser.IntegerLiteral) interface{} {
	v.print(fmt.Sprintf("Integer: %d", node.Value))
	return nil
}

func (v *PrintVisitor) VisitFloatLiteral(node *parser.FloatLiteral) interface{} {
	v.print(fmt.Sprintf("Float: %f", node.Value))
	return nil
}

func (v *PrintVisitor) VisitStringLiteral(node *parser.StringLiteral) interface{} {
	v.print(fmt.Sprintf("String: %q", node.Value))
	return nil
}

func (v *PrintVisitor) VisitBinaryExpression(node *parser.BinaryExpression) interface{} {
	v.print(fmt.Sprintf("Binary Expression: %s", node.Operator))
	v.indent++
	v.print("Left:")
	v.indent++
	parser.Walk(v, node.Left)
	v.indent--
	v.print("Right:")
	v.indent++
	parser.Walk(v, node.Right)
	v.indent--
	v.indent--
	return nil
}

func (v *PrintVisitor) VisitComment(node *parser.Comment) interface{} {
	v.print(fmt.Sprintf("Comment (%s): %s", node.Type, node.Text))
	return nil
}

// StatCountVisitor counts different types of statements
type StatCountVisitor struct {
	parser.BaseVisitor
	GateCallCount    int
	DeclarationCount int
	MeasurementCount int
	IncludeCount     int
}

func (s *StatCountVisitor) VisitGateCall(node *parser.GateCall) interface{} {
	s.GateCallCount++
	return nil
}

func (s *StatCountVisitor) VisitQuantumDeclaration(node *parser.QuantumDeclaration) interface{} {
	s.DeclarationCount++
	return nil
}

func (s *StatCountVisitor) VisitClassicalDeclaration(node *parser.ClassicalDeclaration) interface{} {
	s.DeclarationCount++
	return nil
}

func (s *StatCountVisitor) VisitMeasurement(node *parser.Measurement) interface{} {
	s.MeasurementCount++
	return nil
}

func (s *StatCountVisitor) VisitInclude(node *parser.Include) interface{} {
	s.IncludeCount++
	return nil
}

func main() {
	// Create a simple AST manually for demonstration
	program := &parser.Program{
		BaseNode: parser.BaseNode{
			Position: parser.Position{Line: 1, Column: 1},
		},
		Version: &parser.Version{
			BaseNode: parser.BaseNode{Position: parser.Position{Line: 1, Column: 1}},
			Number:   "3.0",
		},
		Statements: []parser.Statement{
			&parser.Include{
				BaseNode: parser.BaseNode{Position: parser.Position{Line: 2, Column: 1}},
				Path:     "stdgates.qasm",
			},
			&parser.QuantumDeclaration{
				BaseNode:   parser.BaseNode{Position: parser.Position{Line: 4, Column: 1}},
				Type:       "qubit",
				Identifier: "q",
				Size: &parser.IntegerLiteral{
					BaseNode: parser.BaseNode{Position: parser.Position{Line: 4, Column: 7}},
					Value:    2,
				},
			},
			&parser.ClassicalDeclaration{
				BaseNode:   parser.BaseNode{Position: parser.Position{Line: 5, Column: 1}},
				Type:       "bit",
				Identifier: "c",
				Size: &parser.IntegerLiteral{
					BaseNode: parser.BaseNode{Position: parser.Position{Line: 5, Column: 5}},
					Value:    2,
				},
			},
			&parser.GateCall{
				BaseNode: parser.BaseNode{Position: parser.Position{Line: 7, Column: 1}},
				Name:     "h",
				Qubits: []parser.Expression{
					&parser.IndexedIdentifier{
						BaseNode: parser.BaseNode{Position: parser.Position{Line: 7, Column: 3}},
						Name:     "q",
						Index: &parser.IntegerLiteral{
							BaseNode: parser.BaseNode{Position: parser.Position{Line: 7, Column: 5}},
							Value:    0,
						},
					},
				},
			},
			&parser.GateCall{
				BaseNode: parser.BaseNode{Position: parser.Position{Line: 8, Column: 1}},
				Name:     "cx",
				Qubits: []parser.Expression{
					&parser.IndexedIdentifier{
						BaseNode: parser.BaseNode{Position: parser.Position{Line: 8, Column: 4}},
						Name:     "q",
						Index: &parser.IntegerLiteral{
							BaseNode: parser.BaseNode{Position: parser.Position{Line: 8, Column: 6}},
							Value:    0,
						},
					},
					&parser.IndexedIdentifier{
						BaseNode: parser.BaseNode{Position: parser.Position{Line: 8, Column: 10}},
						Name:     "q",
						Index: &parser.IntegerLiteral{
							BaseNode: parser.BaseNode{Position: parser.Position{Line: 8, Column: 12}},
							Value:    1,
						},
					},
				},
			},
			&parser.Measurement{
				BaseNode: parser.BaseNode{Position: parser.Position{Line: 9, Column: 1}},
				Qubit: &parser.Identifier{
					BaseNode: parser.BaseNode{Position: parser.Position{Line: 9, Column: 9}},
					Name:     "q",
				},
				Target: &parser.Identifier{
					BaseNode: parser.BaseNode{Position: parser.Position{Line: 9, Column: 14}},
					Name:     "c",
				},
			},
		},
		Comments: []parser.Comment{
			{
				BaseNode: parser.BaseNode{Position: parser.Position{Line: 3, Column: 1}},
				Text:     "Declare qubits and classical bits",
				Type:     "line",
			},
		},
	}

	fmt.Println("=== AST Visitor Pattern Example ===")
	fmt.Println()

	// Example 1: Pretty printing with PrintVisitor
	fmt.Println("1. Pretty printing AST structure:")
	printVisitor := &PrintVisitor{}
	parser.Walk(printVisitor, program)
	fmt.Print(printVisitor.output.String())
	fmt.Println()

	// Example 2: Counting statements with StatCountVisitor
	fmt.Println("2. Statement statistics:")
	statVisitor := &StatCountVisitor{}

	// Use depth-first visitor to automatically traverse children
	depthFirst := parser.NewDepthFirstVisitor(statVisitor)
	parser.Walk(depthFirst, program)

	fmt.Printf("Gate calls: %d\n", statVisitor.GateCallCount)
	fmt.Printf("Declarations: %d\n", statVisitor.DeclarationCount)
	fmt.Printf("Measurements: %d\n", statVisitor.MeasurementCount)
	fmt.Printf("Includes: %d\n", statVisitor.IncludeCount)
	fmt.Println()

	// Example 3: Walking specific node types
	fmt.Println("3. Walking statements only:")
	for i, stmt := range program.Statements {
		fmt.Printf("Statement %d: %s\n", i+1, stmt.String())
	}
	fmt.Println()

	fmt.Println("Note: This example uses manually created AST nodes.")
	fmt.Println("In real usage, the AST would be generated by parsing QASM source code.")
}
