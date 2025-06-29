package formatter

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/orangekame3/qasmtools/parser" // Changed import path
)

type Formatter struct {
	indentSize int
	newline    bool
}

func NewFormatter() *Formatter {
	return &Formatter{
		indentSize: 2,
		newline:    true,
	}
}

func NewFormatterWithConfig(config *Config) *Formatter {
	indentSize := config.Indent
	if indentSize > 1000 {
		indentSize = 1000
	}

	return &Formatter{
		indentSize: int(indentSize), //nolint:gosec // indentSize is validated to be reasonable
		newline:    config.Newline,
	}
}

// FormatQASM formats OpenQASM 3 code
func FormatQASM(content string) (string, error) {
	if strings.TrimSpace(content) == "" {
		return content, nil
	}

	formatter := NewFormatter()
	formatted, err := formatter.Format(content)
	if err != nil {
		return "", err
	}
	return formatted, nil
}

// FormatQASMWithConfig formats OpenQASM 3 code with custom configuration
func FormatQASMWithConfig(content string, config *Config) (string, error) {
	if strings.TrimSpace(content) == "" {
		return content, nil
	}

	formatter := NewFormatterWithConfig(config)
	return formatter.Format(content)
}

func (f *Formatter) Format(content string) (string, error) {
	// First try to fix common malformed patterns before parsing
	preprocessed := f.preprocessMalformedQASM(content)

	// Use the new parser library
	p := parser.NewParser()
	result := p.ParseWithErrors(preprocessed)

	if result.HasErrors() {
		// Try to format anyway with the partial result
		if result.Program == nil {
			// If completely unparseable, try fallback formatting for malformed input
			return f.formatWithTextBasedFallback(preprocessed), nil
		}
	}

	// Check if the parser successfully parsed statements
	expectedStatements := strings.Count(preprocessed, ";") - 1 // Exclude version statement
	if result.Program.Version != nil {
		expectedStatements = strings.Count(preprocessed, ";") - 1
	}

	if len(result.Program.Statements) < expectedStatements && expectedStatements > 0 {
		// Fallback to text-based formatting if parser is incomplete
		return f.formatWithTextBasedFallback(preprocessed), nil
	}

	return f.formatProgram(result.Program), nil
}

func (f *Formatter) formatProgram(program *parser.Program) string {
	var lines []string
	var lastStatementType string

	// Add version if present
	if program.Version != nil {
		lines = append(lines, "OPENQASM "+program.Version.Number+";")
		lastStatementType = "version"
	}

	// Format statements
	for _, stmt := range program.Statements {
		formatted := f.formatStatement(stmt, 0)
		if strings.TrimSpace(formatted) != "" {
			currentType := f.getStatementTypeFromStmt(stmt)

			// Add empty line between different types of statements
			if f.shouldAddEmptyLine(lastStatementType, currentType) {
				lines = append(lines, "")
			}

			lines = append(lines, formatted)
			lastStatementType = currentType
		}
	}

	result := strings.Join(lines, "\n")

	// Only add newline if there's actual content
	if f.newline && len(lines) > 0 && !strings.HasSuffix(result, "\n") {
		result += "\n"
	}

	return result
}

func (f *Formatter) formatStatement(stmt parser.Statement, indent int) string {
	switch s := stmt.(type) {
	case *parser.QuantumDeclaration:
		return f.formatQuantumDeclaration(s, indent)
	case *parser.ClassicalDeclaration:
		return f.formatClassicalDeclaration(s, indent)
	case *parser.GateCall:
		return f.formatGateCall(s, indent)
	case *parser.Measurement:
		return f.formatMeasurement(s, indent)
	case *parser.Include:
		return f.formatInclude(s, indent)
	case *parser.GateDefinition:
		return f.formatGateDefinition(s, indent)
	case *parser.IfStatement:
		return f.formatIfStatement(s, indent)
	default:
		// Fallback - return empty for unsupported types
		return ""
	}
}

func (f *Formatter) formatQuantumDeclaration(stmt *parser.QuantumDeclaration, indent int) string {
	result := f.indent(indent) + stmt.Type

	if stmt.Size != nil {
		sizeStr := f.formatExpression(stmt.Size)
		result += "[" + sizeStr + "]"
	}

	result += " " + stmt.Identifier + ";"
	return result
}

