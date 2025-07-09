package commands

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/orangekame3/qasmtools/parser"
	"github.com/spf13/cobra"
)

// BenchmarkResult represents the results of a single benchmark run
type BenchmarkResult struct {
	Filename   string        `json:"filename"`
	FileSize   int64         `json:"file_size"`
	Statements int           `json:"statements"`
	ParseTime  time.Duration `json:"parse_time"`
	Success    bool          `json:"success"`
	Error      string        `json:"error,omitempty"`
	Throughput float64       `json:"throughput"` // bytes/second
}

// BenchmarkSuite represents a complete benchmark suite
type BenchmarkSuite struct {
	Name      string            `json:"name"`
	Version   string            `json:"version"`
	Results   []BenchmarkResult `json:"results"`
	Summary   BenchmarkSummary  `json:"summary"`
	Timestamp time.Time         `json:"timestamp"`
}

// BenchmarkSummary provides aggregate statistics
type BenchmarkSummary struct {
	TotalFiles      int           `json:"total_files"`
	SuccessfulFiles int           `json:"successful_files"`
	FailedFiles     int           `json:"failed_files"`
	SuccessRate     float64       `json:"success_rate"`
	TotalTime       time.Duration `json:"total_time"`
	AverageTime     time.Duration `json:"average_time"`
	MinTime         time.Duration `json:"min_time"`
	MaxTime         time.Duration `json:"max_time"`
	MedianTime      time.Duration `json:"median_time"`
	StdDeviation    time.Duration `json:"std_deviation"`
	TotalStatements int           `json:"total_statements"`
	TotalBytes      int64         `json:"total_bytes"`
	AvgThroughput   float64       `json:"avg_throughput"`
}

// NewBenchmarkCommand creates the benchmark command
func NewBenchmarkCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "benchmark [directory]",
		Short: "Run comprehensive parser benchmarks",
		Long:  `Run comprehensive parser benchmarks against QASM files to measure performance and compatibility.`,
		RunE:  runBenchmark,
	}

	// Add flags
	cmd.Flags().Int("runs", 10, "Number of parsing runs per file")
	cmd.Flags().Int("warmup", 3, "Number of warmup runs before benchmarking")
	cmd.Flags().Bool("detailed", false, "Show detailed results for each file")
	cmd.Flags().String("output", "text", "Output format (text, json, csv)")
	cmd.Flags().String("filter", "*.qasm", "File pattern to match")
	cmd.Flags().Bool("recursive", true, "Search subdirectories recursively")
	cmd.Flags().Bool("ignore-errors", false, "Continue benchmarking even if some files fail")
	cmd.Flags().String("suite-name", "qasmtools", "Name of the benchmark suite")

	return cmd
}

