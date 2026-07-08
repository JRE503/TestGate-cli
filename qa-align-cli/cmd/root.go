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

var targetDir string

var rootCmd = &cobra.Command{
	Use:   "qa-align",
	Short: "QA-Align Core Telemetry Engine CLI",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("🚀 Initiating Repository Intelligence Scan inside: %s\n", targetDir)

		// 1. Discover relevant files recursively
		var testFiles []string
		err := filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// Block asset/dependency indexing tracking loops
			if info.IsDir() && (info.Name() == "node_modules" || info.Name() == "venv" || info.Name() == ".git") {
				return filepath.SkipDir
			}
			if !info.IsDir() {
				ext := filepath.Ext(path)
				if ext == ".py" || ext == ".ts" || ext == ".js" || ext == ".java" || ext == ".cpp" || ext == ".c" || ext == ".kt" {
					testFiles = append(testFiles, path)
				}
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("failed codebase tree file discovery crawl: %w", err)
		}

		fmt.Printf("📍 Identified %d matching polyglot test source scripts for processing.\n", len(testFiles))

		var masterPayload []schema.TestCaseMetadata

		// 2. Process discovered test files contextually
		for _, file := range testFiles {
			// Extract plain-text telemetry attributes securely (Read-Only)
			records, err := parser.ParseMetadata(file)
			if err != nil {
				fmt.Printf("⚠️  Skipped processing file [%s]: %v\n", file, err)
				continue
			}

			// Crunch Git commit logs and calculate Risk Metrics
			for i := range records {
				churn := gitops.CalculateFileChurn(targetDir, file, 30)
				records[i].ChangeFrequency = churn
				records[i].RiskScore = risk.ComputeMatrixScore(file, records[i].What, churn)

				masterPayload = append(masterPayload, records[i])
			}
		}

		// 3. Serialize output cleanly to the target JSON schema artifact
		outputFile := "telemetry.json"
		if err := schema.WriteTelemetryPayload(outputFile, masterPayload); err != nil {
			return fmt.Errorf("failed syncing compliance matrix arrays to disk: %w", err)
		}

		fmt.Printf("✨ Extraction sequence successfully compiled! Mapped matrix metrics saved to `%s`.\n", outputFile)
		return nil
	},
}

func init() {
	rootCmd.Flags().StringVarP(&targetDir, "dir", "d", ".", "Target repository directory path to audit")
}

func Execute() error {
	return rootCmd.Execute()
}
