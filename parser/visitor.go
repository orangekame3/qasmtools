package parser

// Visitor interface for AST traversal
type Visitor interface {
	VisitProgram(node *Program) interface{}
	VisitVersion(node *Version) interface{}
	VisitComment(node *Comment) interface{}

	// Statement visitors
	VisitQuantumDeclaration(node *QuantumDeclaration) interface{}
	VisitClassicalDeclaration(node *ClassicalDeclaration) interface{}
	VisitGateCall(node *GateCall) interface{}
	VisitMeasurement(node *Measurement) interface{}
	VisitInclude(node *Include) interface{}
	VisitGateDefinition(node *GateDefinition) interface{}
	VisitIfStatement(node *IfStatement) interface{}
	VisitForStatement(node *ForStatement) interface{}
	VisitWhileStatement(node *WhileStatement) interface{}

	// Expression visitors
	VisitIdentifier(node *Identifier) interface{}
	VisitIndexedIdentifier(node *IndexedIdentifier) interface{}
	VisitRangedIdentifier(node *RangedIdentifier) interface{}
	VisitIntegerLiteral(node *IntegerLiteral) interface{}
	VisitFloatLiteral(node *FloatLiteral) interface{}
	VisitStringLiteral(node *StringLiteral) interface{}
	VisitBooleanLiteral(node *BooleanLiteral) interface{}
	VisitBinaryExpression(node *BinaryExpression) interface{}
	VisitUnaryExpression(node *UnaryExpression) interface{}
	VisitFunctionCall(node *FunctionCall) interface{}
	VisitParenthesizedExpression(node *ParenthesizedExpression) interface{}

	// Other visitors
	VisitModifier(node *Modifier) interface{}
	VisitParameter(node *Parameter) interface{}
}

// BaseVisitor provides default implementations that return nil
type BaseVisitor struct{}

func (v *BaseVisitor) VisitProgram(node *Program) interface{}                           { return nil }
func (v *BaseVisitor) VisitVersion(node *Version) interface{}                           { return nil }
func (v *BaseVisitor) VisitComment(node *Comment) interface{}                           { return nil }
func (v *BaseVisitor) VisitQuantumDeclaration(node *QuantumDeclaration) interface{}     { return nil }
func (v *BaseVisitor) VisitClassicalDeclaration(node *ClassicalDeclaration) interface{} { return nil }
func (v *BaseVisitor) VisitGateCall(node *GateCall) interface{}                         { return nil }
func (v *BaseVisitor) VisitMeasurement(node *Measurement) interface{}                   { return nil }
func (v *BaseVisitor) VisitInclude(node *Include) interface{}                           { return nil }
func (v *BaseVisitor) VisitGateDefinition(node *GateDefinition) interface{}             { return nil }
func (v *BaseVisitor) VisitIfStatement(node *IfStatement) interface{}                   { return nil }
func (v *BaseVisitor) VisitForStatement(node *ForStatement) interface{}                 { return nil }
func (v *BaseVisitor) VisitWhileStatement(node *WhileStatement) interface{}             { return nil }
func (v *BaseVisitor) VisitIdentifier(node *Identifier) interface{}                     { return nil }
func (v *BaseVisitor) VisitIndexedIdentifier(node *IndexedIdentifier) interface{}       { return nil }
func (v *BaseVisitor) VisitRangedIdentifier(node *RangedIdentifier) interface{}         { return nil }
func (v *BaseVisitor) VisitIntegerLiteral(node *IntegerLiteral) interface{}             { return nil }
func (v *BaseVisitor) VisitFloatLiteral(node *FloatLiteral) interface{}                 { return nil }
func (v *BaseVisitor) VisitStringLiteral(node *StringLiteral) interface{}               { return nil }
func (v *BaseVisitor) VisitBooleanLiteral(node *BooleanLiteral) interface{}             { return nil }
func (v *BaseVisitor) VisitBinaryExpression(node *BinaryExpression) interface{}         { return nil }
func (v *BaseVisitor) VisitUnaryExpression(node *UnaryExpression) interface{}           { return nil }
func (v *BaseVisitor) VisitFunctionCall(node *FunctionCall) interface{}                 { return nil }
func (v *BaseVisitor) VisitParenthesizedExpression(node *ParenthesizedExpression) interface{} {
	return nil
}
func (v *BaseVisitor) VisitModifier(node *Modifier) interface{}   { return nil }
func (v *BaseVisitor) VisitParameter(node *Parameter) interface{} { return nil }

