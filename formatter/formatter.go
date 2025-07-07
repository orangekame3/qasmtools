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
	comments   []parser.Comment
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

	// Use the parser library with fallback for compatibility
	p := parser.NewParser()
	result := p.ParseWithErrors(preprocessed)

	if result.HasErrors() {
		// If completely unparseable, use simple text-based formatting
		if result.Program == nil {
			return f.formatWithSimpleFallback(preprocessed), nil
		}
	}

	// Extract comments from the original content
	f.comments = f.extractComments(preprocessed)

	// Check if we have enough statements parsed or use fallback
	hasComments := strings.Contains(preprocessed, "//") || strings.Contains(preprocessed, "/*")

	// Count non-empty lines in input to detect parser issues
	inputLines := strings.Split(strings.TrimSpace(preprocessed), "\n")
	nonEmptyLines := 0
	for _, line := range inputLines {
		if strings.TrimSpace(line) != "" {
			nonEmptyLines++
		}
	}

	// Use fallback if we have fewer parsed statements than non-empty input lines
	// This indicates the parser may have incorrectly combined statements
	// BUT: Account for version statement not being counted in Program.Statements
	expectedStatements := nonEmptyLines
	if result.Program.Version != nil {
		expectedStatements-- // Version is not counted in statements
	}

	if hasComments || len(result.Program.Statements) == 0 ||
		(len(result.Program.Statements) < expectedStatements && expectedStatements > 1) {
		return f.formatWithSimpleFallback(preprocessed), nil
	}

	// Use AST-based formatting for simple cases
	return f.formatProgram(result.Program), nil
}

// extractComments extracts comments from the content and returns them
func (f *Formatter) extractComments(content string) []parser.Comment {
	var comments []parser.Comment
	lines := strings.Split(content, "\n")

	for lineNum, line := range lines {
		// Handle line comments
		if idx := strings.Index(line, "//"); idx != -1 {
			comments = append(comments, parser.Comment{
				BaseNode: parser.BaseNode{
					Position: parser.Position{
						Line:   lineNum + 1,
						Column: idx + 1,
					},
				},
				Text: strings.TrimSpace(line[idx:]),
				Type: "line",
			})
		}

		// Handle block comments (simple implementation)
		if idx := strings.Index(line, "/*"); idx != -1 {
			endIdx := strings.Index(line[idx:], "*/")
			if endIdx != -1 {
				comments = append(comments, parser.Comment{
					BaseNode: parser.BaseNode{
						Position: parser.Position{
							Line:   lineNum + 1,
							Column: idx + 1,
						},
					},
					Text: strings.TrimSpace(line[idx : idx+endIdx+2]),
					Type: "block",
				})
			}
		}
	}

	return comments
}

