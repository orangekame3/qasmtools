package features

import (
	"github.com/tliron/commonlog"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/orangekame3/qasmtools/lint"
)

// DiagnosticsProvider handles linting and diagnostic publishing
type DiagnosticsProvider struct {
	linter *lint.Linter
	log    commonlog.Logger
}

// NewDiagnosticsProvider creates a new diagnostics provider
func NewDiagnosticsProvider(linter *lint.Linter, log commonlog.Logger) *DiagnosticsProvider {
	return &DiagnosticsProvider{
		linter: linter,
		log:    log,
	}
}

// PublishDiagnostics runs linting and publishes diagnostics for a document
func (d *DiagnosticsProvider) PublishDiagnostics(context *glsp.Context, uri protocol.DocumentUri, content string) {
	// Run linting
	violations := d.runLinting(content, string(uri))
	
	// Convert violations to diagnostics
	diagnostics := d.convertViolationsToDiagnostics(violations)
	
	// Publish diagnostics
	d.log.Info("Publishing diagnostics", "uri", uri, "count", len(diagnostics))
	context.Notify(protocol.ServerTextDocumentPublishDiagnostics, &protocol.PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: diagnostics,
	})
}

// runLinting executes the linter on the given content
func (d *DiagnosticsProvider) runLinting(content, filename string) []*lint.Violation {
	if d.linter == nil {
		d.log.Error("Linter not available")
		return nil
	}
	
	violations, err := d.linter.LintContent(content, filename)
	if err != nil {
		d.log.Error("Failed to run linting", "error", err)
		return nil
	}
	
	return violations
}

// convertViolationsToDiagnostics converts lint violations to LSP diagnostics
func (d *DiagnosticsProvider) convertViolationsToDiagnostics(violations []*lint.Violation) []protocol.Diagnostic {
	var diagnostics []protocol.Diagnostic
	
	for _, violation := range violations {
		severity := d.convertSeverity(string(violation.Severity))
		code := protocol.IntegerOrString{Value: violation.Rule.ID}
		source := "qasm-lint"
		
		diagnostic := protocol.Diagnostic{
			Range: protocol.Range{
				Start: protocol.Position{
					Line:      uint32(violation.Line - 1), // LSP uses 0-based line numbers
					Character: uint32(violation.Column - 1), // LSP uses 0-based column numbers
				},
				End: protocol.Position{
					Line:      uint32(violation.Line - 1),
					Character: uint32(violation.Column), // End at next character for highlighting
				},
			},
			Severity: &severity,
			Code:     &code,
			Source:   &source,
			Message:  violation.Message,
		}
		
		// Add documentation URL if available
		if violation.Rule.DocumentationURL != "" {
			diagnostic.CodeDescription = &protocol.CodeDescription{
				HRef: violation.Rule.DocumentationURL,
			}
		}
		
		diagnostics = append(diagnostics, diagnostic)
	}
	
	return diagnostics
}

// convertSeverity converts lint severity to LSP diagnostic severity
func (d *DiagnosticsProvider) convertSeverity(severity string) protocol.DiagnosticSeverity {
	switch severity {
	case "error":
		return protocol.DiagnosticSeverityError
	case "warning":
		return protocol.DiagnosticSeverityWarning
	case "info":
		return protocol.DiagnosticSeverityInformation
	default:
		return protocol.DiagnosticSeverityWarning
	}
}