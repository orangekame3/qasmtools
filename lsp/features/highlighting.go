package features

import (
	"github.com/tliron/commonlog"
	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/orangekame3/qasmtools/highlight"
)

// HighlightingProvider handles semantic token highlighting
type HighlightingProvider struct {
	highlighter *highlight.Highlighter
	log         commonlog.Logger
}

// NewHighlightingProvider creates a new highlighting provider
func NewHighlightingProvider(highlighter *highlight.Highlighter, log commonlog.Logger) *HighlightingProvider {
	return &HighlightingProvider{
		highlighter: highlighter,
		log:         log,
	}
}

// GetSemanticTokens generates semantic tokens for syntax highlighting
func (h *HighlightingProvider) GetSemanticTokens(content string) (*protocol.SemanticTokens, error) {
	if h.highlighter == nil {
		h.log.Error("Highlighter not available")
		return &protocol.SemanticTokens{Data: []uint32{}}, nil
	}
	
	_, err := h.highlighter.Highlight(content)
	if err != nil {
		h.log.Error("Failed to generate semantic tokens", "error", err)
		return &protocol.SemanticTokens{Data: []uint32{}}, err
	}
	
	// Get tokens from highlighter
	tokens := h.highlighter.GetTokens()
	
	// Convert highlight tokens to LSP semantic tokens format
	data := h.convertToSemanticTokensData(tokens)
	
	return &protocol.SemanticTokens{Data: data}, nil
}

// convertToSemanticTokensData converts highlight tokens to LSP format
func (h *HighlightingProvider) convertToSemanticTokensData(tokens []highlight.TokenInfo) []uint32 {
	var data []uint32
	
	var lastLine, lastChar uint32
	
	for _, token := range tokens {
		// Calculate delta values for LSP format
		deltaLine := uint32(token.Line) - lastLine
		deltaChar := uint32(token.Column)
		if deltaLine == 0 {
			deltaChar = uint32(token.Column) - lastChar
		}
		
		// LSP semantic tokens format: [deltaLine, deltaChar, length, tokenType, tokenModifiers]
		data = append(data,
			deltaLine,
			deltaChar,
			uint32(token.Length),
			h.getTokenTypeIndex(token.TypeName),
			0, // No modifiers for now
		)
		
		lastLine = uint32(token.Line)
		lastChar = uint32(token.Column)
	}
	
	return data
}

// getTokenTypeIndex maps token types to indices
func (h *HighlightingProvider) getTokenTypeIndex(tokenType string) uint32 {
	// Map token types to their indices in the legend
	typeMap := map[string]uint32{
		"keyword":           0,
		"operator":          1,
		"identifier":        2,
		"number":            3,
		"string":            4,
		"comment":           5,
		"gate":              6,
		"measurement":       7,
		"register":          8,
		"punctuation":       9,
		"builtin_gate":      10,
		"builtin_quantum":   11,
		"builtin_classical": 12,
		"builtin_constant":  13,
		"access_control":    14,
		"extern":            15,
		"hardware_qubit":    16,
	}
	
	if index, exists := typeMap[tokenType]; exists {
		return index
	}
	
	return 2 // Default to "identifier"
}