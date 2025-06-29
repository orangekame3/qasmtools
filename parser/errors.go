package parser

import (
	"fmt"
	"strings"

	"github.com/antlr4-go/antlr/v4"
)

// ParseError represents parsing errors
type ParseError struct {
	Message  string   `json:"message"`
	Position Position `json:"position"`
	Type     string   `json:"type"` // "syntax", "semantic", "lexer"
	Context  string   `json:"context,omitempty"`
	Code     string   `json:"code,omitempty"`     // Error code
	Severity string   `json:"severity,omitempty"` // "error", "warning", "info"
	Expected []string `json:"expected,omitempty"` // Expected tokens for recovery
	Actual   string   `json:"actual,omitempty"`   // Actual token found
}

func (e *ParseError) Error() string {
	if e.Context != "" {
		return fmt.Sprintf("%s error at line %d, column %d: %s (context: %s)",
			e.Type, e.Position.Line, e.Position.Column, e.Message, e.Context)
	}
	return fmt.Sprintf("%s error at line %d, column %d: %s",
		e.Type, e.Position.Line, e.Position.Column, e.Message)
}

// ParseResult contains parsing results with errors
type ParseResult struct {
	Program *Program     `json:"program,omitempty"`
	Errors  []ParseError `json:"errors,omitempty"`
}

// HasErrors returns true if there are any errors
func (r *ParseResult) HasErrors() bool {
	return len(r.Errors) > 0
}

// ErrorMessages returns all error messages as strings
func (r *ParseResult) ErrorMessages() []string {
	messages := make([]string, len(r.Errors))
	for i, err := range r.Errors {
		messages[i] = err.Error()
	}
	return messages
}

// String returns a formatted string of all errors
func (r *ParseResult) String() string {
	if !r.HasErrors() {
		return "No errors"
	}
	return strings.Join(r.ErrorMessages(), "\n")
}

// ErrorListener implements ANTLR error listener interface
type ErrorListener struct {
	errors []ParseError
}

// NewErrorListener creates a new error listener
func NewErrorListener() *ErrorListener {
	return &ErrorListener{
		errors: make([]ParseError, 0),
	}
}

// GetErrors returns collected errors
func (l *ErrorListener) GetErrors() []ParseError {
	return l.errors
}

// HasErrors returns true if any errors were collected
func (l *ErrorListener) HasErrors() bool {
	return len(l.errors) > 0
}

// SyntaxError implements antlr.ErrorListener interface
func (l *ErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	pos := Position{
		Line:   line,
		Column: column,
		Offset: 0, // Would need to be calculated from token
		Length: 1, // Default length
	}

	// Extract additional information if available
	var actual string
	var expected []string

	if token, ok := offendingSymbol.(antlr.Token); ok {
		actual = token.GetText()
		pos.Offset = token.GetStart()
		pos.Length = len(actual)
	}

	// Try to extract expected tokens from exception
	if e != nil {
		// This would require more sophisticated parsing of the exception
		// For now, just use the message
	}

	l.errors = append(l.errors, ParseError{
		Message:  msg,
		Position: pos,
		Type:     "syntax",
		Severity: "error",
		Expected: expected,
		Actual:   actual,
	})
}

// ReportAmbiguity implements antlr.ErrorListener interface
func (l *ErrorListener) ReportAmbiguity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, exact bool, ambigAlts *antlr.BitSet, configs *antlr.ATNConfigSet) {
	// Optional: Handle ambiguity errors if needed
}

// ReportAttemptingFullContext implements antlr.ErrorListener interface
func (l *ErrorListener) ReportAttemptingFullContext(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, conflictingAlts *antlr.BitSet, configs *antlr.ATNConfigSet) {
	// Optional: Handle full context attempts if needed
}

// ReportContextSensitivity implements antlr.ErrorListener interface
func (l *ErrorListener) ReportContextSensitivity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, prediction int, configs *antlr.ATNConfigSet) {
	// Optional: Handle context sensitivity if needed
}

// NewSyntaxError creates a new syntax error
func NewSyntaxError(message string, pos Position) ParseError {
	return ParseError{
		Message:  message,
		Position: pos,
		Type:     "syntax",
	}
}

// NewSemanticError creates a new semantic error
func NewSemanticError(message string, pos Position) ParseError {
	return ParseError{
		Message:  message,
		Position: pos,
		Type:     "semantic",
	}
}

// NewLexerError creates a new lexer error
func NewLexerError(message string, pos Position) ParseError {
	return ParseError{
		Message:  message,
		Position: pos,
		Type:     "lexer",
	}
}
