package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lucasmodrich/git-worktree-manager/internal/git"
	"github.com/lucasmodrich/git-worktree-manager/internal/ui"
	"github.com/spf13/cobra"
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

	// Reject slashes — they cause directory path issues on all platforms
	if strings.ContainsAny(branchName, `/\`) {
		ui.PrintError(
			fmt.Errorf("branch name %q contains a slash", branchName),
			"Use hyphens instead of slashes (e.g. 'feature-my-thing' not 'feature/my-thing')",
		)
		return
	}

	// Resolve the worktree repo root (works from any subdirectory)
	root, err := findWorktreeRoot()
	if err != nil {
		ui.PrintError(err, "Run this command from within a worktree-managed repository")
		return
	}

	// baseBranch is kept local so it doesn't bleed between invocations
	var baseBranch string
	if len(args) > 1 {
		baseBranch = args[1]
	}

	client := git.NewClient(root)
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

	branchExistsLocal := client.BranchExists(branchName, false)
	branchExistsRemote := client.BranchExists(branchName, true)

	var shouldPush bool

	if branchExistsLocal {
		ui.PrintStatus("📂", "Branch '"+branchName+"' exists locally — creating worktree from it")

		if !branchExistsRemote {
			ui.PrintStatus("⚠️", "Branch '"+branchName+"' not found on remote")

			answer, err := ui.PromptYesNo("☁️  Push branch to remote?", os.Stdin)
			if err != nil {
				ui.PrintError(err, "Invalid input")
				return
			}
			if answer {
				shouldPush = true
			}
		}
	} else if branchExistsRemote {
		ui.PrintStatus("☁️", "Branch '"+branchName+"' exists on remote but not locally")

		answer, err := ui.PromptYesNo("📥 Fetch and create worktree from remote branch?", os.Stdin)
		if err != nil {
			ui.PrintError(err, "Invalid input")
			return
		}
		if !answer {
			ui.PrintStatus("❌", "Cancelled")
			return
		}

		if err := client.CreateBranch(branchName, "origin/"+branchName); err != nil {
			ui.PrintError(err, "Failed to create tracking branch")
			return
		}
	} else {
		if baseBranch == "" {
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

	worktreePath := filepath.Join(root, branchName)
	if err := client.WorktreeAdd(worktreePath, branchName, false); err != nil {
		ui.PrintError(err, "Failed to create worktree")
		return
	}

	if shouldPush {
		ui.PrintStatus("☁️", "Pushing new branch '"+branchName+"' to origin")
		if err := client.Push(branchName, true); err != nil {
			ui.PrintError(err, "Failed to push branch to remote")
			return
		}
	}

	ui.PrintStatus("✅", "Worktree for '"+branchName+"' is ready at: "+worktreePath)
}
