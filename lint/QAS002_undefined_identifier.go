package lint

import (
	"os"
	"regexp"
	"strings"

	"github.com/orangekame3/qasmtools/parser"
)

// UndefinedIdentifierChecker checks for undefined identifier usage (QAS002)
type UndefinedIdentifierChecker struct{}

func (c *UndefinedIdentifierChecker) Check(node parser.Node, context *CheckContext) []*Violation {
	// This method is required by RuleChecker but not used for program-level analysis
	return nil
}

// CheckProgram implements ProgramChecker interface for program-level analysis
func (c *UndefinedIdentifierChecker) CheckProgram(context *CheckContext) []*Violation {
	// Use text-based analysis due to AST parsing issues
	return c.CheckFile(context)
}

// CheckFile performs file-level undefined identifier analysis
func (c *UndefinedIdentifierChecker) CheckFile(context *CheckContext) []*Violation {
	var violations []*Violation

	// Read file content for text-based analysis
	content, err := os.ReadFile(context.File)
	if err != nil {
		return violations
	}

	text := string(content)
	lines := strings.Split(text, "\n")

	// First pass: collect all declared identifiers
	declared := c.findDeclaredIdentifiers(lines)

	// Second pass: find usage of undefined identifiers
	for i, line := range lines {
		// Skip comments and empty lines
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "//") {
			continue
		}

		// Find all identifier usages in this line
		usages := c.findIdentifierUsages(line)
		
		for _, usage := range usages {
			// Check if the identifier is declared
			if !declared[usage.name] {
				violation := &Violation{
					Rule:     nil, // Will be set by the runner
					File:     context.File,
					Line:     i + 1,
					Column:   usage.column,
					NodeName: usage.name,
					Message:  "Identifier '" + usage.name + "' is not declared.",
					Severity: SeverityError,
				}
				violations = append(violations, violation)
			}
		}
	}

	return violations
}

type identifierUsage struct {
	name   string
	column int
}

// findDeclaredIdentifiers finds all declared identifiers in the file
func (c *UndefinedIdentifierChecker) findDeclaredIdentifiers(lines []string) map[string]bool {
	declared := make(map[string]bool)

	// Built-in identifiers that are always available
	builtins := []string{
		// Standard gates (these might be included via stdgates.qasm)
		"h", "x", "y", "z", "s", "t", "cx", "cy", "cz", "ccx", "reset",
		// Standard functions
		"sin", "cos", "tan", "exp", "ln", "sqrt",
		// Constants
		"pi", "euler", "tau",
	}
	
	for _, builtin := range builtins {
		declared[builtin] = true
	}

	// Declaration patterns
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`^\s*qubit(?:\[\s*\d+\s*\])?\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*;`),        // qubit declarations
		regexp.MustCompile(`^\s*bit(?:\[\s*\d+\s*\])?\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*;`),          // bit declarations
		regexp.MustCompile(`^\s*gate\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*[\(\s]`),                     // gate definitions (followed by space or parentheses)
		regexp.MustCompile(`^\s*def\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\(`),                           // function definitions
		regexp.MustCompile(`^\s*include\s+"([^"]+)"`),                                            // include statements
	}

	for _, line := range lines {
		// Skip comments
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "//") {
			continue
		}

		for _, pattern := range patterns {
			matches := pattern.FindStringSubmatch(line)
			if len(matches) > 1 {
				identifier := matches[1]
				declared[identifier] = true
				
				// Special handling for include statements
				if strings.Contains(line, "include") && strings.Contains(identifier, "stdgates") {
					// Add standard gates when stdgates.qasm is included
					stdGates := []string{
						"h", "x", "y", "z", "s", "sdg", "t", "tdg", 
						"rx", "ry", "rz", "p", "cx", "cy", "cz", "swap",
						"ccx", "cswap", "u1", "u2", "u3",
					}
					for _, gate := range stdGates {
						declared[gate] = true
					}
				}
			}
		}
	}

	return declared
}

// findIdentifierUsages finds all identifier usages in a line (excluding declarations)
func (c *UndefinedIdentifierChecker) findIdentifierUsages(line string) []identifierUsage {
	var usages []identifierUsage

	// Skip declaration lines
	if c.isDeclarationLine(line) {
		return usages
	}

	// Remove comments from the line before processing
	codeOnly := c.removeComments(line)
	if strings.TrimSpace(codeOnly) == "" {
		return usages
	}

	// Find all potential identifiers in the line
	identifierPattern := regexp.MustCompile(`\b([a-zA-Z_][a-zA-Z0-9_]*)\b`)
	matches := identifierPattern.FindAllStringSubmatch(codeOnly, -1)
	matchIndices := identifierPattern.FindAllStringIndex(codeOnly, -1)

	for i, match := range matches {
		if len(match) > 1 {
			identifier := match[1]
			column := matchIndices[i][0] + 1 // Convert to 1-based indexing

			// Skip keywords and built-in types
			if c.isKeyword(identifier) {
				continue
			}

			usages = append(usages, identifierUsage{
				name:   identifier,
				column: column,
			})
		}
	}

	return usages
}

// isDeclarationLine checks if a line contains a declaration
func (c *UndefinedIdentifierChecker) isDeclarationLine(line string) bool {
	declarationPatterns := []string{
		`^\s*qubit`,
		`^\s*bit`,
		`^\s*gate\s+`,
		`^\s*def\s+`,
		`^\s*include\s+`,
		`^\s*OPENQASM`,
	}

	for _, pattern := range declarationPatterns {
		matched, _ := regexp.MatchString(pattern, line)
		if matched {
			return true
		}
	}

	return false
}

// isKeyword checks if an identifier is a reserved keyword
func (c *UndefinedIdentifierChecker) isKeyword(identifier string) bool {
	keywords := map[string]bool{
		// OpenQASM keywords
		"OPENQASM": true,
		"include":  true,
		"qubit":    true,
		"bit":      true,
		"gate":     true,
		"def":      true,
		"if":       true,
		"else":     true,
		"for":      true,
		"while":    true,
		"break":    true,
		"continue": true,
		"measure":  true,
		"reset":    true,
		"barrier":  true,
		"delay":    true,
		"stretch":  true,
		"box":      true,
		"let":      true,
		"const":    true,
		"input":    true,
		"output":   true,
		
		// Types
		"int":      true,
		"uint":     true,
		"float":    true,
		"angle":    true,
		"bool":     true,
		"duration": true,
		
		// Boolean values
		"true":  true,
		"false": true,
	}

	return keywords[identifier]
}

// removeComments removes comments from a line, handling string literals properly
func (c *UndefinedIdentifierChecker) removeComments(line string) string {
	// Simple approach: find // and remove everything after it
	// This doesn't handle // inside string literals, but that's rare in QASM
	if idx := strings.Index(line, "//"); idx != -1 {
		return line[:idx]
	}
	return line
}