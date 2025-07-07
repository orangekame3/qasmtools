package lint

import (
	"regexp"
	"strings"
)

// IdentifierDeclaration represents a declared identifier with position
type IdentifierDeclaration struct {
	Name   string
	Column int
}

// Common regex patterns used across multiple checkers
var (
	// Qubit declaration patterns
	QubitDeclarationPattern      = regexp.MustCompile(`^\s*qubit(?:\[\s*\d+\s*\])?\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*;`)
	ArrayQubitDeclarationPattern = regexp.MustCompile(`^\s*qubit\[\s*(\d+)\s*\]\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*;`)

	// Bit declaration patterns
	BitDeclarationPattern      = regexp.MustCompile(`^\s*bit(?:\[\s*\d+\s*\])?\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*;`)
	ArrayBitDeclarationPattern = regexp.MustCompile(`^\s*bit\[\s*(\d+)\s*\]\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*;`)

	// Gate declaration patterns
	GateDeclarationPattern = regexp.MustCompile(`^\s*gate\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*(?:\([^)]*\))?\s+([^{]+)\{`)

	// Identifier patterns
	IdentifierPattern = regexp.MustCompile(`\b([a-zA-Z_][a-zA-Z0-9_]*)\b`)

	// Array access patterns
	ArrayAccessPattern = regexp.MustCompile(`\b([a-zA-Z_][a-zA-Z0-9_]*)\[\s*(\d+)\s*\]`)
)

// removeComments removes comments from a line
func RemoveComments(line string) string {
	if idx := strings.Index(line, "//"); idx != -1 {
		return line[:idx]
	}
	return line
}

// extractRegisterName extracts the base register name from register[index] format
func ExtractRegisterName(register string) string {
	// Pattern to match register[index] or just register
	pattern := regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)(?:\[[^\]]*\])?`)
	matches := pattern.FindStringSubmatch(register)
	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}

// isKeyword checks if a name is a reserved keyword or statement type
func IsKeyword(name string) bool {
	keywords := map[string]bool{
		"OPENQASM": true,
		"include":  true,
		"qubit":    true,
		"bit":      true,
		"int":      true,
		"uint":     true,
		"float":    true,
		"angle":    true,
		"complex":  true,
		"bool":     true,
		"const":    true,
		"def":      true,
		"gate":     true,
		"circuit":  true,
		"measure":  true,
		"reset":    true,
		"barrier":  true,
		"if":       true,
		"else":     true,
		"for":      true,
		"while":    true,
		"break":    true,
		"continue": true,
		"return":   true,
		"input":    true,
		"output":   true,
		// Types
		"duration": true,
		// Boolean values
		"true":  true,
		"false": true,
	}
	return keywords[name]
}

// skipCommentAndEmptyLine checks if a line should be skipped (empty or comment)
func SkipCommentAndEmptyLine(line string) bool {
	trimmedLine := strings.TrimSpace(line)
	return trimmedLine == "" || strings.HasPrefix(trimmedLine, "//")
}

// findIdentifierDeclarations finds all identifier declarations in a line
func FindIdentifierDeclarations(line string) []IdentifierDeclaration {
	var declarations []IdentifierDeclaration

	// Remove comments from the line
	codeOnly := RemoveComments(line)

	// Patterns for different types of declarations
	// Use proper identifier patterns to avoid matching numbers or invalid identifiers
	patterns := []*regexp.Regexp{
		// qubit declarations: qubit name; or qubit[size] name;
		regexp.MustCompile(`\bqubit(?:\[\s*\d+\s*\])?\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*;`),
		// bit declarations: bit name; or bit[size] name;
		regexp.MustCompile(`\bbit(?:\[\s*\d+\s*\])?\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*;`),
		// gate declarations: gate name(...) {...}
		regexp.MustCompile(`\bgate\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*(?:\(|[a-zA-Z_])`),
		// function declarations: def name(...) {...}
		regexp.MustCompile(`\bdef\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\(`),
		// circuit declarations: circuit name(...) {...}
		regexp.MustCompile(`\bcircuit\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\(`),
		// register declarations: int name; float name; etc.
		regexp.MustCompile(`\b(?:int|uint|float|angle|complex|bool)(?:\[\s*\d+\s*\])?\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*;`),
		// const declarations: const name = value;
		regexp.MustCompile(`\bconst\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*=`),
		// input/output declarations: input name; output name;
		regexp.MustCompile(`\b(?:input|output)\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*;`),
	}

	for _, pattern := range patterns {
		matches := pattern.FindAllStringSubmatch(codeOnly, -1)
		indices := pattern.FindAllStringIndex(codeOnly, -1)

		for i, match := range matches {
			if len(match) >= 2 {
				identifierName := match[1]
				// Find the position of the identifier within the match
				matchStart := indices[i][0]
				identifierPos := strings.Index(codeOnly[matchStart:], identifierName)
				if identifierPos != -1 {
					column := matchStart + identifierPos + 1 // Convert to 1-based indexing

					declarations = append(declarations, IdentifierDeclaration{
						Name:   identifierName,
						Column: column,
					})
				}
			}
		}
	}

	return declarations
}

// getBuiltinIdentifiers returns a map of built-in identifiers that are always available
func GetBuiltinIdentifiers() map[string]bool {
	return map[string]bool{
		// Standard gates (these might be included via stdgates.qasm)
		"h": true, "x": true, "y": true, "z": true, "s": true, "t": true,
		"cx": true, "cy": true, "cz": true, "ccx": true, "reset": true,
		"sdg": true, "tdg": true, "rx": true, "ry": true, "rz": true,
		"p": true, "swap": true, "cswap": true, "u1": true, "u2": true, "u3": true,

		// Standard functions
		"sin": true, "cos": true, "tan": true, "exp": true, "ln": true, "sqrt": true,

		// Constants
		"pi": true, "euler": true, "tau": true,
	}
}
