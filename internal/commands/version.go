package commands

import (
	"fmt"

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
	// Display current version and build info
	currentVersion := GetVersion()
	commit := GetCommit()
	date := GetDate()

	if commit != "" && commit != "none" {
		fmt.Printf("gwtm version %s (%s, %s)\n", currentVersion, commit, date)
	} else {
		fmt.Printf("gwtm version %s\n", currentVersion)
	}

	// Check for newer version on GitHub
	ui.PrintStatus("🔍", "Checking for newer version on GitHub...")

	latestVersion, err := version.FetchLatestVersion()
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
		if currentVersion == "dev" {
			ui.PrintStatus("ℹ️", "Dev build — version comparison not available.")
		} else {
			ui.PrintStatus("⚠️", "Unable to compare versions.")
		}
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
