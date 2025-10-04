package commands

import (
	"fmt"
	"path/filepath"

	"github.com/lucasmodrich/git-worktree-manager/internal/git"
	"github.com/lucasmodrich/git-worktree-manager/internal/ui"
	"github.com/spf13/cobra"
)

var (
	baseBranch string
)

var branchCmd = &cobra.Command{
	Use:   "new-branch <branch-name> [base-branch]",
	Short: "Create new branch worktree",
	Long:  "Create a new worktree for a branch. If the branch doesn't exist, it will be created from the base branch (or default branch).",
	Args:  cobra.RangeArgs(1, 2),
	Run:   runBranch,
}

func init() {
	rootCmd.AddCommand(branchCmd)
}

func runBranch(cmd *cobra.Command, args []string) {
	branchName := args[0]

	// Get base branch from args or detect default
	if len(args) > 1 {
		baseBranch = args[1]
	}

	// Verify we're in a worktree-managed repo
	if err := verifyWorktreeRepo(); err != nil {
		ui.PrintError(err, "Run this command from a directory where .git points to .bare")
		return
	}

	// Create git client
	client := git.NewClient(".")
	client.DryRun = GetDryRun()

	if client.DryRun {
		ui.PrintDryRun("Would fetch latest from origin")
		ui.PrintDryRun("Would create new branch '" + branchName + "'")
		ui.PrintDryRun("Would push new branch '" + branchName + "' to origin")
		ui.PrintDryRun("Would create worktree for '" + branchName + "'")
		return
	}

	// Fetch latest from origin
	ui.PrintStatus("📡", "Fetching latest from origin")
	if err := client.Fetch(true, false); err != nil {
		ui.PrintError(err, "Check network connection")
		return
	}

	// Check if branch exists locally
	branchExistsLocal := client.BranchExists(branchName, false)
	branchExistsRemote := client.BranchExists(branchName, true)

	var shouldPush bool

	if branchExistsLocal {
		// Branch exists locally - prompt for confirmation
		ui.PrintStatus("📂", "Branch '"+branchName+"' exists locally — creating worktree from it")

		if !branchExistsRemote {
			// Branch not on remote - prompt to push
			ui.PrintStatus("⚠️", "Branch '"+branchName+"' not found on remote")

			answer, err := ui.PromptYesNo("☁️  Push branch to remote?")
			if err != nil {
				ui.PrintError(err, "Invalid input")
				return
			}

			if answer {
				shouldPush = true
			}
		}
	} else if branchExistsRemote {
		// Branch exists remotely but not locally - prompt to fetch
		ui.PrintStatus("☁️", "Branch '"+branchName+"' exists on remote but not locally")

		answer, err := ui.PromptYesNo("📥 Fetch and create worktree from remote branch?")
		if err != nil {
			ui.PrintError(err, "Invalid input")
			return
		}

		if !answer {
			ui.PrintStatus("❌", "Cancelled")
			return
		}

		// Create tracking branch from remote
		if err := client.CreateBranch(branchName, "origin/"+branchName); err != nil {
			ui.PrintError(err, "Failed to create tracking branch")
			return
		}
	} else {
		// Branch doesn't exist anywhere - create new
		if baseBranch == "" {
			// Detect default branch
			var err error
			baseBranch, err = client.DetectDefaultBranch()
			if err != nil {
				ui.PrintError(err, "Could not detect default branch")
				return
			}
		}

		ui.PrintStatus("🌱", fmt.Sprintf("Creating new branch '%s' from '%s'", branchName, baseBranch))

		if err := client.CreateBranch(branchName, baseBranch); err != nil {
			ui.PrintError(err, "Failed to create branch")
			return
		}

		shouldPush = true
	}

	// Create worktree
	worktreePath := filepath.Join(".", branchName)

	if err := client.WorktreeAdd(worktreePath, branchName, false); err != nil {
		ui.PrintError(err, "Failed to create worktree")
		return
	}

	// Push to remote if needed
	if shouldPush {
		ui.PrintStatus("☁️", "Pushing new branch '"+branchName+"' to origin")

		if err := client.Push(branchName, true); err != nil {
			ui.PrintError(err, "Failed to push branch to remote")
			return
		}
	}

	ui.PrintStatus("✅", "Worktree for '"+branchName+"' is ready")
}
