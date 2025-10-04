package config

import (
	"os"
	"path/filepath"
)

// GetInstallDir returns the installation directory for git-worktree-manager
// Respects GIT_WORKTREE_MANAGER_HOME environment variable, defaults to $HOME/.git-worktree-manager
func GetInstallDir() string {
	if customDir := os.Getenv("GIT_WORKTREE_MANAGER_HOME"); customDir != "" {
		return customDir
	}

	home := os.Getenv("HOME")
	if home == "" {
		// Fallback for Windows
		home = os.Getenv("USERPROFILE")
	}

	return filepath.Join(home, ".git-worktree-manager")
}
