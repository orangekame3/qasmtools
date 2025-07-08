package highlight

import (
	"github.com/orangekame3/qasmtools/parser"
)

// ASTHighlighter provides enhanced AST-based semantic highlighting
type ASTHighlighter struct {
	*Highlighter // Embed existing highlighter
	program      *parser.Program
	symbols      map[string]SymbolInfo
}

// SymbolInfo provides context information about identifiers
type SymbolInfo struct {
	Type        SymbolType
	Declaration parser.Node
	Usages      []parser.Node
	IsUsed      bool
	IsArray     bool
	ArraySize   int
}

// SymbolType represents different types of symbols for semantic highlighting
type SymbolType int

const (
	SymbolUnknown SymbolType = iota
	SymbolQubitDeclaration
	SymbolBitDeclaration
	SymbolGateDefinition
	SymbolParameter
	SymbolUnusedVariable
	SymbolArrayAccess
	SymbolGateCall
	SymbolMeasurementTarget
)

// Enhanced token types for semantic highlighting
const (
	TokenSemanticQubitDecl TokenType = iota + 100
	TokenSemanticBitDecl
	TokenSemanticUnusedVar
	TokenSemanticArrayAccess
	TokenSemanticGateParam
	TokenSemanticMeasureTarget
	TokenSemanticGateDefinition
	TokenSemanticConstantValue
)

// NewASTHighlighter creates a new AST-enhanced highlighter
func NewASTHighlighter() *ASTHighlighter {
	return &ASTHighlighter{
		Highlighter: New(),
		symbols:     make(map[string]SymbolInfo),
	}
}

// NewASTHighlighterWithColorScheme creates an AST highlighter with custom colors
func NewASTHighlighterWithColorScheme(scheme *ColorScheme) *ASTHighlighter {
	return &ASTHighlighter{
		Highlighter: NewWithColorScheme(scheme),
		symbols:     make(map[string]SymbolInfo),
	}
}

// HighlightWithAST performs semantic highlighting using both ANTLR tokens and AST context
func (h *ASTHighlighter) HighlightWithAST(input string) (string, error) {
	// Phase 1: Parse AST for semantic analysis
	p := parser.NewParser()
	result := p.ParseWithErrors(input)

	if result.HasErrors() || result.Program == nil {
		// Fallback to token-based highlighting
		return h.Highlight(input)
	}

	// Phase 2: Analyze AST for semantic information
	h.program = result.Program
	h.analyzeSymbols()

	// Phase 3: Perform enhanced highlighting
	return h.highlightWithSemanticInfo(input)
}

// analyzeSymbols builds symbol table from AST for semantic highlighting
func (h *ASTHighlighter) analyzeSymbols() {
	if h.program == nil {
		return
	}

	// Clear previous analysis
	h.symbols = make(map[string]SymbolInfo)

	// First pass: collect declarations
	for _, stmt := range h.program.Statements {
		h.analyzeDeclaration(stmt)
	}

	// Second pass: find usages and mark unused variables
	for _, stmt := range h.program.Statements {
		h.analyzeUsages(stmt)
	}

	// Third pass: mark unused variables
	h.markUnusedVariables()
}

// analyzeDeclaration processes declaration statements
func (h *ASTHighlighter) analyzeDeclaration(stmt parser.Statement) {
	switch s := stmt.(type) {
	case *parser.QuantumDeclaration:
		symbolType := SymbolQubitDeclaration
		isArray := s.Size != nil
		arraySize := 0
		if isArray {
			arraySize = h.extractArraySize(s.Size)
		}

		h.symbols[s.Identifier] = SymbolInfo{
			Type:        symbolType,
			Declaration: s,
			IsArray:     isArray,
			ArraySize:   arraySize,
			Usages:      make([]parser.Node, 0),
		}

	case *parser.ClassicalDeclaration:
		symbolType := SymbolBitDeclaration
		isArray := s.Size != nil
		arraySize := 0
		if isArray {
			arraySize = h.extractArraySize(s.Size)
		}

		h.symbols[s.Identifier] = SymbolInfo{
			Type:        symbolType,
			Declaration: s,
			IsArray:     isArray,
			ArraySize:   arraySize,
			Usages:      make([]parser.Node, 0),
		}

	case *parser.GateDefinition:
		h.symbols[s.Name] = SymbolInfo{
			Type:        SymbolGateDefinition,
			Declaration: s,
			Usages:      make([]parser.Node, 0),
		}

		// Also analyze parameters
		for _, param := range s.Parameters {
			h.symbols[param.Name] = SymbolInfo{
				Type:        SymbolParameter,
				Declaration: s,
				Usages:      make([]parser.Node, 0),
			}
		}
	}
}

// analyzeUsages finds identifier usages throughout the program
func (h *ASTHighlighter) analyzeUsages(stmt parser.Statement) {
	switch s := stmt.(type) {
	case *parser.GateCall:
		// Mark gate usage
		if info, exists := h.symbols[s.Name]; exists {
			info.Usages = append(info.Usages, s)
			info.IsUsed = true
			h.symbols[s.Name] = info
		}

		// Mark qubit usages
		for _, qubit := range s.Qubits {
			h.analyzeExpressionUsage(qubit)
		}

		// Mark parameter usages
		for _, param := range s.Parameters {
			h.analyzeExpressionUsage(param)
		}

	case *parser.Measurement:
		h.analyzeExpressionUsage(s.Qubit)
		if s.Target != nil {
			h.analyzeExpressionUsage(s.Target)
		}

	case *parser.IfStatement:
		h.analyzeExpressionUsage(s.Condition)
		for _, thenStmt := range s.ThenBody {
			h.analyzeUsages(thenStmt)
		}
		for _, elseStmt := range s.ElseBody {
			h.analyzeUsages(elseStmt)
		}

	case *parser.GateDefinition:
		for _, bodyStmt := range s.Body {
			h.analyzeUsages(bodyStmt)
		}
	}
}

