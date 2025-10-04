package commands

import (
	"fmt"

	"github.com/lucasmodrich/git-worktree-manager/internal/git"
	"github.com/lucasmodrich/git-worktree-manager/internal/ui"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all worktrees",
	Long:  "Display all active git worktrees in the current repository",
	Run:   runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) {
	// Verify we're in a worktree-managed repo
	if err := verifyWorktreeRepo(); err != nil {
		ui.PrintError(err, "Run this command from a directory where .git points to .bare")
		return
	}

	// Create git client
	client := git.NewClient(".")
	client.DryRun = GetDryRun()

	if client.DryRun {
		ui.PrintDryRun("Would list all worktrees")
		return
	}

	// Get worktree list
	worktrees, err := client.WorktreeList()
	if err != nil {
		ui.PrintError(err, "Failed to list worktrees")
		return
	}

	// Display worktrees
	ui.PrintStatus("📋", "Active Git worktrees:")
	for _, wt := range worktrees {
		fmt.Println(wt)
	}
}
