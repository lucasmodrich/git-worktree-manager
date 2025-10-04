package git

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupTestRepo(t *testing.T) (*Client, string) {
	tmpDir := t.TempDir()

	// Create a bare repo
	bareDir := filepath.Join(tmpDir, ".bare")
	os.MkdirAll(bareDir, 0755)

	client := NewClient(bareDir)
	client.ExecGit("init", "--bare")

	// Set up a main worktree with an initial commit
	mainDir := filepath.Join(tmpDir, "main")
	os.MkdirAll(mainDir, 0755)

	mainClient := NewClient(mainDir)
	mainClient.ExecGit("init")
	mainClient.ExecGit("config", "user.name", "Test User")
	mainClient.ExecGit("config", "user.email", "test@example.com")

	// Create initial commit
	testFile := filepath.Join(mainDir, "README.md")
	os.WriteFile(testFile, []byte("# Test Repo\n"), 0644)
	mainClient.ExecGit("add", "README.md")
	mainClient.ExecGit("commit", "-m", "Initial commit")

	return NewClient(tmpDir), tmpDir
}

func TestWorktreeAdd(t *testing.T) {
	client, tmpDir := setupTestRepo(t)

	tests := []struct {
		name     string
		path     string
		branch   string
		track    bool
		wantErr  bool
	}{
		{
			name:     "add new worktree",
			path:     filepath.Join(tmpDir, "feature"),
			branch:   "feature/test",
			track:    true,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.WorktreeAdd(tt.path, tt.branch, tt.track)
			if (err != nil) != tt.wantErr {
				t.Errorf("WorktreeAdd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWorktreeList(t *testing.T) {
	client, tmpDir := setupTestRepo(t)

	// Add a worktree first
	featurePath := filepath.Join(tmpDir, "feature")
	client.WorktreeAdd(featurePath, "feature/test", true)

	worktrees, err := client.WorktreeList()
	if err != nil {
		t.Fatalf("WorktreeList() error = %v", err)
	}

	if len(worktrees) == 0 {
		t.Error("WorktreeList() returned empty list, expected at least one worktree")
	}

	// Check that the worktree paths are present
	var foundFeature bool
	for _, wt := range worktrees {
		if strings.Contains(wt, "feature") {
			foundFeature = true
		}
	}

	if !foundFeature {
		t.Error("WorktreeList() did not include the feature worktree we added")
	}
}

func TestWorktreeRemove(t *testing.T) {
	client, tmpDir := setupTestRepo(t)

	// Add a worktree first
	featurePath := filepath.Join(tmpDir, "feature")
	client.WorktreeAdd(featurePath, "feature/test", true)

	// Now remove it
	err := client.WorktreeRemove(featurePath)
	if err != nil {
		t.Errorf("WorktreeRemove() error = %v", err)
	}

	// Verify it's removed
	worktrees, _ := client.WorktreeList()
	for _, wt := range worktrees {
		if strings.Contains(wt, "feature") {
			t.Error("WorktreeRemove() did not remove the worktree")
		}
	}
}

func TestWorktreePrune(t *testing.T) {
	client, _ := setupTestRepo(t)

	// Prune should not error even if nothing to prune
	err := client.WorktreePrune()
	if err != nil {
		t.Errorf("WorktreePrune() error = %v", err)
	}
}
