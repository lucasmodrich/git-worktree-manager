# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Overview

This repository contains `git-worktree-manager.sh`, a self-updating Bash script that simplifies Git worktree management using a bare clone + worktree workflow. The script provides commands for repository setup, branch creation, worktree management, and self-updating.

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

### Testing
```bash
# Run version comparison tests
./tests/version_compare_tests.sh
```

### Release Process
```bash
# Install dependencies (for local testing of release process)
npm ci

# The actual release is automated via GitHub Actions when pushing to main
# Semantic Release handles versioning, changelog, and GitHub releases
```

### Script Development
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