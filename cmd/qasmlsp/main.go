package main

import (
	"sort"
	"sync"
	"time"

	"github.com/tliron/commonlog"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"

	"github.com/orangekame3/qasmtools/formatter"
	"github.com/orangekame3/qasmtools/highlight"
	_ "github.com/tliron/commonlog/simple"
)

const lsName = "qasm"

var (
	version        string = "0.0.1"
	handler        protocol.Handler
	log            = commonlog.GetLogger("qasm")
	documents      = make(map[protocol.DocumentUri]string)
	lastFormatTime = make(map[protocol.DocumentUri]time.Time)
	formatMutex    sync.RWMutex
)

func main() {
	commonlog.Configure(1, nil) // Lower log level for more verbose output
	log.Info("qasmlsp server starting")

	handler = protocol.Handler{
		Initialize:                     initialize,
		Initialized:                    initialized,
		Shutdown:                       shutdown,
		TextDocumentDidOpen:            textDocumentDidOpen,
		TextDocumentDidChange:          textDocumentDidChange,
		TextDocumentSemanticTokensFull: textDocumentSemanticTokensFull,
		TextDocumentFormatting:         textDocumentFormatting,
	}

	server := server.NewServer(&handler, lsName, true)
	log.Info("Starting QASM Language Server")
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
	capabilities.DocumentFormattingProvider = protocol.True

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
	log.Info("Document content stored", "uri", params.TextDocument.URI, "content_length", len(params.TextDocument.Text))
	return nil
}

func textDocumentDidChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	log.Info("Document changed", "uri", params.TextDocument.URI, "version", params.TextDocument.Version, "changes", len(params.ContentChanges))

	if len(params.ContentChanges) > 0 {
		// Look for full document changes first
		for i, changeInterface := range params.ContentChanges {
			if change, ok := changeInterface.(protocol.TextDocumentContentChangeEvent); ok {
				log.Info("Processing change", "index", i, "text_length", len(change.Text), "has_range", change.Range != nil)

				// Check if this is a full document change (no range specified)
				if change.Range == nil {
					// Full document replacement - this is what we want for format results
					documents[params.TextDocument.URI] = change.Text
					log.Info("Full document replacement applied", "new_length", len(change.Text))
					return nil
				}
			}
		}

		// If we only have incremental changes from formatting, don't update the cache
		// The formatter has already updated the cache, so incremental changes might be inconsistent
		log.Info("Only incremental changes detected - skipping document cache update to preserve formatter consistency")
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

	sort.Slice(tokens, func(i, j int) bool {
		if tokens[i].Line != tokens[j].Line {
			return tokens[i].Line < tokens[j].Line
		}
		return tokens[i].Column < tokens[j].Column
	})

	log.Info("Tokens:", "tokens", h.GetTokens())

	data := make([]uint32, 0, len(tokens)*5)
	var prevLine uint32 = 0
	var prevChar uint32 = 0

	for _, token := range tokens {
		if token.Line <= 0 || token.Column < 0 {
			continue // Skip invalid tokens
		}

		currentLine := uint32(token.Line - 1) // Convert to 0-based
		deltaLine := currentLine - prevLine

		var deltaStart uint32
		if deltaLine == 0 {
			deltaStart = uint32(token.Column) - prevChar
		} else {
			deltaStart = uint32(token.Column)
		}

		tokenTypeIndex, exists := tokenTypeMap[token.TypeName]
		if !exists {
			tokenTypeIndex = tokenTypeMap["identifier"] // Default fallback
		}

		data = append(data,
			deltaLine,
			deltaStart,
			uint32(token.Length),
			uint32(tokenTypeIndex),
			0, // No modifiers
		)

		prevLine = currentLine
		prevChar = uint32(token.Column)
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

func textDocumentFormatting(context *glsp.Context, params *protocol.DocumentFormattingParams) ([]protocol.TextEdit, error) {
	log.Info("=== FORMATTING REQUEST START ===", "uri", params.TextDocument.URI)

	// Prevent rapid consecutive formatting requests
	formatMutex.Lock()
	now := time.Now()
	if lastTime, exists := lastFormatTime[params.TextDocument.URI]; exists {
		if now.Sub(lastTime) < 500*time.Millisecond {
			log.Info("Skipping format request - too soon after last request", "uri", params.TextDocument.URI, "time_since_last", now.Sub(lastTime))
			formatMutex.Unlock()
			return []protocol.TextEdit{}, nil
		}
	}
	lastFormatTime[params.TextDocument.URI] = now
	formatMutex.Unlock()

	log.Info("Available documents", "count", len(documents))
	for uri, content := range documents {
		log.Info("Document in cache", "uri", uri, "length", len(content))
	}

	content, exists := documents[params.TextDocument.URI]
	if !exists {
		log.Error("Document not found in cache", "uri", params.TextDocument.URI, "available_uris", getDocumentURIs())
		return []protocol.TextEdit{}, nil
	}

	log.Info("Found document content", "content_length", len(content))
	if len(content) > 0 {
		preview := content
		if len(content) > 200 {
			preview = content[:200] + "..."
		}
		log.Info("Content preview", "content", preview)
	}

	f := formatter.NewFormatter()
	formatted, err := f.Format(content)
	if err != nil {
		log.Error("Formatting failed", "error", err)
		return []protocol.TextEdit{}, err
	}

	log.Info("Formatting completed", "original_length", len(content), "formatted_length", len(formatted))

	// If content is the same, return empty edits
	if content == formatted {
		log.Info("No formatting changes needed - content is identical")
		return []protocol.TextEdit{}, nil
	}

	log.Info("Content differs, preparing text edit")

	lines := splitLines(content)
	endLine := uint32(len(lines))
	endChar := uint32(0)
	if endLine > 0 {
		endChar = uint32(len(lines[endLine-1]))
	} else {
		endLine = 1
	}

	edit := protocol.TextEdit{
		Range: protocol.Range{
			Start: protocol.Position{Line: 0, Character: 0},
			End:   protocol.Position{Line: endLine, Character: endChar},
		},
		NewText: formatted,
	}

	log.Info("Text edit prepared", "start_line", edit.Range.Start.Line, "start_char", edit.Range.Start.Character,
		"end_line", edit.Range.End.Line, "end_char", edit.Range.End.Character, "newtext_length", len(edit.NewText))

	// Update the document cache with the formatted content to ensure consistency
	documents[params.TextDocument.URI] = formatted
	log.Info("Updated document cache with formatted content", "length", len(formatted))

	log.Info("=== FORMATTING REQUEST END ===")

	return []protocol.TextEdit{edit}, nil
}

func getDocumentURIs() []string {
	uris := make([]string, 0, len(documents))
	for uri := range documents {
		uris = append(uris, string(uri))
	}
	return uris
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func splitLines(text string) []string {
	lines := []string{}
	start := 0
	for i, char := range text {
		if char == '\n' {
			lines = append(lines, text[start:i])
			start = i + 1
		}
	}
	if start < len(text) {
		lines = append(lines, text[start:])
	}
	return lines
}
