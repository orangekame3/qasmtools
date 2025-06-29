package parser

import (
	"context"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	qasm_gen "github.com/orangekame3/qasmtools/parser/gen"
)

// ParseOptions configures the parser behavior
type ParseOptions struct {
	// StrictMode enables strict OpenQASM 3.0 compliance
	StrictMode bool

	// IncludeComments preserves comments in the AST
	IncludeComments bool

	// ErrorRecovery enables error recovery for partial parsing
	ErrorRecovery bool

	// MaxErrors limits the number of errors to collect
	MaxErrors int
}

// DefaultParseOptions returns default parsing options
func DefaultParseOptions() *ParseOptions {
	return &ParseOptions{
		StrictMode:      false,
		IncludeComments: true,
		ErrorRecovery:   true,
		MaxErrors:       100,
	}
}

// Parser represents the main OpenQASM 3.0 parser
type Parser struct {
	options *ParseOptions
	stream  antlr.TokenStream
}

// NewParser creates a new parser with default options
func NewParser() *Parser {
	return &Parser{
		options: DefaultParseOptions(),
	}
}

// NewParserWithOptions creates a parser with custom options
func NewParserWithOptions(opts *ParseOptions) *Parser {
	if opts == nil {
		opts = DefaultParseOptions()
	}
	return &Parser{
		options: opts,
	}
}

// GetTokenStream returns the current token stream
func (p *Parser) GetTokenStream() antlr.TokenStream {
	return p.stream
}

// ParseString parses QASM code from a string
func (p *Parser) ParseString(content string) (*Program, error) {
	result := p.ParseWithErrors(content)
	if result.HasErrors() {
		return result.Program, &result.Errors[0]
	}
	return result.Program, nil
}

// ParseReader parses QASM code from an io.Reader
func (p *Parser) ParseReader(reader io.Reader) (*Program, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return p.ParseString(string(content))
}

// ParseFile parses QASM code from a file
func (p *Parser) ParseFile(filename string) (*Program, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return p.ParseString(string(content))
}

// ParseWithContext parses with context for cancellation
func (p *Parser) ParseWithContext(ctx context.Context, content string) (*Program, error) {
	// Check if context is already cancelled
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// For now, we'll parse directly without context support
	// Future enhancement could implement context-aware parsing
	return p.ParseString(content)
}

// Validate validates QASM syntax without building full AST
func (p *Parser) Validate(content string) error {
	result := p.ParseWithErrors(content)
	if result.HasErrors() {
		return &result.Errors[0]
	}
	return nil
}

// ParseWithErrors returns partial results even with errors
func (p *Parser) ParseWithErrors(content string) *ParseResult {
	// Preprocess content to handle common issues
	content = p.preprocessContent(content)

	// Create input stream
	input := antlr.NewInputStream(content)

	// Create lexer (this will be replaced with generated code)
	lexer := p.createLexer(input)

	// Create error listener for lexer
	lexerErrors := NewErrorListener()
	lexer.RemoveErrorListeners()
	if !p.options.ErrorRecovery {
		lexer.AddErrorListener(lexerErrors)
	}

	// Create token stream
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p.stream = stream

	// Create parser (this will be replaced with generated code)
	parser := p.createParser(stream)

	// Create error listener for parser
	parserErrors := NewErrorListener()
	parser.RemoveErrorListeners()
	if !p.options.ErrorRecovery {
		parser.AddErrorListener(parserErrors)
	}

	// Parse the program
	tree := p.parseProgram(parser)

	// Collect all errors
	allErrors := make([]ParseError, 0)
	allErrors = append(allErrors, lexerErrors.GetErrors()...)
	allErrors = append(allErrors, parserErrors.GetErrors()...)

	// Convert parse tree to AST
	program := p.convertToAST(tree, content)

	// Extract and associate comments if enabled
	if p.options.IncludeComments {
		commentExtractor := NewCommentExtractor(content)
		commentExtractor.ExtractComments(stream)
		commentExtractor.AssociateCommentsWithStatements(program)
	}

	// Add AST builder errors if any
	if visitor, ok := tree.(interface{ GetErrors() []ParseError }); ok {
		allErrors = append(allErrors, visitor.GetErrors()...)
	}

	// Limit errors if specified
	if p.options.MaxErrors > 0 && len(allErrors) > p.options.MaxErrors {
		allErrors = allErrors[:p.options.MaxErrors]
	}

	result := &ParseResult{
		Program: program,
		Errors:  allErrors,
	}

	return result
}

// preprocessContent handles common formatting issues
func (p *Parser) preprocessContent(content string) string {
	// Normalize line endings
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")

	// Ensure content ends with newline
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}

	return content
}

// createLexer creates the ANTLR lexer
// This is a placeholder - will be replaced with generated lexer
func (p *Parser) createLexer(input antlr.CharStream) antlr.Lexer {
	return qasm_gen.Newqasm3Lexer(input)
}

// createParser creates the ANTLR parser
// This is a placeholder - will be replaced with generated parser
func (p *Parser) createParser(stream antlr.TokenStream) antlr.Parser {
	return qasm_gen.Newqasm3Parser(stream)
}

// parseProgram parses the root program
func (p *Parser) parseProgram(parser antlr.Parser) antlr.Tree {
	// Use reflection to call Program() method since qasm3Parser is not exported
	parserValue := reflect.ValueOf(parser)
	programMethod := parserValue.MethodByName("Program")
	if !programMethod.IsValid() {
		return nil
	}

	result := programMethod.Call([]reflect.Value{})
	if len(result) == 0 {
		return nil
	}

	if tree, ok := result[0].Interface().(antlr.Tree); ok {
		return tree
	}
	return nil
}

// convertToAST converts ANTLR parse tree to our AST
func (p *Parser) convertToAST(tree antlr.Tree, content string) *Program {
	if tree == nil {
		return &Program{
			BaseNode: BaseNode{
				Position: Position{Line: 1, Column: 1},
			},
			Statements: make([]Statement, 0),
			Comments:   make([]Comment, 0),
		}
	}

	// Use the AST builder visitor for complete parsing
	visitor := NewASTBuilderVisitor(content)

	// Visit the parse tree to build AST
	result := visitor.VisitProgram(tree.(*qasm_gen.ProgramContext))

	if program, ok := result.(*Program); ok {
		return program
	}

	// Fallback if visitor fails
	return &Program{
		BaseNode: BaseNode{
			Position: Position{Line: 1, Column: 1},
		},
		Statements: make([]Statement, 0),
		Comments:   make([]Comment, 0),
	}
}

// extractVersion extracts version information from parse tree
func (p *Parser) extractVersion(tree antlr.ParseTree) *Version {
	// Look for version declaration in the parse tree
	// This is a simplified implementation
	if tree.GetText() != "" && strings.Contains(tree.GetText(), "3.0") {
		return &Version{
			BaseNode: BaseNode{Position: Position{Line: 1, Column: 1}},
			Number:   "3.0",
		}
	}

	// Recursively check children
	for i := 0; i < tree.GetChildCount(); i++ {
		if child := tree.GetChild(i); child != nil {
			if parseTree, ok := child.(antlr.ParseTree); ok {
				if version := p.extractVersion(parseTree); version != nil {
					return version
				}
			}
		}
	}

	return nil
}

// GetOptions returns the current parser options
func (p *Parser) GetOptions() *ParseOptions {
	return p.options
}

// SetOptions updates the parser options
func (p *Parser) SetOptions(opts *ParseOptions) {
	if opts != nil {
		p.options = opts
	}
}
