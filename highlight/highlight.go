package highlight

import (
	"strings"

	"github.com/antlr4-go/antlr/v4"
	qasm_gen "github.com/orangekame3/qasmtools/parser/gen"
)

// Color represents ANSI color codes for terminal output
type Color string

const (
	Reset   Color = "\033[0m"
	Red     Color = "\033[31m"
	Green   Color = "\033[32m"
	Yellow  Color = "\033[33m"
	Blue    Color = "\033[34m"
	Magenta Color = "\033[35m"
	Cyan    Color = "\033[36m"
	White   Color = "\033[37m"
	Bold    Color = "\033[1m"
)

// TokenInfo represents detailed information about a token
type TokenInfo struct {
	Type     TokenType `json:"type"`
	TypeName string    `json:"type_name"`
	Content  string    `json:"content"`
	Line     int       `json:"line"`
	Column   int       `json:"column"`
	Length   int       `json:"length"`
}

// HighlightResult represents the result of syntax highlighting
type HighlightResult struct {
	Tokens []TokenInfo `json:"tokens"`
}

// ColorScheme defines the colors for different token types
type ColorScheme struct {
	Keyword          Color
	Operator         Color
	Identifier       Color
	Number           Color
	String           Color
	Comment          Color
	Gate             Color
	Measurement      Color
	Register         Color
	Punctuation      Color
	BuiltinGate      Color
	BuiltinQuantum   Color
	BuiltinClassical Color
	BuiltinConstant  Color
	AccessControl    Color
	Extern           Color
	HardwareQubit    Color
}

// DefaultColorScheme returns the default color scheme
func DefaultColorScheme() *ColorScheme {
	return &ColorScheme{
		Keyword:          Yellow,
		Operator:         Red,
		Identifier:       White,
		Number:           Cyan,
		String:           Green,
		Comment:          Blue,
		Gate:            Magenta,
		Measurement:      Red,
		Register:         Yellow,
		Punctuation:      White,
		BuiltinGate:      Magenta,
		BuiltinQuantum:   Red,
		BuiltinClassical: Green,
		BuiltinConstant:  Cyan,
		AccessControl:    Yellow,
		Extern:          Blue,
		HardwareQubit:    Magenta,
	}
}

// TokenType represents the type of syntax token
type TokenType int

const (
	// Token types for syntax highlighting
	TokenKeyword TokenType = iota
	TokenOperator
	TokenIdentifier
	TokenNumber
	TokenString
	TokenComment
	TokenGate
	TokenMeasurement
	TokenRegister
	TokenPunctuation
	TokenBuiltinGate
	TokenBuiltinQuantum
	TokenBuiltinClassical
	TokenBuiltinConstant
	TokenAccessControl
	TokenExtern
	TokenHardwareQubit
)

// Token represents a syntax highlighted token
type Token struct {
	Type    TokenType
	Content string
	Line    int
	Column  int
}

// Highlighter provides syntax highlighting for QASM code
type Highlighter struct {
	tokens      []TokenInfo
	colorScheme *ColorScheme
}

// New creates a new Highlighter instance with default color scheme
func New() *Highlighter {
	return &Highlighter{
		colorScheme: DefaultColorScheme(),
	}
}

// NewWithColorScheme creates a new Highlighter instance with custom color scheme
func NewWithColorScheme(scheme *ColorScheme) *Highlighter {
	return &Highlighter{
		colorScheme: scheme,
	}
}

// Highlight performs syntax highlighting on the given QASM code and returns colored output
func (h *Highlighter) Highlight(input string) (string, error) {
	// Create the input stream
	inputStream := antlr.NewInputStream(input)

	// Create the lexer
	lexer := qasm_gen.Newqasm3Lexer(inputStream)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// Clear previous tokens
	h.tokens = make([]TokenInfo, 0)

	// Walk through all tokens
	for {
		t := stream.LT(1)
		if t.GetTokenType() == antlr.TokenEOF {
			break
		}
		h.processToken(t)
		stream.Consume()
	}

	return h.ColoredString(), nil
}

// processToken converts an ANTLR token to a syntax highlighting token
func (h *Highlighter) processToken(t antlr.Token) {
	tokenType := h.getTokenType(t)
	token := TokenInfo{
		Type:     tokenType,
		TypeName: h.tokenTypeToString(tokenType),
		Content:  t.GetText(),
		Line:     t.GetLine(),
		Column:   t.GetColumn(),
		Length:   len(t.GetText()),
	}
	h.tokens = append(h.tokens, token)
}

// GetTokens returns the tokens for LSP integration
func (h *Highlighter) GetTokens() []TokenInfo {
	return h.tokens
}

// ColoredString returns the colored string representation of the tokens
func (h *Highlighter) ColoredString() string {
	var result strings.Builder
	currentLine := 1
	currentColumn := 0

	for _, token := range h.tokens {
		// Add newlines if needed
		for currentLine < token.Line {
			result.WriteString("\n")
			currentLine++
			currentColumn = 0
		}

		// Add spaces for correct column position
		for currentColumn < token.Column {
			result.WriteString(" ")
			currentColumn++
		}

		// Get color for token type
		color := h.getColorForToken(token.Type)

		// Write colored token
		result.WriteString(string(color))
		result.WriteString(token.Content)
		result.WriteString(string(Reset))

		currentColumn += token.Length
	}

	return result.String()
}

