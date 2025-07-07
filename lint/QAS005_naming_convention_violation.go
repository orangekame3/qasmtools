package lint

import (
	"regexp"
	"strings"

	"github.com/orangekame3/qasmtools/parser"
)

// NamingConventionViolationChecker is the new implementation using BaseChecker framework
type NamingConventionViolationChecker struct {
	*BaseChecker
	namingPattern *regexp.Regexp
}

// NewNamingConventionViolationChecker creates a new NamingConventionViolationChecker
func NewNamingConventionViolationChecker() *NamingConventionViolationChecker {
	return &NamingConventionViolationChecker{
		BaseChecker:   NewBaseChecker("QAS005"),
		namingPattern: regexp.MustCompile(`^[a-z][a-zA-Z0-9_]*$`),
	}
}

// CheckFile performs file-level naming convention analysis
func (c *NamingConventionViolationChecker) CheckFile(context *CheckContext) []*Violation {
	// Use the shared ProcessFileLines function with this checker as LineProcessor
	return ProcessFileLines(context, c)
}

// ProcessLine implements LineProcessor interface for line-by-line processing
func (c *NamingConventionViolationChecker) ProcessLine(line string, lineNum int, context *CheckContext) []*Violation {
	var violations []*Violation

	// Find all identifier-like declarations (including invalid ones)
	identifiers := c.findAllIdentifierDeclarations(line)

	for _, identifier := range identifiers {
		if !c.isValidIdentifierName(identifier.Name) {
			violation := c.NewViolationBuilder().
				WithMessage("Identifier '"+identifier.Name+"' violates naming conventions. Follow pattern: ^[a-z][a-zA-Z0-9_]*$.").
				WithFile(context.File).
				WithPosition(lineNum, identifier.Column).
				WithNodeName(identifier.Name).
				AsWarning().
				Build()
			violations = append(violations, violation)
		}
	}

	return violations
}

// findAllIdentifierDeclarations finds both valid and invalid identifier declarations
func (c *NamingConventionViolationChecker) findAllIdentifierDeclarations(line string) []IdentifierDeclaration {
	var declarations []IdentifierDeclaration

	// Remove comments from the line
	codeOnly := RemoveComments(line)

	// Patterns that capture both valid and invalid identifiers in declarations
	// Use broader patterns that can match invalid identifiers too
	patterns := []*regexp.Regexp{
		// qubit declarations: qubit name; or qubit[size] name; (including invalid names)
		regexp.MustCompile(`\bqubit(?:\[\s*\d+\s*\])?\s+([^\s;]+)\s*;`),
		// bit declarations: bit name; or bit[size] name; (including invalid names)
		regexp.MustCompile(`\bbit(?:\[\s*\d+\s*\])?\s+([^\s;]+)\s*;`),
		// gate declarations: gate name(...) {...} (including invalid names)
		regexp.MustCompile(`\bgate\s+([^\s\(]+)\s*(?:\(|[a-zA-Z0-9_])`),
		// function declarations: def name(...) {...} (including invalid names)
		regexp.MustCompile(`\bdef\s+([^\s\(]+)\s*\(`),
		// circuit declarations: circuit name(...) {...} (including invalid names)
		regexp.MustCompile(`\bcircuit\s+([^\s\(]+)\s*\(`),
		// register declarations: int name; float name; etc. (including invalid names)
		regexp.MustCompile(`\b(?:int|uint|float|angle|complex|bool)(?:\[\s*\d+\s*\])?\s+([^\s;]+)\s*;`),
		// const declarations: const name = value; (including invalid names)
		regexp.MustCompile(`\bconst\s+([^\s=]+)\s*=`),
		// input/output declarations: input name; output name; (including invalid names)
		regexp.MustCompile(`\b(?:input|output)\s+([^\s;]+)\s*;`),
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

// isValidIdentifierName checks if an identifier follows the naming convention
func (c *NamingConventionViolationChecker) isValidIdentifierName(name string) bool {
	// Use the precompiled regex pattern
	return c.namingPattern.MatchString(name)
}

// Check implements RuleChecker interface (required but delegates to CheckProgram)
func (c *NamingConventionViolationChecker) Check(node parser.Node, context *CheckContext) []*Violation {
	return nil
}

// CheckProgram implements ProgramChecker interface
func (c *NamingConventionViolationChecker) CheckProgram(context *CheckContext) []*Violation {
	return c.CheckFile(context)
}
