package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestGetBinaryPath(t *testing.T) {
	installDir := t.TempDir()
	os.Setenv("GIT_WORKTREE_MANAGER_HOME", installDir)
	defer os.Unsetenv("GIT_WORKTREE_MANAGER_HOME")

	binaryName := "gwtm"
	if runtime.GOOS == "windows" {
		binaryName = "gwtm.exe"
	}
	expected := filepath.Join(installDir, binaryName)

	got := GetBinaryPath()
	if got != expected {
		t.Errorf("GetBinaryPath() = %v, want %v", got, expected)
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
			parts:    []string{"/usr", "local", "bin", "gwtm"},
			expected: filepath.Join("/usr", "local", "bin", "gwtm"),
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
