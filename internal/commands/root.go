package commands

import (
	"github.com/spf13/cobra"
)

var (
	// Global flags
	dryRun     bool
	appVersion string
)

var rootCmd = &cobra.Command{
	Use:   "git-worktree-manager",
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

// SetVersion sets the version string (called from main)
func SetVersion(v string) {
	appVersion = v
}

// GetVersion returns the app version
func GetVersion() string {
	return appVersion
}

// GetDryRun returns the dry-run flag value
func GetDryRun() bool {
	return dryRun
}
