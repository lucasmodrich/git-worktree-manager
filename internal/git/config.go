package git

import (
	"fmt"
)

// SetConfig sets a git configuration value
func (c *Client) SetConfig(key, value string) error {
	_, _, err := c.ExecGit("config", key, value)
	if err != nil {
		return fmt.Errorf("failed to set config %s: %w", key, err)
	}

	return nil
}

// ConfigureFetchRefspec configures the fetch refspec to fetch all remote branches
func (c *Client) ConfigureFetchRefspec() error {
	// Set remote.origin.fetch to fetch all branches
	err := c.SetConfig("remote.origin.fetch", "+refs/heads/*:refs/remotes/origin/*")
	if err != nil {
		return fmt.Errorf("failed to configure fetch refspec: %w", err)
	}

	return nil
}

// ConfigureWorktreeSettings configures git settings for worktree management
func (c *Client) ConfigureWorktreeSettings() error {
	settings := map[string]string{
		"push.default":           "current",
		"branch.autosetupmerge":  "always",
		"branch.autosetuprebase": "always",
	}

	for key, value := range settings {
		if err := c.SetConfig(key, value); err != nil {
			return err
		}
	}

	return nil
}
