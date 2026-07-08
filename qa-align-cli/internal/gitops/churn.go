package gitops

import (
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// BuildChurnMap opens the repository once and walks all commits within
// dayWindow, building a map[relFilePath]commitCount for the whole repo.
// This is O(commits) instead of O(files × commits).
func BuildChurnMap(repoRoot string, dayWindow int) map[string]int {
	churnMap := make(map[string]int)

	repo, err := git.PlainOpen(repoRoot)
	if err != nil {
		return churnMap
	}

	logIter, err := repo.Log(&git.LogOptions{})
	if err != nil {
		return churnMap
	}
	defer logIter.Close()

	cutoff := time.Now().AddDate(0, 0, -dayWindow)

	_ = logIter.ForEach(func(c *object.Commit) error {
		if !c.Author.When.After(cutoff) {
			return nil
		}

		// Get file stats for this commit (one diff per commit, not per file)
		stats, err := c.Stats()
		if err != nil {
			return nil
		}
		for _, s := range stats {
			churnMap[s.Name]++
		}
		return nil
	})

	return churnMap
}

// LookupChurn returns the churn count for a file from a pre-built map.
// Returns 1 as a baseline if the file has no recorded changes.
func LookupChurn(churnMap map[string]int, relFilePath string) int {
	if count, ok := churnMap[relFilePath]; ok && count > 0 {
		return count
	}
	return 1
}
