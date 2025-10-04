package commands

import (
	"github.com/lucasmodrich/git-worktree-manager/internal/git"
	"github.com/lucasmodrich/git-worktree-manager/internal/ui"
	"github.com/spf13/cobra"
)

var pruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Prune stale worktrees",
	Long:  "Remove stale worktree references from .git/worktrees",
	Run:   runPrune,
}

func init() {
	rootCmd.AddCommand(pruneCmd)
}

func runPrune(cmd *cobra.Command, args []string) {
	// Verify we're in a worktree-managed repo
	if err := verifyWorktreeRepo(); err != nil {
		ui.PrintError(err, "Run this command from a directory where .git points to .bare")
		return
	}

	// Create git client
	client := git.NewClient(".")
	client.DryRun = GetDryRun()

	if client.DryRun {
		ui.PrintDryRun("Would prune stale worktrees")
		return
	}

	// Prune worktrees
	ui.PrintStatus("🧹", "Pruning stale worktrees...")

	if err := client.WorktreePrune(); err != nil {
		ui.PrintError(err, "Failed to prune worktrees")
		return
	}

	ui.PrintStatus("✅", "Prune complete.")
}
