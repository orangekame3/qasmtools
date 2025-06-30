package lint

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DocumentationGenerator generates markdown documentation for lint rules
type DocumentationGenerator struct {
	outputDir string
}

// NewDocumentationGenerator creates a new documentation generator
func NewDocumentationGenerator(outputDir string) *DocumentationGenerator {
	return &DocumentationGenerator{
		outputDir: outputDir,
	}
}

// GenerateRuleDocumentation generates markdown documentation for a single rule
func (dg *DocumentationGenerator) GenerateRuleDocumentation(rule *Rule) (string, error) {
	var sb strings.Builder

	// Title
	sb.WriteString(fmt.Sprintf("# %s (%s)\n\n", rule.Name, rule.ID))

	// Severity and status
	sb.WriteString(fmt.Sprintf("**Severity:** %s  \n", rule.Level))
	sb.WriteString(fmt.Sprintf("**Category:** %s  \n", strings.Join(rule.Tags, ", ")))
	sb.WriteString(fmt.Sprintf("**Fixable:** %t  \n\n", rule.Fixable))

	// Description
	sb.WriteString("## Description\n\n")
	sb.WriteString(fmt.Sprintf("%s\n\n", rule.Description))

	// Rule details
	sb.WriteString("## Rule Details\n\n")
	sb.WriteString(fmt.Sprintf("This rule checks for %s violations.\n\n", strings.ToLower(rule.Name)))

	// Message format
	sb.WriteString("## Message Format\n\n")
	sb.WriteString(fmt.Sprintf("```\n%s\n```\n\n", rule.Message))

	// Examples section
	sb.WriteString("## Examples\n\n")
	sb.WriteString("### ❌ Incorrect\n\n")
	sb.WriteString(dg.generateIncorrectExample(rule))
	sb.WriteString("\n### ✅ Correct\n\n")
	sb.WriteString(dg.generateCorrectExample(rule))
	sb.WriteString("\n\n")

	// Configuration
	sb.WriteString("## Configuration\n\n")
	sb.WriteString(fmt.Sprintf("- **Enabled by default:** %t\n", rule.Enabled))
	sb.WriteString(fmt.Sprintf("- **Match type:** %s\n", rule.Match.Type))
	sb.WriteString(fmt.Sprintf("- **Match kind:** %s\n\n", rule.Match.Kind))

	// Related rules
	sb.WriteString("## Related Rules\n\n")
	sb.WriteString(dg.generateRelatedRules(rule))

	return sb.String(), nil
}

// GenerateAllDocumentation generates documentation for all rules
func (dg *DocumentationGenerator) GenerateAllDocumentation(rules []*Rule) error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(dg.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate documentation for each rule
	for _, rule := range rules {
		content, err := dg.GenerateRuleDocumentation(rule)
		if err != nil {
			return fmt.Errorf("failed to generate documentation for rule %s: %w", rule.ID, err)
		}

		filename := fmt.Sprintf("%s.md", rule.ID)
		filepath := filepath.Join(dg.outputDir, filename)

		if err := os.WriteFile(filepath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write documentation file %s: %w", filepath, err)
		}
	}

	// Generate index file
	if err := dg.generateIndexFile(rules); err != nil {
		return fmt.Errorf("failed to generate index file: %w", err)
	}

	return nil
}

// generateIndexFile creates an index file listing all rules
func (dg *DocumentationGenerator) generateIndexFile(rules []*Rule) error {
	var sb strings.Builder

	sb.WriteString("# QASM Lint Rules Documentation\n\n")
	sb.WriteString("This directory contains documentation for all available QASM linting rules.\n\n")

	// Rules by severity
	errorRules := []*Rule{}
	warningRules := []*Rule{}
	infoRules := []*Rule{}

	for _, rule := range rules {
		switch rule.Level {
		case SeverityError:
			errorRules = append(errorRules, rule)
		case SeverityWarning:
			warningRules = append(warningRules, rule)
		case SeverityInfo:
			infoRules = append(infoRules, rule)
		}
	}

	// Error rules
	if len(errorRules) > 0 {
		sb.WriteString("## Error Rules\n\n")
		for _, rule := range errorRules {
			sb.WriteString(fmt.Sprintf("- **[%s](%s.md)** - %s\n", rule.ID, rule.ID, rule.Description))
		}
		sb.WriteString("\n")
	}

	// Warning rules
	if len(warningRules) > 0 {
		sb.WriteString("## Warning Rules\n\n")
		for _, rule := range warningRules {
			sb.WriteString(fmt.Sprintf("- **[%s](%s.md)** - %s\n", rule.ID, rule.ID, rule.Description))
		}
		sb.WriteString("\n")
	}

	// Info rules
	if len(infoRules) > 0 {
		sb.WriteString("## Info Rules\n\n")
		for _, rule := range infoRules {
			sb.WriteString(fmt.Sprintf("- **[%s](%s.md)** - %s\n", rule.ID, rule.ID, rule.Description))
		}
		sb.WriteString("\n")
	}

	// All rules table
	sb.WriteString("## All Rules\n\n")
	sb.WriteString("| Rule ID | Name | Severity | Tags | Fixable |\n")
	sb.WriteString("|---------|------|----------|------|---------|\n")

	for _, rule := range rules {
		sb.WriteString(fmt.Sprintf("| [%s](%s.md) | %s | %s | %s | %t |\n",
			rule.ID, rule.ID, rule.Name, rule.Level, strings.Join(rule.Tags, ", "), rule.Fixable))
	}

	indexPath := filepath.Join(dg.outputDir, "README.md")
	return os.WriteFile(indexPath, []byte(sb.String()), 0644)
}

