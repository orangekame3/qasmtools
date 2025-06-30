package main

import (
	"fmt"
	"strings"
	"syscall/js"

	"github.com/orangekame3/qasmtools/formatter"
)

func main() {
	c := make(chan struct{}, 0)

	// Export formatQASM function to JavaScript
	js.Global().Set("formatQASM", js.FuncOf(formatQASM))

	// Keep the program running
	<-c
}

func formatQASM(this js.Value, args []js.Value) interface{} {
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