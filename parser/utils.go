package parser

import (
	"fmt"
	"reflect"
	"strings"
)

// ASTTransformer interface for transforming AST nodes
type ASTTransformer interface {
	Transform(node Node) Node
}

// ValidationError represents semantic validation errors
type ValidationError struct {
	Message  string   `json:"message"`
	Position Position `json:"position"`
	Code     string   `json:"code"`
	Severity string   `json:"severity"`
}

func (ve *ValidationError) Error() string {
	return fmt.Sprintf("validation %s at line %d, column %d: %s (code: %s)",
		ve.Severity, ve.Position.Line, ve.Position.Column, ve.Message, ve.Code)
}

// ToQASM converts AST back to QASM string representation
func (p *Program) ToQASM() string {
	var builder strings.Builder

	// Add version if present
	if p.Version != nil {
		builder.WriteString(fmt.Sprintf("OPENQASM %s;\n", p.Version.Number))
	}

	// Add statements
	for _, stmt := range p.Statements {
		builder.WriteString(p.statementToQASM(stmt))
		builder.WriteString("\n")
	}

	return builder.String()
}

// statementToQASM converts a statement to QASM string
func (p *Program) statementToQASM(stmt Statement) string {
	switch s := stmt.(type) {
	case *QuantumDeclaration:
		return p.quantumDeclToQASM(s)
	case *ClassicalDeclaration:
		return p.classicalDeclToQASM(s)
	case *GateCall:
		return p.gateCallToQASM(s)
	case *Measurement:
		return p.measurementToQASM(s)
	case *Include:
		return fmt.Sprintf("include \"%s\";", s.Path)
	case *GateDefinition:
		return p.gateDefToQASM(s)
	case *IfStatement:
		return p.ifStmtToQASM(s)
	case *ForStatement:
		return p.forStmtToQASM(s)
	case *WhileStatement:
		return p.whileStmtToQASM(s)
	default:
		return "// Unknown statement type"
	}
}

func (p *Program) quantumDeclToQASM(q *QuantumDeclaration) string {
	if q.Size != nil {
		return fmt.Sprintf("%s[%s] %s;", q.Type, p.expressionToQASM(q.Size), q.Identifier)
	}
	return fmt.Sprintf("%s %s;", q.Type, q.Identifier)
}

func (p *Program) classicalDeclToQASM(c *ClassicalDeclaration) string {
	result := c.Type
	if c.Size != nil {
		result += fmt.Sprintf("[%s]", p.expressionToQASM(c.Size))
	}
	result += " " + c.Identifier
	if c.Initializer != nil {
		result += " = " + p.expressionToQASM(c.Initializer)
	}
	return result + ";"
}

func (p *Program) gateCallToQASM(g *GateCall) string {
	result := g.Name
	if len(g.Parameters) > 0 {
		params := make([]string, len(g.Parameters))
		for i, param := range g.Parameters {
			params[i] = p.expressionToQASM(param)
		}
		result += "(" + strings.Join(params, ", ") + ")"
	}

	if len(g.Qubits) > 0 {
		qubits := make([]string, len(g.Qubits))
		for i, qubit := range g.Qubits {
			qubits[i] = p.expressionToQASM(qubit)
		}
		result += " " + strings.Join(qubits, ", ")
	}

	return result + ";"
}

func (p *Program) measurementToQASM(m *Measurement) string {
	result := "measure " + p.expressionToQASM(m.Qubit)
	if m.Target != nil {
		result += " -> " + p.expressionToQASM(m.Target)
	}
	return result + ";"
}

func (p *Program) gateDefToQASM(g *GateDefinition) string {
	result := "gate " + g.Name

	if len(g.Parameters) > 0 {
		params := make([]string, len(g.Parameters))
		for i, param := range g.Parameters {
			params[i] = param.Name
		}
		result += "(" + strings.Join(params, ", ") + ")"
	}

	if len(g.Qubits) > 0 {
		qubits := make([]string, len(g.Qubits))
		for i, qubit := range g.Qubits {
			qubits[i] = qubit.Name
		}
		result += " " + strings.Join(qubits, ", ")
	}

	result += " {\n"
	for _, stmt := range g.Body {
		result += "  " + p.statementToQASM(stmt) + "\n"
	}
	result += "}"

	return result
}

func (p *Program) ifStmtToQASM(i *IfStatement) string {
	result := "if (" + p.expressionToQASM(i.Condition) + ") {\n"
	for _, stmt := range i.ThenBody {
		result += "  " + p.statementToQASM(stmt) + "\n"
	}
	result += "}"

	if len(i.ElseBody) > 0 {
		result += " else {\n"
		for _, stmt := range i.ElseBody {
			result += "  " + p.statementToQASM(stmt) + "\n"
		}
		result += "}"
	}

	return result
}