// generateIncorrectExample generates example code that violates the rule
func (dg *DocumentationGenerator) generateIncorrectExample(rule *Rule) string {
	switch rule.ID {
	case "QAS001":
		return "```qasm\nOPENQASM 3.0;\ninclude \"stdgates.qasm\";\n\nqubit q;\nqubit unused_qubit;  // ❌ Never used\n\nh q;\n```"
	case "QAS002":
		return "```qasm\nOPENQASM 3.0;\ninclude \"stdgates.qasm\";\n\nqubit[2] q;\nbit c;  // ❌ Only 1 bit for 2 measurements\n\nh q;\nmeasure q -> c;\n```"
	case "QAS003":
		return "```qasm\nOPENQASM 3.0;\ninclude \"stdgates.qasm\";\n\nqubit q;\nbit c;\n\nmeasure q -> c;  // ❌ Measuring unaffected qubit\n```"
	case "QAS004":
		return "```qasm\nOPENQASM 3.0;\ninclude \"stdgates.qasm\";\n\nqubit[2] q;\n\nh q[3];  // ❌ Index 3 exceeds array bounds [0,1]\n```"
	case "QAS005":
		return "```qasm\nOPENQASM 3.0;\ninclude \"stdgates.qasm\";\n\nqubit MyQubit;  // ❌ Should start with lowercase\nqubit q_2;      // ❌ Should not start with uppercase\n```"
	default:
		return "```qasm\n// Example code that violates this rule\n```"
	}
}

// generateCorrectExample generates example code that follows the rule
func (dg *DocumentationGenerator) generateCorrectExample(rule *Rule) string {
	switch rule.ID {
	case "QAS001":
		return "```qasm\nOPENQASM 3.0;\ninclude \"stdgates.qasm\";\n\nqubit q;\n\nh q;  // ✅ Qubit is used\n```"
	case "QAS002":
		return "```qasm\nOPENQASM 3.0;\ninclude \"stdgates.qasm\";\n\nqubit[2] q;\nbit[2] c;  // ✅ Sufficient bits for measurements\n\nh q;\nmeasure q -> c;\n```"
	case "QAS003":
		return "```qasm\nOPENQASM 3.0;\ninclude \"stdgates.qasm\";\n\nqubit q;\nbit c;\n\nh q;          // ✅ Affect qubit before measurement\nmeasure q -> c;\n```"
	case "QAS004":
		return "```qasm\nOPENQASM 3.0;\ninclude \"stdgates.qasm\";\n\nqubit[4] q;\n\nh q[3];  // ✅ Index 3 is within bounds [0,3]\n```"
	case "QAS005":
		return "```qasm\nOPENQASM 3.0;\ninclude \"stdgates.qasm\";\n\nqubit myQubit;  // ✅ Starts with lowercase\nqubit q2;       // ✅ Valid naming convention\n```"
	default:
		return "```qasm\n// Example code that follows this rule\n```"
	}
}

// generateRelatedRules generates a list of related rules
func (dg *DocumentationGenerator) generateRelatedRules(rule *Rule) string {
	// This could be enhanced to automatically detect related rules based on tags
	switch rule.ID {
	case "QAS001":
		return "- QAS005 (naming-convention-violation): Both relate to variable declaration best practices\n"
	case "QAS002":
		return "- QAS003 (constant-measured-bit): Both relate to measurement operations\n- QAS004 (exceeding-qubit-limits): Both relate to array bounds checking\n"
	case "QAS003":
		return "- QAS002 (insufficient-classical-bits): Both relate to measurement operations\n"
	case "QAS004":
		return "- QAS002 (insufficient-classical-bits): Both relate to array bounds checking\n"
	case "QAS005":
		return "- QAS001 (unused-qubit): Both relate to variable declaration best practices\n"
	default:
		return "None currently identified.\n"
	}
}