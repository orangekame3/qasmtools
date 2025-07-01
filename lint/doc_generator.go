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
	sb.WriteString(fmt.Sprintf("**Fixable:** %t  \n", rule.Fixable))
	
	// Add specification URL if available
	if rule.SpecificationURL != "" {
		sb.WriteString(fmt.Sprintf("**OpenQASM Specification:** [View Details](%s)  \n", rule.SpecificationURL))
	}
	sb.WriteString("\n")

	// Description
	sb.WriteString("## Description\n\n")
	sb.WriteString(fmt.Sprintf("%s\n\n", rule.Description))

	// Rule details
	sb.WriteString("## Rule Details\n\n")
	sb.WriteString(fmt.Sprintf("This rule checks for %s violations according to OpenQASM 3.0 specifications.\n\n", strings.ReplaceAll(rule.Name, "-", " ")))

	// Message format
	sb.WriteString("## Message Format\n\n")
	sb.WriteString(fmt.Sprintf("```\n%s\n```\n\n", rule.Message))

	// Examples section - use examples from YAML if available
	sb.WriteString("## Examples\n\n")
	if rule.Examples.Incorrect != "" || rule.Examples.Correct != "" {
		// Use examples from YAML
		if rule.Examples.Incorrect != "" {
			sb.WriteString("### ❌ Incorrect\n\n")
			sb.WriteString(fmt.Sprintf("```qasm\n%s\n```\n\n", strings.TrimSpace(rule.Examples.Incorrect)))
		}
		if rule.Examples.Correct != "" {
			sb.WriteString("### ✅ Correct\n\n")
			sb.WriteString(fmt.Sprintf("```qasm\n%s\n```\n\n", strings.TrimSpace(rule.Examples.Correct)))
		}
	} else {
		// Fallback to generated examples
		sb.WriteString("### ❌ Incorrect\n\n")
		sb.WriteString(dg.generateIncorrectExample(rule))
		sb.WriteString("\n### ✅ Correct\n\n")
		sb.WriteString(dg.generateCorrectExample(rule))
		sb.WriteString("\n\n")
	}

	// Configuration
	sb.WriteString("## Configuration\n\n")
	sb.WriteString(fmt.Sprintf("- **Enabled by default:** %t\n", rule.Enabled))
	sb.WriteString(fmt.Sprintf("- **Match type:** %s\n", rule.Match.Type))
	sb.WriteString(fmt.Sprintf("- **Match kind:** %s\n\n", rule.Match.Kind))

	// Related rules
	sb.WriteString("## Related Rules\n\n")
	sb.WriteString(dg.generateRelatedRules(rule))

	// References
	sb.WriteString("## References\n\n")
	if rule.SpecificationURL != "" {
		sb.WriteString(fmt.Sprintf("- [OpenQASM 3.0 Specification](%s)\n", rule.SpecificationURL))
	}
	if rule.DocumentationURL != "" {
		sb.WriteString(fmt.Sprintf("- [Rule Documentation](%s)\n", rule.DocumentationURL))
	}

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
	sb.WriteString("Each rule includes detailed descriptions, examples, and links to the relevant OpenQASM 3.0 specification sections.\n\n")

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
		sb.WriteString("These rules catch critical issues that prevent valid OpenQASM execution:\n\n")
		for _, rule := range errorRules {
			sb.WriteString(fmt.Sprintf("- **[%s](%s.md)** - %s\n", rule.ID, rule.ID, rule.Description))
		}
		sb.WriteString("\n")
	}

	// Warning rules
	if len(warningRules) > 0 {
		sb.WriteString("## Warning Rules\n\n")
		sb.WriteString("These rules identify potential issues and style violations:\n\n")
		for _, rule := range warningRules {
			sb.WriteString(fmt.Sprintf("- **[%s](%s.md)** - %s\n", rule.ID, rule.ID, rule.Description))
		}
		sb.WriteString("\n")
	}

	// Info rules
	if len(infoRules) > 0 {
		sb.WriteString("## Info Rules\n\n")
		sb.WriteString("These rules provide helpful suggestions and best practices:\n\n")
		for _, rule := range infoRules {
			sb.WriteString(fmt.Sprintf("- **[%s](%s.md)** - %s\n", rule.ID, rule.ID, rule.Description))
		}
		sb.WriteString("\n")
	}

	// All rules table
	sb.WriteString("## All Rules Summary\n\n")
	sb.WriteString("| Rule ID | Name | Severity | Tags | Fixable | Specification |\n")
	sb.WriteString("|---------|------|----------|------|---------|---------------|\n")

	for _, rule := range rules {
		specLink := "N/A"
		if rule.SpecificationURL != "" {
			specLink = fmt.Sprintf("[Link](%s)", rule.SpecificationURL)
		}
		sb.WriteString(fmt.Sprintf("| [%s](%s.md) | %s | %s | %s | %t | %s |\n",
			rule.ID, rule.ID, rule.Name, rule.Level, strings.Join(rule.Tags, ", "), rule.Fixable, specLink))
	}

	sb.WriteString("\n## Usage\n\n")
	sb.WriteString("To disable specific rules, use the `--disable` flag:\n\n")
	sb.WriteString("```bash\n")
	sb.WriteString("qasm lint --disable=QAS005,QAS012 input.qasm\n")
	sb.WriteString("```\n\n")
	sb.WriteString("To enable only specific rules, use the `--enable-only` flag:\n\n")
	sb.WriteString("```bash\n")
	sb.WriteString("qasm lint --enable-only=QAS001,QAS002 input.qasm\n")
	sb.WriteString("```\n\n")

	indexPath := filepath.Join(dg.outputDir, "README.md")
	return os.WriteFile(indexPath, []byte(sb.String()), 0644)
}

