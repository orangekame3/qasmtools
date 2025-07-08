package parser

import (
	"strings"

	"github.com/antlr4-go/antlr/v4"
	qasm_gen "github.com/orangekame3/qasmtools/parser/gen"
)

// ASTBuilderVisitor implements ANTLR visitor to build our AST
type ASTBuilderVisitor struct {
	*qasm_gen.Baseqasm3ParserVisitor
	errors []ParseError
}

// NewASTBuilderVisitor creates a new AST builder visitor
func NewASTBuilderVisitor() *ASTBuilderVisitor {
	return &ASTBuilderVisitor{
		Baseqasm3ParserVisitor: &qasm_gen.Baseqasm3ParserVisitor{},
		errors:                 make([]ParseError, 0),
	}
}

// GetErrors returns any errors encountered during AST building
func (v *ASTBuilderVisitor) GetErrors() []ParseError {
	return v.errors
}

// addError adds an error to the error list
func (v *ASTBuilderVisitor) addError(ctx antlr.ParserRuleContext, message string) {
	pos := v.getPosition(ctx)
	v.errors = append(v.errors, ParseError{
		Message:  message,
		Position: pos,
		Type:     "semantic",
		Context:  ctx.GetText(),
	})
}

// getPosition extracts position information from parser context
func (v *ASTBuilderVisitor) getPosition(ctx antlr.ParserRuleContext) Position {
	if ctx == nil {
		return Position{Line: 1, Column: 1}
	}

	start := ctx.GetStart()
	if start == nil {
		return Position{Line: 1, Column: 1}
	}

	return Position{
		Line:   start.GetLine(),
		Column: start.GetColumn() + 1, // ANTLR uses 0-based columns
		Offset: start.GetStart(),
	}
}

// getEndPosition extracts end position information from parser context
func (v *ASTBuilderVisitor) getEndPosition(ctx antlr.ParserRuleContext) Position {
	if ctx == nil {
		return Position{Line: 1, Column: 1}
	}

	stop := ctx.GetStop()
	if stop == nil {
		return v.getPosition(ctx)
	}

	return Position{
		Line:   stop.GetLine(),
		Column: stop.GetColumn() + len(stop.GetText()),
		Offset: stop.GetStop() + 1,
	}
}

// createBaseNode creates a BaseNode with position information
func (v *ASTBuilderVisitor) createBaseNode(ctx antlr.ParserRuleContext) BaseNode {
	return BaseNode{
		Position: v.getPosition(ctx),
		EndPos:   v.getEndPosition(ctx),
	}
}

// VisitProgram builds a Program node using improved ANTLR visitor pattern
func (v *ASTBuilderVisitor) VisitProgram(ctx *qasm_gen.ProgramContext) interface{} {
	program := &Program{
		BaseNode:     v.createBaseNode(ctx),
		Statements:   make([]Statement, 0),
		Comments:     make([]Comment, 0),
		LineComments: make(map[int][]Comment),
	}

	// Parse version if present
	if ctx.Version() != nil {
		program.Version = &Version{
			BaseNode: v.createBaseNode(ctx),
			Number:   "3.0",
		}
	}

	// Visit all statement or scope contexts using proper ANTLR pattern
	for _, statementOrScopeCtx := range ctx.AllStatementOrScope() {
		if statementOrScopeCtx != nil {
			if stmt := v.visitStatementOrScope(statementOrScopeCtx); stmt != nil {
				program.Statements = append(program.Statements, stmt)
			}
		}
	}

	// If no statements were found through ANTLR, fall back to text parsing
	if len(program.Statements) == 0 {
		fullText := ctx.GetText()
		fullText = strings.Replace(fullText, "<EOF>", "", -1)
		if fullText != "" {
			statements := v.parseStatementsFromText(fullText)
			program.Statements = statements
		}
	}

	return program
}

// visitStatementOrScope handles statement or scope contexts
func (v *ASTBuilderVisitor) visitStatementOrScope(ctx qasm_gen.IStatementOrScopeContext) Statement {
	if ctx == nil {
		return nil
	}

	// Try to get the underlying statement context
	if statementCtx := ctx.Statement(); statementCtx != nil {
		return v.visitStatement(statementCtx)
	}

	// Handle scope contexts if needed (for now, return nil)
	return nil
}

