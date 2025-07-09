package main

import (
	"fmt"
	"runtime"
)

// Version information for the QASM CLI tool
var (
	// Version is the current version of the qasm CLI
	Version = "0.2.15"

	// BuildDate is the date when the binary was built
	BuildDate = "unknown"

	// GitCommit is the git commit hash
	GitCommit = "unknown"
)

// GetVersion returns the formatted version string
func GetVersion() string {
	return fmt.Sprintf("%s (commit: %s, built: %s, %s/%s)",
		Version, GitCommit, BuildDate, runtime.GOOS, runtime.GOARCH)
}
