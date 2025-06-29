package parser

import "fmt"

// Position represents source code position
type Position struct {
	Line   int    `json:"line"`
	Column int    `json:"column"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`         // Length of the range
	File   string `json:"file,omitempty"` // Source file name
}

// Node represents any AST node
type Node interface {
	Pos() Position
	End() Position
	String() string
}

// BaseNode provides common functionality for all AST nodes
type BaseNode struct {
	Position Position `json:"position"`
	EndPos   Position `json:"end_position"`
}

func (n *BaseNode) Pos() Position {
	return n.Position
}

func (n *BaseNode) End() Position {
	return n.EndPos
}

// Program represents the root AST node
type Program struct {
	BaseNode
	Version      *Version          `json:"version,omitempty"`
	Statements   []Statement       `json:"statements"`
	Comments     []Comment         `json:"comments,omitempty"`
	LineComments map[int][]Comment `json:"line_comments,omitempty"`
}

func (p *Program) String() string {
	return "Program"
}

// Statement represents any QASM statement
type Statement interface {
	Node
	StatementNode()
}

// Expression interface for all expressions
type Expression interface {
	Node
	ExpressionNode()
}

// Version represents OPENQASM version declaration
type Version struct {
	BaseNode
	Number string `json:"number"`
}

func (v *Version) String() string {
	return "Version: " + v.Number
}

// Comment represents code comments
type Comment struct {
	BaseNode
	Text     string `json:"text"`
	Type     string `json:"type"`     // "line" or "block"
	Attached bool   `json:"attached"` // whether attached to a statement
}

func (c *Comment) String() string {
	return "Comment: " + c.Text
}

// QuantumDeclaration represents qubit declarations
type QuantumDeclaration struct {
	BaseNode
	Type       string     `json:"type"`           // "qubit"
	Size       Expression `json:"size,omitempty"` // for qubit[n]
	Identifier string     `json:"identifier"`
	TypeInfo   *TypeInfo  `json:"type_info,omitempty"` // Type information
}

func (q *QuantumDeclaration) StatementNode() {}
func (q *QuantumDeclaration) String() string {
	return "QuantumDeclaration: " + q.Identifier
}

// ClassicalDeclaration represents classical variable declarations
type ClassicalDeclaration struct {
	BaseNode
	Type        string     `json:"type"`           // "bit", "int", "float", etc.
	Size        Expression `json:"size,omitempty"` // for bit[n], int[32], etc.
	Identifier  string     `json:"identifier"`
	Initializer Expression `json:"initializer,omitempty"`
	TypeInfo    *TypeInfo  `json:"type_info,omitempty"` // Type information
}

func (c *ClassicalDeclaration) StatementNode() {}
func (c *ClassicalDeclaration) String() string {
	return "ClassicalDeclaration: " + c.Identifier
}

// GateCall represents gate applications
type GateCall struct {
	BaseNode
	Name       string       `json:"name"`
	Parameters []Expression `json:"parameters,omitempty"`
	Qubits     []Expression `json:"qubits"`
	Modifiers  []Modifier   `json:"modifiers,omitempty"`
}

func (g *GateCall) StatementNode() {}
func (g *GateCall) String() string {
	return "GateCall: " + g.Name
}

// Modifier represents gate modifiers like inv, ctrl, etc.
type Modifier struct {
	BaseNode
	Type       string       `json:"type"`                 // "inv", "ctrl", "pow"
	Parameters []Expression `json:"parameters,omitempty"` // for parameterized modifiers
}

func (m *Modifier) String() string {
	return "Modifier: " + m.Type
}

// Measurement represents measure statements
type Measurement struct {
	BaseNode
	Qubit  Expression `json:"qubit"`
	Target Expression `json:"target,omitempty"` // for measure q -> c
}

func (m *Measurement) StatementNode() {}
func (m *Measurement) String() string {
	return "Measurement"
}

// Include represents include statements
type Include struct {
	BaseNode
	Path string `json:"path"`
}

func (i *Include) StatementNode() {}
func (i *Include) String() string {
	return "Include: " + i.Path
}

// GateDefinition represents gate definitions
type GateDefinition struct {
	BaseNode
	Name       string      `json:"name"`
	Parameters []Parameter `json:"parameters,omitempty"`
	Qubits     []Parameter `json:"qubits"`
	Body       []Statement `json:"body"`
}

func (g *GateDefinition) StatementNode() {}
func (g *GateDefinition) String() string {
	return "GateDefinition: " + g.Name
}

// Parameter represents function/gate parameters
type Parameter struct {
	BaseNode
	Name string `json:"name"`
	Type string `json:"type,omitempty"`
}

func (p *Parameter) String() string {
	return "Parameter: " + p.Name
}

// IfStatement represents conditional statements
type IfStatement struct {
	BaseNode
	Condition Expression  `json:"condition"`
	ThenBody  []Statement `json:"then_body"`
	ElseBody  []Statement `json:"else_body,omitempty"`
}

func (i *IfStatement) StatementNode() {}
func (i *IfStatement) String() string {
	return "IfStatement"
}

// ForStatement represents for loops
type ForStatement struct {
	BaseNode
	Variable string      `json:"variable"`
	Iterable Expression  `json:"iterable"`
	Body     []Statement `json:"body"`
}

func (f *ForStatement) StatementNode() {}
func (f *ForStatement) String() string {
	return "ForStatement"
}

