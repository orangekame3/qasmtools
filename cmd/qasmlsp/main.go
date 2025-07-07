package main

import (
	"github.com/tliron/commonlog"
	_ "github.com/tliron/commonlog/simple"

	"github.com/orangekame3/qasmtools/lsp/server"
)

var version string = "0.0.1"

func main() {
	commonlog.Configure(1, nil) // Lower log level for more verbose output

	// Create and start the LSP server
	lspServer := server.NewServer(version)
	lspServer.Start()
}