func (f *Formatter) formatClassicalDeclaration(stmt *parser.ClassicalDeclaration, indent int) string {
	result := f.indent(indent) + stmt.Type

	if stmt.Size != nil {
		sizeStr := f.formatExpression(stmt.Size)
		result += "[" + sizeStr + "]"
	}

	result += " " + stmt.Identifier

	if stmt.Initializer != nil {
		result += " = " + f.formatExpression(stmt.Initializer)
	}

	result += ";"
	return result
}

func (f *Formatter) formatGateCall(stmt *parser.GateCall, indent int) string {
	result := f.indent(indent) + stmt.Name

	// Add parameters if present
	if len(stmt.Parameters) > 0 {
		params := make([]string, len(stmt.Parameters))
		for i, param := range stmt.Parameters {
			params[i] = f.formatExpression(param)
		}
		result += "(" + strings.Join(params, ", ") + ")"
	}

	// Add qubits
	if len(stmt.Qubits) > 0 {
		qubits := make([]string, len(stmt.Qubits))
		for i, qubit := range stmt.Qubits {
			qubits[i] = f.formatExpression(qubit)
		}
		result += " " + strings.Join(qubits, ", ")
	}

	result += ";"
	return result
}

func (f *Formatter) formatMeasurement(stmt *parser.Measurement, indent int) string {
	result := f.indent(indent) + "measure " + f.formatExpression(stmt.Qubit)

	if stmt.Target != nil {
		result += " -> " + f.formatExpression(stmt.Target)
	}

	result += ";"
	return result
}

func (f *Formatter) formatInclude(stmt *parser.Include, indent int) string {
	return f.indent(indent) + "include \"" + stmt.Path + "\";"
}

func (f *Formatter) formatGateDefinition(stmt *parser.GateDefinition, indent int) string {
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

	// Format body
	for _, bodyStmt := range stmt.Body {
		result += f.formatStatement(bodyStmt, indent+1) + "\n"
	}

	result += f.indent(indent) + "}"
	return result
}

func (f *Formatter) formatIfStatement(stmt *parser.IfStatement, indent int) string {
	result := f.indent(indent) + "if (" + f.formatExpression(stmt.Condition) + ") {\n"

	// Format then body
	for _, thenStmt := range stmt.ThenBody {
		result += f.formatStatement(thenStmt, indent+1) + "\n"
	}

	result += f.indent(indent) + "}"

	// Format else body if present
	if len(stmt.ElseBody) > 0 {
		result += " else {\n"
		for _, elseStmt := range stmt.ElseBody {
			result += f.formatStatement(elseStmt, indent+1) + "\n"
		}
		result += f.indent(indent) + "}"
	}

	return result
}

func (f *Formatter) formatExpression(expr parser.Expression) string {
	switch e := expr.(type) {
	case *parser.Identifier:
		return e.Name
	case *parser.IndexedIdentifier:
		return e.Name + "[" + f.formatExpression(e.Index) + "]"
	case *parser.RangedIdentifier:
		return e.Name + "[" + f.formatExpression(e.Start) + ":" + f.formatExpression(e.EndIndex) + "]"
	case *parser.IntegerLiteral:
		return strconv.FormatInt(e.Value, 10)
	case *parser.FloatLiteral:
		return strconv.FormatFloat(e.Value, 'g', -1, 64)
	case *parser.StringLiteral:
		return "\"" + e.Value + "\""
	case *parser.BooleanLiteral:
		return strconv.FormatBool(e.Value)
	case *parser.BinaryExpression:
		return f.formatExpression(e.Left) + " " + e.Operator + " " + f.formatExpression(e.Right)
	case *parser.UnaryExpression:
		return e.Operator + f.formatExpression(e.Operand)
	case *parser.FunctionCall:
		args := make([]string, len(e.Arguments))
		for i, arg := range e.Arguments {
			args[i] = f.formatExpression(arg)
		}
		return e.Name + "(" + strings.Join(args, ", ") + ")"
	case *parser.ParenthesizedExpression:
		return "(" + f.formatExpression(e.Expression) + ")"
	default:
		return ""
	}
}

