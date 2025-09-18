# Git Worktree Manager

## 📌 Overview
`git-worktree-manager.sh` is a self-updating shell script for managing Git repositories using a **bare clone + worktree** workflow.
I created this project following frustration using the standard tooling.

It supports:

- **Full setup** from GitHub using `org/repo` shorthand with input validation
- **Branch creation** with automatic remote push
- **Worktree listing**, pruning, and removal with optional remote cleanup
- **Version tracking** and **self-upgrade** with robust error handling
- **Dry-run mode** for safe preview of actions
- **Configurable installation** directory
- **Comprehensive testing** suite
- **Markdown-style help card** for onboarding

---

## 🚀 Installation

By default `git-worktree-manager` will install and update itself into the `$HOME/.git-worktree-manager/` (`~/.git-worktree-manager`) folder.

### Custom Installation Directory

You can customize the installation directory by setting the `GIT_WORKTREE_MANAGER_HOME` environment variable:

```bash
export GIT_WORKTREE_MANAGER_HOME="/opt/git-tools"
./git-worktree-manager.sh --upgrade
```

This will install the script to `/opt/git-tools/` instead of the default location.


To install directly from this GitHub repo, use the following command:
```bash
curl -sSL https://raw.githubusercontent.com/lucasmodrich/git-worktree-manager/refs/heads/main/git-worktree-manager.sh | bash -s -- --upgrade
```

## 🧠 Versioning & Upgrade

- Check version:
  ```bash
  ./git-worktree-manager.sh --version
  ```
- Upgrade to latest from GitHub:
  ```bash
  ./git-worktree-manager.sh --upgrade
  ```


## 📂 Folder Structure

After setup:
```
<repo-name>/
├── .bare/             # Bare repository clone
├── .git               # Points to .bare
└── <default-branch>/  # Initial worktree
```

---

## 🖼 Architecture Diagram

```
                ┌───────────────────────────┐
                │        .bare repo         │
                │  (Git metadata & objects) │
                └─────────────┬─────────────┘
                              │
     ┌────────────────────────┼────────────────────────┐
     │                        │                        │
┌──────────────┐       ┌──────────────┐        ┌──────────────┐
│ main/        │       │ feature-x/   │        │ bugfix-y/    │
│ (worktree)   │       │ (worktree)   │        │ (worktree)   │
└──────────────┘       └──────────────┘        └──────────────┘
```

---

## 🔄 Flow Diagrams

### Full Setup
```
Input: [--dry-run] <org>/<repo>
→ Validate repository format
→ [DRY-RUN] Preview actions OR
→ Create root folder
→ Clone into .bare
→ Point .git to .bare
→ Configure fetch
→ Fetch branches
→ Detect default branch
→ Create worktree
→ Push if new
```

### Branch Creation
```
Input: [--dry-run] --new-branch <branch> [base]
→ [DRY-RUN] Preview actions OR
→ Fetch branches
→ Create worktree
→ Push if new
```

### Branch Removal
```
Input: [--dry-run] --remove <branch> [--remote]
→ [DRY-RUN] Preview actions OR
→ Remove worktree
→ Delete local branch
→ [OPTIONAL] Delete remote branch
```

---

## 🌐 Local ↔ Remote Relationship

```
GitHub Remote (origin)
┌────────────────────┐
│ origin/main        │
│ origin/feature-x   │
└──────────┬─────────┘
           ▼
       .bare repo
           ▼
┌──────────────┐
│ feature-x/   │
└──────────────┘
```

---

## 🚀 Usage

### Make executable
```bash
chmod +x git-worktree-manager.sh
```

---

### Full Setup
```bash
./git-worktree-manager.sh your-org/your-repo
```

---

### Create New Branch
```bash
./git-worktree-manager.sh --new-branch <branch> [base]
```

---

### Remove Worktree + Branch
```bash
# Remove worktree and local branch only
./git-worktree-manager.sh --remove <branch>

# Remove worktree, local branch, AND remote branch
./git-worktree-manager.sh --remove <branch> --remote
```

---

### List Worktrees
```bash
./git-worktree-manager.sh --list
```

---

### Prune Stale Worktrees
```bash
./git-worktree-manager.sh --prune
```

