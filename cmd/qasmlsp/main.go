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
	version   string = "0.0.1"
	handler   protocol.Handler
	log       = commonlog.GetLogger("qasm")
	documents = make(map[protocol.DocumentUri]string)
)

func main() {
	commonlog.Configure(1, nil)

	handler = protocol.Handler{
		Initialize:                     initialize,
		Initialized:                    initialized,
		Shutdown:                       shutdown,
		SetTrace:                       setTrace,
		TextDocumentDidOpen:            textDocumentDidOpen,
		TextDocumentDidChange:          textDocumentDidChange,
		TextDocumentSemanticTokensFull: textDocumentSemanticTokensFull,
	}

	server := server.NewServer(&handler, lsName, false)
	server.RunStdio()
}

func initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	capabilities := handler.CreateServerCapabilities()
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
	documents[params.TextDocument.URI] = params.TextDocument.Text
	return nil
}

func textDocumentDidChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	log.Info("Document changed", "uri", params.TextDocument.URI)
	if len(params.ContentChanges) > 0 {
		if change, ok := params.ContentChanges[0].(protocol.TextDocumentContentChangeEvent); ok {
			documents[params.TextDocument.URI] = change.Text
		}
	}
	return nil
}

func textDocumentSemanticTokensFull(context *glsp.Context, params *protocol.SemanticTokensParams) (*protocol.SemanticTokens, error) {
	log.Info("Semantic tokens requested", "uri", params.TextDocument.URI)

	content, exists := documents[params.TextDocument.URI]
	if !exists {
		log.Error("Document not found", "uri", params.TextDocument.URI)
		return &protocol.SemanticTokens{Data: []uint32{}}, nil
	}

	return computeSemanticTokens(content)
}

func computeSemanticTokens(content string) (*protocol.SemanticTokens, error) {
	h := highlight.New()
	_, err := h.Highlight(content)
	if err != nil {
		return nil, err
	}

	tokens := h.GetTokens()

	// Sort tokens by line and column
	for i := 0; i < len(tokens)-1; i++ {
		for j := i + 1; j < len(tokens); j++ {
			if tokens[i].Line > tokens[j].Line ||
				(tokens[i].Line == tokens[j].Line && tokens[i].Column > tokens[j].Column) {
				tokens[i], tokens[j] = tokens[j], tokens[i]
			}
		}
	}

	data := make([]uint32, 0, len(tokens)*5)
	var prevLine uint32 = 0
	var prevChar uint32 = 0

	for _, token := range tokens {
		if token.Line <= 0 || token.Column < 0 {
			continue // Skip invalid tokens
		}

		// Calculate delta line and delta start
		currentLine := uint32(token.Line - 1) // Convert to 0-based
		deltaLine := currentLine - prevLine

		var deltaStart uint32
		if deltaLine == 0 {
			deltaStart = uint32(token.Column) - prevChar
		} else {
			deltaStart = uint32(token.Column)
		}

		// Get token type index
		tokenTypeIndex, exists := tokenTypeMap[token.TypeName]
		if !exists {
			tokenTypeIndex = tokenTypeMap["identifier"] // Default fallback
		}

		// Add token data (5 values per token)
		data = append(data,
			deltaLine,
			deltaStart,
			uint32(token.Length),
			uint32(tokenTypeIndex),
			0, // No modifiers
		)

		prevLine = currentLine
		if deltaLine == 0 {
			prevChar = uint32(token.Column) + uint32(token.Length)
		} else {
			prevChar = uint32(token.Column) + uint32(token.Length)
		}
	}

	return &protocol.SemanticTokens{
		Data: data,
	}, nil
}

// tokenTypeMap maps token type names to their indices in the legend array
// This MUST match the order in the SemanticTokensLegend TokenTypes array
var tokenTypeMap = map[string]int{
	"keyword":           0,  // "keyword"
	"operator":          1,  // "operator"
	"identifier":        2,  // "identifier"
	"number":            3,  // "number"
	"string":            4,  // "string"
	"comment":           5,  // "comment"
	"gate":              6,  // "gate"
	"measurement":       7,  // "measurement"
	"register":          8,  // "register"
	"punctuation":       9,  // "punctuation"
	"builtin_gate":      10, // "builtin_gate"
	"builtin_quantum":   11, // "builtin_quantum"
	"builtin_classical": 12, // "builtin_classical"
	"builtin_constant":  13, // "builtin_constant"
	"access_control":    14, // "access_control"
	"extern":            15, // "extern"
	"hardware_qubit":    16, // "hardware_qubit"
}
