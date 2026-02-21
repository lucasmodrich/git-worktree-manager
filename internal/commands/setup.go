package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

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

	url, repoName, err := parseRepoSpec(repoSpec)
	if err != nil {
		ui.PrintError(err, "Examples: acme/webapp, git@gitlab.com:org/repo.git, https://github.com/org/repo")
		return
	}

	// Compute the absolute path for the new project directory
	cwd, err := os.Getwd()
	if err != nil {
		ui.PrintError(err, "Failed to determine current directory")
		return
	}
	repoDir := filepath.Join(cwd, repoName)

	// Fail early if the directory already exists
	if _, err := os.Stat(repoDir); !os.IsNotExist(err) {
		ui.PrintError(
			fmt.Errorf("directory %q already exists", repoName),
			"Run setup in a different directory or choose a different name",
		)
		return
	}

	// All git operations use repoDir as the working directory
	client := git.NewClient(repoDir)
	client.DryRun = GetDryRun()

	if client.DryRun {
		ui.PrintDryRun("Would create project root: " + repoDir)
		ui.PrintDryRun("Would clone bare repository into .bare")
		ui.PrintDryRun("Would create .git file pointing to .bare")
		ui.PrintDryRun("Would configure Git for auto remote tracking")
		ui.PrintDryRun("Would fetch all remote branches")
		ui.PrintDryRun("Would create initial worktree for default branch")
		return
	}

	// Create project root directory
	ui.PrintStatus("📂", "Creating project root: "+repoName)
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		ui.PrintError(err, "Failed to create project directory")
		return
	}

	// Remove the partially-created directory if any subsequent step fails
	cleanup := true
	defer func() {
		if cleanup {
			os.RemoveAll(repoDir)
		}
	}()

	bareDir := filepath.Join(repoDir, ".bare")

	ui.PrintStatus("📦", "Cloning bare repository into .bare")
	if err := client.Clone(url, bareDir, true); err != nil {
		ui.PrintError(err, "Check network connection and verify repository URL is accessible")
		return
	}

	ui.PrintStatus("📝", "Creating .git file pointing to .bare")
	gitFile := filepath.Join(repoDir, ".git")
	if err := os.WriteFile(gitFile, []byte("gitdir: ./.bare"), 0644); err != nil {
		ui.PrintError(err, "Failed to create .git file")
		return
	}

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

	ui.PrintStatus("📡", "Fetching all remote branches")
	if err := client.Fetch(true, false); err != nil {
		ui.PrintError(err, "Failed to fetch remote branches")
		return
	}

	defaultBranch, err := client.DetectDefaultBranch()
	if err != nil {
		ui.PrintError(err, "Could not detect default branch")
		return
	}

	ui.PrintStatus("🌱", "Creating initial worktree for branch: "+defaultBranch)
	worktreePath := filepath.Join(repoDir, defaultBranch)
	if err := client.WorktreeAdd(worktreePath, defaultBranch, false); err != nil {
		ui.PrintError(err, "Failed to create worktree for default branch")
		return
	}

	cleanup = false // all steps succeeded — keep the directory
	ui.PrintStatus("✅", fmt.Sprintf("Setup complete! cd %s/%s to start working.", repoName, defaultBranch))
}

// parseRepoSpec accepts the following formats and returns the clone URL and repo name:
//   - org/repo            → expanded to git@github.com:org/repo.git
//   - git@<host>:<path>   → used as-is (any SSH host)
//   - https?://<host>/... → used as-is
//
// The repo name is always derived from the last path component (without .git suffix).
func parseRepoSpec(spec string) (url, repoName string, err error) {
	// HTTPS or HTTP URL
	if strings.Contains(spec, "://") {
		name := repoNameFromPath(strings.TrimSuffix(spec, ".git"), "/")
		if name == "" {
			return "", "", fmt.Errorf("cannot determine repository name from URL %q", spec)
		}
		return spec, name, nil
	}

	// Generic SSH URL: git@<host>:<path>[.git]
	if strings.HasPrefix(spec, "git@") {
		colonIdx := strings.Index(spec, ":")
		if colonIdx == -1 {
			return "", "", fmt.Errorf("invalid SSH URL %q: missing ':'", spec)
		}
		name := repoNameFromPath(strings.TrimSuffix(spec[colonIdx+1:], ".git"), "/")
		if name == "" {
			return "", "", fmt.Errorf("cannot determine repository name from SSH URL %q", spec)
		}
		return spec, name, nil
	}

	// org/repo shorthand — expand to GitHub SSH
	orgRepoRegex := regexp.MustCompile(`^([a-zA-Z0-9_.-]+)/([a-zA-Z0-9_.-]+)$`)
	if matches := orgRepoRegex.FindStringSubmatch(spec); matches != nil {
		org, repo := matches[1], matches[2]
		return fmt.Sprintf("git@github.com:%s/%s.git", org, repo), repo, nil
	}

	return "", "", fmt.Errorf("invalid repository format %q\nExamples: org/repo, git@github.com:org/repo.git, https://github.com/org/repo", spec)
}

// repoNameFromPath returns the last non-empty segment of a slash-separated path.
func repoNameFromPath(path, sep string) string {
	parts := strings.Split(path, sep)
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] != "" {
			return parts[i]
		}
	}
	return ""
}
