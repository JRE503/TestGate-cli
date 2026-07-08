package parser

import (
	"bufio"
	"os"
	"qa-align-cli/internal/schema"
	"regexp"
	"strings"
)

var (
	whatRegex = regexp.MustCompile(`(?:\[What\]:|@test|\{\s*@code\s+@test\s*\})\s*(.+)`)
	whyRegex  = regexp.MustCompile(`(?:\[Why\]:|@description|\{\s*@code\s+@description\s*\})\s*(.+)`)
	refRegex  = regexp.MustCompile(`(?:\[Reference\]:|@issue|@reference|\{\s*@code\s+@issue\s*\})\s*(\S+)`)
	// funcRegex matches actual function/method declarations only — anchored to
	// non-comment line starts. Requires the keyword at the beginning of trimmed content.
	// Group 1: def/function/void style
	// Group 2: test*() name style
	// Group 3: ESP-IDF Unity TEST_CASE("name", "[tag]") style
	funcRegex    = regexp.MustCompile(`^\s*(?:(?:public\s+(?:static\s+)?void|private\s+void|protected\s+void|async\s+function|function|def)\s+(\w+)|(\w*[Tt]est\w*)\s*\(|TEST_CASE\s*\(\s*"([^"]+)")`)
	commentRegex = regexp.MustCompile(`^\s*(?://|#|\*|/\*)`)
)

// ParseMetadata inspects target file paths and strips away documentation tokens
func ParseMetadata(filePath string) ([]schema.TestCaseMetadata, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cases []schema.TestCaseMetadata
	scanner := bufio.NewScanner(file)

	var currentWhat, currentWhy, currentRef string

	for scanner.Scan() {
		line := scanner.Text()

		// Match tag values inside comment/annotation lines only
		if matches := whatRegex.FindStringSubmatch(line); len(matches) > 1 {
			currentWhat = strings.TrimSpace(matches[1])
		}
		if matches := whyRegex.FindStringSubmatch(line); len(matches) > 1 {
			currentWhy = strings.TrimSpace(matches[1])
		}
		if matches := refRegex.FindStringSubmatch(line); len(matches) > 1 {
			currentRef = strings.TrimSpace(matches[1])
		}

		// Skip comment lines for function detection to prevent false positives
		// e.g. "// @test Creates..." must NOT trigger a function record
		if commentRegex.MatchString(line) {
			continue
		}

		// Match actual function/method declaration lines
		if funcMatches := funcRegex.FindStringSubmatch(line); len(funcMatches) > 0 {
			// Group 1 = named via def/function/void
			// Group 2 = test*( pattern
			// Group 3 = TEST_CASE("name", ...) — ESP-IDF Unity style
			methodName := funcMatches[1]
			if methodName == "" {
				methodName = funcMatches[2]
			}
			if methodName == "" {
				methodName = funcMatches[3]
			}
			if methodName == "" {
				continue
			}

			if currentWhat != "" || currentRef != "" {
				if currentRef == "" {
					currentRef = "UNMAPPED"
				}
				cases = append(cases, schema.TestCaseMetadata{
					TestMethod:    methodName,
					FilePath:      filePath,
					What:          currentWhat,
					Why:           currentWhy,
					RequirementID: currentRef,
				})
				// Refresh memory boundaries for the next block
				currentWhat, currentWhy, currentRef = "", "", ""
			}
		}
	}

	return cases, scanner.Err()
}