// visitStatement handles individual statement contexts using proper ANTLR pattern
func (v *ASTBuilderVisitor) visitStatement(ctx qasm_gen.IStatementContext) Statement {
	if ctx == nil {
		return nil
	}

	// We'll use createBaseNode for position information

	// Use proper ANTLR context methods instead of text parsing

	// Check for quantum declaration statement
	if quantumDeclCtx := ctx.QuantumDeclarationStatement(); quantumDeclCtx != nil {
		return v.visitQuantumDeclarationStatement(quantumDeclCtx)
	}

	// Check for classical declaration statement
	if classicalDeclCtx := ctx.ClassicalDeclarationStatement(); classicalDeclCtx != nil {
		return v.visitClassicalDeclarationStatement(classicalDeclCtx)
	}

	// Check for include statement
	if includeCtx := ctx.IncludeStatement(); includeCtx != nil {
		return v.visitIncludeStatement(includeCtx)
	}

	// Check for gate call statement
	if gateCallCtx := ctx.GateCallStatement(); gateCallCtx != nil {
		return v.visitGateCallStatement(gateCallCtx)
	}

	// Check for measurement assignment statement (measure -> classical)
	if measureCtx := ctx.MeasureArrowAssignmentStatement(); measureCtx != nil {
		return v.visitMeasureArrowAssignmentStatement(measureCtx)
	}

	// Check for expression statement (other expressions)
	if exprCtx := ctx.ExpressionStatement(); exprCtx != nil {
		return v.visitExpressionStatement(exprCtx)
	}

	// Check for assignment statement (other assignments)
	if assignCtx := ctx.AssignmentStatement(); assignCtx != nil {
		return v.visitAssignmentStatement(assignCtx)
	}

	return nil
}

// visitQuantumDeclarationStatement handles quantum declarations using ANTLR context
func (v *ASTBuilderVisitor) visitQuantumDeclarationStatement(ctx qasm_gen.IQuantumDeclarationStatementContext) Statement {
	if ctx == nil {
		return nil
	}

	// Position will be handled by createBaseNode

	// Get the identifier name
	var identifier string
	if idNode := ctx.Identifier(); idNode != nil {
		identifier = idNode.GetText()
	}

	// Get the qubit type information
	var size Expression
	var declType = "qubit"

	if qubitType := ctx.QubitType(); qubitType != nil {
		// Get designator (array size) if present
		if designator := qubitType.Designator(); designator != nil {
			// Parse the designator to get size
			if sizeExpr := v.visitDesignator(designator); sizeExpr != nil {
				size = sizeExpr
			}
		}
	}

	return &QuantumDeclaration{
		BaseNode:   v.createBaseNode(ctx),
		Type:       declType,
		Size:       size,
		Identifier: identifier,
	}
}

// visitDesignator handles array designators like [2] in qubit[2]
func (v *ASTBuilderVisitor) visitDesignator(ctx qasm_gen.IDesignatorContext) Expression {
	if ctx == nil {
		return nil
	}

	// For now, handle simple integer designators
	// This is a simplified implementation - full implementation would handle expressions
	text := ctx.GetText()
	// Remove brackets and parse as integer
	if len(text) >= 3 && text[0] == '[' && text[len(text)-1] == ']' {
		sizeStr := text[1 : len(text)-1]
		if val := parseInt(sizeStr); val > 0 {
			return &IntegerLiteral{Value: val}
		}
	}

	return nil
}

// visitClassicalDeclarationStatement handles classical declarations using ANTLR context
func (v *ASTBuilderVisitor) visitClassicalDeclarationStatement(ctx qasm_gen.IClassicalDeclarationStatementContext) Statement {
	if ctx == nil {
		return nil
	}

	// Get the identifier name
	var identifier string
	if idNode := ctx.Identifier(); idNode != nil {
		identifier = idNode.GetText()
	}

	// Get type information
	var size Expression
	var declType string = "bit" // default

	// Check for scalar type (bit, int, etc.)
	if scalarType := ctx.ScalarType(); scalarType != nil {
		declType = scalarType.GetText()
	}

	// Check for array type (bit[n], int[n], etc.)
	if arrayType := ctx.ArrayType(); arrayType != nil {
		// Parse array type to get base type and size
		text := arrayType.GetText()
		if idx := strings.Index(text, "["); idx != -1 {
			declType = text[:idx]
			if endIdx := strings.Index(text, "]"); endIdx > idx {
				sizeStr := text[idx+1 : endIdx]
				if val := parseInt(sizeStr); val > 0 {
					size = &IntegerLiteral{Value: val}
				}
			}
		}
	}

	return &ClassicalDeclaration{
		BaseNode:    v.createBaseNode(ctx),
		Type:        declType,
		Size:        size,
		Identifier:  identifier,
		Initializer: nil, // TODO: Handle initializers if needed
	}
}