func (f *Formatter) formatProgram(program *parser.Program) string {
	var lines []string
	var lastStatementType string

	// Add version if present
	if program.Version != nil {
		lines = append(lines, "OPENQASM "+program.Version.Number+";")
		lastStatementType = "version"
	}

	// Format statements with comments integration
	for _, stmt := range program.Statements {
		// Check for comments that belong to this statement line
		stmtLine := stmt.Pos().Line
		var lineComments []parser.Comment
		for _, comment := range f.comments {
			if comment.Pos().Line == stmtLine {
				lineComments = append(lineComments, comment)
			}
		}

		formatted := f.formatStatement(stmt, 0, program)
		if strings.TrimSpace(formatted) != "" {
			currentType := f.getStatementTypeFromStmt(stmt)

			// Add empty line between different types of statements
			// Only add if the last line is not already empty
			if f.shouldAddEmptyLine(lastStatementType, currentType) {
				if len(lines) == 0 || lines[len(lines)-1] != "" {
					lines = append(lines, "")
				}
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

func (f *Formatter) formatStatement(stmt parser.Statement, indent int, program *parser.Program) string {
	// Add leading comments
	var result string
	if comments := parser.GetLeadingComments(program, stmt); len(comments) > 0 {
		for _, comment := range comments {
			result += f.indent(indent) + "//" + comment.Text + "\n"
		}
	}

	// Format the statement
	formatted := f.formatStatementContent(stmt, indent, program)

	// Add trailing comments
	if comments := parser.GetTrailingComments(program, stmt); len(comments) > 0 {
		formatted += " // " + comments[0].Text
	}

	return result + formatted
}

func (f *Formatter) formatStatementContent(stmt parser.Statement, indent int, program *parser.Program) string {
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
		return f.formatGateDefinition(s, indent, program)
	case *parser.IfStatement:
		return f.formatIfStatement(s, indent, program)
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

func (f *Formatter) formatGateDefinition(stmt *parser.GateDefinition, indent int, program *parser.Program) string {
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
		formatted := f.formatStatementContent(bodyStmt, indent+1, program)
		if strings.TrimSpace(formatted) != "" {
			result += f.indent(indent+1) + formatted + "\n"
		}
	}

	result += f.indent(indent) + "}"
	return result
}

func (f *Formatter) formatIfStatement(stmt *parser.IfStatement, indent int, program *parser.Program) string {
	result := f.indent(indent) + "if (" + f.formatExpression(stmt.Condition) + ") {\n"

	// Format then body
	for _, thenStmt := range stmt.ThenBody {
		formatted := f.formatStatementContent(thenStmt, indent+1, program)
		if strings.TrimSpace(formatted) != "" {
			result += f.indent(indent+1) + formatted + "\n"
		}
	}

	result += f.indent(indent) + "}"

	// Format else body if present
	if len(stmt.ElseBody) > 0 {
		result += " else {\n"
		for _, elseStmt := range stmt.ElseBody {
			formatted := f.formatStatementContent(elseStmt, indent+1, program)
			if strings.TrimSpace(formatted) != "" {
				result += f.indent(indent+1) + formatted + "\n"
			}
		}
		result += f.indent(indent) + "}"
	}

	return result
}

func (f *Formatter) formatExpression(expr parser.Expression) string {
	if expr == nil {
		return ""
	}
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
		left := f.formatExpression(e.Left)
		right := f.formatExpression(e.Right)
		// Add spaces around all binary operators according to spec
		return left + " " + e.Operator + " " + right
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
	case *parser.TimingExpression:
		// Handle timing expressions like 100ns
		value := f.formatExpression(e.Value)
		unit := strings.TrimSpace(e.Unit)
		return value + unit
	case *parser.DelayExpression:
		// Handle delay expressions like delay[100ns]
		timing := f.formatExpression(e.Timing)
		return "delay[" + timing + "]"
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

	// According to spec.yaml: Do not insert empty lines between variable declarations unless grouping is intentional
	// Remove the rule that adds empty lines between declarations and gate calls/measurements

	// Add empty line before gate definitions (but not if it's the first statement)
	if currentType == "gate_definition" && lastType != "gate_definition" && lastType != "" {
		return true
	}

	// Add empty line after gate definitions
	if lastType == "gate_definition" && currentType != "gate_definition" {
		return true
	}

	// Add empty line after block endings (like gate definitions) only if there are more statements
	if lastType == "block_end" && currentType != "block_end" && currentType != "" && currentType != "other" {
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

	// Additional validation for unknown/invalid statements
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "/*") {
			continue
		}

		// Check if line matches known QASM patterns
		if !isValidQASMStatement(line) {
			return fmt.Errorf("invalid QASM statement at line %d: %s", i+1, line)
		}
	}

	return nil
}

// isValidQASMStatement checks if a statement follows valid QASM syntax
func isValidQASMStatement(line string) bool {
	line = strings.TrimSpace(line)

	// Check for specific known invalid patterns first
	invalidPatterns := []string{
		`^invalid\s+`, // Lines starting with "invalid"
		`^unknown\s+`, // Lines starting with "unknown"
		`^bad\s+`,     // Lines starting with "bad"
	}

	for _, pattern := range invalidPatterns {
		matched, _ := regexp.MatchString(pattern, line)
		if matched {
			return false
		}
	}

	// Known valid patterns
	validPatterns := []string{
		`^OPENQASM\s+[\d\.]+\s*;$`, // OPENQASM version
		`^include\s+"[^"]+"\s*;$`,  // include statements
		`^(qubit|bit|int|float|bool)(\[[^\]]+\])?\s+[a-zA-Z_][a-zA-Z0-9_]*(\s*=\s*[^;]+)?\s*;$`, // declarations
		`^measure\s+[^;]+\s*;$`, // measure statements
		`^gate\s+[a-zA-Z_][a-zA-Z0-9_]*\s*(\([^)]*\))?\s+[^{]+\s*\{$`, // gate definitions
		`^if\s*\([^)]+\)\s*\{$`,                    // if statements
		`^\}(\s*else\s*\{)?$`,                      // closing braces and else
		`^[a-zA-Z_][a-zA-Z0-9_]*\s*=\s*[^;]+\s*;$`, // assignments
	}

	// Check for known gate names (common quantum gates)
	knownGates := []string{
		"h", "x", "y", "z", "s", "t", "cx", "cy", "cz", "ccx", "rx", "ry", "rz",
		"p", "cp", "u", "u1", "u2", "u3", "swap", "iswap", "cswap", "toffoli",
		"fredkin", "rxx", "ryy", "rzz", "cphase", "crx", "cry", "crz", "cu", "cu1", "cu2", "cu3",
	}

	// Check if it's a gate call
	for _, gate := range knownGates {
		gatePattern := fmt.Sprintf(`^%s(\([^)]*\))?\s+[a-zA-Z_][a-zA-Z0-9_]*(\[[^\]]*\])?(\s*,\s*[a-zA-Z_][a-zA-Z0-9_]*(\[[^\]]*\])?)*\s*;$`, gate)
		matched, _ := regexp.MatchString(gatePattern, line)
		if matched {
			return true
		}
	}

	for _, pattern := range validPatterns {
		matched, _ := regexp.MatchString(pattern, line)
		if matched {
			return true
		}
	}

	return false
}

// preprocessMalformedQASM fixes common malformed patterns before parsing
func (f *Formatter) preprocessMalformedQASM(content string) string {
	// Normalize line endings
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")

	// Split compound statements (multiple statements on one line) into separate lines first
	content = f.splitCompoundStatements(content)

	// Split into lines
	lines := strings.Split(content, "\n")
	var processed []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			processed = append(processed, "") // Preserve empty lines
			continue
		}

		// Check if this is a comment line that should preserve spacing
		if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "/*") || strings.Contains(trimmed, "*/") {
			// Preserve original spacing for comment lines
			line = f.fixMalformedLine(line)
		} else {
			// Trim and process non-comment lines
			line = f.fixMalformedLine(trimmed)
		}
		processed = append(processed, line)
	}

	// Join lines and ensure trailing newline
	result := strings.Join(processed, "\n")
	if !strings.HasSuffix(result, "\n") {
		result += "\n"
	}

	return result
}

