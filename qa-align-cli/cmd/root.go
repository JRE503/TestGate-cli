package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"qa-align-cli/internal/gitops"
	"qa-align-cli/internal/parser"
	"qa-align-cli/internal/risk"
	"qa-align-cli/internal/schema"

	"github.com/spf13/cobra"
)

var (
	targetDir          string
	outputFile         string
	includeUnannotated bool
)

// skipDirs covers dependency, build, and artifact directories across
// embedded (ESP-IDF, PlatformIO, CMake), web, and Python ecosystems.
var skipDirs = map[string]bool{
	"node_modules":    true,
	"venv":            true,
	".git":            true,
	"build":           true,
	"dist":            true,
	".pio":            true,
	"managed_components": true,
	"vendor":          true,
	"third_party":     true,
	"extern":          true,
	"__pycache__":     true,
	".cache":          true,
	"coverage":        true,
	"CMakeFiles":      true,
	"cmake-build-debug": true,
	".gradle":         true,
	"target":          true, // Maven/Cargo
}

var rootCmd = &cobra.Command{
	Use:   "qa-align",
	Short: "QA-Align Core Telemetry Engine CLI",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("🚀 Initiating Repository Intelligence Scan inside: %s\n", targetDir)

		// 1. Discover relevant files recursively, skipping dependency/build dirs
		var testFiles []string
		err := filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				if skipDirs[info.Name()] {
					return filepath.SkipDir
				}
				return nil
			}
			ext := filepath.Ext(path)
			if ext == ".py" || ext == ".ts" || ext == ".js" || ext == ".java" ||
				ext == ".cpp" || ext == ".c" || ext == ".kt" {
				testFiles = append(testFiles, path)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("failed codebase tree file discovery crawl: %w", err)
		}

		fmt.Printf("📍 Identified %d matching polyglot test source scripts for processing.\n", len(testFiles))

		// BUG FIX: initialize as empty slice, not nil — prevents JSON `null` output
		masterPayload := []schema.TestCaseMetadata{}

		// 2. Process discovered test files
		for _, file := range testFiles {
			// BUG FIX: compute repo-relative path for git log — absolute paths
			// never match git history entries, causing churn to always return 1
			relPath, err := filepath.Rel(targetDir, file)
			if err != nil {
				relPath = file // fallback gracefully
			}

			records, err := parser.ParseMetadata(file, includeUnannotated)
			if err != nil {
				fmt.Printf("⚠️  Skipped processing file [%s]: %v\n", file, err)
				continue
			}

			for i := range records {
				churn := gitops.CalculateFileChurn(targetDir, relPath, 30)
				records[i].ChangeFrequency = churn
				records[i].RiskScore = risk.ComputeMatrixScore(file, records[i].What, churn)
				masterPayload = append(masterPayload, records[i])
			}
		}

		// 3. Warn clearly when nothing was mapped — never silently succeed
		if len(masterPayload) == 0 {
			fmt.Println("")
			fmt.Println("⚠️  0 test cases mapped.")
			fmt.Println("   Files were scanned but no annotated tests were found.")
			fmt.Println("   → Add // @test, // @description, // @issue above test functions.")
			fmt.Println("   → Or run with --include-unannotated to capture all test functions.")
			fmt.Println("")
		}

		// 4. Resolve output path — default to telemetry.json inside target dir
		if outputFile == "" {
			outputFile = filepath.Join(targetDir, "telemetry.json")
		}

		if err := schema.WriteTelemetryPayload(outputFile, masterPayload); err != nil {
			return fmt.Errorf("failed syncing compliance matrix arrays to disk: %w", err)
		}

		if len(masterPayload) > 0 {
			fmt.Printf("✨ %d test cases mapped → `%s`\n", len(masterPayload), outputFile)
		}
		return nil
	},
}

func init() {
	rootCmd.Flags().StringVarP(&targetDir, "dir", "d", ".", "Target repository directory path to audit")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output path for telemetry.json (default: <dir>/telemetry.json)")
	rootCmd.Flags().BoolVar(&includeUnannotated, "include-unannotated", false, "Capture all detected test functions, even without annotations")
}

func Execute() error {
	return rootCmd.Execute()
}
