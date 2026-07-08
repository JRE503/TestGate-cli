package schema

import (
	"encoding/json"
	"os"
)

// TestCaseMetadata holds all telemetry fields for a single test case
type TestCaseMetadata struct {
	TestMethod      string `json:"test_method"`
	FilePath        string `json:"file_path"`
	What            string `json:"what"`
	Why             string `json:"why"`
	RequirementID   string `json:"requirement_id"`
	ChangeFrequency int    `json:"change_frequency_30_days"`
	RiskScore       int    `json:"calculated_risk_score"`
}

// WriteTelemetryPayload materializes metadata into a single schema json asset
func WriteTelemetryPayload(outputPath string, records []TestCaseMetadata) error {
	jsonData, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(outputPath, jsonData, 0644)
}
