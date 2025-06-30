package main

import (
	"fmt"
	"strings"
	"syscall/js"

	"github.com/orangekame3/qasmtools/formatter"
	"github.com/orangekame3/qasmtools/highlight"
)

func main() {
	// Defer panic recovery
	defer func() {
		if r := recover(); r != nil {
			// Log the panic to JavaScript console for debugging
			js.Global().Get("console").Call("error", fmt.Sprintf("WASM panic: %v", r))
		}
	}()

	// Export formatQASM function to JavaScript
	formatFunc := js.FuncOf(formatQASM)
	js.Global().Set("formatQASM", formatFunc)
	
	// Export highlightQASM function to JavaScript  
	// Temporarily disabled to debug issues
	highlightFunc := js.FuncOf(highlightQASM)
	js.Global().Set("highlightQASM", highlightFunc)

	// Signal that functions are ready
	js.Global().Set("qasmToolsReady", js.ValueOf(true))

	// Log successful initialization
	js.Global().Get("console").Call("log", "QASM Tools WASM initialized successfully")

	// Keep the program running indefinitely
	select {}
}

func formatQASM(this js.Value, args []js.Value) interface{} {
	defer func() {
		if r := recover(); r != nil {
			js.Global().Get("console").Call("error", fmt.Sprintf("formatQASM panic: %v", r))
		}
	}()

	if len(args) != 1 {
		return map[string]interface{}{
			"success": false,
			"error":   "Expected exactly one argument (QASM code)",
		}
	}

	qasmCode := args[0].String()
	if strings.TrimSpace(qasmCode) == "" {
		return map[string]interface{}{
			"success": false,
			"error":   "Input QASM code is empty",
		}
	}

	// Create formatter with default configuration
	f := formatter.NewFormatter()

	// Format the QASM code
	formatted, err := f.Format(qasmCode)
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"error":   fmt.Sprintf("Failed to format QASM: %v", err),
		}
	}

	return map[string]interface{}{
		"success":   true,
		"formatted": formatted,
	}
}

func highlightQASM(this js.Value, args []js.Value) interface{} {
	defer func() {
		if r := recover(); r != nil {
			js.Global().Get("console").Call("error", fmt.Sprintf("highlightQASM panic: %v", r))
		}
	}()

	if len(args) != 1 {
		return map[string]interface{}{
			"success": false,
			"error":   "Expected exactly one argument (QASM code)",
		}
	}

	qasmCode := args[0].String()
	if strings.TrimSpace(qasmCode) == "" {
		return map[string]interface{}{
			"success": false,
			"error":   "Input QASM code is empty",
		}
	}

	// Create highlighter
	h := highlight.New()

	// Highlight the QASM code to get tokens
	_, err := h.Highlight(qasmCode)
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"error":   fmt.Sprintf("Failed to highlight QASM: %v", err),
		}
	}

	// Get tokens for JavaScript
	tokens := h.GetTokens()

	return map[string]interface{}{
		"success": true,
		"tokens":  convertTokensToJS(tokens),
	}
}

// convertTokensToJS converts tokens to a JavaScript-friendly format
func convertTokensToJS(tokens []highlight.TokenInfo) []map[string]interface{} {
	result := make([]map[string]interface{}, len(tokens))
	for i, token := range tokens {
		result[i] = map[string]interface{}{
			"type":     token.TypeName,
			"content":  token.Content,
			"line":     token.Line,
			"column":   token.Column,
			"length":   token.Length,
		}
	}
	return result
}