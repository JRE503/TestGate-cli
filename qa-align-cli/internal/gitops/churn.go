package gitops

import (
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// CalculateFileChurn parses local git log references contextually
func CalculateFileChurn(repoRoot, targetFilePath string, dayWindow int) int {
	repo, err := git.PlainOpen(repoRoot)
	if err != nil {
		return 1 // Default baseline assignment fallback if un-versioned workspace
	}

	logIterator, err := repo.Log(&git.LogOptions{FileName: &targetFilePath})
	if err != nil {
		return 1
	}
	defer logIterator.Close()

	churnCount := 0
	cutoffDate := time.Now().AddDate(0, 0, -dayWindow)

	_ = logIterator.ForEach(func(commit *object.Commit) error {
		if commit.Author.When.After(cutoffDate) {
			churnCount++
		}
		return nil
	})

	if churnCount == 0 {
		return 1
	}
	return churnCount
}
