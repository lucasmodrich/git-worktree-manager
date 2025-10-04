package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/lucasmodrich/git-worktree-manager/internal/git"
	"github.com/lucasmodrich/git-worktree-manager/internal/ui"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup <org>/<repo>",
	Short: "Full repository setup",
	Long:  "Clone a repository as a bare repo and create initial worktree for the default branch",
	Args:  cobra.ExactArgs(1),
	Run:   runSetup,
}

func init() {
	rootCmd.AddCommand(setupCmd)
}

func runSetup(cmd *cobra.Command, args []string) {
	repoSpec := args[0]

	// Parse repository specification
	url, repoName, err := parseRepoSpec(repoSpec)
	if err != nil {
		ui.PrintError(err, "Examples: acme/webapp, user123/my-project")
		return
	}

	// Check if .bare already exists
	bareDir := ".bare"
	if _, err := os.Stat(bareDir); !os.IsNotExist(err) {
		ui.PrintError(fmt.Errorf(".bare directory already exists in current directory"),
			"Remove existing .bare directory or run setup in a different directory")
		return
	}

	// Create git client
	client := git.NewClient(".")
	client.DryRun = GetDryRun()

	if client.DryRun {
		ui.PrintDryRun("Would create project root: " + repoName)
		ui.PrintDryRun("Would clone bare repository into .bare")
		ui.PrintDryRun("Would create .git file pointing to .bare")
		ui.PrintDryRun("Would configure Git for auto remote tracking")
		ui.PrintDryRun("Would fetch all remote branches")
		ui.PrintDryRun("Would create initial worktree for default branch")
		return
	}

	// Create project root directory
	ui.PrintStatus("📂", "Creating project root: "+repoName)
	if err := os.MkdirAll(repoName, 0755); err != nil {
		ui.PrintError(err, "Failed to create project directory")
		return
	}

	// Change to project directory
	if err := os.Chdir(repoName); err != nil {
		ui.PrintError(err, "Failed to change to project directory")
		return
	}

	// Clone bare repository
	ui.PrintStatus("📦", "Cloning bare repository into .bare")
	if err := client.Clone(url, bareDir, true); err != nil {
		ui.PrintError(err, "Check network connection and verify repository URL is accessible")
		return
	}

	// Create .git file pointing to .bare
	ui.PrintStatus("📝", "Creating .git file pointing to .bare")
	if err := os.WriteFile(".git", []byte("gitdir: ./"+bareDir), 0644); err != nil {
		ui.PrintError(err, "Failed to create .git file")
		return
	}

	// Configure git settings
	ui.PrintStatus("⚙️", "Configuring Git for auto remote tracking")
	if err := client.ConfigureWorktreeSettings(); err != nil {
		ui.PrintError(err, "Failed to configure git settings")
		return
	}

	ui.PrintStatus("🔧", "Ensuring all remote branches are fetched")
	if err := client.ConfigureFetchRefspec(); err != nil {
		ui.PrintError(err, "Failed to configure fetch refspec")
		return
	}

	// Fetch all remote branches
	ui.PrintStatus("📡", "Fetching all remote branches")
	if err := client.Fetch(true, false); err != nil {
		ui.PrintError(err, "Failed to fetch remote branches")
		return
	}

	// Detect default branch
	defaultBranch, err := client.DetectDefaultBranch()
	if err != nil {
		ui.PrintError(err, "Could not detect default branch")
		return
	}

	// Create initial worktree for default branch
	ui.PrintStatus("🌱", "Creating initial worktree for branch: "+defaultBranch)
	worktreePath := filepath.Join(".", defaultBranch)
	if err := client.WorktreeAdd(worktreePath, defaultBranch, false); err != nil {
		ui.PrintError(err, "Failed to create worktree for default branch")
		return
	}

	ui.PrintStatus("✅", "Setup complete!")
}

// parseRepoSpec parses org/repo or git@github.com:org/repo.git format
func parseRepoSpec(spec string) (url, repoName string, err error) {
	// Match org/repo format
	orgRepoRegex := regexp.MustCompile(`^([a-zA-Z0-9_-]+)/([a-zA-Z0-9_-]+)$`)
	if matches := orgRepoRegex.FindStringSubmatch(spec); matches != nil {
		org := matches[1]
		repo := matches[2]
		url = fmt.Sprintf("git@github.com:%s/%s.git", org, repo)
		repoName = repo
		return url, repoName, nil
	}

	// Match git@github.com:org/repo.git format
	sshRegex := regexp.MustCompile(`^git@github\.com:([a-zA-Z0-9_-]+)/([a-zA-Z0-9_-]+)\.git$`)
	if matches := sshRegex.FindStringSubmatch(spec); matches != nil {
		repo := matches[2]
		url = spec
		repoName = repo
		return url, repoName, nil
	}

	return "", "", fmt.Errorf("invalid repository format. Expected: org/repo or git@github.com:org/repo.git")
}