---

### Show Version
```bash
./git-worktree-manager.sh --version
```

---

### Upgrade Script
```bash
./git-worktree-manager.sh --upgrade
```

---

### Help Card
```bash
./git-worktree-manager.sh --help
./git-worktree-manager.sh -h
```

---

### Dry-run Mode

Preview actions without executing them:

```bash
# Preview repository setup
./git-worktree-manager.sh --dry-run acme/webapp

# Preview branch creation
./git-worktree-manager.sh --dry-run --new-branch feature/test

# Preview branch removal
./git-worktree-manager.sh --dry-run --remove feature/old-branch --remote
```

---

## 📖 Example Workflows

### Basic Workflow
```bash
# Setup repository
./git-worktree-manager.sh acme/webapp

# Create feature branch
./git-worktree-manager.sh --new-branch feature/login-page

# Work in the branch
cd feature/login-page
git add .
git commit -m "Add login page"
git push

# Clean up (local only)
cd ..
./git-worktree-manager.sh --remove feature/login-page
```

### Advanced Workflow with Dry-run and Remote Cleanup
```bash
# Preview setup first
./git-worktree-manager.sh --dry-run acme/webapp

# Actually setup
./git-worktree-manager.sh acme/webapp

# Preview branch creation
./git-worktree-manager.sh --dry-run --new-branch feature/advanced-feature

# Create the branch
./git-worktree-manager.sh --new-branch feature/advanced-feature

# Work and commit
cd feature/advanced-feature
git add .
git commit -m "Implement advanced feature"
git push

# Complete cleanup including remote branch
cd ..
./git-worktree-manager.sh --remove feature/advanced-feature --remote
```

### Custom Installation Directory
```bash
# Set custom installation directory
export GIT_WORKTREE_MANAGER_HOME="/opt/dev-tools"

# Install to custom location
curl -sSL https://raw.githubusercontent.com/lucasmodrich/git-worktree-manager/refs/heads/main/git-worktree-manager.sh | bash -s -- --upgrade

# Script is now available in /opt/dev-tools/
/opt/dev-tools/git-worktree-manager.sh --version
```

---

## 🧪 Testing

The script includes a comprehensive test suite to ensure reliability:

```bash
# Run all tests
./tests/run_all_tests.sh

# Run individual test suites
./tests/version_compare_tests.sh      # Semantic version comparison
./tests/input_validation_tests.sh    # Repository format validation
./tests/dry_run_tests.sh             # Dry-run functionality
```

Current test coverage: **33 tests, 100% passing** ✅

---

## 🔒 Security & Reliability

- **Input validation**: Repository formats are validated against known patterns
- **Error handling**: Comprehensive error checking for network operations
- **Safe operations**: Dry-run mode allows preview before execution
- **No command injection**: All user inputs are properly sanitized
- **Atomic downloads**: Upgrade operations use temporary files for safety

---

## 🔧 Configuration

### Environment Variables

- `GIT_WORKTREE_MANAGER_HOME`: Custom installation directory (default: `$HOME/.git-worktree-manager`)

### Git Alias Setup

For convenience, add this alias to your `~/.gitconfig`:

```ini
[alias]
    wtm = "!bash $HOME/.git-worktree-manager/git-worktree-manager.sh"
```

Then use: `git wtm --help`

> **Note**: The `--help` flag doesn't work when called via `git` as it invokes Git's built-in help.

---

## ✅ Benefits

- **Disk-efficient** multi-branch development
- **No detached HEADs** - each branch has its own working directory
- **Safe operations** with dry-run mode and input validation
- **Easy onboarding** with comprehensive help and examples
- **Self-updating** and version-aware with robust error handling
- **GitHub-native** workflow with SSH and HTTPS support
- **Configurable** installation and behavior
- **Well-tested** with comprehensive test suite
- **Ubuntu-optimized** for reliable bash script execution

---

## 🛠 Requirements

- **Bash 4.0+** (Ubuntu default)
- **Git 2.0+** with worktree support
- **curl** for self-update functionality
- **SSH access** to GitHub (recommended) or HTTPS

---

## 📝 Version History

See [CHANGELOG.md](CHANGELOG.md) for detailed version history and release notes.
