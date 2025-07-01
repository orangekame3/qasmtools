package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/tliron/commonlog"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"

	"github.com/orangekame3/qasmtools/formatter"
	"github.com/orangekame3/qasmtools/highlight"
	"github.com/orangekame3/qasmtools/lint"
	_ "github.com/tliron/commonlog/simple"
)

const lsName = "qasm"

var (
	version            string = "0.0.1"
	handler            protocol.Handler
	log                = commonlog.GetLogger("qasm")
	documents          = make(map[protocol.DocumentUri]string)
	lastFormatTime     = make(map[protocol.DocumentUri]time.Time)
	formatMutex        sync.RWMutex
	linter             *lint.Linter
	recentlyFormatted  = make(map[protocol.DocumentUri]time.Time)
	formattingMutex    sync.RWMutex
)

func main() {
	commonlog.Configure(1, nil) // Lower log level for more verbose output
	log.Info("qasmlsp server starting")

	// Initialize linter
	linter = lint.NewLinter("")
	if err := linter.LoadRules(); err != nil {
		log.Error("Failed to load linting rules", "error", err)
	} else {
		rules := linter.GetRules()
		log.Info("Linting rules loaded successfully", "rule_count", len(rules))
		for i, rule := range rules {
			log.Info("Rule loaded", "index", i, "id", rule.ID, "name", rule.Name, "enabled", rule.Enabled)
		}
	}

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

	// Run linting on newly opened document
	if linter != nil {
		log.Info("Running linting on opened document", "uri", params.TextDocument.URI)
		violations := runLinting(params.TextDocument.Text, string(params.TextDocument.URI))
		log.Info("Linting completed", "violations_count", len(violations))
		publishDiagnostics(context, params.TextDocument.URI, violations)
	}

	return nil
}

func textDocumentDidChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	log.Info("Document changed", "uri", params.TextDocument.URI, "version", params.TextDocument.Version, "changes", len(params.ContentChanges))

	// Check if this document was recently formatted - if so, ignore incremental changes for a short period
	formattingMutex.RLock()
	if lastFormat, exists := recentlyFormatted[params.TextDocument.URI]; exists {
		if time.Since(lastFormat) < 2*time.Second {
			log.Info("Ignoring changes shortly after formatting to prevent feedback loop", "time_since_format", time.Since(lastFormat))
			formattingMutex.RUnlock()
			return nil
		}
	}
	formattingMutex.RUnlock()

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

					// Run linting on updated document content
					if linter != nil {
						log.Info("Running linting on changed document", "uri", params.TextDocument.URI)
						violations := runLinting(change.Text, string(params.TextDocument.URI))
						log.Info("Linting completed", "violations_count", len(violations))
						publishDiagnostics(context, params.TextDocument.URI, violations)
					}
					return nil
				}
			}
		}

		// Handle incremental changes - apply them to the document and run linting
		currentContent, exists := documents[params.TextDocument.URI]
		if !exists {
			log.Error("Document not found for incremental changes", "uri", params.TextDocument.URI)
			return nil
		}

		// Apply incremental changes
		updatedContent := applyIncrementalChanges(currentContent, params.ContentChanges)
		documents[params.TextDocument.URI] = updatedContent
		log.Info("Applied incremental changes", "new_length", len(updatedContent))

		// Run linting on updated content
		if linter != nil {
			log.Info("Running linting on incrementally changed document", "uri", params.TextDocument.URI)
			violations := runLinting(updatedContent, string(params.TextDocument.URI))
			log.Info("Linting completed", "violations_count", len(violations))
			publishDiagnostics(context, params.TextDocument.URI, violations)
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

	// Mark document as recently formatted to prevent feedback loop
	formattingMutex.Lock()
	recentlyFormatted[params.TextDocument.URI] = time.Now()
	formattingMutex.Unlock()

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

// convertViolationsToDiagnostics converts lint violations to LSP diagnostics
func convertViolationsToDiagnostics(violations []*lint.Violation) []protocol.Diagnostic {
	var diagnostics []protocol.Diagnostic

	for _, violation := range violations {
		severity := protocol.DiagnosticSeverityInformation
		switch violation.Severity {
		case lint.SeverityError:
			severity = protocol.DiagnosticSeverityError
		case lint.SeverityWarning:
			severity = protocol.DiagnosticSeverityWarning
		case lint.SeverityInfo:
			severity = protocol.DiagnosticSeverityInformation
		}

		// Convert 1-based line/column to 0-based for LSP
		line := uint32(violation.Line - 1)
		column := uint32(violation.Column)

		// Create code as IntegerOrString
		code := protocol.IntegerOrString{Value: violation.Rule.ID}
		source := "qasm-lint"

		diagnostic := protocol.Diagnostic{
			Range: protocol.Range{
				Start: protocol.Position{Line: line, Character: column},
				End:   protocol.Position{Line: line, Character: column + 10}, // Estimate end position
			},
			Severity: &severity,
			Code:     &code,
			Source:   &source,
			Message:  violation.Message,
		}

		diagnostics = append(diagnostics, diagnostic)
	}

	return diagnostics
}

// publishDiagnostics sends diagnostics to the client
func publishDiagnostics(context *glsp.Context, uri protocol.DocumentUri, violations []*lint.Violation) {
	diagnostics := convertViolationsToDiagnostics(violations)

	params := protocol.PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: diagnostics,
	}

	context.Notify(protocol.ServerTextDocumentPublishDiagnostics, params)
	log.Info("Published diagnostics", "uri", uri, "count", len(diagnostics))
}

