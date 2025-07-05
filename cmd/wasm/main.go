package main

import (
	"fmt"
	"strconv"
	"strings"
	"syscall/js"

	"github.com/orangekame3/qasmtools/formatter"
	"github.com/orangekame3/qasmtools/highlight"
	"github.com/orangekame3/qasmtools/lint"
)

// Global variables
var (
	globalLinter  *lint.Linter
	formatFunc    js.Func
	highlightFunc js.Func
	lintFunc      js.Func
	stopFunc      js.Func
)

func main() {
	// Defer panic recovery
	defer func() {
		if r := recover(); r != nil {
			// Log the panic to JavaScript console for debugging
			js.Global().Get("console").Call("error", fmt.Sprintf("WASM panic: %v", r))
		}
	}()

	// Initialize global linter with empty path to use embedded rules
	globalLinter = lint.NewLinter("")
	if err := globalLinter.LoadRules(); err != nil {
		js.Global().Get("console").Call("error", fmt.Sprintf("Failed to load rules: %v", err))
		return
	}

	// Initialize and export global functions
	formatFunc = js.FuncOf(formatQASM)
	js.Global().Set("formatQASM", formatFunc)

	highlightFunc = js.FuncOf(highlightQASM)
	js.Global().Set("highlightQASM", highlightFunc)

	lintFunc = js.FuncOf(lintQASM)
	js.Global().Set("lintQASM", lintFunc)

	// Signal that functions are ready
	js.Global().Set("qasmToolsReady", js.ValueOf(true))

	// Log successful initialization
	js.Global().Get("console").Call("log", "QASM Tools WASM initialized successfully")

	// Keep the program running
	select {}
}

func formatQASM(this js.Value, args []js.Value) interface{} {
	defer func() {
		if r := recover(); r != nil {
			js.Global().Get("console").Call("error", fmt.Sprintf("formatQASM panic: %v", r))
		}
	}()

	if len(args) < 1 || len(args) > 2 {
		return map[string]interface{}{
			"success": false,
			"error":   "Expected 1-2 arguments (QASM code, optional unescape flag)",
		}
	}

	qasmCode := args[0].String()
	js.Global().Get("console").Call("log", "Debug: Raw QASM code:", qasmCode)
	if strings.TrimSpace(qasmCode) == "" {
		return map[string]interface{}{
			"success": false,
			"error":   "Input QASM code is empty",
		}
	}

	// Check if unescape is requested (second argument)
	unescape := false
	if len(args) > 1 && !args[1].IsNull() && !args[1].IsUndefined() {
		unescape = args[1].Bool()
	}

	// Apply unescaping if requested
	if unescape {
		unescaped, err := strconv.Unquote(qasmCode)
		if err != nil {
			return map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("Failed to unescape input: %v", err),
			}
		}
		qasmCode = unescaped
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
			"type":    token.TypeName,
			"content": token.Content,
			"line":    token.Line,
			"column":  token.Column,
			"length":  token.Length,
		}
	}
	return result
}

