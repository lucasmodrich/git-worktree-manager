package commands

import (
	"github.com/lucasmodrich/git-worktree-manager/internal/git"
	"github.com/lucasmodrich/git-worktree-manager/internal/ui"
	"github.com/spf13/cobra"
)

var (
	removeRemote bool
)

var removeCmd = &cobra.Command{
	Use:   "remove <branch>",
	Short: "Remove worktree and branch",
	Long:  "Remove a worktree and its associated local branch. Use --remote to also delete the remote branch.",
	Args:  cobra.ExactArgs(1),
	Run:   runRemove,
}

func init() {
	removeCmd.Flags().BoolVar(&removeRemote, "remote", false, "Also delete remote branch")
	rootCmd.AddCommand(removeCmd)
}

func runRemove(cmd *cobra.Command, args []string) {
	branchName := args[0]

	// Verify we're in a worktree-managed repo
	if err := verifyWorktreeRepo(); err != nil {
		ui.PrintError(err, "Run this command from a directory where .git points to .bare")
		return
	}

	// Create git client
	client := git.NewClient(".")
	client.DryRun = GetDryRun()

	if client.DryRun {
		ui.PrintDryRun("Would remove worktree '" + branchName + "'")
		ui.PrintDryRun("Would delete local branch '" + branchName + "'")
		if removeRemote {
			ui.PrintDryRun("Would delete remote branch 'origin/" + branchName + "'")
		}
		return
	}

	// Remove worktree
	ui.PrintStatus("🗑", "Removing worktree '"+branchName+"'")

	// Construct worktree path (handle feature/ prefixes etc)
	worktreePath := branchName

	if err := client.WorktreeRemove(worktreePath); err != nil {
		ui.PrintError(err, "Use --list to see available worktrees and branches")
		return
	}

	// Delete local branch
	ui.PrintStatus("🧨", "Deleting local branch '"+branchName+"'")

	if err := client.DeleteBranch(branchName, false); err != nil {
		ui.PrintError(err, "Branch may have already been deleted")
		// Continue anyway - not fatal
	}

	// Delete remote branch if requested
	if removeRemote {
		ui.PrintStatus("☁️", "Deleting remote branch 'origin/"+branchName+"'")

		if err := client.DeleteRemoteBranch(branchName); err != nil {
			ui.PrintError(err, "Remote branch may not exist or network issue")
			return
		}
	}

	ui.PrintStatus("✅", "Removal complete.")
}