// generateIncorrectExample generates example code that violates the rule
func (dg *DocumentationGenerator) generateIncorrectExample(rule *Rule) string {
	switch rule.ID {
	case "QAS001":
		return "```qasm\nOPENQASM 3.0;\ninclude \"stdgates.qasm\";\n\nqubit q;\nqubit unused_qubit;  // ❌ Never used\n\nh q;\n```"
	case "QAS002":
		return "```qasm\nOPENQASM 3.0;\n\nh q[0];  // ❌ q is not declared\n```"
	case "QAS003":
		return "```qasm\nOPENQASM 3.0;\ninclude \"stdgates.qasm\";\n\nqubit q;\nbit c;\n\nmeasure q -> c;  // ❌ Measuring unaffected qubit\n```"
	case "QAS004":
		return "```qasm\nOPENQASM 3.0;\ninclude \"stdgates.qasm\";\n\nqubit[2] q;\n\nh q[2];  // ❌ Index 2 exceeds array bounds [0,1]\n```"
	case "QAS005":
		return "```qasm\nOPENQASM 3.0;\ninclude \"stdgates.qasm\";\n\nqubit MyQubit;  // ❌ Should start with lowercase\n```"
	case "QAS006":
		return "```qasm\nOPENQASM 3.0;\n\ngate mygate(a, b) { /* ... */ }\nqubit[2] q1;\nqubit[3] q2;\nmygate q1, q2;  // ❌ Size mismatch: 2 vs 3\n```"
	case "QAS007":
		return "```qasm\nOPENQASM 3.0;\n\ngate mygate(a) {\n  h a[0];  // ❌ Index access on parameter\n}\n```"
	case "QAS008":
		return "```qasm\nOPENQASM 3.0;\n\ndef foo() {\n  qubit q;  // ❌ Local qubit declaration\n}\n```"
	case "QAS009":
		return "```qasm\nOPENQASM 3.0;\n\nbreak;  // ❌ Outside of loop\n```"
	case "QAS010":
		return "```qasm\nOPENQASM 3.0;\n\ngate g(a) {\n  measure a -> c;  // ❌ Non-unitary operation\n}\n```"
	case "QAS011":
		return "```qasm\nOPENQASM 3.0;\n\nqubit __reserved_name;  // ❌ Uses reserved prefix\n```"
	case "QAS012":
		return "```qasm\nOPENQASM 3.0;\n\nqubit myQubit;  // ❌ Should use snake_case\ngate MyGate(q) { h q; }  // ❌ Should use snake_case\n```"
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
		return "```qasm\nOPENQASM 3.0;\ninclude \"stdgates.qasm\";\n\nqubit q;\nh q[0];  // ✅ q is declared\n```"
	case "QAS003":
		return "```qasm\nOPENQASM 3.0;\ninclude \"stdgates.qasm\";\n\nqubit q;\nbit c;\n\nh q;          // ✅ Affect qubit before measurement\nmeasure q -> c;\n```"
	case "QAS004":
		return "```qasm\nOPENQASM 3.0;\ninclude \"stdgates.qasm\";\n\nqubit[4] q;\n\nh q[3];  // ✅ Index 3 is within bounds [0,3]\n```"
	case "QAS005":
		return "```qasm\nOPENQASM 3.0;\ninclude \"stdgates.qasm\";\n\nqubit myQubit;  // ✅ Starts with lowercase\n```"
	case "QAS006":
		return "```qasm\nOPENQASM 3.0;\n\ngate mygate(a, b) { /* ... */ }\nqubit[2] q1;\nqubit[2] q2;\nmygate q1, q2;  // ✅ Both registers size 2\n```"
	case "QAS007":
		return "```qasm\nOPENQASM 3.0;\n\ngate mygate(a) {\n  h a;  // ✅ Use parameter directly\n}\n```"
	case "QAS008":
		return "```qasm\nOPENQASM 3.0;\n\nqubit q;  // ✅ Global scope declaration\n\ndef foo() {\n  h q;  // ✅ Use global qubit\n}\n```"
	case "QAS009":
		return "```qasm\nOPENQASM 3.0;\n\nfor i in [0:1:10] {\n  if (i == 5) break;  // ✅ Inside loop\n}\n```"
	case "QAS010":
		return "```qasm\nOPENQASM 3.0;\n\ngate g(a) {\n  h a;  // ✅ Unitary operation only\n}\n```"
	case "QAS011":
		return "```qasm\nOPENQASM 3.0;\n\nqubit my_qubit;  // ✅ Normal identifier\n```"
	case "QAS012":
		return "```qasm\nOPENQASM 3.0;\n\nqubit my_qubit;  // ✅ snake_case\ngate my_gate(q) { h q; }  // ✅ snake_case\n```"
	default:
		return "```qasm\n// Example code that follows this rule\n```"
	}
}

