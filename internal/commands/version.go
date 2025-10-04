package commands

import (
	"fmt"
	"io"
	"net/http"
	"strings"

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
	fmt.Printf("git-worktree-manager version %s\n", currentVersion)

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
		ui.PrintStatus("⬇️", fmt.Sprintf("Run 'git-worktree-manager upgrade' to upgrade to version %s.", latestVersion))
	} else {
		fmt.Printf("%s <= %s\n", latestVersion, currentVersion)
		ui.PrintStatus("✅", "You already have the latest version.")
	}
}

func fetchLatestVersion() (string, error) {
	url := "https://raw.githubusercontent.com/lucasmodrich/git-worktree-manager/refs/heads/main/VERSION"

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch latest version: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	return strings.TrimSpace(string(body)), nil
}
