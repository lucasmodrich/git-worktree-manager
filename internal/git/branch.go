package git

import (
	"fmt"
	"strings"
)

// BranchExists checks if a branch exists locally or remotely
func (c *Client) BranchExists(name string, remote bool) bool {
	var args []string
	if remote {
		args = []string{"branch", "-r", "--list", fmt.Sprintf("origin/%s", name)}
	} else {
		args = []string{"branch", "--list", name}
	}

	stdout, _, err := c.ExecGit(args...)
	if err != nil {
		return false
	}

	return strings.TrimSpace(stdout) != ""
}

// CreateBranch creates a new branch from the specified base branch
func (c *Client) CreateBranch(name, baseBranch string) error {
	args := []string{"branch", name}
	if baseBranch != "" {
		args = append(args, baseBranch)
	}

	_, _, err := c.ExecGit(args...)
	if err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}

	return nil
}

// DeleteBranch deletes a local branch
func (c *Client) DeleteBranch(name string, force bool) error {
	flag := "-d"
	if force {
		flag = "-D"
	}

	_, _, err := c.ExecGit("branch", flag, name)
	if err != nil {
		return fmt.Errorf("failed to delete branch: %w", err)
	}

	return nil
}

// DeleteRemoteBranch deletes a remote branch
func (c *Client) DeleteRemoteBranch(name string) error {
	_, _, err := c.ExecGit("push", "origin", "--delete", name)
	if err != nil {
		return fmt.Errorf("failed to delete remote branch: %w", err)
	}

	return nil
}
