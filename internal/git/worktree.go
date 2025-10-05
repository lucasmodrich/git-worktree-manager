package git

import (
	"fmt"
	"strings"
)

// WorktreeAdd creates a new worktree at the specified path for the given branch
func (c *Client) WorktreeAdd(path, branch string, track bool) error {
	args := []string{"worktree", "add"}

	if track {
		// Create new branch: git worktree add -b <branch> <path>
		args = append(args, "-b", branch, path)
	} else {
		// Checkout existing branch: git worktree add <path> <branch>
		args = append(args, path, branch)
	}

	_, stderr, err := c.ExecGit(args...)
	if err != nil {
		return fmt.Errorf("failed to add worktree: %w", err)
	}

	if stderr != "" && !strings.Contains(stderr, "Preparing worktree") {
		// Git sometimes outputs to stderr even on success
		return fmt.Errorf("worktree add warnings: %s", stderr)
	}

	return nil
}

// WorktreeList returns a list of all worktrees
func (c *Client) WorktreeList() ([]string, error) {
	stdout, _, err := c.ExecGit("worktree", "list")
	if err != nil {
		return nil, fmt.Errorf("failed to list worktrees: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(stdout), "\n")
	var worktrees []string
	for _, line := range lines {
		if line != "" {
			worktrees = append(worktrees, line)
		}
	}

	return worktrees, nil
}

// WorktreeRemove removes the worktree at the specified path
func (c *Client) WorktreeRemove(path string) error {
	_, _, err := c.ExecGit("worktree", "remove", path)
	if err != nil {
		return fmt.Errorf("failed to remove worktree: %w", err)
	}

	return nil
}

// WorktreePrune prunes stale worktree references
func (c *Client) WorktreePrune() error {
	_, _, err := c.ExecGit("worktree", "prune")
	if err != nil {
		return fmt.Errorf("failed to prune worktrees: %w", err)
	}

	return nil
}