// splitCompoundStatements splits lines with multiple statements into separate lines
func (f *Formatter) splitCompoundStatements(content string) string {
	lines := strings.Split(content, "\n")
	var processedLines []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			processedLines = append(processedLines, line)
			continue
		}

		// Skip lines that are only comments
		if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "/*") {
			processedLines = append(processedLines, line)
			continue
		}

		// Count semicolons to determine if this line has multiple statements
		semicolonCount := strings.Count(line, ";")

		// If there's only 0 or 1 semicolon, keep the line as-is
		if semicolonCount <= 1 {
			processedLines = append(processedLines, line)
			continue
		}

		// Split compound statements on lines with multiple semicolons
		var statements []string
		parts := strings.Split(line, ";")

		for i, part := range parts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}

			// Check if this part starts with a comment - if so, attach it to the previous statement
			if strings.HasPrefix(part, "//") && len(statements) > 0 {
				// This is a trailing comment, attach it to the previous statement
				statements[len(statements)-1] += "; " + part
			} else {
				// This is a regular statement, add semicolon back if not the last empty part
				if i < len(parts)-1 || strings.TrimSpace(parts[len(parts)-1]) != "" {
					part += ";"
				}
				statements = append(statements, part)
			}
		}

		processedLines = append(processedLines, statements...)
	}

	return strings.Join(processedLines, "\n")
}

