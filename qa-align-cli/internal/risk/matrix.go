package risk

import "strings"

// ComputeMatrixScore calculates priorities based on churn and keyword patterns
func ComputeMatrixScore(filePath, descriptiveText string, churnRate int) int {
	impact := 2 // Default baseline medium score assignment

	// Check directory patterns and keyword heuristics
	lowText := strings.ToLower(descriptiveText + filePath)
	if strings.Contains(lowText, "auth") || strings.Contains(lowText, "billing") || strings.Contains(lowText, "security") {
		impact = 3 // High visibility safety boundary tier
	} else if strings.Contains(lowText, "utils") || strings.Contains(lowText, "styles") {
		impact = 1 // Cosmetic tier
	}

	frequency := 1
	if churnRate > 10 {
		frequency = 3
	} else if churnRate > 2 {
		frequency = 2
	}

	probability := 2 // Structural complexity mid-point anchor

	return impact * frequency * probability
}