func lintQASM(this js.Value, args []js.Value) (result interface{}) {
	defer func() {
		if r := recover(); r != nil {
			js.Global().Get("console").Call("error", fmt.Sprintf("lintQASM panic: %v", r))
			result = map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("Internal error: %v", r),
			}
		}
	}()

	if len(args) != 1 {
		return map[string]interface{}{
			"success": false,
			"error":   "Expected exactly one argument (QASM code)",
		}
	}

	qasmCode := args[0].String()
	js.Global().Get("console").Call("log", "Debug: Raw QASM code:", qasmCode)
	js.Global().Get("console").Call("log", "Debug: QASM code length:", len(qasmCode))
	if strings.TrimSpace(qasmCode) == "" {
		return map[string]interface{}{
			"success":    true,
			"violations": []interface{}{},
		}
	}

	// Format the code before linting
	f := formatter.NewFormatter()
	formatted, err := f.Format(qasmCode)
	if err != nil {
		js.Global().Get("console").Call("error", fmt.Sprintf("Format error: %v", err))
		return map[string]interface{}{
			"success": false,
			"error":   fmt.Sprintf("Failed to format QASM: %v", err),
		}
	}
	js.Global().Get("console").Call("log", "Formatted code:", formatted)
	qasmCode = formatted
	// Use global linter instance
	js.Global().Get("console").Call("log", "Linting code:", qasmCode)

	// Check if rules are loaded
	rules := globalLinter.GetRules()
	js.Global().Get("console").Call("log", "Loaded rules:", len(rules))
	for _, rule := range rules {
		js.Global().Get("console").Call("log", "Rule:", rule.ID, "Enabled:", rule.Enabled)
	}

	// Add debug logs
	js.Global().Get("console").Call("log", "Debug: Starting lint process")

	js.Global().Get("console").Call("log", "Debug: Before LintContent")
	violations, err := globalLinter.LintContent(qasmCode, "<stdin>")
	js.Global().Get("console").Call("log", "Debug: After LintContent, err:", err)
	if err != nil {
		js.Global().Get("console").Call("error", fmt.Sprintf("Lint error: %v", err))
		return map[string]interface{}{
			"success": false,
			"error":   fmt.Sprintf("Failed to lint QASM: %v", err),
		}
	}
	js.Global().Get("console").Call("log", "Debug: Violations length:", len(violations))
	if violations == nil {
		js.Global().Get("console").Call("log", "Debug: Violations is nil")
		violations = []*lint.Violation{}
	}
	for i, v := range violations {
		js.Global().Get("console").Call("log", "Debug: Processing violation", i)
		if v == nil {
			js.Global().Get("console").Call("log", "Debug: Violation is nil")
			continue
		}
		if v.Rule == nil {
			js.Global().Get("console").Call("log", "Debug: Rule is nil")
			continue
		}
		js.Global().Get("console").Call("log", "Debug: Violation details:", fmt.Sprintf("%+v", v))
		js.Global().Get("console").Call("log", "Debug: Rule details:", fmt.Sprintf("%+v", v.Rule))
		js.Global().Get("console").Call("log", "Violation:", v.Rule.ID, v.Message)
	}

	// Convert violations to JavaScript array
	jsViolations := js.Global().Get("Array").New(len(violations))
	js.Global().Get("console").Call("log", "Found violations:", len(violations))
	for i, v := range violations {
		js.Global().Get("console").Call("log", fmt.Sprintf("Violation %d: %s at line %d, column %d", i+1, v.Message, v.Line, v.Column))
		jsViolation := js.Global().Get("Object").New()
		jsViolation.Set("file", "<stdin>") // Always use <stdin> for consistency
		jsViolation.Set("line", v.Line)
		jsViolation.Set("column", v.Column)
		jsViolation.Set("severity", "error") // すべての違反をerrorとして扱う
		jsViolation.Set("rule_id", v.Rule.ID)
		jsViolation.Set("message", v.Message)
		jsViolation.Set("documentation_url", fmt.Sprintf("https://github.com/orangekame3/qasmtools/blob/main/docs/rules/%s.md", v.Rule.ID))

		// Add rule details
		jsRuleDetails := js.Global().Get("Object").New()
		jsRuleDetails.Set("name", v.Rule.Name)
		jsRuleDetails.Set("description", v.Rule.Description)
		// Convert tags array to JavaScript array
		jsTags := js.Global().Get("Array").New(len(v.Rule.Tags))
		for i, tag := range v.Rule.Tags {
			jsTags.SetIndex(i, tag)
		}
		jsRuleDetails.Set("tags", jsTags)
		jsRuleDetails.Set("fixable", v.Rule.Fixable)
		jsRuleDetails.Set("specification_url", v.Rule.SpecificationURL)

		jsExamples := js.Global().Get("Object").New()
		jsExamples.Set("incorrect", v.Rule.Examples.Incorrect)
		jsExamples.Set("correct", v.Rule.Examples.Correct)
		jsRuleDetails.Set("examples", jsExamples)

		jsViolation.Set("rule_details", jsRuleDetails)
		jsViolations.SetIndex(i, jsViolation)
	}

	// Create result object
	jsResult := js.Global().Get("Object").New()
	jsResult.Set("success", len(violations) == 0)
	jsResult.Set("violations", jsViolations)

	jsSummary := js.Global().Get("Object").New()
	jsSummary.Set("total", len(violations))
	jsSummary.Set("errors", len(violations)) // すべての違反をエラーとして扱う
	jsSummary.Set("warnings", 0)
	jsSummary.Set("info", 0)
	jsResult.Set("summary", jsSummary)

	result = jsResult
	return result
}
