package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/lucasmodrich/git-worktree-manager/internal/ui"
	"github.com/lucasmodrich/git-worktree-manager/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version and check for updates",
	Long:  "Display current version and check GitHub for newer versions",
	Run:   runVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func runVersion(cmd *cobra.Command, args []string) {
	// Display current version
	currentVersion := GetVersion()
	fmt.Printf("gwtm version %s\n", currentVersion)

	// Check for newer version on GitHub
	ui.PrintStatus("🔍", "Checking for newer version on GitHub...")

	latestVersion, err := fetchLatestVersion()
	if err != nil {
		ui.PrintError(err, "Could not check for updates. Try again later.")
		return
	}

	ui.PrintStatus("🔢", fmt.Sprintf("Local version: %s", currentVersion))
	ui.PrintStatus("🌐", fmt.Sprintf("Remote version: %s", latestVersion))

	// Compare versions
	currentVer, err1 := version.ParseVersion(currentVersion)
	latestVer, err2 := version.ParseVersion(latestVersion)

	if err1 != nil || err2 != nil {
		ui.PrintStatus("⚠️", "Unable to compare versions")
		return
	}

	if latestVer.GreaterThan(currentVer) {
		fmt.Printf("%s > %s\n", latestVersion, currentVersion)
		ui.PrintStatus("⬇️", fmt.Sprintf("Run 'gwtm upgrade' to upgrade to version %s.", latestVersion))
	} else {
		fmt.Printf("%s <= %s\n", latestVersion, currentVersion)
		ui.PrintStatus("✅", "You already have the latest version.")
	}
}

// fetchLatestVersion queries the GitHub Releases API for the latest published release version.
func fetchLatestVersion() (string, error) {
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
