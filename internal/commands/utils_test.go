package commands

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindWorktreeRoot(t *testing.T) {
	// Create a temp dir structure simulating a worktree-managed repo:
	//   tmpDir/.git  (file, not directory)
	//   tmpDir/branch/subdir/
	// findWorktreeRoot called from branch/subdir should return tmpDir.
	tmpDir := t.TempDir()

	// Create .git as a file (worktree pointer)
	gitFile := filepath.Join(tmpDir, ".git")
	if err := os.WriteFile(gitFile, []byte("gitdir: ./.bare"), 0644); err != nil {
		t.Fatalf("failed to create .git file: %v", err)
	}

	// Create a nested subdirectory
	subDir := filepath.Join(tmpDir, "branch", "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("failed to create subdirectory: %v", err)
	}

	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}
	defer os.Chdir(origDir) //nolint:errcheck

	if err := os.Chdir(subDir); err != nil {
		t.Fatalf("failed to chdir to subDir: %v", err)
	}

	root, err := findWorktreeRoot()
	if err != nil {
		t.Fatalf("findWorktreeRoot() error = %v", err)
	}

	// Use os.SameFile to handle potential symlinks in temp paths
	rootInfo, err := os.Stat(root)
	if err != nil {
		t.Fatalf("failed to stat root %q: %v", root, err)
	}
	tmpDirInfo, err := os.Stat(tmpDir)
	if err != nil {
		t.Fatalf("failed to stat tmpDir %q: %v", tmpDir, err)
	}
	if !os.SameFile(rootInfo, tmpDirInfo) {
		t.Errorf("findWorktreeRoot() = %q, want %q", root, tmpDir)
	}
}

func TestFindWorktreeRoot_NotFound(t *testing.T) {
	tmpDir := t.TempDir()

	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}
	defer os.Chdir(origDir) //nolint:errcheck

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}

	_, err = findWorktreeRoot()
	if err == nil {
		t.Error("findWorktreeRoot() expected error when no .git file present, got nil")
	}
}

func TestFindWorktreeRoot_GitDirectory(t *testing.T) {
	// A .git directory (standard git repo) must NOT satisfy findWorktreeRoot
	tmpDir := t.TempDir()

	gitDir := filepath.Join(tmpDir, ".git")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatalf("failed to create .git directory: %v", err)
	}

	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}
	defer os.Chdir(origDir) //nolint:errcheck

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}

	_, err = findWorktreeRoot()
	if err == nil {
		t.Error("findWorktreeRoot() should not match a .git directory, expected error")
	}
}
