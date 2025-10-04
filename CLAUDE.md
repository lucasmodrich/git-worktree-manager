# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Overview

This repository contains the Git Worktree Manager, available in two implementations:
- **Go CLI: `gwtm`** (Primary) - Fast, cross-platform compiled binary with enhanced UX
- **Bash Script: `git-worktree-manager.sh`** (Legacy) - Single-file shell script for maximum portability

Both implementations simplify Git worktree management using a bare clone + worktree workflow, providing commands for repository setup, branch creation, worktree management, and self-updating.

## Key Architecture

### Core Components

1. **Main Script**: `git-worktree-manager.sh` - A monolithic Bash script (400+ lines) containing all functionality
   - Uses semantic versioning with self-upgrade capability
   - Script version is hardcoded at line 5: `SCRIPT_VERSION="1.1.7"`
   - Self-installs to `$HOME/.git-worktree-manager/`

2. **Workflow Structure**: Creates and manages the following directory structure:
   ```
   <repo-name>/
   ├── .bare/             # Bare repository clone
   ├── .git               # Points to .bare
   └── <branches>/        # Individual worktrees for each branch
   ```

3. **Version Management**:
   - Implements custom semantic version comparison function (`version_gt()`)
   - Fetches latest version from GitHub for upgrade checks
   - VERSION file tracks current release version

## Common Commands

### Building the Go CLI
```bash
# Build the binary (creates 'gwtm' executable)
make build

# Or use go directly
go build -o gwtm ./cmd/git-worktree-manager

# Install to $GOPATH/bin or $HOME/.git-worktree-manager/
make install

# Test the binary
./gwtm --help
./gwtm version
```

### Testing
```bash
# Run Go tests
go test ./...

# Run version comparison tests (Bash)
./tests/version_compare_tests.sh
```

### Release Process
```bash
# Install dependencies (for local testing of release process)
npm ci

# The actual release is automated via GitHub Actions when pushing to main
# Semantic Release handles versioning, changelog, and GitHub releases
# GoReleaser builds multi-platform binaries (Linux, macOS, Windows)
```

### Bash Script Development
```bash
# Make script executable
chmod +x git-worktree-manager.sh

# Test script locally
./git-worktree-manager.sh --help

# Version check
./git-worktree-manager.sh --version

# Self-upgrade (fetches from GitHub main branch)
./git-worktree-manager.sh --upgrade
```

## Release Configuration

The project uses semantic-release for automated versioning:
- **release.config.js**: Configures semantic-release plugins
- Automatically updates:
  - SCRIPT_VERSION in git-worktree-manager.sh
  - VERSION file
  - CHANGELOG.md
- Creates GitHub releases with the script as an asset
- GitHub Actions workflow (`.github/workflows/release.yml`) triggers on main branch pushes

## Important Implementation Details

1. **Version Comparison Logic**: The `version_gt()` function (lines 17-111) implements full semantic versioning comparison including prerelease precedence rules

2. **Self-Upgrade Mechanism**:
   - Fetches script from `https://raw.githubusercontent.com/lucasmodrich/git-worktree-manager/refs/heads/main/git-worktree-manager.sh`
   - Compares versions before replacing
   - Preserves executable permissions

3. **Error Handling**: Uses `set -e` for strict error handling throughout the script

4. **Git Worktree Operations**:
   - Full setup from GitHub using org/repo shorthand
   - Branch creation with automatic remote push
   - Worktree listing, pruning, and removal
   - Configures fetch to include all remote refs

## Development Notes

- The script is designed to be self-contained with no external dependencies beyond standard Unix tools and Git
- All functionality is in a single file for easy distribution
- The script assumes `$HOME/.git-worktree-manager/` as its installation directory (hardcoded)
- Uses bash-specific features, requires bash shell (not sh-compatible)
