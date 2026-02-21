package version

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// FetchLatestVersion queries the GitHub Releases API for the latest published release version.
// Returns the version number without a leading 'v' (e.g. "1.4.0").
func FetchLatestVersion() (string, error) {
	const apiURL = "https://api.github.com/repos/lucasmodrich/git-worktree-manager/releases/latest"

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "gwtm")
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch latest release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("failed to parse GitHub response: %w", err)
	}

	if release.TagName == "" {
		return "", fmt.Errorf("no releases found")
	}

	// Strip leading 'v' so callers work with bare version numbers (e.g. "1.3.0")
	return strings.TrimPrefix(release.TagName, "v"), nil
}
