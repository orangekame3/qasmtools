package features

import (
	"strings"

	"github.com/tliron/commonlog"
	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/orangekame3/qasmtools/formatter"
)

// FormattingProvider handles document formatting
type FormattingProvider struct {
	formatter *formatter.Formatter
	log       commonlog.Logger
}

// NewFormattingProvider creates a new formatting provider
func NewFormattingProvider(formatter *formatter.Formatter, log commonlog.Logger) *FormattingProvider {
	return &FormattingProvider{
		formatter: formatter,
		log:       log,
	}
}

// FormatDocument formats a document and returns the formatted content
func (f *FormattingProvider) FormatDocument(content string) (string, error) {
	if f.formatter == nil {
		f.log.Error("Formatter not available")
		return content, nil
	}
	
	formatted, err := f.formatter.Format(content)
	if err != nil {
		f.log.Error("Failed to format document", "error", err)
		return content, err
	}
	
	return formatted, nil
}

// CreateTextEdit creates a text edit that replaces the entire document
func (f *FormattingProvider) CreateTextEdit(original, formatted string) []protocol.TextEdit {
	if original == formatted {
		return []protocol.TextEdit{}
	}
	
	// Count lines in original content
	lines := strings.Split(original, "\n")
	endLine := len(lines) - 1
	endChar := len(lines[endLine])
	
	return []protocol.TextEdit{
		{
			Range: protocol.Range{
				Start: protocol.Position{Line: 0, Character: 0},
				End:   protocol.Position{Line: uint32(endLine), Character: uint32(endChar)},
			},
			NewText: formatted,
		},
	}
}