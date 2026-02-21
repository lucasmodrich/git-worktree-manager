package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetInstallDir(t *testing.T) {
	tests := []struct {
		name    string
		envVar  string
		envVal  string
		wantDir string
	}{
		{
			name:    "with GIT_WORKTREE_MANAGER_HOME set",
			envVar:  "GIT_WORKTREE_MANAGER_HOME",
			envVal:  "/custom/path",
			wantDir: "/custom/path",
		},
		{
			name:    "without env var uses HOME default",
			envVar:  "",
			envVal:  "",
			wantDir: "", // Will use $HOME/.git-worktree-manager
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original env
			origEnv := os.Getenv("GIT_WORKTREE_MANAGER_HOME")
			defer os.Setenv("GIT_WORKTREE_MANAGER_HOME", origEnv)

			// Set test env
			if tt.envVar != "" {
				os.Setenv(tt.envVar, tt.envVal)
			} else {
				os.Unsetenv("GIT_WORKTREE_MANAGER_HOME")
			}

			got := GetInstallDir()

			// If no custom path, verify it uses HOME/.git-worktree-manager
			if tt.wantDir == "" {
				home := os.Getenv("HOME")
				if home == "" {
					home = os.Getenv("USERPROFILE")
				}
				if home == "" {
					t.Skip("HOME not set, skipping default path test")
				}
				expected := filepath.Join(home, ".git-worktree-manager")
				if got != expected {
					t.Errorf("GetInstallDir() = %v, want %v", got, expected)
				}
			} else {
				if got != tt.wantDir {
					t.Errorf("GetInstallDir() = %v, want %v", got, tt.wantDir)
				}
			}
		})
	}
}
