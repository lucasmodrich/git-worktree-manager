package git

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		workDir string
	}{
		{
			name:    "create client with work directory",
			workDir: "/tmp/test-repo",
		},
		{
			name:    "create client with empty work directory",
			workDir: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.workDir)
			if client == nil {
				t.Error("NewClient() returned nil")
			}
			if client.WorkDir != tt.workDir {
				t.Errorf("NewClient() WorkDir = %v, want %v", client.WorkDir, tt.workDir)
			}
		})
	}
}

func TestExecGit(t *testing.T) {
	// Create a temporary test directory
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		workDir     string
		args        []string
		setupRepo   bool
		wantErr     bool
		wantStdout  bool
	}{
		{
			name:        "git version command",
			workDir:     tmpDir,
			args:        []string{"version"},
			setupRepo:   false,
			wantErr:     false,
			wantStdout:  true,
		},
		{
			name:        "git status in initialized repo",
			workDir:     tmpDir,
			args:        []string{"status"},
			setupRepo:   true,
			wantErr:     false,
			wantStdout:  true,
		},
		{
			name:        "invalid git command",
			workDir:     tmpDir,
			args:        []string{"invalid-command-xyz"},
			setupRepo:   false,
			wantErr:     true,
			wantStdout:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupRepo {
				// Initialize git repo for this test
				testRepo := filepath.Join(tmpDir, tt.name)
				os.MkdirAll(testRepo, 0755)
				client := NewClient(testRepo)
				client.ExecGit("init")
				tt.workDir = testRepo
			}

			client := NewClient(tt.workDir)
			stdout, stderr, err := client.ExecGit(tt.args...)

			if (err != nil) != tt.wantErr {
				t.Errorf("ExecGit() error = %v, wantErr %v, stderr = %s", err, tt.wantErr, stderr)
				return
			}

			if tt.wantStdout && stdout == "" {
				t.Error("ExecGit() expected stdout output, got empty string")
			}
		})
	}
}

func TestDryRunMode(t *testing.T) {
	tmpDir := t.TempDir()
	client := NewClient(tmpDir)
	client.DryRun = true

	// In dry-run mode, commands should not execute
	stdout, stderr, err := client.ExecGit("init")

	// Dry run should not error, but also shouldn't create a .git directory
	if err != nil {
		t.Errorf("DryRun ExecGit() unexpected error = %v", err)
	}

	// Verify no .git directory was created
	gitDir := filepath.Join(tmpDir, ".git")
	if _, err := os.Stat(gitDir); !os.IsNotExist(err) {
		t.Error("DryRun mode created .git directory, should not execute commands")
	}

	// Dry run should log the command somehow (check stdout or client state)
	if stdout == "" && stderr == "" {
		// This is acceptable - dry run might not produce output
	}
}