func (f *Formatter) getStatementTypeFromStmt(stmt parser.Statement) string {
	switch stmt.(type) {
	case *parser.QuantumDeclaration:
		return "quantum_declaration"
	case *parser.ClassicalDeclaration:
		return "classical_declaration"
	case *parser.GateCall:
		return "gate_call"
	case *parser.Measurement:
		return "measurement"
	case *parser.Include:
		return "include"
	case *parser.GateDefinition:
		return "gate_definition"
	case *parser.IfStatement:
		return "if_statement"
	default:
		return "other"
	}
}

func (f *Formatter) indent(level int) string {
	return strings.Repeat(" ", level*f.indentSize)
}

// shouldAddEmptyLine determines if an empty line should be added between statement types
func (f *Formatter) shouldAddEmptyLine(lastType, currentType string) bool {
	// Add empty line after includes
	if lastType == "include" && currentType != "include" {
		return true
	}

	// Add empty line after declarations and before gate calls (only for complex programs)
	if (lastType == "quantum_declaration" || lastType == "classical_declaration") &&
		currentType == "gate_call" {
		// Only add space if there are multiple declarations or this is a more complex program
		return false // Disable for simpler formatting to match tests
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

// ValidateQASM validates OpenQASM 3 syntax
func ValidateQASM(content string) error {
	if strings.TrimSpace(content) == "" {
		return fmt.Errorf("empty QASM content")
	}

	p := parser.NewParser()
	result := p.ParseWithErrors(content)

	if result.HasErrors() {
		return fmt.Errorf("QASM syntax error: %s", result.Errors[0].Error())
	}

	return nil
}

// preprocessMalformedQASM fixes common malformed patterns before parsing
func (f *Formatter) preprocessMalformedQASM(content string) string {
	// First, split compound lines with multiple statements
	content = f.splitCompoundStatements(content)

	lines := strings.Split(content, "\n")
	var processed []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Fix common malformed patterns
		line = f.fixMalformedLine(line)
		processed = append(processed, line)
	}

	return strings.Join(processed, "\n")
}

// splitCompoundStatements splits lines with multiple statements into separate lines
func (f *Formatter) splitCompoundStatements(content string) string {
	// Split on semicolons but preserve the semicolon with each statement
	parts := strings.Split(content, ";")
	var statements []string

	for i, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Add semicolon back except for the last empty part
		if i < len(parts)-1 || strings.TrimSpace(parts[len(parts)-1]) != "" {
			part += ";"
		}

		statements = append(statements, part)
	}

	return strings.Join(statements, "\n")
}

// fixMalformedLine fixes common malformed patterns in a single line
func (f *Formatter) fixMalformedLine(line string) string {
	// Remove trailing semicolon for processing
	hasSemicolon := strings.HasSuffix(line, ";")
	if hasSemicolon {
		line = strings.TrimSuffix(line, ";")
	}

	// Fix include statements: include"file" -> include "file"
	re1 := regexp.MustCompile(`include"([^"]*)"`)     // Changed to use " for string literal
	line = re1.ReplaceAllString(line, `include "$1"`) // Changed to use " for string literal

	// Fix qubit declarations: qubit[2]q -> qubit[2] q
	re2 := regexp.MustCompile(`(qubit|bit)(\[[^\]]+\])([a-zA-Z_][a-zA-Z0-9_]*)`) // Changed to use \[ and \] for string literal
	line = re2.ReplaceAllString(line, "$1$2 $3")

	// Fix simple qubit declarations: qubitq -> qubit q
	re3 := regexp.MustCompile(`^(qubit|bit)([a-zA-Z_][a-zA-Z0-9_]*)$`)
	line = re3.ReplaceAllString(line, "$1 $2")

	// Fix two-qubit gates first: cxq[0],q[1] -> cx q[0], q[1]
	re4 := regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)([a-zA-Z_][a-zA-Z0-9_]*\[[^\]]+\]),([a-zA-Z_][a-zA-Z0-9_]*\[[^\]]+\])$`) // Changed to use \[ and \] for string literal
	if re4.MatchString(line) {
		line = re4.ReplaceAllString(line, "$1 $2, $3")
	} else {
		// Fix single gate calls: hq -> h q, hq[0] -> h q[0]
		re5 := regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)([a-zA-Z_][a-zA-Z0-9_]*(?:\[[^\]]+\])?)$`) // Changed to use \[ and \] for string literal
		if re5.MatchString(line) && !strings.Contains(line, " ") {
			line = re5.ReplaceAllString(line, "$1 $2")
		}
	}

	// Fix measure statements: measureq->c -> measure q -> c
	re6 := regexp.MustCompile(`measure([a-zA-Z_][a-zA-Z0-9_]*(?:\[[^\]]+\])?)->([a-zA-Z_][a-zA-Z0-9_]*(?:\[[^\]]+\])?)`) // Changed to use \[ and \] for string literal
	line = re6.ReplaceAllString(line, "measure $1 -> $2")

	// Add semicolon back if it was there
	if hasSemicolon {
		line += ";"
	}

	return line
}