// Walk traverses AST with visitor using dispatch pattern
func Walk(visitor Visitor, node Node) interface{} {
	if node == nil {
		return nil
	}

	switch n := node.(type) {
	case *Program:
		return visitor.VisitProgram(n)
	case *Version:
		return visitor.VisitVersion(n)
	case *Comment:
		return visitor.VisitComment(n)
	case *QuantumDeclaration:
		return visitor.VisitQuantumDeclaration(n)
	case *ClassicalDeclaration:
		return visitor.VisitClassicalDeclaration(n)
	case *GateCall:
		return visitor.VisitGateCall(n)
	case *Measurement:
		return visitor.VisitMeasurement(n)
	case *Include:
		return visitor.VisitInclude(n)
	case *GateDefinition:
		return visitor.VisitGateDefinition(n)
	case *IfStatement:
		return visitor.VisitIfStatement(n)
	case *ForStatement:
		return visitor.VisitForStatement(n)
	case *WhileStatement:
		return visitor.VisitWhileStatement(n)
	case *Identifier:
		return visitor.VisitIdentifier(n)
	case *IndexedIdentifier:
		return visitor.VisitIndexedIdentifier(n)
	case *RangedIdentifier:
		return visitor.VisitRangedIdentifier(n)
	case *IntegerLiteral:
		return visitor.VisitIntegerLiteral(n)
	case *FloatLiteral:
		return visitor.VisitFloatLiteral(n)
	case *StringLiteral:
		return visitor.VisitStringLiteral(n)
	case *BooleanLiteral:
		return visitor.VisitBooleanLiteral(n)
	case *BinaryExpression:
		return visitor.VisitBinaryExpression(n)
	case *UnaryExpression:
		return visitor.VisitUnaryExpression(n)
	case *FunctionCall:
		return visitor.VisitFunctionCall(n)
	case *ParenthesizedExpression:
		return visitor.VisitParenthesizedExpression(n)
	case *Modifier:
		return visitor.VisitModifier(n)
	case *Parameter:
		return visitor.VisitParameter(n)
	default:
		// Unknown node type
		return nil
	}
}

// WalkStatements traverses a slice of statements
func WalkStatements(visitor Visitor, statements []Statement) []interface{} {
	results := make([]interface{}, len(statements))
	for i, stmt := range statements {
		results[i] = Walk(visitor, stmt)
	}
	return results
}

// WalkExpressions traverses a slice of expressions
func WalkExpressions(visitor Visitor, expressions []Expression) []interface{} {
	results := make([]interface{}, len(expressions))
	for i, expr := range expressions {
		results[i] = Walk(visitor, expr)
	}
	return results
}

// DepthFirstVisitor provides depth-first traversal with automatic child visiting
type DepthFirstVisitor struct {
	BaseVisitor
	visitor Visitor
}

// NewDepthFirstVisitor creates a visitor that automatically traverses children
func NewDepthFirstVisitor(visitor Visitor) *DepthFirstVisitor {
	return &DepthFirstVisitor{visitor: visitor}
}

func (d *DepthFirstVisitor) VisitProgram(node *Program) interface{} {
	result := d.visitor.VisitProgram(node)
	if node.Version != nil {
		Walk(d, node.Version)
	}
	WalkStatements(d, node.Statements)
	for _, comment := range node.Comments {
		Walk(d, &comment)
	}
	return result
}

func (d *DepthFirstVisitor) VisitGateCall(node *GateCall) interface{} {
	result := d.visitor.VisitGateCall(node)
	WalkExpressions(d, node.Parameters)
	WalkExpressions(d, node.Qubits)
	for _, modifier := range node.Modifiers {
		Walk(d, &modifier)
	}
	return result
}

func (d *DepthFirstVisitor) VisitMeasurement(node *Measurement) interface{} {
	result := d.visitor.VisitMeasurement(node)
	Walk(d, node.Qubit)
	if node.Target != nil {
		Walk(d, node.Target)
	}
	return result
}

func (d *DepthFirstVisitor) VisitGateDefinition(node *GateDefinition) interface{} {
	result := d.visitor.VisitGateDefinition(node)
	for _, param := range node.Parameters {
		Walk(d, &param)
	}
	for _, qubit := range node.Qubits {
		Walk(d, &qubit)
	}
	WalkStatements(d, node.Body)
	return result
}

