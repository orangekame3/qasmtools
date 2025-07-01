package formatter

// Config holds the formatting configuration
type Config struct {
	Write    bool
	Check    bool
	Indent   uint
	Newline  bool
	Verbose  bool
	Diff     bool
	Unescape bool
}
