package config

import (
	"path/filepath"
)

// GetBinaryPath returns the full path to the git-worktree-manager binary
func GetBinaryPath() string {
	installDir := GetInstallDir()
	return filepath.Join(installDir, "git-worktree-manager")
}
