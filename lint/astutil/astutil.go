// Package astutil provides utility functions for working with QASM AST structures
package astutil

import (
	"strings"

	"github.com/orangekame3/qasmtools/parser"
)

// VisitAllNodes traverses all nodes in the AST and calls the visitor function for each
func VisitAllNodes(node parser.Node, visitor func(parser.Node)) {
	if node == nil {
		return
	}

	visitor(node)

	switch n := node.(type) {
	case *parser.Program:
		if n.Version != nil {
			VisitAllNodes(n.Version, visitor)
		}
		for _, stmt := range n.Statements {
			if stmt != nil {
				VisitAllNodes(stmt, visitor)
			}
		}
		for _, comment := range n.Comments {
			VisitAllNodes(&comment, visitor)
		}

	case *parser.QuantumDeclaration:
		// Check if n is a valid pointer before accessing fields
		if n == nil {
			return
		}
		if n.Size != nil {
			VisitAllNodes(n.Size, visitor)
		}

	case *parser.ClassicalDeclaration:
		if n == nil {
			return
		}
		if n.Size != nil {
			VisitAllNodes(n.Size, visitor)
		}
		if n.Initializer != nil {
			VisitAllNodes(n.Initializer, visitor)
		}

	case *parser.GateCall:
		if n == nil {
			return
		}
		for _, param := range n.Parameters {
			VisitAllNodes(param, visitor)
		}
		for _, qubit := range n.Qubits {
			VisitAllNodes(qubit, visitor)
		}
		for _, modifier := range n.Modifiers {
			VisitAllNodes(&modifier, visitor)
		}

	case *parser.Measurement:
		if n == nil {
			return
		}
		if n.Qubit != nil {
			VisitAllNodes(n.Qubit, visitor)
		}
		if n.Target != nil {
			VisitAllNodes(n.Target, visitor)
		}

	case *parser.GateDefinition:
		for _, param := range n.Parameters {
			VisitAllNodes(&param, visitor)
		}
		for _, qubit := range n.Qubits {
			VisitAllNodes(&qubit, visitor)
		}
		for _, stmt := range n.Body {
			VisitAllNodes(stmt, visitor)
		}

	case *parser.IfStatement:
		VisitAllNodes(n.Condition, visitor)
		for _, stmt := range n.ThenBody {
			VisitAllNodes(stmt, visitor)
		}
		for _, stmt := range n.ElseBody {
			VisitAllNodes(stmt, visitor)
		}

	case *parser.ForStatement:
		VisitAllNodes(n.Iterable, visitor)
		for _, stmt := range n.Body {
			VisitAllNodes(stmt, visitor)
		}

	case *parser.WhileStatement:
		VisitAllNodes(n.Condition, visitor)
		for _, stmt := range n.Body {
			VisitAllNodes(stmt, visitor)
		}

	case *parser.IndexedIdentifier:
		VisitAllNodes(n.Index, visitor)

	case *parser.RangedIdentifier:
		VisitAllNodes(n.Start, visitor)
		VisitAllNodes(n.EndIndex, visitor)

	case *parser.BinaryExpression:
		VisitAllNodes(n.Left, visitor)
		VisitAllNodes(n.Right, visitor)

	case *parser.UnaryExpression:
		VisitAllNodes(n.Operand, visitor)

	case *parser.FunctionCall:
		for _, arg := range n.Arguments {
			VisitAllNodes(arg, visitor)
		}

	case *parser.ParenthesizedExpression:
		VisitAllNodes(n.Expression, visitor)

	case *parser.TimingExpression:
		VisitAllNodes(n.Value, visitor)

	case *parser.DelayExpression:
		VisitAllNodes(n.Timing, visitor)

	case *parser.Modifier:
		for _, param := range n.Parameters {
			VisitAllNodes(param, visitor)
		}

	// For leaf nodes (literals, identifiers, etc.), no further traversal needed
	}
}

// FindNodesByType finds all nodes of a specific type in the AST
func FindNodesByType[T parser.Node](program *parser.Program, nodeType T) []T {
	var results []T

	VisitAllNodes(program, func(node parser.Node) {
		if typedNode, ok := node.(T); ok {
			results = append(results, typedNode)
		}
	})

	return results
}

// FindDeclarations finds all declarations of a specific type
func FindDeclarations(program *parser.Program) *Declarations {
	declarations := &Declarations{
		Quantum:   make([]*parser.QuantumDeclaration, 0),
		Classical: make([]*parser.ClassicalDeclaration, 0),
		Gates:     make([]*parser.GateDefinition, 0),
	}

	VisitAllNodes(program, func(node parser.Node) {
		switch n := node.(type) {
		case *parser.QuantumDeclaration:
			if n != nil {
				declarations.Quantum = append(declarations.Quantum, n)
			}
		case *parser.ClassicalDeclaration:
			if n != nil {
				declarations.Classical = append(declarations.Classical, n)
			}
		case *parser.GateDefinition:
			if n != nil {
				declarations.Gates = append(declarations.Gates, n)
			}
		}
	})

	return declarations
}

// Declarations holds categorized declarations
type Declarations struct {
	Quantum   []*parser.QuantumDeclaration
	Classical []*parser.ClassicalDeclaration
	Gates     []*parser.GateDefinition
}

