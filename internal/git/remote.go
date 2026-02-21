package git

import (
	"fmt"
	"strings"
)

// Clone clones a repository to the specified target directory
func (c *Client) Clone(url, target string, bare bool) error {
	args := []string{"clone"}

	if bare {
		args = append(args, "--bare")
	}

	args = append(args, url, target)

	_, _, err := c.ExecGit(args...)
	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	return nil
}

// Fetch fetches from the remote repository
func (c *Client) Fetch(all, prune bool) error {
	args := []string{"fetch"}

	if all {
		args = append(args, "--all")
	}

	if prune {
		args = append(args, "--prune")
	}

	_, _, err := c.ExecGit(args...)
	if err != nil {
		return fmt.Errorf("failed to fetch: %w", err)
	}

	return nil
}

// Push pushes the specified branch to the remote
func (c *Client) Push(branch string, setUpstream bool) error {
	args := []string{"push"}

	if setUpstream {
		args = append(args, "-u", "origin", branch)
	} else {
		args = append(args, "origin", branch)
	}

	_, _, err := c.ExecGit(args...)
	if err != nil {
		return fmt.Errorf("failed to push: %w", err)
	}

	return nil
}

// DetectDefaultBranch detects the default branch of the remote repository
func (c *Client) DetectDefaultBranch() (string, error) {
	// Try to get the default branch from symbolic-ref
	stdout, _, err := c.ExecGit("symbolic-ref", "refs/remotes/origin/HEAD")
	if err == nil && stdout != "" {
		// Parse output like "refs/remotes/origin/main"
		parts := strings.Split(strings.TrimSpace(stdout), "/")
		if len(parts) > 0 {
			return parts[len(parts)-1], nil
		}
	}

	// Fallback: try to detect from remote show
	stdout, _, err = c.ExecGit("remote", "show", "origin")
	if err != nil {
		// Last fallback: check if main or master exists locally
		if c.BranchExists("main", false) {
			return "main", nil
		}
		if c.BranchExists("master", false) {
			return "master", nil
		}
		return "", fmt.Errorf("failed to detect default branch: %w", err)
	}

	// Parse "HEAD branch: main" from output
	lines := strings.Split(stdout, "\n")
	for _, line := range lines {
		if strings.Contains(line, "HEAD branch:") {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1]), nil
			}
		}
	}

	// Fallback to common defaults
	if c.BranchExists("main", false) {
		return "main", nil
	}
	if c.BranchExists("master", false) {
		return "master", nil
	}

	return "", fmt.Errorf("could not detect default branch")
}
