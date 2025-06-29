package main

import (
	"github.com/tliron/commonlog"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"

	"github.com/orangekame3/qasmtools/highlight"
	_ "github.com/tliron/commonlog/simple"
)

const lsName = "qasm"

var (
	version string = "0.0.1"
	handler protocol.Handler
	log     = commonlog.GetLogger("qasm")
)

func main() {
	commonlog.Configure(1, nil)

	handler = protocol.Handler{
		Initialize:           initialize,
		Initialized:         initialized,
		Shutdown:           shutdown,
		SetTrace:           setTrace,
		TextDocumentDidOpen: textDocumentDidOpen,
		TextDocumentDidChange: textDocumentDidChange,
	}

	server := server.NewServer(&handler, lsName, false)
	server.RunStdio()
}

func initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	capabilities := handler.CreateServerCapabilities()
	capabilities.SemanticTokensProvider = &protocol.SemanticTokensOptions{
		Legend: protocol.SemanticTokensLegend{
			TokenTypes:     []string{"keyword", "operator", "string", "number", "type", "variable", "function"},
			TokenModifiers: []string{},
		},
		Full: protocol.True,
	}

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    lsName,
			Version: &version,
		},
	}, nil
}

func initialized(context *glsp.Context, params *protocol.InitializedParams) error {
	return nil
}

func shutdown(context *glsp.Context) error {
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func setTrace(context *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}

func textDocumentDidOpen(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	log.Info("Document opened", "uri", params.TextDocument.URI)
	return updateHighlighting(context, params.TextDocument.URI, params.TextDocument.Text)
}

func textDocumentDidChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	log.Info("Document changed", "uri", params.TextDocument.URI)
	if len(params.ContentChanges) > 0 {
		if change, ok := params.ContentChanges[0].(protocol.TextDocumentContentChangeEvent); ok {
			return updateHighlighting(context, params.TextDocument.URI, change.Text)
		}
	}
	return nil
}

func updateHighlighting(context *glsp.Context, uri protocol.DocumentUri, content string) error {
	log.Info("Updating highlighting", "uri", uri)
	tokens, err := computeSemanticTokens(content)
	if err != nil {
		return err
	}

	// Send semantic tokens to the client
	context.Notify("textDocument/semanticTokens/full", tokens)
	return nil
}

func computeSemanticTokens(content string) (*protocol.SemanticTokens, error) {
	h := highlight.New()
	_, err := h.Highlight(content)
	if err != nil {
		return nil, err
	}

	tokens := h.GetTokens()
	data := make([]uint32, 0, len(tokens)*5)

	var prevLine uint32 = 0
	var prevChar uint32 = 0

	for _, token := range tokens {
		// Calculate delta line and delta start
		deltaLine := uint32(token.Line - 1) - prevLine
		var deltaStart uint32
		if deltaLine == 0 {
			deltaStart = uint32(token.Column) - prevChar
		} else {
			deltaStart = uint32(token.Column)
		}

		// Get token type
		tokenType := uint32(tokenTypeMap[token.TypeName])

		// Add token data
		data = append(data,
			deltaLine,
			deltaStart,
			uint32(token.Length),
			tokenType,
			0, // No modifiers
		)

		prevLine = uint32(token.Line - 1)
		prevChar = uint32(token.Column)
	}

	return &protocol.SemanticTokens{
		Data: data,
	}, nil
}

var tokenTypeMap = map[string]int{
	"keyword":           0,
	"operator":          1,
	"identifier":        2,
	"number":            3,
	"string":            4,
	"comment":          5,
	"gate":             6,
	"measurement":      7,
	"register":         8,
	"punctuation":      9,
	"builtin_gate":     10,
	"builtin_quantum":  11,
	"builtin_classical": 12,
	"builtin_constant": 13,
	"access_control":   14,
	"extern":           15,
	"hardware_qubit":   16,
}