func (p *Program) forStmtToQASM(f *ForStatement) string {
	result := fmt.Sprintf("for %s in %s {\n", f.Variable, p.expressionToQASM(f.Iterable))
	for _, stmt := range f.Body {
		result += "  " + p.statementToQASM(stmt) + "\n"
	}
	result += "}"
	return result
}

func (p *Program) whileStmtToQASM(w *WhileStatement) string {
	result := "while (" + p.expressionToQASM(w.Condition) + ") {\n"
	for _, stmt := range w.Body {
		result += "  " + p.statementToQASM(stmt) + "\n"
	}
	result += "}"
	return result
}

func (p *Program) expressionToQASM(expr Expression) string {
	switch e := expr.(type) {
	case *Identifier:
		return e.Name
	case *IndexedIdentifier:
		return fmt.Sprintf("%s[%s]", e.Name, p.expressionToQASM(e.Index))
	case *RangedIdentifier:
		return fmt.Sprintf("%s[%s:%s]", e.Name, p.expressionToQASM(e.Start), p.expressionToQASM(e.EndIndex))
	case *IntegerLiteral:
		return fmt.Sprintf("%d", e.Value)
	case *FloatLiteral:
		return fmt.Sprintf("%g", e.Value)
	case *StringLiteral:
		return fmt.Sprintf("\"%s\"", e.Value)
	case *BooleanLiteral:
		if e.Value {
			return "true"
		}
		return "false"
	case *BinaryExpression:
		return fmt.Sprintf("(%s %s %s)", p.expressionToQASM(e.Left), e.Operator, p.expressionToQASM(e.Right))
	case *UnaryExpression:
		return fmt.Sprintf("%s%s", e.Operator, p.expressionToQASM(e.Operand))
	case *FunctionCall:
		args := make([]string, len(e.Arguments))
		for i, arg := range e.Arguments {
			args[i] = p.expressionToQASM(arg)
		}
		return fmt.Sprintf("%s(%s)", e.Name, strings.Join(args, ", "))
	case *ParenthesizedExpression:
		return "(" + p.expressionToQASM(e.Expression) + ")"
	default:
		return "/* unknown expression */"
	}
}

// Validate performs semantic validation on the AST
func (p *Program) Validate() []ValidationError {
	// Basic validation rules
	validator := &BasicValidator{errors: make([]ValidationError, 0)}

	// Validate version
	if p.Version == nil {
		validator.addError("missing OPENQASM version declaration", Position{Line: 1, Column: 1}, "V001", "error")
	}

	// Validate statements
	symbolTable := make(map[string]*TypeInfo)
	for _, stmt := range p.Statements {
		validator.validateStatement(stmt, symbolTable)
	}

	return validator.errors
}

// BasicValidator implements basic semantic validation
type BasicValidator struct {
	errors []ValidationError
}

func (bv *BasicValidator) addError(message string, pos Position, code string, severity string) {
	bv.errors = append(bv.errors, ValidationError{
		Message:  message,
		Position: pos,
		Code:     code,
		Severity: severity,
	})
}

func (bv *BasicValidator) validateStatement(stmt Statement, symbolTable map[string]*TypeInfo) {
	switch s := stmt.(type) {
	case *QuantumDeclaration:
		bv.validateQuantumDeclaration(s, symbolTable)
	case *ClassicalDeclaration:
		bv.validateClassicalDeclaration(s, symbolTable)
	case *GateCall:
		bv.validateGateCall(s, symbolTable)
	case *Measurement:
		bv.validateMeasurement(s, symbolTable)
		// Add more statement validations as needed
	}
}

func (bv *BasicValidator) validateQuantumDeclaration(q *QuantumDeclaration, symbolTable map[string]*TypeInfo) {
	// Check for duplicate declarations
	if _, exists := symbolTable[q.Identifier]; exists {
		bv.addError(fmt.Sprintf("identifier '%s' already declared", q.Identifier), q.Position, "V002", "error")
		return
	}

	// Add to symbol table
	typeInfo := &TypeInfo{
		Kind: q.Type,
	}

	if q.Size != nil {
		// Validate size expression and extract dimensions
		// This would need more sophisticated expression evaluation
		typeInfo.Dimensions = []int{1} // Placeholder
	}

	symbolTable[q.Identifier] = typeInfo
}

