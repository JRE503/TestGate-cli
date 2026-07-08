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

	// funcRegex captures test function declarations across all supported languages.
	// Group 1: keyword-prefixed style (def/function/void/async function)
	// Group 2: test*() name style — catches most framework naming conventions
	// Group 3: ESP-IDF Unity + Catch2 TEST_CASE("name", "[tag]") macro
	// Group 4: Catch2 SECTION("name") — nested test blocks inside TEST_CASE
	funcRegex = regexp.MustCompile(
		`^\s*(?:` +
			`(?:public\s+(?:static\s+)?void|private\s+void|protected\s+void|async\s+function|function|def)\s+(\w+)` +
			`|(\w*[Tt]est\w*)\s*\(` +
			`|(?:TEST_CASE|TEST_F|TEST)\s*\(\s*"([^"]+)"` +
			`|SECTION\s*\(\s*"([^"]+)"` +
			`)`,
	)

	// commentRegex guards funcRegex from triggering on annotation comment lines.
	// e.g. "// @test Creates..." must NOT create a function record.
	commentRegex = regexp.MustCompile(`^\s*(?://|#|\*|/\*)`)
)

// ParseMetadata scans a source file and extracts test case telemetry.
// When includeUnannotated is true, test functions without annotations are
// included in the output with sentinel values (what="UNANNOTATED").
func ParseMetadata(filePath string, includeUnannotated bool) ([]schema.TestCaseMetadata, error) {
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

		// Collect annotation tags from comment lines
		if matches := whatRegex.FindStringSubmatch(line); len(matches) > 1 {
			currentWhat = strings.TrimSpace(matches[1])
		}
		if matches := whyRegex.FindStringSubmatch(line); len(matches) > 1 {
			currentWhy = strings.TrimSpace(matches[1])
		}
		if matches := refRegex.FindStringSubmatch(line); len(matches) > 1 {
			currentRef = strings.TrimSpace(matches[1])
		}

		// Skip comment lines for function detection — prevents false positives
		// where "// @test Creates..." would be treated as a function declaration
		if commentRegex.MatchString(line) {
			continue
		}

		// Detect function/test declarations
		funcMatches := funcRegex.FindStringSubmatch(line)
		if len(funcMatches) == 0 {
			continue
		}

		// Resolve method name from whichever capture group matched
		methodName := funcMatches[1]
		if methodName == "" {
			methodName = funcMatches[2]
		}
		if methodName == "" {
			methodName = funcMatches[3] // TEST_CASE / TEST_F / TEST
		}
		if methodName == "" {
			methodName = funcMatches[4] // SECTION
		}
		if methodName == "" {
			continue
		}

		hasAnnotation := currentWhat != "" || currentRef != ""

		if hasAnnotation {
			// Fully annotated record
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
			currentWhat, currentWhy, currentRef = "", "", ""
		} else if includeUnannotated {
			// Unannotated record — emit with sentinel values when flag is set
			cases = append(cases, schema.TestCaseMetadata{
				TestMethod:    methodName,
				FilePath:      filePath,
				What:          "UNANNOTATED",
				Why:           "",
				RequirementID: "UNMAPPED",
			})
			// Don't reset — there were no annotations to consume
		} else {
			// No annotation and flag not set — reset and skip
			currentWhat, currentWhy, currentRef = "", "", ""
		}
	}

	return cases, scanner.Err()
}