// getColorForToken returns the appropriate color for a token type
func (h *Highlighter) getColorForToken(t TokenType) Color {
	switch t {
	case TokenKeyword:
		return h.colorScheme.Keyword
	case TokenOperator:
		return h.colorScheme.Operator
	case TokenIdentifier:
		return h.colorScheme.Identifier
	case TokenNumber:
		return h.colorScheme.Number
	case TokenString:
		return h.colorScheme.String
	case TokenComment:
		return h.colorScheme.Comment
	case TokenGate:
		return h.colorScheme.Gate
	case TokenMeasurement:
		return h.colorScheme.Measurement
	case TokenRegister:
		return h.colorScheme.Register
	case TokenPunctuation:
		return h.colorScheme.Punctuation
	case TokenBuiltinGate:
		return h.colorScheme.BuiltinGate
	case TokenBuiltinQuantum:
		return h.colorScheme.BuiltinQuantum
	case TokenBuiltinClassical:
		return h.colorScheme.BuiltinClassical
	case TokenBuiltinConstant:
		return h.colorScheme.BuiltinConstant
	case TokenAccessControl:
		return h.colorScheme.AccessControl
	case TokenExtern:
		return h.colorScheme.Extern
	case TokenHardwareQubit:
		return h.colorScheme.HardwareQubit
	default:
		return Reset
	}
}

// getTokenType determines the token type for syntax highlighting
func (h *Highlighter) getTokenType(t antlr.Token) TokenType {
	switch t.GetTokenType() {
	// Keywords
	case 1, 2, // OPENQASM, INCLUDE
		4, 5, 6, // DEF, CAL, DEFCAL
		10, 11, // LET, BREAK
		12, 13, 14, // CONTINUE, IF, ELSE
		15, 16, 17, // END, RETURN, FOR
		18, 19: // WHILE, IN
		return TokenKeyword

	// Types and registers
	case 31, 32, 33, // QREG, QUBIT, CREG
		34, 35: // BOOL, BIT
		return TokenRegister

	// Builtin gates and quantum operations
	case 53: // RESET
		return TokenBuiltinQuantum

	// Operators
	case 69, 71, 72, // PLUS, MINUS, ASTERISK
		74, 67, 68, // SLASH, EQUALS, ARROW
		84, 86, // EqualityOperator, ComparisonOperator
		63: // COLON
		return TokenOperator

	// Gates and extern
	case 7: // GATE token
		return TokenGate
	case 8: // EXTERN
		return TokenExtern

	// Measurements
	case 54: // MEASURE
		return TokenMeasurement

	// Numbers
	case 92, 93, 96: // DecimalIntegerLiteral, HexIntegerLiteral, FloatLiteral
		return TokenNumber

	// Strings
	case 106: // StringLiteral
		return TokenString

	// Comments
	case 101, 102: // LineComment, BlockComment
		return TokenComment

	// Access control
	case 28, 29, 30: // CONST, READONLY, MUTABLE
		return TokenAccessControl

	// Builtin constants
	case 88: // IMAG
		return TokenBuiltinConstant

	// Builtin classical functions
	case 42: // VOID
		return TokenBuiltinClassical

	// Hardware qubits (special identifiers starting with $)
	case 95: // HardwareQubit
		return TokenHardwareQubit

	// Regular identifiers
	case 94: // Identifier
		// Check if it's a builtin gate (U, CX)
		content := t.GetText()
		if content == "U" || content == "CX" {
			return TokenBuiltinGate
		}
		// Check if it's a builtin classical function
		if content == "sin" || content == "cos" || content == "tan" ||
			content == "exp" || content == "ln" || content == "sqrt" {
			return TokenBuiltinClassical
		}
		// Check if it's a builtin constant
		if content == "pi" || content == "tau" || content == "euler" {
			return TokenBuiltinConstant
		}
		return TokenIdentifier

	// Punctuation
	case 64, 66, // SEMICOLON, COMMA
		61, 62, // LPAREN, RPAREN
		57, 58, // LBRACKET, RBRACKET
		59, 60, // LBRACE, RBRACE
		83: // EXCLAMATION_POINT
		return TokenPunctuation

	default:
		return TokenIdentifier
	}
}

// tokenTypeToString converts a TokenType to its string representation
func (h *Highlighter) tokenTypeToString(t TokenType) string {
	switch t {
	case TokenKeyword:
		return "keyword"
	case TokenOperator:
		return "operator"
	case TokenIdentifier:
		return "identifier"
	case TokenNumber:
		return "number"
	case TokenString:
		return "string"
	case TokenComment:
		return "comment"
	case TokenGate:
		return "gate"
	case TokenMeasurement:
		return "measurement"
	case TokenRegister:
		return "register"
	case TokenPunctuation:
		return "punctuation"
	case TokenBuiltinGate:
		return "builtin_gate"
	case TokenBuiltinQuantum:
		return "builtin_quantum"
	case TokenBuiltinClassical:
		return "builtin_classical"
	case TokenBuiltinConstant:
		return "builtin_constant"
	case TokenAccessControl:
		return "access_control"
	case TokenExtern:
		return "extern"
	case TokenHardwareQubit:
		return "hardware_qubit"
	default:
		return "unknown"
	}
}