// fixMalformedLine fixes common malformed patterns in a single line
func (f *Formatter) fixMalformedLine(line string) string {
	// Handle comments - preserve all comment lines without processing
	if strings.HasPrefix(line, "//") {
		return line
	}
	if strings.HasPrefix(line, "/*") {
		return line // This includes both single-line and multi-line comment starts
	}
	if strings.Contains(line, "*/") {
		return line // This handles multi-line comment ends and lines inside multi-line comments
	}

	// Remove trailing semicolon for processing
	hasSemicolon := strings.HasSuffix(line, ";")
	if hasSemicolon {
		line = strings.TrimSuffix(line, ";")
	}

	// Fix include statements: include"file" -> include "file"
	re1 := regexp.MustCompile(`include"([^"]*)"`)
	line = re1.ReplaceAllString(line, `include "$1"`)

	// Fix classical/quantum declarations with array sizes and assignments
	// bit[2]c=0 -> bit[2] c = 0
	// int[32]x=5+3*2 -> int[32] x = 5 + 3 * 2
	re2 := regexp.MustCompile(`^((?:bit|int|qubit)(?:\[[^\]]+\])?)([a-zA-Z_][a-zA-Z0-9_]*)=(.+)$`)
	if re2.MatchString(line) {
		line = re2.ReplaceAllString(line, "$1 $2 = $3")
	}

	// Fix declarations without assignment: bit[2]c -> bit[2] c
	re3 := regexp.MustCompile(`^((?:bit|int|qubit)(?:\[[^\]]+\])?)([a-zA-Z_][a-zA-Z0-9_]*)$`)
	if re3.MatchString(line) {
		line = re3.ReplaceAllString(line, "$1 $2")
	}

	// Fix simple declarations: bitc -> bit c
	re4 := regexp.MustCompile(`^(bit|int|qubit)([a-zA-Z_][a-zA-Z0-9_]*)$`)
	if re4.MatchString(line) {
		line = re4.ReplaceAllString(line, "$1 $2")
	}

	// Fix array indices early - remove spaces inside brackets before operator processing
	re_brackets := regexp.MustCompile(`\[\s*([^\]]+?)\s*\]`)
	line = re_brackets.ReplaceAllString(line, "[$1]")

	// Fix binary operators with simplified rules
	// Handle specific common patterns directly

	// Fix comparison operators in if statements
	if strings.Contains(line, "if") {
		// Use word boundaries and check for existing spacing to avoid duplicating spaces
		line = regexp.MustCompile(`\s*==\s*`).ReplaceAllString(line, " == ")
		line = regexp.MustCompile(`\s*!=\s*`).ReplaceAllString(line, " != ")
		line = regexp.MustCompile(`\s*<=\s*`).ReplaceAllString(line, " <= ")
		line = regexp.MustCompile(`\s*>=\s*`).ReplaceAllString(line, " >= ")
	}

	// Skip assignment operator processing for now
	// if strings.Contains(line, "=") && !strings.Contains(line, "==") && !strings.Contains(line, "!=") && !strings.Contains(line, "if") {
	//	line = regexp.MustCompile(`([a-zA-Z0-9_\])])\s*=\s*([a-zA-Z0-9_\[\(])`).ReplaceAllString(line, "$1 = $2")
	// }

	// Fix arithmetic operators (but preserve function parameters like pi/2)
	if !strings.Contains(line, "(") {
		line = regexp.MustCompile(`([a-zA-Z0-9_\])])\s*\+\s*([a-zA-Z0-9_\[\(])`).ReplaceAllString(line, "$1 + $2")
		line = regexp.MustCompile(`([a-zA-Z0-9_\])])\s*-\s*([a-zA-Z0-9_\[\(])`).ReplaceAllString(line, "$1 - $2")
		line = regexp.MustCompile(`([a-zA-Z0-9_\])])\s*\*\s*([a-zA-Z0-9_\[\(])`).ReplaceAllString(line, "$1 * $2")
		line = regexp.MustCompile(`([a-zA-Z0-9_\])])\s*/\s*([a-zA-Z0-9_\[\(])`).ReplaceAllString(line, "$1 / $2")
	}

	// Fix gate calls with parameters
	re6 := regexp.MustCompile(`([a-zA-Z_][a-zA-Z0-9_]*)\(([^)]+)\)([a-zA-Z_][a-zA-Z0-9_]*(?:\[[^\]]+\])?)`)
	if re6.MatchString(line) {
		line = re6.ReplaceAllString(line, "$1($2) $3")
	}

	// Fix two-qubit gates: cxq[0],q[1] -> cx q[0], q[1]
	re7 := regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)([a-zA-Z_][a-zA-Z0-9_]*\[[^\]]+\]),([a-zA-Z_][a-zA-Z0-9_]*\[[^\]]+\])$`)
	if re7.MatchString(line) {
		line = re7.ReplaceAllString(line, "$1 $2, $3")
	} else {
		// Fix single gate calls: hq -> h q, hq[0] -> h q[0]
		re8 := regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)([a-zA-Z_][a-zA-Z0-9_]*(?:\[[^\]]+\])?)$`)
		if re8.MatchString(line) && !strings.Contains(line, " ") {
			line = re8.ReplaceAllString(line, "$1 $2")
		}
	}

	// Fix measure statements: measureq->c -> measure q -> c
	re9 := regexp.MustCompile(`measure([a-zA-Z_][a-zA-Z0-9_]*(?:\[[^\]]+\])?)->([a-zA-Z_][a-zA-Z0-9_]*(?:\[[^\]]+\])?)`)
	line = re9.ReplaceAllString(line, "measure $1 -> $2")

	// Fix timing expressions
	re10 := regexp.MustCompile(`\[\s*(\d+)\s*(ns|us|ms)\s*\]`)
	line = re10.ReplaceAllString(line, "[$1$2]")

	// Fix if statements: if(condition){ -> if (condition) { (after operator processing)
	re_if := regexp.MustCompile(`^if\s*\(([^)]+)\)\s*\{`)
	if re_if.MatchString(line) {
		line = re_if.ReplaceAllString(line, "if ($1) {")
	}

	// Fix else statements: }else{ -> } else {
	re_else := regexp.MustCompile(`\}\s*else\s*\{`)
	line = re_else.ReplaceAllString(line, "} else {")

	// Add semicolon back if it was there (but not for control flow statements and block definitions)
	trimmed := strings.TrimSpace(line)
	if hasSemicolon && !strings.HasPrefix(line, "//") && !strings.HasPrefix(line, "/*") &&
		!strings.HasPrefix(trimmed, "if") && !strings.HasPrefix(trimmed, "gate") &&
		!strings.HasSuffix(trimmed, "{") && !strings.HasSuffix(trimmed, "}") &&
		!strings.Contains(trimmed, "} else {") {
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
	// Handle assignments with proper spacing around = (but avoid comparison operators)
	if strings.Contains(text, "==") || strings.Contains(text, "!=") ||
		strings.Contains(text, "<=") || strings.Contains(text, ">=") ||
		strings.Contains(text, "if") {
		return text // Don't process
	}
	re := regexp.MustCompile(`\s*=\s*`)
	return re.ReplaceAllString(text, " = ")
}

