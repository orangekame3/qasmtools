package formatter

import (
	"strconv"
	"strings"

	"github.com/orangekame3/qasmtools/parser"
)

// shouldUseASTFormatting determines whether to use AST-based or text-based formatting
func (f *Formatter) shouldUseASTFormatting(result *parser.ParseResult, content string) bool {
	// Use AST formatting when:
	// 1. Parse was successful with no errors
	// 2. All statements were parsed correctly
	// 3. No comments present (comments require text-based handling)
	// 4. No incomplete parsing indicators

	if result.HasErrors() || result.Program == nil {
		return false
	}

	// Check for comments
	hasComments := strings.Contains(content, "//") || strings.Contains(content, "/*")
	if hasComments {
		return false
	}

	// Count non-empty lines to detect parsing completeness
	inputLines := strings.Split(strings.TrimSpace(content), "\n")
	nonEmptyLines := 0
	for _, line := range inputLines {
		if strings.TrimSpace(line) != "" {
			nonEmptyLines++
		}
	}

	// Account for version statement not being counted in Program.Statements
	expectedStatements := nonEmptyLines
	if result.Program.Version != nil {
		expectedStatements--
	}

	// Check if all statements were parsed
	if len(result.Program.Statements) < expectedStatements && expectedStatements > 1 {
		return false
	}

	// Check for incomplete parsing (lost initializers, parameters, etc.)
	if f.hasIncompleteParsingIndicators(content, result.Program) {
		return false
	}

	// All conditions met - use AST formatting
	return true
}

// formatWithAST provides enhanced AST-based formatting
func (f *Formatter) formatWithAST(program *parser.Program) string {
	var lines []string
	var lastStatementType string

	// Add version if present
	if program.Version != nil {
		lines = append(lines, "OPENQASM "+program.Version.Number+";")
		lastStatementType = "version"
	}

	// Format statements with enhanced AST formatting
	for _, stmt := range program.Statements {
		formatted := f.formatStatementWithAST(stmt, 0)
		if strings.TrimSpace(formatted) != "" {
			currentType := f.getStatementTypeFromStmt(stmt)

			// Add empty line between different types of statements (AST-based logic)
			if f.shouldAddEmptyLineAST(lastStatementType, currentType) {
				if len(lines) == 0 || lines[len(lines)-1] != "" {
					lines = append(lines, "")
				}
			}

			lines = append(lines, formatted)
			lastStatementType = currentType
		}
	}

	return f.joinWithNewline(lines)
}

// formatStatementWithAST provides enhanced AST-based statement formatting
func (f *Formatter) formatStatementWithAST(stmt parser.Statement, indent int) string {
	// Enhanced AST-based formatting with better spacing and structure
	switch s := stmt.(type) {
	case *parser.QuantumDeclaration:
		return f.formatQuantumDeclarationAST(s, indent)
	case *parser.ClassicalDeclaration:
		return f.formatClassicalDeclarationAST(s, indent)
	case *parser.GateCall:
		return f.formatGateCallAST(s, indent)
	case *parser.Measurement:
		return f.formatMeasurementAST(s, indent)
	case *parser.Include:
		return f.formatIncludeAST(s, indent)
	case *parser.GateDefinition:
		return f.formatGateDefinitionAST(s, indent)
	case *parser.IfStatement:
		return f.formatIfStatementAST(s, indent)
	default:
		// Fallback to existing formatting
		return f.formatStatementContent(stmt, indent, nil)
	}
}

// shouldAddEmptyLineAST determines when to add empty lines between statements (AST-based logic)
func (f *Formatter) shouldAddEmptyLineAST(lastType, currentType string) bool {
	// Add empty lines between different statement types for better readability
	if lastType == "" {
		return false
	}

	// Add empty line after version statement (except when followed by include)
	if lastType == "version" && currentType != "include" {
		return true
	}

	// Add empty line after includes
	if lastType == "include" && currentType != "include" {
		return true
	}

	// Add empty line before gate definitions
	if currentType == "gate_definition" && lastType != "gate_definition" {
		return true
	}

	// Add empty line after gate definitions
	if lastType == "gate_definition" && currentType != "gate_definition" {
		return true
	}

	return false
}

// joinWithNewline joins lines and handles final newline based on formatter settings
func (f *Formatter) joinWithNewline(lines []string) string {
	result := strings.Join(lines, "\n")
	if f.newline && len(lines) > 0 && !strings.HasSuffix(result, "\n") {
		result += "\n"
	}
	return result
}

// Enhanced AST-based formatting methods with improved regex-free approach

// formatQuantumDeclarationAST formats quantum declarations using pure AST approach
func (f *Formatter) formatQuantumDeclarationAST(stmt *parser.QuantumDeclaration, indent int) string {
	result := f.indent(indent) + stmt.Type

	if stmt.Size != nil {
		sizeStr := f.formatExpressionAST(stmt.Size)
		result += "[" + sizeStr + "]"
	}

	result += " " + stmt.Identifier + ";"
	return result
}

// formatClassicalDeclarationAST formats classical declarations using pure AST approach
func (f *Formatter) formatClassicalDeclarationAST(stmt *parser.ClassicalDeclaration, indent int) string {
	result := f.indent(indent) + stmt.Type

	if stmt.Size != nil {
		sizeStr := f.formatExpressionAST(stmt.Size)
		result += "[" + sizeStr + "]"
	}

	result += " " + stmt.Identifier

	if stmt.Initializer != nil {
		result += " = " + f.formatExpressionAST(stmt.Initializer)
	}

	result += ";"
	return result
}

