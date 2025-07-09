package lint

import (
	"sync"
	"time"

	"github.com/orangekame3/qasmtools/lint/ast"
	"github.com/orangekame3/qasmtools/parser"
)

// PerformanceStats tracks linting performance metrics
type PerformanceStats struct {
	TotalFiles      int
	TotalTime       time.Duration
	ASTRulesUsed    int
	TextRulesUsed   int
	ParseTime       time.Duration
	AnalysisTime    time.Duration
	ViolationsFound int
}

// LinterWithMetrics wraps the linter with performance tracking
type LinterWithMetrics struct {
	*Linter
	stats PerformanceStats
	mutex sync.RWMutex
}

// NewLinterWithMetrics creates a new performance-tracked linter
func NewLinterWithMetrics(rulesDir string) *LinterWithMetrics {
	return &LinterWithMetrics{
		Linter: NewLinter(rulesDir),
		stats:  PerformanceStats{},
	}
}

// GetStats returns current performance statistics
func (l *LinterWithMetrics) GetStats() PerformanceStats {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.stats
}

// ResetStats resets performance statistics
func (l *LinterWithMetrics) ResetStats() {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.stats = PerformanceStats{}
}

// LintFileWithMetrics lints a file with performance tracking
func (l *LinterWithMetrics) LintFileWithMetrics(filename string) ([]*Violation, error) {
	start := time.Now()
	defer func() {
		l.mutex.Lock()
		l.stats.TotalFiles++
		l.stats.TotalTime += time.Since(start)
		l.mutex.Unlock()
	}()

	violations, err := l.LintFile(filename)
	if err != nil {
		return nil, err
	}

	l.mutex.Lock()
	l.stats.ViolationsFound += len(violations)
	l.mutex.Unlock()

	return violations, nil
}

// BatchLinter provides optimized batch linting capabilities
type BatchLinter struct {
	*LinterWithMetrics
	concurrency int
}

// NewBatchLinter creates a new batch linter with configurable concurrency
func NewBatchLinter(rulesDir string, concurrency int) *BatchLinter {
	if concurrency <= 0 {
		concurrency = 4 // Default to 4 workers
	}
	return &BatchLinter{
		LinterWithMetrics: NewLinterWithMetrics(rulesDir),
		concurrency:       concurrency,
	}
}

// LintFilesParallel lints multiple files in parallel
func (l *BatchLinter) LintFilesParallel(filenames []string) ([]*Violation, error) {
	// Load rules once
	if err := l.LoadRules(); err != nil {
		return nil, err
	}

	type result struct {
		violations []*Violation
		err        error
	}

	jobs := make(chan string, len(filenames))
	results := make(chan result, len(filenames))

	// Start workers
	for i := 0; i < l.concurrency; i++ {
		go func() {
			for filename := range jobs {
				violations, err := l.LintFileWithMetrics(filename)
				results <- result{violations: violations, err: err}
			}
		}()
	}

	// Send jobs
	for _, filename := range filenames {
		jobs <- filename
	}
	close(jobs)

	// Collect results
	var allViolations []*Violation
	var firstError error

	for i := 0; i < len(filenames); i++ {
		res := <-results
		if res.err != nil && firstError == nil {
			firstError = res.err
		}
		allViolations = append(allViolations, res.violations...)
	}

	return allViolations, firstError
}

// OptimizedLinter provides additional optimizations for AST-based linting
type OptimizedLinter struct {
	*Linter
	astRuleCache map[string]ast.ASTRule
	programCache map[string]*parser.Program
	cacheMutex   sync.RWMutex
}

// NewOptimizedLinter creates a new optimized linter with caching
func NewOptimizedLinter(rulesDir string) *OptimizedLinter {
	return &OptimizedLinter{
		Linter:       NewLinter(rulesDir),
		astRuleCache: make(map[string]ast.ASTRule),
		programCache: make(map[string]*parser.Program),
	}
}

// LintFileOptimized lints a file with caching optimizations
func (l *OptimizedLinter) LintFileOptimized(filename string, content string) ([]*Violation, error) {
	// Check program cache first
	l.cacheMutex.RLock()
	program, cached := l.programCache[filename]
	l.cacheMutex.RUnlock()

	if !cached {
		// Parse and cache
		p := parser.NewParser()
		result := p.ParseWithErrors(content)
		if result.HasErrors() {
			return nil, &result.Errors[0]
		}
		program = result.Program

		l.cacheMutex.Lock()
		l.programCache[filename] = program
		l.cacheMutex.Unlock()
	}

	// Build usage map
	usageMap := l.buildUsageMap(program)

	// Create check context
	context := &CheckContext{
		File:     filename,
		Content:  content,
		Program:  program,
		UsageMap: usageMap,
	}

	var allViolations []*Violation

	// Run optimized AST rules
	for _, rule := range l.rules {
		if !rule.Enabled {
			continue
		}

		var violations []*Violation

		// Use cached AST rules for better performance
		if l.useAST {
			l.cacheMutex.RLock()
			astRule, exists := l.astRuleCache[rule.ID]
			l.cacheMutex.RUnlock()

			if !exists {
				astRule = CreateASTRule(rule.ID)
				l.cacheMutex.Lock()
				l.astRuleCache[rule.ID] = astRule
				l.cacheMutex.Unlock()
			}

			if astRule != nil {
				astCtx := l.convertToASTContext(context)
				astViolations := astRule.CheckAST(program, astCtx)
				violations = l.convertASTViolations(astViolations)
			}
		}

		// Set rule reference for each violation
		for _, violation := range violations {
			violation.Rule = rule
			violation.Severity = rule.Level
		}

		allViolations = append(allViolations, violations...)
	}

	return allViolations, nil
}

// ClearCache clears the internal caches
func (l *OptimizedLinter) ClearCache() {
	l.cacheMutex.Lock()
	defer l.cacheMutex.Unlock()
	l.astRuleCache = make(map[string]ast.ASTRule)
	l.programCache = make(map[string]*parser.Program)
}

// GetCacheStats returns cache statistics
func (l *OptimizedLinter) GetCacheStats() (astRules int, programs int) {
	l.cacheMutex.RLock()
	defer l.cacheMutex.RUnlock()
	return len(l.astRuleCache), len(l.programCache)
}
