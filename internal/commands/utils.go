package commands

import (
	"fmt"
	"os"
	"path/filepath"
)

// findWorktreeRoot walks up from the current directory to find the root of a
// worktree-managed repository. A worktree-managed repo has a .git file (not
// directory) at its root. Returns the absolute path to the root, or an error
// if no such root is found.
func findWorktreeRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	for {
		info, err := os.Stat(filepath.Join(dir, ".git"))
		if err == nil && !info.IsDir() {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break // reached filesystem root
		}
		dir = parent
	}

	return "", fmt.Errorf("not in a worktree-managed repository")
}