// generateRelatedRules generates a list of related rules
func (dg *DocumentationGenerator) generateRelatedRules(rule *Rule) string {
	switch rule.ID {
	case "QAS001":
		return "- [QAS005](QAS005.md) (naming-convention-violation): Both relate to declaration best practices\n- [QAS008](QAS008.md) (qubit-declared-in-local-scope): Both relate to qubit declaration\n"
	case "QAS002":
		return "- [QAS004](QAS004.md) (out-of-bounds-index): Both relate to identifier usage validation\n- [QAS008](QAS008.md) (qubit-declared-in-local-scope): Both relate to scope management\n"
	case "QAS003":
		return "- [QAS010](QAS010.md) (invalid-instruction-in-gate): Both relate to measurement operations\n"
	case "QAS004":
		return "- [QAS002](QAS002.md) (undefined-identifier): Both relate to identifier validation\n- [QAS006](QAS006.md) (gate-register-size-mismatch): Both relate to array bounds checking\n"
	case "QAS005":
		return "- [QAS001](QAS001.md) (unused-qubit): Both relate to declaration best practices\n- [QAS012](QAS012.md) (snake-case-required): Both relate to naming conventions\n"
	case "QAS006":
		return "- [QAS004](QAS004.md) (out-of-bounds-index): Both relate to size validation\n- [QAS007](QAS007.md) (gate-parameter-indexing): Both relate to gate parameter handling\n"
	case "QAS007":
		return "- [QAS006](QAS006.md) (gate-register-size-mismatch): Both relate to gate parameter handling\n- [QAS010](QAS010.md) (invalid-instruction-in-gate): Both relate to gate definition constraints\n"
	case "QAS008":
		return "- [QAS001](QAS001.md) (unused-qubit): Both relate to qubit management\n- [QAS002](QAS002.md) (undefined-identifier): Both relate to scope validation\n"
	case "QAS009":
		return "- [QAS010](QAS010.md) (invalid-instruction-in-gate): Both relate to syntax constraints\n"
	case "QAS010":
		return "- [QAS003](QAS003.md) (constant-measured-bit): Both relate to measurement operations\n- [QAS007](QAS007.md) (gate-parameter-indexing): Both relate to gate definition constraints\n- [QAS009](QAS009.md) (illegal-break-continue): Both relate to syntax constraints\n"
	case "QAS011":
		return "- [QAS005](QAS005.md) (naming-convention-violation): Both relate to identifier naming\n- [QAS012](QAS012.md) (snake-case-required): Both relate to naming standards\n"
	case "QAS012":
		return "- [QAS005](QAS005.md) (naming-convention-violation): Both relate to naming conventions\n- [QAS011](QAS011.md) (reserved-prefix-usage): Both relate to naming standards\n"
	default:
		return "None currently identified.\n"
	}
}