package commands

import (
	"fmt"
	"os"
)

// verifyWorktreeRepo checks if we're in a worktree-managed repository
// A worktree-managed repo has a .git file (not directory) pointing to .bare
func verifyWorktreeRepo() error {
	// Check if .git exists and is a file (not a directory)
	gitPath := ".git"
	info, err := os.Stat(gitPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("not in a worktree-managed repository")
		}
		return fmt.Errorf("failed to check .git: %w", err)
	}

	// In a worktree-managed repo, .git should be a file, not a directory
	if info.IsDir() {
		return fmt.Errorf("not in a worktree-managed repository")
	}

	return nil
}
