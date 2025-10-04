package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Client is a wrapper for executing git commands
type Client struct {
	WorkDir string // Working directory for git commands
	DryRun  bool   // If true, log commands without executing
}

// NewClient creates a new git client
func NewClient(workDir string) *Client {
	return &Client{
		WorkDir: workDir,
		DryRun:  false,
	}
}

// ExecGit executes a git command and returns stdout, stderr, and error
func (c *Client) ExecGit(args ...string) (stdout, stderr string, err error) {
	if c.DryRun {
		// In dry-run mode, just log the command and return success
		cmdStr := "git " + strings.Join(args, " ")
		return fmt.Sprintf("[DRY-RUN] Would execute: %s", cmdStr), "", nil
	}

	cmd := exec.Command("git", args...)
	if c.WorkDir != "" {
		cmd.Dir = c.WorkDir
	}

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err = cmd.Run()

	stdout = outBuf.String()
	stderr = errBuf.String()

	if err != nil {
		// Enhance error with stderr output
		if stderr != "" {
			err = fmt.Errorf("git command failed: %w\nstderr: %s", err, stderr)
		} else {
			err = fmt.Errorf("git command failed: %w", err)
		}
	}

	return stdout, stderr, err
}
