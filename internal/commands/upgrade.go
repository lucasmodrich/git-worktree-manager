package commands

import (
	"fmt"

	"github.com/lucasmodrich/git-worktree-manager/internal/ui"
	"github.com/lucasmodrich/git-worktree-manager/internal/version"
	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade to latest version",
	Long:  "Download and install the latest version of git-worktree-manager from GitHub",
	Run:   runUpgrade,
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}

func runUpgrade(cmd *cobra.Command, args []string) {
	currentVersion := GetVersion()

	// Check for newer version on GitHub
	ui.PrintStatus("🔍", "Checking for newer version on GitHub...")

	latestVersion, err := fetchLatestVersion()
	if err != nil {
		ui.PrintError(err, "Could not check for updates. Try again later.")
		return
	}

	// Parse and compare versions
	currentVer, err1 := version.ParseVersion(currentVersion)
	latestVer, err2 := version.ParseVersion(latestVersion)

	if err1 != nil || err2 != nil {
		ui.PrintError(fmt.Errorf("unable to parse versions"), "Version format issue")
		return
	}

	if !latestVer.GreaterThan(currentVer) {
		ui.PrintStatus("✅", "You already have the latest version.")
		return
	}

	// Perform upgrade
	ui.PrintStatus("⬇️", fmt.Sprintf("Upgrading to version %s...", latestVersion))

	if err := version.UpgradeToLatest(currentVersion, latestVersion); err != nil {
		ui.PrintError(err, "Upgrade failed. Try again or download manually from GitHub releases")
		return
	}

	ui.PrintStatus("✓", "Binary downloaded")
	ui.PrintStatus("✓", "Checksum verified")
	ui.PrintStatus("✓", "README.md downloaded")
	ui.PrintStatus("✓", "VERSION downloaded")
	ui.PrintStatus("✓", "LICENSE downloaded")
	ui.PrintStatus("✅", fmt.Sprintf("Upgrade complete. Now running version %s.", latestVersion))
}