// visitIncludeStatement handles include statements using ANTLR context
func (v *ASTBuilderVisitor) visitIncludeStatement(ctx qasm_gen.IIncludeStatementContext) Statement {
	if ctx == nil {
		return nil
	}

	// Get the string literal for the include path
	var path string
	if stringLit := ctx.StringLiteral(); stringLit != nil {
		pathText := stringLit.GetText()
		// Remove quotes
		if len(pathText) >= 2 && pathText[0] == '"' && pathText[len(pathText)-1] == '"' {
			path = pathText[1 : len(pathText)-1]
		}
	}

	return &Include{
		BaseNode: v.createBaseNode(ctx),
		Path:     path,
	}
}

// visitGateCallStatement handles gate call statements using ANTLR context
func (v *ASTBuilderVisitor) visitGateCallStatement(ctx qasm_gen.IGateCallStatementContext) Statement {
	if ctx == nil {
		return nil
	}

	// Get gate name
	var gateName string
	if idNode := ctx.Identifier(); idNode != nil {
		gateName = idNode.GetText()
	}

	// Get parameters (in parentheses)
	var parameters []Expression
	if exprList := ctx.ExpressionList(); exprList != nil {
		parameters = v.visitExpressionList(exprList)
	}

	// Get qubits (gate operands)
	var qubits []Expression
	if operandList := ctx.GateOperandList(); operandList != nil {
		qubits = v.visitGateOperandList(operandList)
	}

	// Get modifiers (like ctrl, inv, etc.)
	var modifiers []Modifier
	for _, modCtx := range ctx.AllGateModifier() {
		if mod := v.visitGateModifier(modCtx); mod != nil {
			modifiers = append(modifiers, *mod)
		}
	}

	return &GateCall{
		BaseNode:   v.createBaseNode(ctx),
		Name:       gateName,
		Parameters: parameters,
		Qubits:     qubits,
		Modifiers:  modifiers,
	}
}

// visitMeasureArrowAssignmentStatement handles measurement assignments (measure -> classical)
func (v *ASTBuilderVisitor) visitMeasureArrowAssignmentStatement(ctx qasm_gen.IMeasureArrowAssignmentStatementContext) Statement {
	if ctx == nil {
		return nil
	}

	// Get the measurement expression (what we're measuring)
	var qubit Expression
	if measureExpr := ctx.MeasureExpression(); measureExpr != nil {
		qubit = v.visitMeasureExpression(measureExpr)
	}

	// Get the target (where we store the result)
	var target Expression
	if targetCtx := ctx.IndexedIdentifier(); targetCtx != nil {
		target = v.visitIndexedIdentifier(targetCtx)
	}

	return &Measurement{
		BaseNode: v.createBaseNode(ctx),
		Qubit:    qubit,
		Target:   target,
	}
}

// visitExpressionStatement handles expression statements (other expressions)
func (v *ASTBuilderVisitor) visitExpressionStatement(ctx qasm_gen.IExpressionStatementContext) Statement {
	// For now, we don't handle general expressions as statements
	// Most relevant expressions (gate calls, measurements) have their own statement types
	return nil
}

// visitAssignmentStatement handles assignment statements (other assignments)
func (v *ASTBuilderVisitor) visitAssignmentStatement(ctx qasm_gen.IAssignmentStatementContext) Statement {
	// For now, we handle measurements separately
	// Other assignments might be classical variable assignments
	return nil
}

// visitGateOperandList handles list of gate operands (qubits)
func (v *ASTBuilderVisitor) visitGateOperandList(ctx qasm_gen.IGateOperandListContext) []Expression {
	if ctx == nil {
		return nil
	}

	var operands []Expression
	// For now, use simplified parsing - proper implementation would parse each operand
	text := ctx.GetText()
	// Split by commas and create identifiers
	if text != "" {
		parts := strings.Split(text, ",")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" {
				operands = append(operands, v.parseQubitExpression(part))
			}
		}
	}

	return operands
}

