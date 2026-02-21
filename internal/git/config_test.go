package git

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupConfigTestRepo(t *testing.T) (*Client, string) {
	tmpDir := t.TempDir()

	client := NewClient(tmpDir)
	client.ExecGit("init")

	return client, tmpDir
}

func TestSetConfig(t *testing.T) {
	client, _ := setupConfigTestRepo(t)

	tests := []struct {
		name    string
		key     string
		value   string
		wantErr bool
	}{
		{
			name:    "set user name",
			key:     "user.name",
			value:   "Test User",
			wantErr: false,
		},
		{
			name:    "set user email",
			key:     "user.email",
			value:   "test@example.com",
			wantErr: false,
		},
		{
			name:    "set push default",
			key:     "push.default",
			value:   "current",
			wantErr: false,
		},
		{
			name:    "set branch autosetupmerge",
			key:     "branch.autosetupmerge",
			value:   "always",
			wantErr: false,
		},
		{
			name:    "set branch autosetuprebase",
			key:     "branch.autosetuprebase",
			value:   "always",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.SetConfig(tt.key, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetConfig() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				// Verify the config was set
				stdout, _, _ := client.ExecGit("config", "--get", tt.key)
				got := strings.TrimSpace(stdout)
				if got != tt.value {
					t.Errorf("SetConfig() value = %s, want %s", got, tt.value)
				}
			}
		})
	}
}

func TestConfigureFetchRefspec(t *testing.T) {
	client, tmpDir := setupConfigTestRepo(t)

	// Add a remote first
	remoteDir := filepath.Join(filepath.Dir(tmpDir), "remote.git")
	os.MkdirAll(remoteDir, 0755)
	remoteClient := NewClient(remoteDir)
	remoteClient.ExecGit("init", "--bare")

	client.ExecGit("remote", "add", "origin", remoteDir)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "configure fetch refspec",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.ConfigureFetchRefspec()
			if (err != nil) != tt.wantErr {
				t.Errorf("ConfigureFetchRefspec() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				// Verify the refspec was configured
				stdout, _, _ := client.ExecGit("config", "--get", "remote.origin.fetch")
				got := strings.TrimSpace(stdout)
				if !strings.Contains(got, "refs/heads/*") {
					t.Errorf("ConfigureFetchRefspec() did not set correct refspec, got %s", got)
				}
			}
		})
	}
}

func TestConfigureWorktreeSettings(t *testing.T) {
	client, _ := setupConfigTestRepo(t)

	// This is a convenience function that should set multiple config values
	err := client.ConfigureWorktreeSettings()
	if err != nil {
		t.Fatalf("ConfigureWorktreeSettings() error = %v", err)
	}

	// Verify all required settings were configured
	requiredConfigs := map[string]string{
		"push.default":              "current",
		"branch.autosetupmerge":     "always",
		"branch.autosetuprebase":    "always",
	}

	for key, expectedValue := range requiredConfigs {
		stdout, _, _ := client.ExecGit("config", "--get", key)
		got := strings.TrimSpace(stdout)
		if got != expectedValue {
			t.Errorf("ConfigureWorktreeSettings() %s = %s, want %s", key, got, expectedValue)
		}
	}
}