func runBenchmark(cmd *cobra.Command, args []string) error {
	// Get flags
	runs, _ := cmd.Flags().GetInt("runs")
	warmup, _ := cmd.Flags().GetInt("warmup")
	detailed, _ := cmd.Flags().GetBool("detailed")
	outputFormat, _ := cmd.Flags().GetString("output")
	filter, _ := cmd.Flags().GetString("filter")
	recursive, _ := cmd.Flags().GetBool("recursive")
	ignoreErrors, _ := cmd.Flags().GetBool("ignore-errors")
	suiteName, _ := cmd.Flags().GetString("suite-name")

	// Determine directory to benchmark
	directory := "."
	if len(args) > 0 {
		directory = args[0]
	}

	// Find QASM files
	files, err := findQASMFiles(directory, filter, recursive)
	if err != nil {
		return fmt.Errorf("failed to find QASM files: %w", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("no QASM files found in directory: %s", directory)
	}

	fmt.Printf("🔍 Found %d QASM files to benchmark\n", len(files))
	fmt.Printf("⚙️  Configuration: %d runs per file, %d warmup runs\n", runs, warmup)
	fmt.Printf("📊 Starting benchmark suite...\n\n")

	// Create benchmark suite
	suite := &BenchmarkSuite{
		Name:      suiteName,
		Version:   "1.0.0", // TODO: Get from version
		Results:   make([]BenchmarkResult, 0, len(files)),
		Timestamp: time.Now(),
	}

	// Run benchmarks
	for i, file := range files {
		fmt.Printf("📄 [%d/%d] %s", i+1, len(files), filepath.Base(file))

		result := benchmarkFile(file, runs, warmup)
		suite.Results = append(suite.Results, result)

		if result.Success {
			fmt.Printf(" ✅ %v (%d statements)\n", result.ParseTime, result.Statements)
		} else {
			fmt.Printf(" ❌ %s\n", result.Error)
			if !ignoreErrors {
				return fmt.Errorf("benchmarking failed for %s: %s", file, result.Error)
			}
		}
	}

	// Calculate summary statistics
	suite.Summary = calculateSummary(suite.Results)

	// Output results
	switch outputFormat {
	case "json":
		return outputBenchmarkJSON(suite)
	case "csv":
		return outputCSV(suite)
	default:
		return outputText(suite, detailed)
	}
}

func findQASMFiles(directory, pattern string, recursive bool) ([]string, error) {
	var files []string

	if recursive {
		err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() {
				if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
					files = append(files, path)
				}
			}
			return nil
		})
		return files, err
	}

	// Non-recursive search
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			if matched, _ := filepath.Match(pattern, entry.Name()); matched {
				files = append(files, filepath.Join(directory, entry.Name()))
			}
		}
	}

	return files, nil
}

func benchmarkFile(filename string, runs, warmup int) BenchmarkResult {
	result := BenchmarkResult{
		Filename: filename,
		Success:  false,
	}

	// Get file size
	if stat, err := os.Stat(filename); err == nil {
		result.FileSize = stat.Size()
	}

	// Read file content
	content, err := os.ReadFile(filename)
	if err != nil {
		result.Error = fmt.Sprintf("Failed to read file: %v", err)
		return result
	}

	contentStr := string(content)

	// Warmup runs
	p := parser.NewParser()
	for i := 0; i < warmup; i++ {
		_, _ = p.ParseString(contentStr)
	}

	// Benchmark runs
	var durations []time.Duration
	var program *parser.Program
	var lastErr error

	for i := 0; i < runs; i++ {
		p := parser.NewParser()
		start := time.Now()
		prog, parseErr := p.ParseString(contentStr)
		duration := time.Since(start)

		if parseErr != nil {
			lastErr = parseErr
			continue
		}

		durations = append(durations, duration)
		if i == 0 {
			program = prog
		}
	}

	if len(durations) == 0 {
		result.Error = fmt.Sprintf("All parsing attempts failed: %v", lastErr)
		return result
	}

	// Calculate average time
	var totalTime time.Duration
	for _, d := range durations {
		totalTime += d
	}
	avgTime := totalTime / time.Duration(len(durations))

	// Set results
	result.Success = true
	result.ParseTime = avgTime
	result.Throughput = float64(result.FileSize) / avgTime.Seconds()

	if program != nil {
		result.Statements = len(program.Statements)
	}

	return result
}

func calculateSummary(results []BenchmarkResult) BenchmarkSummary {
	summary := BenchmarkSummary{
		TotalFiles: len(results),
		MinTime:    time.Duration(math.MaxInt64),
		MaxTime:    0,
	}

	var successfulTimes []time.Duration
	var totalThroughput float64

	for _, result := range results {
		summary.TotalBytes += result.FileSize
		summary.TotalStatements += result.Statements

		if result.Success {
			summary.SuccessfulFiles++
			summary.TotalTime += result.ParseTime
			successfulTimes = append(successfulTimes, result.ParseTime)
			totalThroughput += result.Throughput

			if result.ParseTime < summary.MinTime {
				summary.MinTime = result.ParseTime
			}
			if result.ParseTime > summary.MaxTime {
				summary.MaxTime = result.ParseTime
			}
		} else {
			summary.FailedFiles++
		}
	}

	if summary.SuccessfulFiles > 0 {
		summary.SuccessRate = float64(summary.SuccessfulFiles) / float64(summary.TotalFiles) * 100
		summary.AverageTime = summary.TotalTime / time.Duration(summary.SuccessfulFiles)
		summary.AvgThroughput = totalThroughput / float64(summary.SuccessfulFiles)

		// Calculate median
		sort.Slice(successfulTimes, func(i, j int) bool {
			return successfulTimes[i] < successfulTimes[j]
		})
		mid := len(successfulTimes) / 2
		if len(successfulTimes)%2 == 0 {
			summary.MedianTime = (successfulTimes[mid-1] + successfulTimes[mid]) / 2
		} else {
			summary.MedianTime = successfulTimes[mid]
		}

		// Calculate standard deviation
		var variance float64
		for _, t := range successfulTimes {
			diff := t.Seconds() - summary.AverageTime.Seconds()
			variance += diff * diff
		}
		variance /= float64(len(successfulTimes))
		summary.StdDeviation = time.Duration(math.Sqrt(variance) * float64(time.Second))
	}

	return summary
}