// WhileStatement represents while loops
type WhileStatement struct {
	BaseNode
	Condition Expression  `json:"condition"`
	Body      []Statement `json:"body"`
}

func (w *WhileStatement) StatementNode() {}
func (w *WhileStatement) String() string {
	return "WhileStatement"
}

// Expression implementations

// Identifier represents variable references
type Identifier struct {
	BaseNode
	Name string `json:"name"`
}

func (i *Identifier) ExpressionNode() {}
func (i *Identifier) String() string {
	return "Identifier: " + i.Name
}

// IndexedIdentifier represents array access like q[0]
type IndexedIdentifier struct {
	BaseNode
	Name  string     `json:"name"`
	Index Expression `json:"index"`
}

func (i *IndexedIdentifier) ExpressionNode() {}
func (i *IndexedIdentifier) String() string {
	return "IndexedIdentifier: " + i.Name
}

// RangedIdentifier represents range access like q[0:2]
type RangedIdentifier struct {
	BaseNode
	Name     string     `json:"name"`
	Start    Expression `json:"start"`
	EndIndex Expression `json:"end"`
}

func (r *RangedIdentifier) ExpressionNode() {}
func (r *RangedIdentifier) String() string {
	return "RangedIdentifier: " + r.Name
}

// IntegerLiteral represents integer constants
type IntegerLiteral struct {
	BaseNode
	Value int64 `json:"value"`
}

func (i *IntegerLiteral) ExpressionNode() {}
func (i *IntegerLiteral) String() string {
	return "IntegerLiteral"
}

// FloatLiteral represents floating-point constants
type FloatLiteral struct {
	BaseNode
	Value float64 `json:"value"`
}

func (f *FloatLiteral) ExpressionNode() {}
func (f *FloatLiteral) String() string {
	return "FloatLiteral"
}

// StringLiteral represents string constants
type StringLiteral struct {
	BaseNode
	Value string `json:"value"`
}

func (s *StringLiteral) ExpressionNode() {}
func (s *StringLiteral) String() string {
	return "StringLiteral: " + s.Value
}

// BooleanLiteral represents boolean constants
type BooleanLiteral struct {
	BaseNode
	Value bool `json:"value"`
}

func (b *BooleanLiteral) ExpressionNode() {}
func (b *BooleanLiteral) String() string {
	return "BooleanLiteral"
}

// BinaryExpression represents binary operations
type BinaryExpression struct {
	BaseNode
	Left     Expression `json:"left"`
	Operator string     `json:"operator"`
	Right    Expression `json:"right"`
}

func (b *BinaryExpression) ExpressionNode() {}
func (b *BinaryExpression) String() string {
	return "BinaryExpression: " + b.Operator
}

// UnaryExpression represents unary operations
type UnaryExpression struct {
	BaseNode
	Operator string     `json:"operator"`
	Operand  Expression `json:"operand"`
}

func (u *UnaryExpression) ExpressionNode() {}
func (u *UnaryExpression) String() string {
	return "UnaryExpression: " + u.Operator
}

// FunctionCall represents function calls
type FunctionCall struct {
	BaseNode
	Name      string       `json:"name"`
	Arguments []Expression `json:"arguments"`
}

func (f *FunctionCall) ExpressionNode() {}
func (f *FunctionCall) String() string {
	return "FunctionCall: " + f.Name
}

// ParenthesizedExpression represents expressions in parentheses
type ParenthesizedExpression struct {
	BaseNode
	Expression Expression `json:"expression"`
}

func (p *ParenthesizedExpression) ExpressionNode() {}
func (p *ParenthesizedExpression) String() string {
	return "ParenthesizedExpression"
}

// TimingExpression represents timing expressions like 100ns
type TimingExpression struct {
	BaseNode
	Value Expression `json:"value"`
	Unit  string     `json:"unit"` // "ns", "us", "ms", etc.
}

func (t *TimingExpression) ExpressionNode() {}
func (t *TimingExpression) String() string {
	return "TimingExpression"
}

// DelayExpression represents delay expressions like delay[100ns]
type DelayExpression struct {
	BaseNode
	Timing Expression `json:"timing"`
}

func (d *DelayExpression) ExpressionNode() {}
func (d *DelayExpression) String() string {
	return "DelayExpression"
}

// TypeInfo represents type information for AST nodes
type TypeInfo struct {
	Kind        string   `json:"kind"`                  // "qubit", "bit", "int", "float", "bool"
	Dimensions  []int    `json:"dimensions"`            // Array dimensions
	Constraints []string `json:"constraints,omitempty"` // Type constraints
	BitWidth    int      `json:"bit_width,omitempty"`   // For int/float types
}

// IsArray returns true if this type represents an array
func (ti *TypeInfo) IsArray() bool {
	return len(ti.Dimensions) > 0
}

// ArraySize returns the total size of the array (product of all dimensions)
func (ti *TypeInfo) ArraySize() int {
	if !ti.IsArray() {
		return 1
	}

	size := 1
	for _, dim := range ti.Dimensions {
		size *= dim
	}
	return size
}

// String returns a string representation of the type
func (ti *TypeInfo) String() string {
	result := ti.Kind

	if ti.BitWidth > 0 {
		result += fmt.Sprintf("[%d]", ti.BitWidth)
	}

	for _, dim := range ti.Dimensions {
		result += fmt.Sprintf("[%d]", dim)
	}

	return result
}
