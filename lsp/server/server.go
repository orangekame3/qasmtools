package server

import (
	"github.com/tliron/commonlog"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"

	"github.com/orangekame3/qasmtools/formatter"
	"github.com/orangekame3/qasmtools/highlight"
	"github.com/orangekame3/qasmtools/lint"
	"github.com/orangekame3/qasmtools/lsp/features"
	"github.com/orangekame3/qasmtools/lsp/protocol/documents"
)

const LSName = "qasm"

// Server represents the QASM Language Server
type Server struct {
	handler     protocol.Handler
	log         commonlog.Logger
	docManager  *documents.Manager
	diagnostics *features.DiagnosticsProvider
	formatting  *features.FormattingProvider
	highlighting *features.HighlightingProvider
}

// NewServer creates a new QASM Language Server
func NewServer(version string) *Server {
	log := commonlog.GetLogger("qasm")
	
	// Initialize components
	linter := lint.NewLinter("")
	if err := linter.LoadRules(); err != nil {
		log.Error("Failed to load linting rules", "error", err)
	} else {
		rules := linter.GetRules()
		log.Info("Linting rules loaded successfully", "rule_count", len(rules))
	}

	formatter := formatter.NewFormatter()
	highlighter := highlight.New()
	
	// Create document manager
	docManager := documents.NewManager()
	
	// Create feature providers
	diagnostics := features.NewDiagnosticsProvider(linter, log)
	formattingProvider := features.NewFormattingProvider(formatter, log)
	highlightingProvider := features.NewHighlightingProvider(highlighter, log)
	
	server := &Server{
		log:          log,
		docManager:   docManager,
		diagnostics:  diagnostics,
		formatting:   formattingProvider,
		highlighting: highlightingProvider,
	}
	
	// Setup protocol handlers
	server.setupHandlers(version)
	
	return server
}

// Start starts the language server
func (s *Server) Start() {
	s.log.Info("Starting QASM Language Server")
	server := server.NewServer(&s.handler, LSName, true)
	server.RunStdio()
}

// setupHandlers configures the protocol handlers
func (s *Server) setupHandlers(version string) {
	s.handler = protocol.Handler{
		Initialize:                     s.initialize(version),
		Initialized:                    s.initialized,
		Shutdown:                       s.shutdown,
		TextDocumentDidOpen:            s.textDocumentDidOpen,
		TextDocumentDidChange:          s.textDocumentDidChange,
		TextDocumentSemanticTokensFull: s.textDocumentSemanticTokensFull,
		TextDocumentFormatting:         s.textDocumentFormatting,
	}
}

// Protocol handler implementations

func (s *Server) initialize(version string) func(*glsp.Context, *protocol.InitializeParams) (any, error) {
	return func(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
		capabilities := s.handler.CreateServerCapabilities()
		
		// Configure semantic tokens
		capabilities.SemanticTokensProvider = &protocol.SemanticTokensOptions{
			Legend: protocol.SemanticTokensLegend{
				TokenTypes: []string{
					"keyword", "operator", "identifier", "number", "string", "comment",
					"gate", "measurement", "register", "punctuation", "builtin_gate",
					"builtin_quantum", "builtin_classical", "builtin_constant",
					"access_control", "extern", "hardware_qubit",
				},
				TokenModifiers: []string{},
			},
			Full: protocol.True,
		}
		
		capabilities.DocumentFormattingProvider = protocol.True
		
		return protocol.InitializeResult{
			Capabilities: capabilities,
			ServerInfo: &protocol.InitializeResultServerInfo{
				Name:    LSName,
				Version: &version,
			},
		}, nil
	}
}

func (s *Server) initialized(context *glsp.Context, params *protocol.InitializedParams) error {
	return nil
}

func (s *Server) shutdown(context *glsp.Context) error {
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func (s *Server) textDocumentDidOpen(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	uri := params.TextDocument.URI
	content := params.TextDocument.Text
	
	s.log.Info("Document opened", "uri", uri, "length", len(content))
	
	// Update document manager
	s.docManager.DidOpen(uri, content)
	
	// Run diagnostics
	s.diagnostics.PublishDiagnostics(context, uri, content)
	
	return nil
}

func (s *Server) textDocumentDidChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	uri := params.TextDocument.URI
	
	// Convert changes to proper type
	var changes []protocol.TextDocumentContentChangeEvent
	for _, change := range params.ContentChanges {
		if changeEvent, ok := change.(protocol.TextDocumentContentChangeEvent); ok {
			changes = append(changes, changeEvent)
		}
	}
	
	// Update document manager
	content := s.docManager.DidChange(uri, changes)
	
	s.log.Info("Document changed", "uri", uri, "length", len(content))
	
	// Run diagnostics
	s.diagnostics.PublishDiagnostics(context, uri, content)
	
	return nil
}

func (s *Server) textDocumentSemanticTokensFull(context *glsp.Context, params *protocol.SemanticTokensParams) (*protocol.SemanticTokens, error) {
	uri := params.TextDocument.URI
	content, exists := s.docManager.GetContent(uri)
	if !exists {
		s.log.Error("Document not found for semantic tokens", "uri", uri)
		return &protocol.SemanticTokens{Data: []uint32{}}, nil
	}
	
	return s.highlighting.GetSemanticTokens(content)
}

func (s *Server) textDocumentFormatting(context *glsp.Context, params *protocol.DocumentFormattingParams) ([]protocol.TextEdit, error) {
	uri := params.TextDocument.URI
	content, exists := s.docManager.GetContent(uri)
	if !exists {
		s.log.Error("Document not found for formatting", "uri", uri)
		return []protocol.TextEdit{}, nil
	}
	
	formatted, err := s.formatting.FormatDocument(content)
	if err != nil {
		s.log.Error("Failed to format document", "uri", uri, "error", err)
		return []protocol.TextEdit{}, err
	}
	
	// Update document cache with formatted content
	s.docManager.UpdateContent(uri, formatted)
	
	return s.formatting.CreateTextEdit(content, formatted), nil
}