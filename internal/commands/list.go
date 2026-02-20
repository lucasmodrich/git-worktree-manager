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
	root, err := findWorktreeRoot()
	if err != nil {
		ui.PrintError(err, "Run this command from within a worktree-managed repository")
		return
	}

	client := git.NewClient(root)
	client.DryRun = GetDryRun()

	if client.DryRun {
		ui.PrintDryRun("Would list all worktrees")
		return
	}

	worktrees, err := client.WorktreeList()
	if err != nil {
		ui.PrintError(err, "Failed to list worktrees")
		return
	}

	ui.PrintStatus("📋", "Active Git worktrees:")
	for _, wt := range worktrees {
		fmt.Println(wt)
	}
}