// visitExpressionList handles parameter lists
func (v *ASTBuilderVisitor) visitExpressionList(ctx qasm_gen.IExpressionListContext) []Expression {
	if ctx == nil {
		return nil
	}

	// For now, return empty - parameters are less critical for basic linting
	return make([]Expression, 0)
}

// visitGateModifier handles gate modifiers
func (v *ASTBuilderVisitor) visitGateModifier(ctx qasm_gen.IGateModifierContext) *Modifier {
	if ctx == nil {
		return nil
	}

	// For now, return empty modifier - modifiers are less critical for basic linting
	return &Modifier{
		Type: "unknown",
	}
}

// visitMeasureExpression handles measurement expressions
func (v *ASTBuilderVisitor) visitMeasureExpression(ctx qasm_gen.IMeasureExpressionContext) Expression {
	if ctx == nil {
		return nil
	}

	// For now, use simplified parsing
	text := ctx.GetText()
	// Remove "measure" prefix and parse as qubit expression
	if strings.HasPrefix(text, "measure") {
		qubitText := strings.TrimSpace(text[7:]) // Remove "measure"
		return v.parseQubitExpression(qubitText)
	}

	return nil
}

// visitIndexedIdentifier handles indexed identifiers (like c[0])
func (v *ASTBuilderVisitor) visitIndexedIdentifier(ctx qasm_gen.IIndexedIdentifierContext) Expression {
	if ctx == nil {
		return nil
	}

	// For now, use simplified parsing
	text := ctx.GetText()
	return v.parseQubitExpression(text)
}

// isGateCall checks if the text represents a gate call
func (v *ASTBuilderVisitor) isGateCall(text string) bool {
	// Common gates
	commonGates := []string{"h", "x", "y", "z", "cx", "cnot", "s", "t", "rx", "ry", "rz", "u1", "u2", "u3"}

	for _, gate := range commonGates {
		if strings.HasPrefix(strings.TrimSpace(text), gate+" ") ||
			strings.HasPrefix(strings.TrimSpace(text), gate+"\t") {
			return true
		}
	}

	return false
}

// parseStatementsFromText provides simple text-based parsing
func (v *ASTBuilderVisitor) parseStatementsFromText(fullText string) []Statement {
	statements := make([]Statement, 0)

	lines := strings.Split(fullText, "\n")
	for lineNum, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "OPENQASM") {
			continue
		}

		// Simple statement parsing
		if stmt := v.parseSimpleStatement(line, lineNum+1); stmt != nil {
			statements = append(statements, stmt)
		}
	}

	return statements
}

// parseSimpleStatement parses a single statement
func (v *ASTBuilderVisitor) parseSimpleStatement(line string, lineNum int) Statement {
	line = strings.TrimSuffix(strings.TrimSpace(line), ";")

	// Parse different statement types
	if strings.Contains(line, "qubit") {
		return v.parseQuantumDeclaration(line, lineNum)
	}

	if (strings.Contains(line, "bit") || strings.Contains(line, "int")) && !strings.Contains(line, "qubit") {
		return v.parseClassicalDeclaration(line, lineNum)
	}

	if strings.Contains(line, "include") {
		path := "stdgates.inc"
		if start := strings.Index(line, "\""); start != -1 {
			if end := strings.LastIndex(line, "\""); end > start {
				path = line[start+1 : end]
			}
		}
		return &Include{
			BaseNode: BaseNode{Position: Position{Line: lineNum, Column: 1}},
			Path:     path,
		}
	}

	if strings.Contains(line, "measure") {
		return v.parseMeasurementStatement(line, lineNum)
	}

	// Check for gate calls
	if v.isGateCall(line) {
		return v.parseGateCall(line, lineNum)
	}

	return nil
}

