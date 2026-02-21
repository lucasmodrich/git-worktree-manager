# Git Worktree Manager (`gwtm`)

## 📌 Overview

`gwtm` is a CLI tool for managing Git repositories using a **bare clone + worktree** workflow. It removes the friction of multi-branch development by giving each branch its own dedicated working directory — no stashing, no switching, no detached HEADs.

Features:
- **Full setup** from GitHub using `org/repo` shorthand
- **Branch creation** with automatic remote push
- **Worktree listing**, pruning, and removal with optional remote cleanup
- **Self-upgrade** with checksum verification
- **Dry-run mode** to preview any operation before executing it
- **Configurable installation** directory via environment variable
- Works from **any subdirectory** within a managed repository

---

## 🚀 Installation

### Download Pre-built Binary (Recommended)

```bash
# Linux amd64
curl -L https://github.com/lucasmodrich/git-worktree-manager/releases/latest/download/gwtm_Linux_x86_64 -o gwtm
chmod +x gwtm
sudo mv gwtm /usr/local/bin/

# macOS Apple Silicon (M1/M2/M3)
curl -L https://github.com/lucasmodrich/git-worktree-manager/releases/latest/download/gwtm_Darwin_arm64 -o gwtm
chmod +x gwtm
sudo mv gwtm /usr/local/bin/

# macOS Intel
curl -L https://github.com/lucasmodrich/git-worktree-manager/releases/latest/download/gwtm_Darwin_x86_64 -o gwtm
chmod +x gwtm
sudo mv gwtm /usr/local/bin/

# Windows (PowerShell)
# Download gwtm_Windows_x86_64.exe from the releases page, rename to gwtm.exe, and add to PATH
```

### Build from Source

Requires Go 1.25.1 or later.

```bash
git clone https://github.com/lucasmodrich/git-worktree-manager.git
cd git-worktree-manager
go build -o gwtm ./cmd/git-worktree-manager
# Move the binary somewhere on your PATH, e.g.:
sudo mv gwtm /usr/local/bin/
```

### Self-Upgrade

```bash
gwtm upgrade
```

---

## 📂 Repository Structure

After running `gwtm setup`, the following structure is created:

```
<repo-name>/
├── .bare/             # Bare repository clone (Git objects & metadata)
├── .git               # File pointing to .bare
└── <default-branch>/  # Initial worktree, ready to work in
```

Additional worktrees are added as sibling directories:

```
<repo-name>/
├── .bare/
├── .git
├── main/
├── feature-auth/
└── bugfix-crash/
```

---

## 🖼 Architecture

```
                ┌───────────────────────────┐
                │        .bare repo         │
                │  (Git metadata & objects) │
                └─────────────┬─────────────┘
                              │
     ┌────────────────────────┼────────────────────────┐
     │                        │                        │
┌──────────────┐       ┌──────────────┐        ┌──────────────┐
│ main/        │       │ feature-auth/│        │ bugfix-crash/│
│ (worktree)   │       │ (worktree)   │        │ (worktree)   │
└──────────────┘       └──────────────┘        └──────────────┘
```

---

## 🚀 Usage

### Help

```bash
gwtm --help
gwtm <command> --help
```

### Full Setup

Clone a repository and create the initial worktree in one step:

```bash
gwtm setup your-org/your-repo
# also accepts full SSH URL:
gwtm setup git@github.com:your-org/your-repo.git
```

### Create a Branch Worktree

Creates a new branch (or checks out an existing one) and adds a worktree for it. Use hyphens in branch names — slashes are not supported as they conflict with directory paths.

```bash
# New branch from the default branch
gwtm new-branch feature-login

# New branch from a specific base
gwtm new-branch bugfix-crash main

# Check out an existing remote branch
gwtm new-branch feature-login    # detects it exists on remote and prompts
```

### Remove a Worktree and Branch

```bash
# Remove worktree and local branch
gwtm remove feature-login

# Remove worktree, local branch, and remote branch
gwtm remove feature-login --remote
```

### List Worktrees

```bash
gwtm list
```

### Prune Stale Worktrees

```bash
gwtm prune
```

### Version

```bash
gwtm version      # shows local version and checks for updates
gwtm --version    # quick version output (no network check)
```

### Upgrade

```bash
gwtm upgrade
```

### Dry-Run Mode

Preview any operation without executing it:

```bash
gwtm --dry-run setup acme/webapp
gwtm --dry-run new-branch feature-payments
gwtm --dry-run remove feature-payments --remote
```

---

## 📖 Example Workflows

### Start a new feature

```bash
# Set up the repository (once)
gwtm setup acme/webapp
cd acme/webapp

# Create a feature branch and worktree
gwtm new-branch feature-login
cd feature-login

# Do your work
git add .
git commit -m "feat: add login page"
git push

# When done, go back up and clean up
cd ..
gwtm remove feature-login --remote
```

### Pick up a colleague's branch

```bash
cd acme/webapp
gwtm new-branch feature-payments   # detects remote branch, prompts to check out
cd feature-payments
```

### Preview before acting

```bash
gwtm --dry-run setup acme/webapp
gwtm --dry-run remove old-branch --remote
```

---

## 🔧 Configuration

### Environment Variables

| Variable | Default | Description |
|---|---|---|
| `GIT_WORKTREE_MANAGER_HOME` | `$HOME/.git-worktree-manager` | Installation directory for `gwtm upgrade` |

### Git Alias (optional)

```ini
# ~/.gitconfig
[alias]
    wt = "!gwtm"
```

Then use: `git wt list`, `git wt new-branch feature-x`, etc.

---

## 🧪 Testing

```bash
# Run Go tests
go test ./...

# Run with coverage
go test -coverprofile=coverage.txt ./...
go tool cover -html=coverage.txt
```

---

## 🔒 Security & Reliability

- **Input validation**: Repository formats are validated; branch names with path separators are rejected
- **Checksum verification**: `gwtm upgrade` verifies SHA-256 checksums before replacing the binary
- **Atomic upgrades**: New binary downloaded to a temp file and moved into place only after verification
- **Dry-run mode**: Preview any destructive operation before executing it
- **Cleanup on failure**: `gwtm setup` removes the partial directory if any step fails

---

## ✅ Benefits

- **Disk-efficient** multi-branch development — shared object store, no redundant copies
- **No detached HEADs** — each branch has its own working directory
- **No stashing** — switch contexts by changing directories
- **Cross-platform** — runs on Linux, macOS, and Windows
- **Self-updating** with version comparison and checksum verification

---

## 🛠 Requirements

- **Git 2.5+** (worktree support)
- **SSH access** to GitHub (recommended) or HTTPS

The `gwtm` binary is statically compiled with no additional runtime dependencies.

---

## 📝 Version History

See [CHANGELOG.md](CHANGELOG.md) for detailed version history and release notes.
