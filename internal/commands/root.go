package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Global flags
	dryRun bool

	// Build info — set via SetBuildInfo from main
	appVersion string
	appCommit  string
	appDate    string
)

var rootCmd = &cobra.Command{
	Use:   "gwtm",
	Short: "Git worktree manager - Simplify git worktree workflows",
	Long: `🛠 Git Worktree Manager — A tool to simplify git worktree management

Supports:
  - Full repository setup from GitHub
  - Branch and worktree creation
  - Worktree listing and removal
  - Version management and self-upgrade
  - Dry-run mode for all destructive operations`,
}

func init() {
	// Global persistent flags
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Preview actions without executing")
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

// SetBuildInfo sets the version, commit hash, and build date (called from main).
// Also populates rootCmd.Version so Cobra's built-in --version flag works.
func SetBuildInfo(version, commit, date string) {
	appVersion = version
	appCommit = commit
	appDate = date
	if commit != "" && commit != "none" {
		rootCmd.Version = fmt.Sprintf("%s (%s, %s)", version, commit, date)
	} else {
		rootCmd.Version = version
	}
}

// GetVersion returns the semver version string
func GetVersion() string {
	return appVersion
}

// GetCommit returns the git commit hash
func GetCommit() string {
	return appCommit
}

// GetDate returns the build date
func GetDate() string {
	return appDate
}

// GetDryRun returns the dry-run flag value
func GetDryRun() bool {
	return dryRun
}