// formatWithSimpleFallback provides basic text-based formatting
func (f *Formatter) formatWithSimpleFallback(content string) string {
	// First normalize the content and split compound statements
	content = f.splitCompoundStatements(content)
	lines := strings.Split(content, "\n")
	var result []string
	var lastStatementType string
	indentLevel := 0

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			// Don't add trailing empty lines when newline is false
			if i == len(lines)-1 && !f.newline {
				continue
			}
			result = append(result, "")
			continue
		}

		// Handle indentation for blocks
		isClosingBrace := strings.HasPrefix(trimmed, "}")
		if isClosingBrace {
			indentLevel--
		}

		// Determine current statement type
		currentType := f.getStatementType(trimmed)

		// Add empty line between different types of statements
		if f.shouldAddEmptyLineBetween(lastStatementType, currentType) {
			if len(result) > 0 && result[len(result)-1] != "" {
				result = append(result, "")
			}
		}

		// Apply basic formatting
		// For comments, preserve original line with indentation
		lineToFormat := trimmed
		if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "/*") || strings.HasSuffix(trimmed, "*/") {
			lineToFormat = line
		}

		formatted := f.formatLineBasic(lineToFormat)
		if formatted != "" {
			// Apply indentation only for non-comment lines
			if indentLevel > 0 && !isClosingBrace &&
				!strings.HasPrefix(trimmed, "//") && !strings.HasPrefix(trimmed, "/*") && !strings.HasSuffix(trimmed, "*/") {
				formatted = strings.Repeat(" ", indentLevel*f.indentSize) + formatted
			}

			result = append(result, formatted)
			lastStatementType = currentType
		}

		// Increase indent after opening braces
		if strings.HasSuffix(trimmed, "{") {
			indentLevel++
		}
	}

	output := strings.Join(result, "\n")
	if f.newline && len(result) > 0 && !strings.HasSuffix(output, "\n") {
		output += "\n"
	}

	return output
}