func outputText(suite *BenchmarkSuite, detailed bool) error {
	fmt.Printf("\n📊 Benchmark Results Summary\n")
	fmt.Printf("==========================================\n")
	fmt.Printf("Suite: %s v%s\n", suite.Name, suite.Version)
	fmt.Printf("Timestamp: %s\n", suite.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Printf("==========================================\n\n")

	s := suite.Summary
	fmt.Printf("📁 Files: %d total, %d successful, %d failed\n", s.TotalFiles, s.SuccessfulFiles, s.FailedFiles)
	fmt.Printf("✅ Success Rate: %.1f%%\n", s.SuccessRate)
	fmt.Printf("📊 Total Bytes: %d\n", s.TotalBytes)
	fmt.Printf("📝 Total Statements: %d\n", s.TotalStatements)
	fmt.Printf("\n⏱️  Timing Statistics:\n")
	fmt.Printf("   Total Time: %v\n", s.TotalTime)
	fmt.Printf("   Average Time: %v\n", s.AverageTime)
	fmt.Printf("   Median Time: %v\n", s.MedianTime)
	fmt.Printf("   Min Time: %v\n", s.MinTime)
	fmt.Printf("   Max Time: %v\n", s.MaxTime)
	fmt.Printf("   Std Deviation: %v\n", s.StdDeviation)
	fmt.Printf("\n🚀 Throughput: %.2f bytes/sec (average)\n", s.AvgThroughput)

	if detailed {
		fmt.Printf("\n📋 Detailed Results:\n")
		fmt.Printf("%-30s %10s %8s %12s %12s\n", "File", "Size", "Stmts", "Time", "Throughput")
		fmt.Printf("%-30s %10s %8s %12s %12s\n", strings.Repeat("-", 30), strings.Repeat("-", 10), strings.Repeat("-", 8), strings.Repeat("-", 12), strings.Repeat("-", 12))

		for _, result := range suite.Results {
			filename := filepath.Base(result.Filename)
			if len(filename) > 30 {
				filename = filename[:27] + "..."
			}

			if result.Success {
				fmt.Printf("%-30s %10d %8d %12v %12.0f\n",
					filename, result.FileSize, result.Statements, result.ParseTime, result.Throughput)
			} else {
				fmt.Printf("%-30s %10d %8s %12s %12s\n",
					filename, result.FileSize, "ERROR", "-", "-")
			}
		}
	}

	return nil
}

func outputBenchmarkJSON(suite *BenchmarkSuite) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(suite)
}

func outputCSV(suite *BenchmarkSuite) error {
	fmt.Printf("filename,file_size,statements,parse_time_ns,success,throughput,error\n")
	for _, result := range suite.Results {
		errorStr := ""
		if !result.Success {
			errorStr = strings.ReplaceAll(result.Error, "\"", "\"\"")
		}

		fmt.Printf("%s,%d,%d,%d,%t,%.2f,\"%s\"\n",
			result.Filename,
			result.FileSize,
			result.Statements,
			result.ParseTime.Nanoseconds(),
			result.Success,
			result.Throughput,
			errorStr,
		)
	}
	return nil
}
