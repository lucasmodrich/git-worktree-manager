package git

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupTestRepo(t *testing.T) (*Client, string, string) {
	tmpDir := t.TempDir()

	// Step 1: Create a temporary normal repo to get initial content
	setupDir := filepath.Join(tmpDir, "setup")
	os.MkdirAll(setupDir, 0755)

	setupClient := NewClient(setupDir)
	setupClient.ExecGit("init")
	setupClient.ExecGit("config", "user.name", "Test User")
	setupClient.ExecGit("config", "user.email", "test@example.com")

	// Create initial commit
	testFile := filepath.Join(setupDir, "README.md")
	os.WriteFile(testFile, []byte("# Test Repo\n"), 0644)
	setupClient.ExecGit("add", "README.md")
	setupClient.ExecGit("commit", "-m", "Initial commit")

	// Detect default branch
	output, _, _ := setupClient.ExecGit("branch", "--show-current")
	defaultBranch := strings.TrimSpace(output)
	if defaultBranch == "" {
		output, _, _ = setupClient.ExecGit("rev-parse", "--abbrev-ref", "HEAD")
		defaultBranch = strings.TrimSpace(output)
	}

	// Step 2: Create bare repo from the setup repo
	bareDir := filepath.Join(tmpDir, ".bare")
	setupClient.ExecGit("clone", "--bare", setupDir, bareDir)

	// Step 3: Create first worktree from bare repo
	mainDir := filepath.Join(tmpDir, defaultBranch)
	bareClient := NewClient(bareDir)
	bareClient.ExecGit("worktree", "add", mainDir, defaultBranch)

	// Return client pointing to the project root (tmpDir) which contains .bare/
	return NewClient(tmpDir), tmpDir, defaultBranch
}

func TestWorktreeAdd(t *testing.T) {
	_, tmpDir, defaultBranch := setupTestRepo(t)

	// Use bare repo client for worktree commands
	bareDir := filepath.Join(tmpDir, ".bare")
	bareClient := NewClient(bareDir)

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
			track:    false, // Match production usage: branch pre-created, track=false
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Pre-create the branch (matches production workflow in branch.go)
			bareClient.ExecGit("branch", tt.branch, defaultBranch)

			err := bareClient.WorktreeAdd(tt.path, tt.branch, tt.track)
			if (err != nil) != tt.wantErr {
				t.Errorf("WorktreeAdd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWorktreeList(t *testing.T) {
	_, tmpDir, defaultBranch := setupTestRepo(t)

	// Use bare repo client
	bareDir := filepath.Join(tmpDir, ".bare")
	bareClient := NewClient(bareDir)

	// Add a worktree first (match production workflow: create branch, then worktree)
	featurePath := filepath.Join(tmpDir, "feature")
	bareClient.ExecGit("branch", "feature/test", defaultBranch)
	bareClient.WorktreeAdd(featurePath, "feature/test", false)

	worktrees, err := bareClient.WorktreeList()
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
	_, tmpDir, defaultBranch := setupTestRepo(t)

	// Use bare repo client
	bareDir := filepath.Join(tmpDir, ".bare")
	bareClient := NewClient(bareDir)

	// Add a worktree first (match production workflow: create branch, then worktree)
	featurePath := filepath.Join(tmpDir, "feature")
	bareClient.ExecGit("branch", "feature/test", defaultBranch)
	bareClient.WorktreeAdd(featurePath, "feature/test", false)

	// Now remove it
	err := bareClient.WorktreeRemove(featurePath)
	if err != nil {
		t.Errorf("WorktreeRemove() error = %v", err)
	}

	// Verify it's removed
	worktrees, _ := bareClient.WorktreeList()
	for _, wt := range worktrees {
		if strings.Contains(wt, "feature") {
			t.Error("WorktreeRemove() did not remove the worktree")
		}
	}
}

func TestWorktreePrune(t *testing.T) {
	_, tmpDir, _ := setupTestRepo(t)

	// Use bare repo client
	bareDir := filepath.Join(tmpDir, ".bare")
	bareClient := NewClient(bareDir)

	// Prune should not error even if nothing to prune
	err := bareClient.WorktreePrune()
	if err != nil {
		t.Errorf("WorktreePrune() error = %v", err)
	}
}