// formatLineBasic applies basic formatting to a line
func (f *Formatter) formatLineBasic(line string) string {
	// Handle comments - preserve as-is without adding semicolons
	trimmed := strings.TrimSpace(line)
	if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "/*") || strings.HasSuffix(trimmed, "*/") {
		return line
	}

	// Handle trailing comments: extract comment, format statement, then re-attach
	var trailingComment string
	if idx := strings.Index(line, "//"); idx != -1 {
		trailingComment = strings.TrimSpace(line[idx:])
		line = strings.TrimSpace(line[:idx])
	}

	// Apply basic operator spacing and formatting to non-comment part
	if line != "" {
		line = f.formatDeclaration(line)
		line = f.formatMeasure(line)
		line = f.formatOperators(line)
		line = f.formatCommas(line)

		// Ensure semicolon for statements (but not for control structures)
		if !strings.HasSuffix(line, ";") && !strings.HasSuffix(line, "{") && !strings.HasSuffix(line, "}") &&
			!strings.HasPrefix(line, "if") && !strings.HasPrefix(line, "gate") && line != "" {
			line += ";"
		}
	}

	// Re-attach trailing comment
	if trailingComment != "" {
		if line != "" {
			line += " " + trailingComment
		} else {
			line = trailingComment
		}
	}

	return line
}

// formatDeclaration applies formatting to type declarations
func (f *Formatter) formatDeclaration(line string) string {
	// Handle array declarations: bit[2]c -> bit[2] c
	re := regexp.MustCompile(`([a-zA-Z]+\[[^\]]+\])([a-zA-Z_][a-zA-Z0-9_]*)`)
	line = re.ReplaceAllString(line, "$1 $2")

	// Handle simple declarations: bitc -> bit c
	re = regexp.MustCompile(`^(bit|int|float|qubit)([a-zA-Z_][a-zA-Z0-9_]*)`)
	line = re.ReplaceAllString(line, "$1 $2")

	return line
}