// runLinting runs linting on document content and returns violations
func runLinting(content string, filename string) []*lint.Violation {
	if linter == nil {
		log.Error("Linter not initialized")
		return nil
	}

	log.Info("Starting lint process", "filename", filename, "content_length", len(content))
	log.Info("Linter state", "linter_ptr", fmt.Sprintf("%p", linter), "rules_count", len(linter.GetRules()))

	violations, err := linter.LintContent(content, filename)
	if err != nil {
		log.Error("Linting failed", "error", err)
		return nil
	}

	log.Info("Lint process completed", "violations_found", len(violations))
	for i, violation := range violations {
		log.Info("Violation found", "index", i, "rule", violation.Rule.ID, "message", violation.Message, "line", violation.Line, "column", violation.Column)
	}

	return violations
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

// applyIncrementalChanges applies a series of incremental changes to document content
func applyIncrementalChanges(content string, changes []interface{}) string {
	lines := splitLines(content)
	
	// Sort changes by position (reverse order to maintain positions)
	var textChanges []protocol.TextDocumentContentChangeEvent
	for _, changeInterface := range changes {
		if change, ok := changeInterface.(protocol.TextDocumentContentChangeEvent); ok && change.Range != nil {
			textChanges = append(textChanges, change)
		}
	}
	
	// Apply changes in reverse order to maintain line/column positions
	for i := len(textChanges) - 1; i >= 0; i-- {
		change := textChanges[i]
		lines = applyTextChange(lines, change)
	}
	
	return strings.Join(lines, "\n")
}

// applyTextChange applies a single text change to the lines array
func applyTextChange(lines []string, change protocol.TextDocumentContentChangeEvent) []string {
	if change.Range == nil {
		// Full document replacement
		return splitLines(change.Text)
	}
	
	startLine := int(change.Range.Start.Line)
	startChar := int(change.Range.Start.Character)
	endLine := int(change.Range.End.Line)
	endChar := int(change.Range.End.Character)
	
	// Ensure we don't go out of bounds
	if startLine >= len(lines) {
		// Extend lines if necessary
		for len(lines) <= startLine {
			lines = append(lines, "")
		}
	}
	
	if endLine >= len(lines) {
		for len(lines) <= endLine {
			lines = append(lines, "")
		}
	}
	
	newText := change.Text
	newLines := splitLines(newText)
	
	// Build the result
	var result []string
	
	// Add lines before the change
	result = append(result, lines[:startLine]...)
	
	// Handle the change
	if startLine == endLine {
		// Single line change
		line := lines[startLine]
		if startChar > len(line) {
			startChar = len(line)
		}
		if endChar > len(line) {
			endChar = len(line)
		}
		
		newLine := line[:startChar] + strings.Join(newLines, "\n") + line[endChar:]
		if len(newLines) == 1 {
			result = append(result, newLine)
		} else {
			// Multi-line replacement in single line
			firstNewLine := line[:startChar] + newLines[0]
			result = append(result, firstNewLine)
			result = append(result, newLines[1:len(newLines)-1]...)
			lastNewLine := newLines[len(newLines)-1] + line[endChar:]
			result = append(result, lastNewLine)
		}
	} else {
		// Multi-line change
		startLineText := lines[startLine]
		endLineText := lines[endLine]
		
		if startChar > len(startLineText) {
			startChar = len(startLineText)
		}
		if endChar > len(endLineText) {
			endChar = len(endLineText)
		}
		
		if len(newLines) == 1 {
			// Replace multiple lines with single line
			newLine := startLineText[:startChar] + newLines[0] + endLineText[endChar:]
			result = append(result, newLine)
		} else {
			// Replace multiple lines with multiple lines
			firstLine := startLineText[:startChar] + newLines[0]
			result = append(result, firstLine)
			result = append(result, newLines[1:len(newLines)-1]...)
			lastLine := newLines[len(newLines)-1] + endLineText[endChar:]
			result = append(result, lastLine)
		}
	}
	
	// Add lines after the change
	result = append(result, lines[endLine+1:]...)
	
	return result
}
