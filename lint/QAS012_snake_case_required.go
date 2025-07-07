package lint

import (
	"regexp"
	"strings"

	"github.com/orangekame3/qasmtools/parser"
)

// SnakeCaseRequiredChecker is the new implementation using BaseChecker framework
type SnakeCaseRequiredChecker struct {
	*BaseChecker
	snakeCasePattern *regexp.Regexp
}

// NewSnakeCaseRequiredChecker creates a new SnakeCaseRequiredChecker
func NewSnakeCaseRequiredChecker() *SnakeCaseRequiredChecker {
	return &SnakeCaseRequiredChecker{
		BaseChecker:      NewBaseChecker("QAS012"),
		snakeCasePattern: regexp.MustCompile(`^[a-z][a-z0-9_]*[a-z0-9]$|^[a-z]$`),
	}
}

// CheckFile performs file-level snake_case naming analysis
func (c *SnakeCaseRequiredChecker) CheckFile(context *CheckContext) []*Violation {
	// Use the shared ProcessFileLines function with this checker as LineProcessor
	return ProcessFileLines(context, c)
}

// ProcessLine implements LineProcessor interface for line-by-line processing
func (c *SnakeCaseRequiredChecker) ProcessLine(line string, lineNum int, context *CheckContext) []*Violation {
	var violations []*Violation

	// Use shared utility to find identifier declarations, but filter for specific types
	identifiers := c.findSnakeCaseTargetDeclarations(line)

	for _, identifier := range identifiers {
		if !c.isSnakeCase(identifier.Name) {
			violation := c.NewViolationBuilder().
				WithMessage("Identifier '"+identifier.Name+"' should be written in snake_case.").
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

// findSnakeCaseTargetDeclarations finds declarations that should follow snake_case
// This is more specific than the general FindIdentifierDeclarations
func (c *SnakeCaseRequiredChecker) findSnakeCaseTargetDeclarations(line string) []IdentifierDeclaration {
	var declarations []IdentifierDeclaration

	// Remove comments using shared utility
	codeOnly := RemoveComments(line)

	// Check more specific patterns first to avoid duplicates
	// Check array declarations first
	if matches := ArrayQubitDeclarationPattern.FindAllStringSubmatch(codeOnly, -1); len(matches) > 0 {
		indices := ArrayQubitDeclarationPattern.FindAllStringIndex(codeOnly, -1)
		for i, match := range matches {
			if len(match) > 2 {
				identifierName := match[2] // Variable name is in group 2 for array declarations
				matchStart := indices[i][0]
				identifierPos := strings.Index(codeOnly[matchStart:], identifierName)
				if identifierPos != -1 {
					column := matchStart + identifierPos + 1
					declarations = append(declarations, IdentifierDeclaration{
						Name:   identifierName,
						Column: column,
					})
				}
			}
		}
	} else if matches := QubitDeclarationPattern.FindAllStringSubmatch(codeOnly, -1); len(matches) > 0 {
		// Check single qubit declarations only if no array declarations found
		indices := QubitDeclarationPattern.FindAllStringIndex(codeOnly, -1)
		for i, match := range matches {
			if len(match) > 1 {
				identifierName := match[1]
				matchStart := indices[i][0]
				identifierPos := strings.Index(codeOnly[matchStart:], identifierName)
				if identifierPos != -1 {
					column := matchStart + identifierPos + 1
					declarations = append(declarations, IdentifierDeclaration{
						Name:   identifierName,
						Column: column,
					})
				}
			}
		}
	}

	// Check array bit declarations
	if matches := ArrayBitDeclarationPattern.FindAllStringSubmatch(codeOnly, -1); len(matches) > 0 {
		indices := ArrayBitDeclarationPattern.FindAllStringIndex(codeOnly, -1)
		for i, match := range matches {
			if len(match) > 2 {
				identifierName := match[2] // Variable name is in group 2 for array declarations
				matchStart := indices[i][0]
				identifierPos := strings.Index(codeOnly[matchStart:], identifierName)
				if identifierPos != -1 {
					column := matchStart + identifierPos + 1
					declarations = append(declarations, IdentifierDeclaration{
						Name:   identifierName,
						Column: column,
					})
				}
			}
		}
	} else if matches := BitDeclarationPattern.FindAllStringSubmatch(codeOnly, -1); len(matches) > 0 {
		// Check single bit declarations only if no array declarations found
		indices := BitDeclarationPattern.FindAllStringIndex(codeOnly, -1)
		for i, match := range matches {
			if len(match) > 1 {
				identifierName := match[1]
				matchStart := indices[i][0]
				identifierPos := strings.Index(codeOnly[matchStart:], identifierName)
				if identifierPos != -1 {
					column := matchStart + identifierPos + 1
					declarations = append(declarations, IdentifierDeclaration{
						Name:   identifierName,
						Column: column,
					})
				}
			}
		}
	}

	// Check other declaration types
	otherPatterns := []struct {
		pattern   *regexp.Regexp
		nameIndex int
	}{
		{GateDeclarationPattern, 1}, // gate declarations: name in group 1
		// circuit declarations: circuit name(...) {...}
		{regexp.MustCompile(`\bcircuit\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\(`), 1},
	}

	for _, patternInfo := range otherPatterns {
		matches := patternInfo.pattern.FindAllStringSubmatch(codeOnly, -1)
		indices := patternInfo.pattern.FindAllStringIndex(codeOnly, -1)

		for i, match := range matches {
			if len(match) > patternInfo.nameIndex {
				identifierName := match[patternInfo.nameIndex]
				matchStart := indices[i][0]
				identifierPos := strings.Index(codeOnly[matchStart:], identifierName)
				if identifierPos != -1 {
					column := matchStart + identifierPos + 1
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

// isSnakeCase checks if an identifier follows snake_case naming convention
func (c *SnakeCaseRequiredChecker) isSnakeCase(name string) bool {
	// Check basic pattern using precompiled regex
	if !c.snakeCasePattern.MatchString(name) {
		return false
	}

	// Additional checks for proper snake_case:
	// 1. No consecutive underscores
	if strings.Contains(name, "__") {
		return false
	}

	// 2. No ending with underscore (unless single character)
	if len(name) > 1 && strings.HasSuffix(name, "_") {
		return false
	}

	// 3. No starting with underscore (handled by pattern but double-check)
	if strings.HasPrefix(name, "_") {
		return false
	}

	return true
}

// Check implements RuleChecker interface (required but delegates to CheckProgram)
func (c *SnakeCaseRequiredChecker) Check(node parser.Node, context *CheckContext) []*Violation {
	return nil
}

// CheckProgram implements ProgramChecker interface
func (c *SnakeCaseRequiredChecker) CheckProgram(context *CheckContext) []*Violation {
	return c.CheckFile(context)
}