// formatMeasure applies formatting to measurement statements
func (f *Formatter) formatMeasure(line string) string {
	// Add space after measure keyword: measureq -> measure q
	re := regexp.MustCompile(`^measure([a-zA-Z_][a-zA-Z0-9_\[\]]*)`)
	line = re.ReplaceAllString(line, "measure $1")

	return line
}

// formatCommas adds proper spacing after commas
func (f *Formatter) formatCommas(line string) string {
	// Add space after comma: q[0],q[1] -> q[0], q[1]
	re := regexp.MustCompile(`,\s*`)
	line = re.ReplaceAllString(line, ", ")

	return line
}

// getStatementType determines the type of a statement
func (f *Formatter) getStatementType(line string) string {
	line = strings.TrimSpace(line)

	if strings.HasPrefix(line, "OPENQASM") {
		return "version"
	}
	if strings.HasPrefix(line, "include") {
		return "include"
	}
	if strings.Contains(line, "qubit") || strings.Contains(line, "bit") || strings.Contains(line, "int") {
		return "declaration"
	}
	if strings.HasPrefix(line, "gate") {
		return "gate_definition"
	}
	if strings.HasPrefix(line, "}") {
		// Distinguish between } else { and plain }
		if strings.Contains(line, "else") {
			return "else_block"
		}
		return "block_end"
	}
	if strings.Contains(line, "measure") {
		return "measurement"
	}
	if strings.HasPrefix(line, "//") || strings.HasPrefix(line, "/*") {
		return "comment"
	}

	// Check if it's a gate call
	gates := []string{"h", "x", "y", "z", "cx", "cnot", "rx", "ry", "rz", "s", "t"}
	for _, gate := range gates {
		if strings.HasPrefix(line, gate+" ") || strings.HasPrefix(line, gate+";") {
			return "gate_call"
		}
	}

	return "other"
}

// shouldAddEmptyLineBetween determines if an empty line should be added between statement types
func (f *Formatter) shouldAddEmptyLineBetween(lastType, currentType string) bool {
	if lastType == "" {
		return false
	}

	// Add empty line after version statements (but not before comments)
	if lastType == "version" && currentType != "include" && currentType != "version" && currentType != "comment" {
		return true
	}

	// Add empty line after include statements
	if lastType == "include" && currentType != "include" && currentType != "version" {
		return true
	}

	// Add empty line after gate definitions (block_end), but not after } else {
	if lastType == "block_end" && (currentType == "gate_call" || currentType == "measurement" || currentType == "declaration") {
		return true
	}

	// Add empty line between declaration blocks and gate definitions (but not single declarations and gate calls)
	if lastType == "declaration" && currentType == "gate_definition" {
		return true
	}

	return false
}

// formatOperators applies basic operator formatting
func (f *Formatter) formatOperators(line string) string {
	// Skip if it's a comment line
	if strings.HasPrefix(strings.TrimSpace(line), "//") || strings.HasPrefix(strings.TrimSpace(line), "/*") {
		return line
	}

	// Apply operator spacing (order matters - longer operators first)
	// Handle -> specifically first
	line = regexp.MustCompile(`\s*->\s*`).ReplaceAllString(line, " -> ")

	// Handle == specifically to prevent interference with =
	line = regexp.MustCompile(`\s*==\s*`).ReplaceAllString(line, " == ")

	// Handle other comparison operators
	operators := []struct{ old, new string }{
		{"!=", " != "},
		{"<=", " <= "},
		{">=", " >= "},
		{"<", " < "},
	}

	for _, op := range operators {
		pattern := regexp.MustCompile(`\s*` + regexp.QuoteMeta(op.old) + `\s*`)
		line = pattern.ReplaceAllString(line, op.new)
	}

	// Handle = but not when it's part of == (negative lookbehind/lookahead)
	line = regexp.MustCompile(`([^=])\s*=\s*([^=])`).ReplaceAllString(line, "$1 = $2")

	// Handle > but not when it's part of ->
	line = regexp.MustCompile(`([^-])\s*>\s*`).ReplaceAllString(line, "$1 > ")

	return line
}