// formatGateCallAST formats gate calls using pure AST approach
func (f *Formatter) formatGateCallAST(stmt *parser.GateCall, indent int) string {
	result := f.indent(indent) + stmt.Name

	// Add parameters if present (with proper spacing)
	if len(stmt.Parameters) > 0 {
		params := make([]string, len(stmt.Parameters))
		for i, param := range stmt.Parameters {
			params[i] = f.formatExpressionAST(param)
		}
		result += "(" + strings.Join(params, ", ") + ")"
	}

	// Add qubits (with proper spacing)
	if len(stmt.Qubits) > 0 {
		qubits := make([]string, len(stmt.Qubits))
		for i, qubit := range stmt.Qubits {
			qubits[i] = f.formatExpressionAST(qubit)
		}
		result += " " + strings.Join(qubits, ", ")
	}

	result += ";"
	return result
}

// formatMeasurementAST formats measurements using pure AST approach
func (f *Formatter) formatMeasurementAST(stmt *parser.Measurement, indent int) string {
	result := f.indent(indent) + "measure " + f.formatExpressionAST(stmt.Qubit)

	if stmt.Target != nil {
		result += " -> " + f.formatExpressionAST(stmt.Target)
	}

	result += ";"
	return result
}

// formatIncludeAST formats include statements using pure AST approach
func (f *Formatter) formatIncludeAST(stmt *parser.Include, indent int) string {
	return f.indent(indent) + "include \"" + stmt.Path + "\";"
}

// formatGateDefinitionAST formats gate definitions using pure AST approach
func (f *Formatter) formatGateDefinitionAST(stmt *parser.GateDefinition, indent int) string {
	result := f.indent(indent) + "gate " + stmt.Name

	// Add parameters if present
	if len(stmt.Parameters) > 0 {
		params := make([]string, len(stmt.Parameters))
		for i, param := range stmt.Parameters {
			params[i] = param.Name
		}
		result += "(" + strings.Join(params, ", ") + ")"
	}

	// Add qubits
	if len(stmt.Qubits) > 0 {
		qubits := make([]string, len(stmt.Qubits))
		for i, qubit := range stmt.Qubits {
			qubits[i] = qubit.Name
		}
		result += " " + strings.Join(qubits, ", ")
	}

	result += " {\n"

	// Format body statements
	for _, bodyStmt := range stmt.Body {
		formatted := f.formatStatementWithAST(bodyStmt, indent+1)
		if strings.TrimSpace(formatted) != "" {
			result += formatted + "\n"
		}
	}

	result += f.indent(indent) + "}"
	return result
}

// formatIfStatementAST formats if statements using pure AST approach
func (f *Formatter) formatIfStatementAST(stmt *parser.IfStatement, indent int) string {
	result := f.indent(indent) + "if (" + f.formatExpressionAST(stmt.Condition) + ") {\n"

	// Format then body
	for _, thenStmt := range stmt.ThenBody {
		formatted := f.formatStatementWithAST(thenStmt, indent+1)
		if strings.TrimSpace(formatted) != "" {
			result += formatted + "\n"
		}
	}

	result += f.indent(indent) + "}"

	// Format else body if present
	if len(stmt.ElseBody) > 0 {
		result += " else {\n"
		for _, elseStmt := range stmt.ElseBody {
			formatted := f.formatStatementWithAST(elseStmt, indent+1)
			if strings.TrimSpace(formatted) != "" {
				result += formatted + "\n"
			}
		}
		result += f.indent(indent) + "}"
	}

	return result
}

// formatExpressionAST provides enhanced expression formatting with proper operator spacing
func (f *Formatter) formatExpressionAST(expr parser.Expression) string {
	if expr == nil {
		return ""
	}

	switch e := expr.(type) {
	case *parser.Identifier:
		return e.Name
	case *parser.IndexedIdentifier:
		return e.Name + "[" + f.formatExpressionAST(e.Index) + "]"
	case *parser.RangedIdentifier:
		return e.Name + "[" + f.formatExpressionAST(e.Start) + ":" + f.formatExpressionAST(e.EndIndex) + "]"
	case *parser.IntegerLiteral:
		return strconv.FormatInt(e.Value, 10)
	case *parser.FloatLiteral:
		// Handle float formatting properly
		return strconv.FormatFloat(e.Value, 'g', -1, 64)
	case *parser.StringLiteral:
		return "\"" + e.Value + "\""
	case *parser.BooleanLiteral:
		return strconv.FormatBool(e.Value)
	case *parser.BinaryExpression:
		left := f.formatExpressionAST(e.Left)
		right := f.formatExpressionAST(e.Right)
		// Enhanced operator spacing - no regex needed!
		return left + " " + e.Operator + " " + right
	case *parser.UnaryExpression:
		return e.Operator + f.formatExpressionAST(e.Operand)
	case *parser.FunctionCall:
		args := make([]string, len(e.Arguments))
		for i, arg := range e.Arguments {
			args[i] = f.formatExpressionAST(arg)
		}
		return e.Name + "(" + strings.Join(args, ", ") + ")"
	case *parser.ParenthesizedExpression:
		return "(" + f.formatExpressionAST(e.Expression) + ")"
	case *parser.TimingExpression:
		// Handle timing expressions like 100ns - no regex needed!
		value := f.formatExpressionAST(e.Value)
		unit := strings.TrimSpace(e.Unit)
		return value + unit
	case *parser.DelayExpression:
		// Handle delay expressions like delay[100ns] - no regex needed!
		timing := f.formatExpressionAST(e.Timing)
		return "delay[" + timing + "]"
	default:
		// Fallback to original formatExpression for compatibility
		return f.formatExpression(expr)
	}
}
