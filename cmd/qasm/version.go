package main

// Version information for the QASM CLI tool
var (
	// Version is the current version of the qasm CLI
	Version = "0.2.13"

	// BuildDate is the date when the binary was built
	BuildDate = "unknown"

	// GitCommit is the git commit hash
	GitCommit = "unknown"
)

// GetVersion returns the formatted version string
func GetVersion() string {
	if BuildDate != "unknown" && GitCommit != "unknown" {
		return Version + " (built " + BuildDate + ", commit " + GitCommit + ")"
	}
	return Version
}
