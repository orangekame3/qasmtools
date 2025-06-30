package lint

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed rules/*.yaml
var embeddedRulesFS embed.FS

// RuleLoader loads rules from YAML files
type RuleLoader struct {
	rulesDir    string
	useEmbedded bool
}

// NewRuleLoader creates a new rule loader
func NewRuleLoader(rulesDir string) *RuleLoader {
	// Use embedded rules if no custom directory specified or default directory
	useEmbedded := rulesDir == "" || rulesDir == "lint/rules"
	
	return &RuleLoader{
		rulesDir:    rulesDir,
		useEmbedded: useEmbedded,
	}
}

// LoadRules loads all rules from the rules directory
func (l *RuleLoader) LoadRules() ([]*Rule, error) {
	if l.useEmbedded {
		return l.loadEmbeddedRules()
	}
	return l.loadFileSystemRules()
}

// loadEmbeddedRules loads rules from embedded files
func (l *RuleLoader) loadEmbeddedRules() ([]*Rule, error) {
	var rules []*Rule

	err := fs.WalkDir(embeddedRulesFS, "rules", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || (!strings.HasSuffix(path, ".yaml") && !strings.HasSuffix(path, ".yml")) {
			return nil
		}

		data, err := embeddedRulesFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read embedded rule %s: %w", path, err)
		}

		rule, err := l.parseRuleYAML(data)
		if err != nil {
			return fmt.Errorf("failed to parse embedded rule %s: %w", path, err)
		}

		if rule.Enabled {
			rules = append(rules, rule)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return rules, nil
}

// loadFileSystemRules loads rules from filesystem
func (l *RuleLoader) loadFileSystemRules() ([]*Rule, error) {
	var rules []*Rule

	err := filepath.Walk(l.rulesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(info.Name(), ".yaml") && !strings.HasSuffix(info.Name(), ".yml") {
			return nil
		}

		rule, err := l.loadRuleFromFile(path)
		if err != nil {
			return fmt.Errorf("failed to load rule from %s: %w", path, err)
		}

		if rule.Enabled {
			rules = append(rules, rule)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return rules, nil
}

// loadRuleFromFile loads a single rule from a YAML file
func (l *RuleLoader) loadRuleFromFile(filename string) (*Rule, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return l.parseRuleYAML(data)
}

// parseRuleYAML parses YAML data into a Rule struct
func (l *RuleLoader) parseRuleYAML(data []byte) (*Rule, error) {
	var rule Rule
	err := yaml.Unmarshal(data, &rule)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Validate required fields
	if rule.ID == "" {
		return nil, fmt.Errorf("rule ID is required")
	}
	if rule.Name == "" {
		return nil, fmt.Errorf("rule name is required")
	}
	if rule.Message == "" {
		return nil, fmt.Errorf("rule message is required")
	}

	// Set defaults
	if rule.Level == "" {
		rule.Level = SeverityWarning
	}

	return &rule, nil
}

// LoadRule loads a specific rule by ID
func (l *RuleLoader) LoadRule(ruleID string) (*Rule, error) {
	rules, err := l.LoadRules()
	if err != nil {
		return nil, err
	}

	for _, rule := range rules {
		if rule.ID == ruleID {
			return rule, nil
		}
	}

	return nil, fmt.Errorf("rule %s not found", ruleID)
}