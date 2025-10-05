package git

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupBranchTestRepo(t *testing.T) (*Client, string, string) {
	tmpDir := t.TempDir()

	client := NewClient(tmpDir)
	client.ExecGit("init")
	client.ExecGit("config", "user.name", "Test User")
	client.ExecGit("config", "user.email", "test@example.com")

	// Create initial commit
	testFile := filepath.Join(tmpDir, "README.md")
	os.WriteFile(testFile, []byte("# Test\n"), 0644)
	client.ExecGit("add", "README.md")
	client.ExecGit("commit", "-m", "Initial commit")

	// Detect which default branch was created (main or master)
	output, _, _ := client.ExecGit("branch", "--show-current")
	defaultBranch := strings.TrimSpace(output)
	if defaultBranch == "" {
		// Fallback for older git versions
		output, _, _ = client.ExecGit("rev-parse", "--abbrev-ref", "HEAD")
		defaultBranch = strings.TrimSpace(output)
	}

	return client, tmpDir, defaultBranch
}

func TestBranchExists(t *testing.T) {
	client, _, defaultBranch := setupBranchTestRepo(t)

	tests := []struct {
		name     string
		branch   string
		remote   bool
		setup    func()
		want     bool
	}{
		{
			name:   "default branch exists locally",
			branch: defaultBranch,
			remote: false,
			setup:  func() {},
			want:   true,
		},
		{
			name:   "non-existent branch",
			branch: "feature/does-not-exist",
			remote: false,
			setup:  func() {},
			want:   false,
		},
		{
			name:   "created branch exists",
			branch: "feature/test",
			remote: false,
			setup: func() {
				client.CreateBranch("feature/test", "")
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			got := client.BranchExists(tt.branch, tt.remote)
			if got != tt.want {
				t.Errorf("BranchExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateBranch(t *testing.T) {
	client, _, defaultBranch := setupBranchTestRepo(t)

	tests := []struct {
		name       string
		branchName string
		baseBranch string
		wantErr    bool
	}{
		{
			name:       "create branch from default",
			branchName: "feature/new",
			baseBranch: defaultBranch,
			wantErr:    false,
		},
		{
			name:       "create branch with empty base (use current)",
			branchName: "feature/current",
			baseBranch: "",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.CreateBranch(tt.branchName, tt.baseBranch)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateBranch() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				// Verify branch was created
				if !client.BranchExists(tt.branchName, false) {
					t.Errorf("CreateBranch() branch %s was not created", tt.branchName)
				}
			}
		})
	}
}

func TestDeleteBranch(t *testing.T) {
	client, _, defaultBranch := setupBranchTestRepo(t)

	// Create a branch first
	client.CreateBranch("feature/to-delete", defaultBranch)

	tests := []struct {
		name    string
		branch  string
		force   bool
		wantErr bool
	}{
		{
			name:    "delete branch normally",
			branch:  "feature/to-delete",
			force:   false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.DeleteBranch(tt.branch, tt.force)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteBranch() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				// Verify branch was deleted
				if client.BranchExists(tt.branch, false) {
					t.Errorf("DeleteBranch() branch %s still exists", tt.branch)
				}
			}
		})
	}
}

func TestDeleteRemoteBranch(t *testing.T) {
	client, _, _ := setupBranchTestRepo(t)

	tests := []struct {
		name    string
		branch  string
		wantErr bool
	}{
		{
			name:    "delete remote branch (will fail without remote)",
			branch:  "feature/test",
			wantErr: true, // Expected to fail since we don't have a real remote
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.DeleteRemoteBranch(tt.branch)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteRemoteBranch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