// analyzeExpressionUsage analyzes expression nodes for identifier usage
func (h *ASTHighlighter) analyzeExpressionUsage(expr parser.Expression) {
	if expr == nil {
		return
	}

	switch e := expr.(type) {
	case *parser.Identifier:
		if info, exists := h.symbols[e.Name]; exists {
			info.Usages = append(info.Usages, e)
			info.IsUsed = true
			h.symbols[e.Name] = info
		}

	case *parser.IndexedIdentifier:
		if info, exists := h.symbols[e.Name]; exists {
			info.Usages = append(info.Usages, e)
			info.IsUsed = true
			h.symbols[e.Name] = info
		}
		h.analyzeExpressionUsage(e.Index)

	case *parser.BinaryExpression:
		h.analyzeExpressionUsage(e.Left)
		h.analyzeExpressionUsage(e.Right)

	case *parser.UnaryExpression:
		h.analyzeExpressionUsage(e.Operand)

	case *parser.FunctionCall:
		for _, arg := range e.Arguments {
			h.analyzeExpressionUsage(arg)
		}

	case *parser.ParenthesizedExpression:
		h.analyzeExpressionUsage(e.Expression)
	}
}

// markUnusedVariables identifies unused variables for special highlighting
func (h *ASTHighlighter) markUnusedVariables() {
	for name, info := range h.symbols {
		if !info.IsUsed && (info.Type == SymbolQubitDeclaration || info.Type == SymbolBitDeclaration) {
			info.Type = SymbolUnusedVariable
			h.symbols[name] = info
		}
	}
}

// extractArraySize extracts array size from size expression
func (h *ASTHighlighter) extractArraySize(sizeExpr parser.Expression) int {
	if intLit, ok := sizeExpr.(*parser.IntegerLiteral); ok {
		return int(intLit.Value)
	}
	return 0
}

// highlightWithSemanticInfo performs enhanced highlighting with semantic information
func (h *ASTHighlighter) highlightWithSemanticInfo(input string) (string, error) {
	// First get basic token-level highlighting
	_, err := h.Highlight(input)
	if err != nil {
		return "", err
	}

	// Enhance tokens with semantic information
	h.enhanceTokensWithSemantics()

	return h.ColoredString(), nil
}

// enhanceTokensWithSemantics enhances existing tokens with AST semantic information
func (h *ASTHighlighter) enhanceTokensWithSemantics() {
	for i, token := range h.tokens {
		if token.Type == TokenIdentifier {
			// Check if this identifier has semantic information
			if info, exists := h.symbols[token.Content]; exists {
				// Enhance token type based on semantic analysis
				h.tokens[i].Type = h.getSemanticTokenType(info, token)
				h.tokens[i].TypeName = h.tokenTypeToString(h.tokens[i].Type)
			}
		}
	}
}

// getSemanticTokenType determines enhanced token type based on semantic analysis
func (h *ASTHighlighter) getSemanticTokenType(info SymbolInfo, token TokenInfo) TokenType {
	switch info.Type {
	case SymbolQubitDeclaration:
		if h.isDeclarationPosition(token, info.Declaration) {
			return TokenSemanticQubitDecl
		}
		if h.isArrayAccessPosition(token) {
			return TokenSemanticArrayAccess
		}
		return TokenRegister

	case SymbolBitDeclaration:
		if h.isDeclarationPosition(token, info.Declaration) {
			return TokenSemanticBitDecl
		}
		if h.isArrayAccessPosition(token) {
			return TokenSemanticArrayAccess
		}
		return TokenRegister

	case SymbolUnusedVariable:
		return TokenSemanticUnusedVar

	case SymbolGateDefinition:
		if h.isDeclarationPosition(token, info.Declaration) {
			return TokenSemanticGateDefinition
		}
		return TokenBuiltinGate

	case SymbolParameter:
		return TokenSemanticGateParam

	default:
		return token.Type
	}
}

// isDeclarationPosition checks if token is at declaration position
func (h *ASTHighlighter) isDeclarationPosition(token TokenInfo, declaration parser.Node) bool {
	if declaration == nil {
		return false
	}

	declPos := declaration.Pos()
	// This is a simplified check - in practice you'd need more sophisticated position matching
	return token.Line == declPos.Line
}

// isArrayAccessPosition checks if this identifier is used in array access context
func (h *ASTHighlighter) isArrayAccessPosition(token TokenInfo) bool {
	// Look for following '[' token to detect array access
	for _, t := range h.tokens {
		if t.Line == token.Line &&
			t.Column == token.Column+token.Length &&
			t.Content == "[" {
			return true
		}
	}
	return false
}

// Enhanced color scheme for semantic highlighting
func EnhancedColorScheme() *ColorScheme {
	scheme := DefaultColorScheme()

	// Add enhanced colors for semantic highlighting
	// These would be extended in the ColorScheme struct
	return scheme
}

// Override tokenTypeToString for semantic types
func (h *ASTHighlighter) tokenTypeToString(t TokenType) string {
	switch t {
	case TokenSemanticQubitDecl:
		return "semantic_qubit_declaration"
	case TokenSemanticBitDecl:
		return "semantic_bit_declaration"
	case TokenSemanticUnusedVar:
		return "semantic_unused_variable"
	case TokenSemanticArrayAccess:
		return "semantic_array_access"
	case TokenSemanticGateParam:
		return "semantic_gate_parameter"
	case TokenSemanticGateDefinition:
		return "semantic_gate_definition"
	default:
		return h.Highlighter.tokenTypeToString(t)
	}
}
