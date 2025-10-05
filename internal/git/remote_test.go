package git

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupRemoteTestRepo(t *testing.T) (*Client, string, string) {
	tmpDir := t.TempDir()

	// Create a "remote" bare repository
	remoteDir := filepath.Join(tmpDir, "remote.git")
	os.MkdirAll(remoteDir, 0755)
	remoteClient := NewClient(remoteDir)
	remoteClient.ExecGit("init", "--bare")

	// Create a local repository
	localDir := filepath.Join(tmpDir, "local")
	os.MkdirAll(localDir, 0755)
	localClient := NewClient(localDir)
	localClient.ExecGit("init")
	localClient.ExecGit("config", "user.name", "Test User")
	localClient.ExecGit("config", "user.email", "test@example.com")

	// Create initial commit
	testFile := filepath.Join(localDir, "README.md")
	os.WriteFile(testFile, []byte("# Test\n"), 0644)
	localClient.ExecGit("add", "README.md")
	localClient.ExecGit("commit", "-m", "Initial commit")

	// Detect which default branch was created (main or master)
	output, _, _ := localClient.ExecGit("branch", "--show-current")
	defaultBranch := strings.TrimSpace(output)
	if defaultBranch == "" {
		// Fallback for older git versions
		output, _, _ = localClient.ExecGit("rev-parse", "--abbrev-ref", "HEAD")
		defaultBranch = strings.TrimSpace(output)
	}

	// Add remote
	localClient.ExecGit("remote", "add", "origin", remoteDir)

	return localClient, localDir, defaultBranch
}

func TestClone(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a source repository to clone from
	sourceDir := filepath.Join(tmpDir, "source")
	os.MkdirAll(sourceDir, 0755)
	sourceClient := NewClient(sourceDir)
	sourceClient.ExecGit("init")
	sourceClient.ExecGit("config", "user.name", "Test User")
	sourceClient.ExecGit("config", "user.email", "test@example.com")

	testFile := filepath.Join(sourceDir, "README.md")
	os.WriteFile(testFile, []byte("# Source\n"), 0644)
	sourceClient.ExecGit("add", "README.md")
	sourceClient.ExecGit("commit", "-m", "Initial commit")

	tests := []struct {
		name    string
		url     string
		target  string
		bare    bool
		wantErr bool
	}{
		{
			name:    "clone non-bare repository",
			url:     sourceDir,
			target:  filepath.Join(tmpDir, "clone1"),
			bare:    false,
			wantErr: false,
		},
		{
			name:    "clone bare repository",
			url:     sourceDir,
			target:  filepath.Join(tmpDir, "clone2"),
			bare:    true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient("")
			err := client.Clone(tt.url, tt.target, tt.bare)
			if (err != nil) != tt.wantErr {
				t.Errorf("Clone() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				// Verify clone succeeded
				if tt.bare {
					configFile := filepath.Join(tt.target, "config")
					if _, err := os.Stat(configFile); os.IsNotExist(err) {
						t.Error("Clone() bare repository does not have config file")
					}
				} else {
					gitDir := filepath.Join(tt.target, ".git")
					if _, err := os.Stat(gitDir); os.IsNotExist(err) {
						t.Error("Clone() non-bare repository does not have .git directory")
					}
				}
			}
		})
	}
}

func TestFetch(t *testing.T) {
	client, _, defaultBranch := setupRemoteTestRepo(t)

	// Push to remote first
	client.ExecGit("push", "-u", "origin", defaultBranch)

	tests := []struct {
		name    string
		all     bool
		prune   bool
		wantErr bool
	}{
		{
			name:    "fetch all",
			all:     true,
			prune:   false,
			wantErr: false,
		},
		{
			name:    "fetch with prune",
			all:     true,
			prune:   true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Fetch(tt.all, tt.prune)
			if (err != nil) != tt.wantErr {
				t.Errorf("Fetch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPush(t *testing.T) {
	client, _, defaultBranch := setupRemoteTestRepo(t)

	tests := []struct {
		name        string
		branch      string
		setUpstream bool
		wantErr     bool
	}{
		{
			name:        "push default branch with upstream",
			branch:      defaultBranch,
			setUpstream: true,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Push(tt.branch, tt.setUpstream)
			if (err != nil) != tt.wantErr {
				t.Errorf("Push() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDetectDefaultBranch(t *testing.T) {
	client, localDir, defaultBranch := setupRemoteTestRepo(t)

	// Push to create remote tracking
	client.Push(defaultBranch, true)

	// Detect default branch
	branch, err := client.DetectDefaultBranch()
	if err != nil {
		t.Fatalf("DetectDefaultBranch() error = %v", err)
	}

	// Should be either "main" or "master"
	if branch != "main" && branch != "master" {
		t.Errorf("DetectDefaultBranch() = %s, want 'main' or 'master'", branch)
	}

	// Test with a fresh clone that has remote tracking
	cloneDir := filepath.Join(filepath.Dir(localDir), "clone-for-detect")
	cloneClient := NewClient("")
	cloneClient.Clone(client.WorkDir, cloneDir, false)

	cloneClient.WorkDir = cloneDir
	branch2, err2 := cloneClient.DetectDefaultBranch()
	if err2 != nil {
		t.Fatalf("DetectDefaultBranch() on clone error = %v", err2)
	}

	if !strings.Contains(branch2, "main") && !strings.Contains(branch2, "master") {
		t.Errorf("DetectDefaultBranch() on clone = %s, want to contain 'main' or 'master'", branch2)
	}
}