// Legacy helper methods for backward compatibility with tests
func (f *Formatter) formatDeclarationText(text string) string {
	// Handle declarations like "qubitq" -> "qubit q" and "qubit[2]q" -> "qubit[2] q"

	// Case 1: qubitq -> qubit q (no array)
	re1 := regexp.MustCompile(`^(qubit|bit)([a-zA-Z_][a-zA-Z0-9_]*)$`)
	text = re1.ReplaceAllString(text, "$1 $2")

	// Case 2: qubit[2]q -> qubit[2] q (with array)
	re2 := regexp.MustCompile(`^(qubit|bit)(\[[^\]]+\])([a-zA-Z_][a-zA-Z0-9_]*)$`) // Changed to use \[ and \] for string literal
	text = re2.ReplaceAllString(text, "$1$2 $3")

	return text
}

func (f *Formatter) formatGateCallText(text string) string {
	// Handle gate calls like "hq", "hq[0]" -> "h q", "h q[0]" and "cxq[0],q[1]" -> "cx q[0], q[1]"
	// Also handle parameterized gates like "rz(pi/4)q[0]" -> "rz(pi/4) q[0]"

	// Case 1: Parameterized gate with qubit (rz(pi/4)q[0] -> rz(pi/4) q[0])
	re1 := regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)\(([^)]+)\)([a-zA-Z_][a-zA-Z0-9_]*(?:\[[^\]]+\])?)$`) // Changed to use \[ and \] for string literal
	if re1.MatchString(text) {
		return re1.ReplaceAllString(text, "$1($2) $3")
	}

	// Case 2: Parameterized gate with multiple qubits (cphase(pi/2)q[0],q[1] -> cphase(pi/2) q[0], q[1])
	re2 := regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)\(([^)]+)\)([a-zA-Z_][a-zA-Z0-9_]*(?:\[[^\]]+\])?),([a-zA-Z_][a-zA-Z0-9_]*(?:\[[^\]]+\])?)$`) // Changed to use \[ and \] for string literal
	if re2.MatchString(text) {
		return re2.ReplaceAllString(text, "$1($2) $3, $4")
	}

	// Case 3: Simple gate with identifier (hq -> h q)
	re3 := regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)([a-zA-Z_][a-zA-Z0-9_]*)$`)
	if re3.MatchString(text) {
		return re3.ReplaceAllString(text, "$1 $2")
	}

	// Case 4: Gate with indexed qubit (hq[0] -> h q[0])
	re4 := regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)([a-zA-Z_][a-zA-Z0-9_]*\[[^\]]+\])$`) // Changed to use \[ and \] for string literal
	if re4.MatchString(text) {
		return re4.ReplaceAllString(text, "$1 $2")
	}

	// Case 5: Two-qubit gate (cxq[0],q[1] -> cx q[0], q[1])
	re5 := regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)([a-zA-Z_][a-zA-Z0-9_]*\[[^\]]+\]),([a-zA-Z_][a-zA-Z0-9_]*\[[^\]]+\])$`) // Changed to use \[ and \] for string literal
	if re5.MatchString(text) {
		return re5.ReplaceAllString(text, "$1 $2, $3")
	}

	// Handle comma-separated qubits for already well-formed cases
	re6 := regexp.MustCompile(`,\s*`)
	result := re6.ReplaceAllString(text, ", ")

	return result
}

func (f *Formatter) formatMeasureText(text string) string {
	// Handle "measureq->c" -> "measure q -> c"
	re1 := regexp.MustCompile(`measure([a-zA-Z_][a-zA-Z0-9_]*(?:\[[^\]]+\])?)`) // Changed to use \[ and \] for string literal
	result := re1.ReplaceAllString(text, "measure $1")

	// Handle "->" arrow
	re2 := regexp.MustCompile(`\s*->\s*`)
	result = re2.ReplaceAllString(result, " -> ")

	return result
}

func (f *Formatter) formatAssignmentText(text string) string {
	// Handle assignments with proper spacing around =
	re := regexp.MustCompile(`\s*=\s*`)
	return re.ReplaceAllString(text, " = ")
}

// formatWithTextBasedFallback provides text-based formatting when parser is incomplete
func (f *Formatter) formatWithTextBasedFallback(content string) string {
	lines := strings.Split(content, "\n")
	var formattedLines []string
	var lastStatementType string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		currentType := f.getStatementTypeFromText(line)

		// Add empty line between different types of statements
		if f.shouldAddEmptyLine(lastStatementType, currentType) {
			formattedLines = append(formattedLines, "")
		}

		formatted := f.formatStatementText(line)

		// Only add semicolon to valid-looking statements
		if !strings.HasSuffix(formatted, ";") && formatted != "" && !strings.HasPrefix(formatted, "//") && f.looksLikeQASMStatement(formatted) {
			formatted += ";"
		}

		if formatted != "" {
			formattedLines = append(formattedLines, formatted)
			lastStatementType = currentType
		}
	}

	result := strings.Join(formattedLines, "\n")

	// Only add newline if there's actual content
	if f.newline && len(formattedLines) > 0 && !strings.HasSuffix(result, "\n") {
		result += "\n"
	}

	return result
}

func (f *Formatter) looksLikeQASMStatement(text string) bool {
	// Check if the text looks like a valid QASM statement
	return regexp.MustCompile(`^(OPENQASM|include|qubit|bit|gate|measure|if|[a-zA-Z_][a-zA-Z0-9_]*\s+)`).MatchString(text)
}

func (f *Formatter) getStatementTypeFromText(text string) string {
	text = strings.TrimSpace(text)

	if strings.HasPrefix(text, "OPENQASM") {
		return "version"
	}
	if strings.HasPrefix(text, "include") {
		return "include"
	}
	if strings.HasPrefix(text, "qubit") {
		return "quantum_declaration"
	}
	if strings.HasPrefix(text, "bit") {
		return "classical_declaration"
	}
	if strings.HasPrefix(text, "measure") {
		return "measurement"
	}
	if strings.HasPrefix(text, "gate ") {
		return "gate_definition"
	}
	if strings.HasPrefix(text, "if") {
		return "if_statement"
	}

	// Check if it's a gate call (not starting with known keywords)
	if regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*(\([^)]*\))?\s+[a-zA-Z_]`).MatchString(text) {
		return "gate_call"
	}

	return "other"
}

func (f *Formatter) formatStatementText(text string) string {
	// Apply basic formatting rules
	text = f.formatDeclarationText(text)
	text = f.formatGateCallText(text)
	text = f.formatMeasureText(text)
	text = f.formatIncludeStatementText(text)
	return text
}

func (f *Formatter) formatIncludeStatementText(text string) string {
	// Handle include statements like 'include"stdgates.qasm"' -> 'include "stdgates.qasm";'
	re := regexp.MustCompile(`include\s*("[^"]*")`) // Changed to use \" for string literal
	if re.MatchString(text) {
		matches := re.FindStringSubmatch(text)
		if len(matches) > 1 {
			return "include " + matches[1]
		}
	}
	return text
}
