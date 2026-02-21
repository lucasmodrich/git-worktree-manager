package commands

import (
	"path/filepath"

	"github.com/lucasmodrich/git-worktree-manager/internal/git"
	"github.com/lucasmodrich/git-worktree-manager/internal/ui"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove <branch>",
	Short: "Remove worktree and branch",
	Long:  "Remove a worktree and its associated local branch. Use --remote to also delete the remote branch.",
	Args:  cobra.ExactArgs(1),
	Run:   runRemove,
}

func init() {
	removeCmd.Flags().Bool("remote", false, "Also delete remote branch")
	rootCmd.AddCommand(removeCmd)
}

func runRemove(cmd *cobra.Command, args []string) {
	branchName := args[0]
	removeRemote, _ := cmd.Flags().GetBool("remote")

	root, err := findWorktreeRoot()
	if err != nil {
		ui.PrintError(err, "Run this command from within a worktree-managed repository")
		return
	}

	client := git.NewClient(root)
	client.DryRun = GetDryRun()

	if client.DryRun {
		ui.PrintDryRun("Would remove worktree '" + branchName + "'")
		ui.PrintDryRun("Would delete local branch '" + branchName + "'")
		if removeRemote {
			ui.PrintDryRun("Would delete remote branch 'origin/" + branchName + "'")
		}
		return
	}

	worktreePath := filepath.Join(root, branchName)
	ui.PrintStatus("🗑", "Removing worktree '"+branchName+"'")
	if err := client.WorktreeRemove(worktreePath); err != nil {
		ui.PrintError(err, "Use 'gwtm list' to see available worktrees")
		return
	}

	ui.PrintStatus("🧨", "Deleting local branch '"+branchName+"'")
	if err := client.DeleteBranch(branchName, false); err != nil {
		ui.PrintError(err, "Branch may have already been deleted")
		// Continue anyway — worktree was already removed
	}

	if removeRemote {
		ui.PrintStatus("☁️", "Deleting remote branch 'origin/"+branchName+"'")
		if err := client.DeleteRemoteBranch(branchName); err != nil {
			ui.PrintError(err, "Remote branch may not exist or network issue")
			return
		}
	}

	ui.PrintStatus("✅", "Removal complete.")
}