// parseQuantumDeclaration parses quantum variable declarations
func (v *ASTBuilderVisitor) parseQuantumDeclaration(line string, lineNum int) *QuantumDeclaration {
	// Remove semicolon and parse: qubit[2] q or qubit q
	line = strings.TrimSuffix(strings.TrimSpace(line), ";")
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return nil
	}

	var size Expression
	var identifier string
	declType := "qubit"

	// Parse type with optional array size
	typePart := parts[0]
	if strings.Contains(typePart, "[") {
		// Type with array size like "qubit[2]"
		openBracket := strings.Index(typePart, "[")
		closeBracket := strings.Index(typePart, "]")
		if openBracket >= 0 && closeBracket > openBracket {
			declType = typePart[:openBracket]
			sizeStr := typePart[openBracket+1 : closeBracket]
			if val := parseInt(sizeStr); val > 0 {
				size = &IntegerLiteral{Value: val}
			}
		}
		// Get identifier from remaining parts
		if len(parts) >= 2 {
			identifier = strings.TrimSpace(parts[1])
		}
	} else {
		// Simple declaration like "qubit q" or "qubit unused_q"
		declType = parts[0]
		if len(parts) >= 2 {
			identifier = strings.TrimSpace(parts[1])
		}
	}

	return &QuantumDeclaration{
		BaseNode:   BaseNode{Position: Position{Line: lineNum, Column: 1}},
		Type:       declType,
		Size:       size,
		Identifier: identifier,
	}
}

// parseClassicalDeclaration parses classical variable declarations
func (v *ASTBuilderVisitor) parseClassicalDeclaration(line string, lineNum int) *ClassicalDeclaration {
	// Parse: bit[2]c=0 or int[32] x = 5+3*2 (handle cases without spaces)

	var size Expression
	var identifier string
	var initializer Expression
	declType := "bit" // default

	// First, handle assignment if present
	assignParts := strings.SplitN(line, "=", 2)
	declarationPart := strings.TrimSpace(assignParts[0])
	if len(assignParts) > 1 {
		initStr := strings.TrimSpace(assignParts[1])
		initializer = v.parseExpression(initStr)
	}

	// Parse the declaration part: bit[2]c or qubit[2] q
	if strings.Contains(declarationPart, "[") {
		// Type with array size like "bit[2]c" or "int[32]x"
		openBracket := strings.Index(declarationPart, "[")
		closeBracket := strings.Index(declarationPart, "]")
		if openBracket >= 0 && closeBracket > openBracket {
			declType = declarationPart[:openBracket]
			sizeStr := declarationPart[openBracket+1 : closeBracket]
			if val := parseInt(sizeStr); val > 0 {
				size = &IntegerLiteral{Value: val}
			}
			// Everything after ] is the identifier
			if closeBracket+1 < len(declarationPart) {
				identifier = strings.TrimSpace(declarationPart[closeBracket+1:])
			}
		}
	} else {
		// Simple type: try to split on spaces first, then parse as single token
		parts := strings.Fields(declarationPart)
		if len(parts) >= 2 {
			declType = parts[0]
			identifier = parts[1]
		} else if len(parts) == 1 {
			// Single token like "bitc" - try to extract type and identifier
			token := parts[0]
			if strings.HasPrefix(token, "bit") && len(token) > 3 {
				declType = "bit"
				identifier = token[3:]
			} else if strings.HasPrefix(token, "int") && len(token) > 3 {
				declType = "int"
				identifier = token[3:]
			} else {
				declType = token
			}
		}
	}

	return &ClassicalDeclaration{
		BaseNode:    BaseNode{Position: Position{Line: lineNum, Column: 1}},
		Type:        declType,
		Size:        size,
		Identifier:  identifier,
		Initializer: initializer,
	}
}

// parseExpression parses a simple expression string
func (v *ASTBuilderVisitor) parseExpression(expr string) Expression {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return nil
	}

	// Handle binary expressions with + and * (simple precedence)
	if strings.Contains(expr, "+") {
		parts := strings.Split(expr, "+")
		if len(parts) == 2 {
			left := v.parseExpression(strings.TrimSpace(parts[0]))
			right := v.parseExpression(strings.TrimSpace(parts[1]))
			return &BinaryExpression{
				Left:     left,
				Operator: "+",
				Right:    right,
			}
		}
	}

	if strings.Contains(expr, "*") {
		parts := strings.Split(expr, "*")
		if len(parts) == 2 {
			left := v.parseExpression(strings.TrimSpace(parts[0]))
			right := v.parseExpression(strings.TrimSpace(parts[1]))
			return &BinaryExpression{
				Left:     left,
				Operator: "*",
				Right:    right,
			}
		}
	}

	// Handle integer literals
	if val := parseInt(expr); val != 0 || expr == "0" {
		return &IntegerLiteral{Value: val}
	}

	// Handle identifiers
	return &Identifier{Name: expr}
}

