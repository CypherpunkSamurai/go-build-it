package utils

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
)

// GitCloneTemp - Clones a repository to a temporary directory
func GitCloneTemp(gitUrl string) (string, error) {
	// tempDir
	tempDir, err := os.MkdirTemp("", "go-build-")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}
	// clone
	_, err = git.PlainClone(tempDir, false, &git.CloneOptions{
		URL: gitUrl,
		// SingleBranch: true,
		// Progress:     os.Stdout,
	})
	if err != nil {
		return "", fmt.Errorf("failed to clone repository: %w", err)
	}
	// return tempDir
	return tempDir, nil
}