// GetUsages finds all usages of a given identifier in the program
func GetUsages(program *parser.Program, identifierName string) []parser.Node {
	var usages []parser.Node

	VisitAllNodes(program, func(node parser.Node) {
		switch n := node.(type) {
		case *parser.Identifier:
			if n.Name == identifierName {
				usages = append(usages, node)
			}
		case *parser.IndexedIdentifier:
			if n.Name == identifierName {
				usages = append(usages, node)
			}
		case *parser.RangedIdentifier:
			if n.Name == identifierName {
				usages = append(usages, node)
			}
		case *parser.GateCall:
			// Check if gate name matches
			if n.Name == identifierName {
				usages = append(usages, node)
			}
			// Check qubits in gate calls
			for _, qubit := range n.Qubits {
				if id, ok := qubit.(*parser.Identifier); ok && id.Name == identifierName {
					usages = append(usages, node)
				}
				if idx, ok := qubit.(*parser.IndexedIdentifier); ok && idx.Name == identifierName {
					usages = append(usages, node)
				}
			}
		case *parser.Measurement:
			// Check qubit and target in measurements
			if id, ok := n.Qubit.(*parser.Identifier); ok && id.Name == identifierName {
				usages = append(usages, node)
			}
			if id, ok := n.Target.(*parser.Identifier); ok && id.Name == identifierName {
				usages = append(usages, node)
			}
		}
	})

	return usages
}

// ExtractBaseName extracts the base identifier name, removing array access
func ExtractBaseName(name string) string {
	if idx := strings.Index(name, "["); idx != -1 {
		return name[:idx]
	}
	return name
}

// IsSnakeCase checks if an identifier follows snake_case naming convention
func IsSnakeCase(name string) bool {
	if name == "" {
		return false
	}

	// Must start with lowercase letter
	if name[0] < 'a' || name[0] > 'z' {
		return false
	}

	// Can only contain lowercase letters, digits, and underscores
	for i, char := range name {
		if !((char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '_') {
			return false
		}
		
		// No consecutive underscores
		if char == '_' && i < len(name)-1 && name[i+1] == '_' {
			return false
		}
	}

	// Cannot end with underscore (unless single character)
	if len(name) > 1 && name[len(name)-1] == '_' {
		return false
	}

	return true
}

// HasReservedPrefix checks if an identifier starts with a reserved prefix
func HasReservedPrefix(name string) bool {
	return strings.HasPrefix(name, "__")
}

// IsArrayAccess checks if an expression represents array access
func IsArrayAccess(expr parser.Expression) bool {
	_, ok := expr.(*parser.IndexedIdentifier)
	return ok
}

// GetArraySize attempts to extract array size from a declaration
func GetArraySize(sizeExpr parser.Expression) (int, bool) {
	if intLit, ok := sizeExpr.(*parser.IntegerLiteral); ok {
		return int(intLit.Value), true
	}
	return 0, false
}

// IsInGateDefinition checks if a node is within a gate definition
func IsInGateDefinition(program *parser.Program, targetNode parser.Node) bool {
	var found bool

	VisitAllNodes(program, func(node parser.Node) {
		if found {
			return
		}

		if gateDef, ok := node.(*parser.GateDefinition); ok {
			// Check if targetNode is within this gate's body
			for _, stmt := range gateDef.Body {
				if isNodeWithin(stmt, targetNode) {
					found = true
					return
				}
			}
		}
	})

	return found
}

// isNodeWithin checks if a target node is within a subtree rooted at root
func isNodeWithin(root, target parser.Node) bool {
	if root == target {
		return true
	}

	var found bool
	VisitAllNodes(root, func(node parser.Node) {
		if node == target {
			found = true
		}
	})

	return found
}

// GetIdentifierUsages collects all identifier usages in the program
func GetIdentifierUsages(program *parser.Program) map[string][]parser.Node {
	usageMap := make(map[string][]parser.Node)

	VisitAllNodes(program, func(node parser.Node) {
		switch n := node.(type) {
		case *parser.Identifier:
			usageMap[n.Name] = append(usageMap[n.Name], node)
		case *parser.IndexedIdentifier:
			// Record both the indexed access and the base name
			usageMap[n.Name] = append(usageMap[n.Name], node)
		case *parser.RangedIdentifier:
			usageMap[n.Name] = append(usageMap[n.Name], node)
		case *parser.GateCall:
			// Record gate name usage
			usageMap[n.Name] = append(usageMap[n.Name], node)
		}
	})

	return usageMap
}

// GetDeclaredIdentifiers returns all identifiers that are declared in the program
func GetDeclaredIdentifiers(program *parser.Program) map[string]parser.Node {
	declared := make(map[string]parser.Node)

	VisitAllNodes(program, func(node parser.Node) {
		switch n := node.(type) {
		case *parser.QuantumDeclaration:
			declared[n.Identifier] = node
		case *parser.ClassicalDeclaration:
			declared[n.Identifier] = node
		case *parser.GateDefinition:
			declared[n.Name] = node
		}
	})

	return declared
}

// IsValidInstruction checks if an instruction is valid within a gate definition
func IsValidInstruction(stmt parser.Statement) bool {
	switch stmt.(type) {
	case *parser.GateCall:
		return true // Gate calls are valid
	case *parser.Measurement:
		return false // Measurements are not allowed in gates
	case *parser.IfStatement:
		return false // Control flow generally not allowed in gates
	case *parser.ForStatement:
		return false
	case *parser.WhileStatement:
		return false
	default:
		return true // Default to allowing other statements
	}
}