// parseInt parses a string to int64
func parseInt(s string) int64 {
	var result int64
	for _, r := range s {
		if r >= '0' && r <= '9' {
			result = result*10 + int64(r-'0')
		} else {
			return 0
		}
	}
	return result
}

// parseIncludeStatement parses include statements
func (v *ASTBuilderVisitor) parseIncludeStatement(text string, lineNum int) Statement {
	// Extract the file path from include "path"
	text = strings.TrimSpace(text)
	text = strings.TrimSuffix(text, ";")

	if strings.HasPrefix(text, "include") {
		pathPart := strings.TrimSpace(text[7:]) // Remove "include"
		if len(pathPart) >= 2 && pathPart[0] == '"' && pathPart[len(pathPart)-1] == '"' {
			path := pathPart[1 : len(pathPart)-1] // Remove quotes
			return &Include{
				BaseNode: BaseNode{Position: Position{Line: lineNum, Column: 1}},
				Path:     path,
			}
		}
	}

	return nil
}

// parseMeasurementStatement parses measurement statements
func (v *ASTBuilderVisitor) parseMeasurementStatement(text string, lineNum int) Statement {
	// Parse "measure qubit -> classical" or "measure qubit"
	text = strings.TrimSpace(text)
	text = strings.TrimSuffix(text, ";")

	if strings.HasPrefix(text, "measure") {
		measurePart := strings.TrimSpace(text[7:]) // Remove "measure"

		measurement := &Measurement{
			BaseNode: BaseNode{Position: Position{Line: lineNum, Column: 1}},
		}

		// Check for -> notation
		if strings.Contains(measurePart, "->") {
			parts := strings.Split(measurePart, "->")
			if len(parts) == 2 {
				qubitPart := strings.TrimSpace(parts[0])
				targetPart := strings.TrimSpace(parts[1])

				measurement.Qubit = v.parseQubitExpression(qubitPart)
				measurement.Target = v.parseQubitExpression(targetPart)
			}
		} else {
			measurement.Qubit = v.parseQubitExpression(measurePart)
		}

		return measurement
	}

	return nil
}

// parseGateCall parses gate call statements
func (v *ASTBuilderVisitor) parseGateCall(text string, lineNum int) Statement {
	text = strings.TrimSpace(text)
	text = strings.TrimSuffix(text, ";")

	tokens := strings.Fields(text)
	if len(tokens) < 2 {
		return nil
	}

	gateCall := &GateCall{
		BaseNode:   BaseNode{Position: Position{Line: lineNum, Column: 1}},
		Name:       tokens[0],
		Parameters: make([]Expression, 0),
		Qubits:     make([]Expression, 0),
	}

	// Parse qubits (simple approach - split by comma and space)
	qubitsPart := strings.Join(tokens[1:], " ")
	qubits := strings.Split(qubitsPart, ",")

	for _, qubit := range qubits {
		qubit = strings.TrimSpace(qubit)
		if qubit != "" {
			gateCall.Qubits = append(gateCall.Qubits, v.parseQubitExpression(qubit))
		}
	}

	return gateCall
}

// parseQubitExpression parses qubit expressions like q, q[0], etc.
func (v *ASTBuilderVisitor) parseQubitExpression(text string) Expression {
	text = strings.TrimSpace(text)

	// Check for array access q[index]
	if strings.Contains(text, "[") && strings.Contains(text, "]") {
		bracketStart := strings.Index(text, "[")
		bracketEnd := strings.Index(text, "]")

		if bracketStart > 0 && bracketEnd > bracketStart {
			name := text[:bracketStart]
			indexStr := text[bracketStart+1 : bracketEnd]

			return &IndexedIdentifier{
				BaseNode: BaseNode{Position: Position{Line: 1, Column: 1}},
				Name:     name,
				Index:    v.parseExpression(indexStr),
			}
		}
	}

	// Simple identifier
	return &Identifier{
		BaseNode: BaseNode{Position: Position{Line: 1, Column: 1}},
		Name:     text,
	}
}