func (bv *BasicValidator) validateClassicalDeclaration(c *ClassicalDeclaration, symbolTable map[string]*TypeInfo) {
	// Check for duplicate declarations
	if _, exists := symbolTable[c.Identifier]; exists {
		bv.addError(fmt.Sprintf("identifier '%s' already declared", c.Identifier), c.Position, "V003", "error")
		return
	}

	// Add to symbol table
	typeInfo := &TypeInfo{
		Kind: c.Type,
	}

	if c.Size != nil {
		// Validate size expression
		typeInfo.Dimensions = []int{1} // Placeholder
	}

	symbolTable[c.Identifier] = typeInfo
}

func (bv *BasicValidator) validateGateCall(g *GateCall, symbolTable map[string]*TypeInfo) {
	// Validate that all qubit arguments are declared
	for _, qubit := range g.Qubits {
		if id, ok := qubit.(*Identifier); ok {
			if _, exists := symbolTable[id.Name]; !exists {
				bv.addError(fmt.Sprintf("undefined identifier '%s'", id.Name), g.Position, "V004", "error")
			}
		}
	}
}

func (bv *BasicValidator) validateMeasurement(m *Measurement, symbolTable map[string]*TypeInfo) {
	// Validate qubit being measured
	if id, ok := m.Qubit.(*Identifier); ok {
		if typeInfo, exists := symbolTable[id.Name]; exists {
			if typeInfo.Kind != "qubit" {
				bv.addError(fmt.Sprintf("cannot measure non-qubit '%s'", id.Name), m.Position, "V005", "error")
			}
		} else {
			bv.addError(fmt.Sprintf("undefined identifier '%s'", id.Name), m.Position, "V006", "error")
		}
	}
}

// Transform applies a transformation to the AST
func (p *Program) Transform(transformer ASTTransformer) *Program {
	if result := transformer.Transform(p); result != nil {
		if transformed, ok := result.(*Program); ok {
			return transformed
		}
	}
	return p
}

// FindByType finds all nodes of a specific type in the AST
func (p *Program) FindByType(nodeType reflect.Type) []Node {
	visitor := &TypeFinderVisitor{
		targetType: nodeType,
		results:    make([]Node, 0),
	}

	Walk(visitor, p)
	return visitor.results
}

// TypeFinderVisitor finds nodes of a specific type
type TypeFinderVisitor struct {
	BaseVisitor
	targetType reflect.Type
	results    []Node
}

func (tfv *TypeFinderVisitor) visitNode(node Node) {
	if reflect.TypeOf(node) == tfv.targetType {
		tfv.results = append(tfv.results, node)
	}
}

func (tfv *TypeFinderVisitor) VisitProgram(node *Program) interface{} {
	tfv.visitNode(node)
	return nil
}

func (tfv *TypeFinderVisitor) VisitQuantumDeclaration(node *QuantumDeclaration) interface{} {
	tfv.visitNode(node)
	return nil
}

func (tfv *TypeFinderVisitor) VisitClassicalDeclaration(node *ClassicalDeclaration) interface{} {
	tfv.visitNode(node)
	return nil
}

func (tfv *TypeFinderVisitor) VisitGateCall(node *GateCall) interface{} {
	tfv.visitNode(node)
	return nil
}

// FindByPosition finds the AST node at a specific position
func (p *Program) FindByPosition(pos Position) Node {
	visitor := &PositionFinderVisitor{
		targetPos: pos,
		result:    nil,
	}

	Walk(visitor, p)
	return visitor.result
}

// PositionFinderVisitor finds node at a specific position
type PositionFinderVisitor struct {
	BaseVisitor
	targetPos Position
	result    Node
}

func (pfv *PositionFinderVisitor) visitNode(node Node) {
	nodePos := node.Pos()
	nodeEnd := node.End()

	// Check if target position is within this node's range
	if pfv.targetPos.Line >= nodePos.Line && pfv.targetPos.Line <= nodeEnd.Line {
		if pfv.targetPos.Line == nodePos.Line && pfv.targetPos.Column >= nodePos.Column {
			pfv.result = node
		} else if pfv.targetPos.Line == nodeEnd.Line && pfv.targetPos.Column <= nodeEnd.Column {
			pfv.result = node
		} else if pfv.targetPos.Line > nodePos.Line && pfv.targetPos.Line < nodeEnd.Line {
			pfv.result = node
		}
	}
}

func (pfv *PositionFinderVisitor) VisitProgram(node *Program) interface{} {
	pfv.visitNode(node)
	return nil
}

func (pfv *PositionFinderVisitor) VisitQuantumDeclaration(node *QuantumDeclaration) interface{} {
	pfv.visitNode(node)
	return nil
}

func (pfv *PositionFinderVisitor) VisitClassicalDeclaration(node *ClassicalDeclaration) interface{} {
	pfv.visitNode(node)
	return nil
}

func (pfv *PositionFinderVisitor) VisitGateCall(node *GateCall) interface{} {
	pfv.visitNode(node)
	return nil
}
