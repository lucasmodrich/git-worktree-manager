package commands

import (
	"testing"
)

func TestParseRepoSpec(t *testing.T) {
	tests := []struct {
		name         string
		spec         string
		wantURL      string
		wantRepoName string
		wantErr      bool
	}{
		{
			name:         "org/repo shorthand",
			spec:         "acme/webapp",
			wantURL:      "git@github.com:acme/webapp.git",
			wantRepoName: "webapp",
		},
		{
			name:         "org/repo with dots",
			spec:         "my.org/tool.name",
			wantURL:      "git@github.com:my.org/tool.name.git",
			wantRepoName: "tool.name",
		},
		{
			name:         "GitHub SSH URL",
			spec:         "git@github.com:acme/webapp.git",
			wantURL:      "git@github.com:acme/webapp.git",
			wantRepoName: "webapp",
		},
		{
			name:         "generic SSH URL non-GitHub host",
			spec:         "git@gitlab.com:acme/webapp.git",
			wantURL:      "git@gitlab.com:acme/webapp.git",
			wantRepoName: "webapp",
		},
		{
			name:         "HTTPS URL with .git suffix",
			spec:         "https://github.com/acme/webapp.git",
			wantURL:      "https://github.com/acme/webapp.git",
			wantRepoName: "webapp",
		},
		{
			name:         "HTTPS URL without .git suffix",
			spec:         "https://github.com/acme/webapp",
			wantURL:      "https://github.com/acme/webapp",
			wantRepoName: "webapp",
		},
		{
			name:         "HTTPS URL non-GitHub host",
			spec:         "https://gitlab.com/org/repo.git",
			wantURL:      "https://gitlab.com/org/repo.git",
			wantRepoName: "repo",
		},
		{
			name:    "plain name with no slash is invalid",
			spec:    "not-a-repo",
			wantErr: true,
		},
		{
			name:    "too many slashes is invalid",
			spec:    "a/b/c",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotURL, gotRepoName, err := parseRepoSpec(tt.spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRepoSpec(%q) error = %v, wantErr %v", tt.spec, err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if gotURL != tt.wantURL {
				t.Errorf("parseRepoSpec(%q) url = %q, want %q", tt.spec, gotURL, tt.wantURL)
			}
			if gotRepoName != tt.wantRepoName {
				t.Errorf("parseRepoSpec(%q) repoName = %q, want %q", tt.spec, gotRepoName, tt.wantRepoName)
			}
		})
	}
}