func (d *DepthFirstVisitor) VisitIfStatement(node *IfStatement) interface{} {
	result := d.visitor.VisitIfStatement(node)
	Walk(d, node.Condition)
	WalkStatements(d, node.ThenBody)
	WalkStatements(d, node.ElseBody)
	return result
}

func (d *DepthFirstVisitor) VisitForStatement(node *ForStatement) interface{} {
	result := d.visitor.VisitForStatement(node)
	Walk(d, node.Iterable)
	WalkStatements(d, node.Body)
	return result
}

func (d *DepthFirstVisitor) VisitWhileStatement(node *WhileStatement) interface{} {
	result := d.visitor.VisitWhileStatement(node)
	Walk(d, node.Condition)
	WalkStatements(d, node.Body)
	return result
}

func (d *DepthFirstVisitor) VisitIndexedIdentifier(node *IndexedIdentifier) interface{} {
	result := d.visitor.VisitIndexedIdentifier(node)
	Walk(d, node.Index)
	return result
}

func (d *DepthFirstVisitor) VisitRangedIdentifier(node *RangedIdentifier) interface{} {
	result := d.visitor.VisitRangedIdentifier(node)
	Walk(d, node.Start)
	Walk(d, node.EndIndex)
	return result
}

func (d *DepthFirstVisitor) VisitBinaryExpression(node *BinaryExpression) interface{} {
	result := d.visitor.VisitBinaryExpression(node)
	Walk(d, node.Left)
	Walk(d, node.Right)
	return result
}

func (d *DepthFirstVisitor) VisitUnaryExpression(node *UnaryExpression) interface{} {
	result := d.visitor.VisitUnaryExpression(node)
	Walk(d, node.Operand)
	return result
}

func (d *DepthFirstVisitor) VisitFunctionCall(node *FunctionCall) interface{} {
	result := d.visitor.VisitFunctionCall(node)
	WalkExpressions(d, node.Arguments)
	return result
}

func (d *DepthFirstVisitor) VisitParenthesizedExpression(node *ParenthesizedExpression) interface{} {
	result := d.visitor.VisitParenthesizedExpression(node)
	Walk(d, node.Expression)
	return result
}

func (d *DepthFirstVisitor) VisitModifier(node *Modifier) interface{} {
	result := d.visitor.VisitModifier(node)
	WalkExpressions(d, node.Parameters)
	return result
}

// Delegate other methods to wrapped visitor
func (d *DepthFirstVisitor) VisitVersion(node *Version) interface{} {
	return d.visitor.VisitVersion(node)
}
func (d *DepthFirstVisitor) VisitComment(node *Comment) interface{} {
	return d.visitor.VisitComment(node)
}
func (d *DepthFirstVisitor) VisitQuantumDeclaration(node *QuantumDeclaration) interface{} {
	return d.visitor.VisitQuantumDeclaration(node)
}
func (d *DepthFirstVisitor) VisitClassicalDeclaration(node *ClassicalDeclaration) interface{} {
	return d.visitor.VisitClassicalDeclaration(node)
}
func (d *DepthFirstVisitor) VisitInclude(node *Include) interface{} {
	return d.visitor.VisitInclude(node)
}
func (d *DepthFirstVisitor) VisitIdentifier(node *Identifier) interface{} {
	return d.visitor.VisitIdentifier(node)
}
func (d *DepthFirstVisitor) VisitIntegerLiteral(node *IntegerLiteral) interface{} {
	return d.visitor.VisitIntegerLiteral(node)
}
func (d *DepthFirstVisitor) VisitFloatLiteral(node *FloatLiteral) interface{} {
	return d.visitor.VisitFloatLiteral(node)
}
func (d *DepthFirstVisitor) VisitStringLiteral(node *StringLiteral) interface{} {
	return d.visitor.VisitStringLiteral(node)
}
func (d *DepthFirstVisitor) VisitBooleanLiteral(node *BooleanLiteral) interface{} {
	return d.visitor.VisitBooleanLiteral(node)
}
func (d *DepthFirstVisitor) VisitParameter(node *Parameter) interface{} {
	return d.visitor.VisitParameter(node)
}
