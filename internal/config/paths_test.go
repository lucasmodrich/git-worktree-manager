package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetBinaryPath(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() string
		cleanup func()
	}{
		{
			name: "binary path in install directory",
			setup: func() string {
				// Set a custom install directory
				os.Setenv("GIT_WORKTREE_MANAGER_HOME", "/tmp/test-install")
				return "/tmp/test-install/git-worktree-manager"
			},
			cleanup: func() {
				os.Unsetenv("GIT_WORKTREE_MANAGER_HOME")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := tt.setup()
			defer tt.cleanup()

			got := GetBinaryPath()
			if got != expected {
				t.Errorf("GetBinaryPath() = %v, want %v", got, expected)
			}
		})
	}
}

func TestPathJoinCrossPlatform(t *testing.T) {
	tests := []struct {
		name     string
		parts    []string
		expected string
	}{
		{
			name:     "simple path join",
			parts:    []string{"/home", "user", ".git-worktree-manager"},
			expected: filepath.Join("/home", "user", ".git-worktree-manager"),
		},
		{
			name:     "path with binary name",
			parts:    []string{"/usr", "local", "bin", "git-worktree-manager"},
			expected: filepath.Join("/usr", "local", "bin", "git-worktree-manager"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filepath.Join(tt.parts...)
			if got != tt.expected {
				t.Errorf("filepath.Join() = %v, want %v", got, tt.expected)
			}
		})
	}
}
