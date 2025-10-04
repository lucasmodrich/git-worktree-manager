# CLI Interface Contract

**Feature**: 002-go-cli-redesign
**Purpose**: Define exact CLI interface for Go implementation (must match Bash version)
**Date**: 2025-10-03

---

## Global Flags

These flags apply to all commands:

```
--dry-run          Preview actions without executing (shows what would happen)
--help, -h         Show help for command
```

---

## Commands

### 1. Full Repository Setup

**Syntax**:
```bash
git-worktree-manager <org>/<repo>
git-worktree-manager git@github.com:<org>/<repo>.git
```

**Arguments**:
- `<org>/<repo>`: GitHub organization/repository shorthand (required)
- OR `git@github.com:<org>/<repo>.git`: Full SSH URL (alternative format)

**Behavior**:
1. Validate `.bare` does not exist (FR-002a)
2. Create project directory named after repo
3. Clone bare repository to `.bare/`
4. Create `.git` file pointing to `.bare`
5. Configure git settings (push.default, autosetupmerge, etc.)
6. Detect default branch from remote
7. Create initial worktree for default branch

**Exit Codes**:
- `0`: Success
- `1`: Invalid repository format, .bare already exists
- `2`: Network failure, git command failed

**Example Output**:
```
📂 Creating project root: my-repo
📦 Cloning bare repository into .bare
📝 Creating .git file pointing to .bare
⚙️  Configuring Git for auto remote tracking
🔧 Ensuring all remote branches are fetched
📡 Fetching all remote branches
🌱 Creating initial worktree for branch: main
✅ Setup complete!
```

**Contract Test**: Verify identical directory structure and git config as Bash version

---

### 2. Create New Branch Worktree

**Syntax**:
```bash
git-worktree-manager --new-branch <branch-name> [base-branch]
```

**Flags**:
- `--new-branch <branch-name>`: Branch name to create (required)

**Arguments**:
- `[base-branch]`: Branch to create from (optional, defaults to detected default branch)

**Behavior**:
1. Verify running in worktree-managed repo (FR-046)
2. Fetch latest from origin
3. If branch exists locally: Prompt for confirmation (FR-013a)
   - If user confirms: Check if branch exists remotely
   - If not remote: Prompt to push (FR-013b)
4. If branch doesn't exist: Create from base branch, track origin
5. Create worktree directory
6. Push new branch to remote (if newly created)
7. Display worktree list

**Exit Codes**:
- `0`: Success
- `1`: Not in worktree repo, invalid branch name, user declined prompt
- `2`: Git operation failed, network failure

**Example Output (new branch)**:
```
📡 Fetching latest from origin
🌱 Creating new branch 'feature/new-ui' from 'main'
☁️  Pushing new branch 'feature/new-ui' to origin
✅ Worktree for 'feature/new-ui' is ready
```

**Example Output (existing branch)**:
```
📡 Fetching latest from origin
📂 Branch 'feature/new-ui' exists locally — creating worktree from it
⚠️  Branch 'feature/new-ui' not found on remote
☁️  Push branch to remote? [y/N]: y
☁️  Pushing branch 'feature/new-ui' to origin
✅ Worktree for 'feature/new-ui' is ready
```

**Contract Test**: Verify prompt behavior matches specification

---

### 3. List Worktrees

**Syntax**:
```bash
git-worktree-manager --list
```

**Flags**: None

**Behavior**:
1. Verify running in worktree-managed repo (FR-046)
2. Execute `git worktree list`
3. Display output

**Exit Codes**:
- `0`: Success
- `1`: Not in worktree repo

**Example Output**:
```
📋 Active Git worktrees:
/path/to/repo/.bare   (bare)
/path/to/repo/main    abc1234 [main]
/path/to/repo/feature def5678 [feature/new-ui]
```

**Contract Test**: Verify output format matches `git worktree list` exactly

---

### 4. Remove Worktree and Branch

**Syntax**:
```bash
git-worktree-manager --remove <branch-name> [--remote]
```

**Flags**:
- `--remove <branch-name>`: Branch/worktree to remove (required)
- `--remote`: Also delete remote branch (optional)

**Behavior**:
1. Verify running in worktree-managed repo (FR-046)
2. Validate worktree exists (FR-018)
3. Remove worktree directory
4. Delete local branch
5. If `--remote` specified: Delete remote branch (FR-016)

**Exit Codes**:
- `0`: Success
- `1`: Worktree not found, not in worktree repo
- `2`: Git operation failed, network failure (if --remote)

**Example Output (local only)**:
```
🗑 Removing worktree 'feature/old-code'
🧨 Deleting local branch 'feature/old-code'
✅ Removal complete.
```

**Example Output (with --remote)**:
```
🗑 Removing worktree 'feature/old-code'
🧨 Deleting local branch 'feature/old-code'
☁️  Deleting remote branch 'origin/feature/old-code'
✅ Removal complete.
```

**Contract Test**: Verify worktree and branch removed from git state

---

### 5. Prune Stale Worktrees

**Syntax**:
```bash
git-worktree-manager --prune
```

**Flags**: None

**Behavior**:
1. Verify running in worktree-managed repo (FR-046)
2. Execute `git worktree prune`
3. Display completion message

**Exit Codes**:
- `0`: Success
- `1`: Not in worktree repo
- `2`: Git operation failed

**Example Output**:
```
🧹 Pruning stale worktrees...
✅ Prune complete.
```

**Contract Test**: Verify stale entries removed from `.git/worktrees/`

---

### 6. Show Version

**Syntax**:
```bash
git-worktree-manager --version
```

**Flags**: None

**Behavior**:
1. Display current version
2. Check GitHub for newer version (FR-021)
3. If newer version available: Show upgrade message
4. If on latest: Confirm

**Exit Codes**:
- `0`: Success (regardless of update availability)

**Example Output (upgrade available)**:
```
git-worktree-manager.sh version 1.3.0
🔍 Checking for newer version on GitHub...
🔢 Local version: 1.3.0
🌐 Remote version: 1.4.0
1.4.0 > 1.3.0
⬇️  Run 'git-worktree-manager --upgrade' to upgrade to version 1.4.0.
```

**Example Output (latest)**:
```
git-worktree-manager.sh version 1.3.0
🔍 Checking for newer version on GitHub...
🔢 Local version: 1.3.0
🌐 Remote version: 1.3.0
1.3.0 <= 1.3.0
✅ You already have the latest version.
```

**Contract Test**: Verify version comparison produces identical result to Bash

---

### 7. Self-Upgrade

**Syntax**:
```bash
git-worktree-manager --upgrade
```

**Flags**: None

**Behavior**:
1. Check for newer version on GitHub (FR-021, FR-027)
2. If not newer: Display message and exit
3. Download binary for current OS/arch from GitHub Releases (FR-026)
4. Download README, VERSION, LICENSE files
5. Verify checksum against checksums.txt (FR-045)
6. Replace current executable atomically
7. Preserve executable permissions (FR-028)

**Exit Codes**:
- `0`: Success (upgrade completed or already on latest)
- `2`: Network failure, download failed, checksum mismatch

**Example Output (upgrade performed)**:
```
🔍 Checking for newer version on GitHub...
⬇️  Upgrading to version 1.4.0...
✓ Binary downloaded
✓ Checksum verified
✓ README.md downloaded
✓ VERSION downloaded
✓ LICENSE downloaded
✅ Upgrade complete. Now running version 1.4.0.
```

**Example Output (already latest)**:
```
🔍 Checking for newer version on GitHub...
✅ You already have the latest version.
```

**Contract Test**: Verify binary replaced without breaking functionality

---

### 8. Help

**Syntax**:
```bash
git-worktree-manager --help
git-worktree-manager -h
```

**Flags**: None (or any command with `-h`)

**Behavior**:
1. Display comprehensive usage information (FR-029)
2. Show all commands with descriptions
3. Show global options
4. Show examples

**Exit Codes**:
- `0`: Success

**Example Output**:
```
🛠 Git Worktree Manager — Help Card

Usage:
  git-worktree-manager <org>/<repo>                     # Full setup from GitHub
  git-worktree-manager --new-branch <branch> [base]     # Create new branch worktree
  git-worktree-manager --remove <branch> [--remote]     # Remove worktree and local branch
  git-worktree-manager --list                           # List active worktrees
  git-worktree-manager --prune                          # Prune stale worktrees
  git-worktree-manager --version                        # Show script version
  git-worktree-manager --upgrade                        # Upgrade to latest version
  git-worktree-manager --help (-h)                      # Show this help card

Global Options:
  --dry-run                           # Preview actions without executing

Examples:
  git-worktree-manager acme/webapp
  git-worktree-manager --new-branch feature/login-page
  git-worktree-manager --remove feature/login-page --remote
  git-worktree-manager --dry-run --new-branch feature/test

Notes:
  - Run from repo root (where .git points to .bare)
  - New branches are pushed to GitHub automatically
  - Use --remote with --remove to also delete the remote branch
  - Installation directory: $GIT_WORKTREE_MANAGER_HOME or $HOME/.git-worktree-manager
```

**Contract Test**: Verify all commands documented, matches Bash version

---

## Dry-Run Mode Behavior

When `--dry-run` is specified (FR-030, FR-031):

**All Commands**:
- Prefix output with `🔍 [DRY-RUN]`
- Perform validation but skip execution
- Show what would happen without making changes
- Skip prompts (log what would be asked)

**Example**:
```bash
git-worktree-manager --dry-run --new-branch feature/test
```

**Output**:
```
🔍 [DRY-RUN] Would fetch latest from origin
🔍 [DRY-RUN] Would create new branch 'feature/test' from 'main'
🔍 [DRY-RUN] Would push new branch 'feature/test' to origin
🔍 [DRY-RUN] Would list all worktrees
```

**Contract Test**: Verify no filesystem or git changes occur in dry-run mode

---

## Error Message Contract

All errors must follow this format (FR-032, FR-033):

```
❌ <Error description>
💡 <Actionable guidance on how to fix>
```

**Examples**:
```
❌ .bare directory already exists in current directory
💡 Remove existing .bare directory or run setup in a different directory

❌ Not in a worktree-managed repository
💡 Run this command from a directory where .git points to .bare

❌ Branch 'feature/old' not found
💡 Use --list to see available worktrees and branches

❌ Git command failed: git clone failed
💡 Check network connection and verify repository URL is accessible
```

**Contract Test**: Verify error messages include guidance and correct emoji

---

## Exit Code Contract

| Code | Meaning | Examples |
|------|---------|----------|
| 0 | Success | Operation completed, version check (even if update available) |
| 1 | User Error | Invalid input, not in worktree repo, .bare exists, user declined prompt |
| 2 | System Error | Network failure, git command failed, permission denied |

**Contract Test**: Verify exit codes match specification for all scenarios

---

## Progress Indicator Contract

Long-running operations must show progress (FR-034):

**Clone Operation**:
```
📦 Cloning bare repository into .bare
[git clone output with progress bar]
```

**Fetch Operation**:
```
📡 Fetching all remote branches
[git fetch output with progress bar]
```

**Contract Test**: Verify progress indicators appear for operations >1 second

---

## Environment Variable Contract

**GIT_WORKTREE_MANAGER_HOME** (FR-039):
- If set: Use as installation base directory
- If unset: Default to `$HOME/.git-worktree-manager`
- Must be absolute path
- Used for storing binary, README, VERSION, LICENSE

**Contract Test**: Verify both implementations respect environment variable

---

## Compatibility Requirements

The Go implementation must:
1. ✅ Accept identical command syntax as Bash version
2. ✅ Produce output with same emoji and message structure
3. ✅ Return same exit codes for same scenarios
4. ✅ Support same global flags (--dry-run, --help)
5. ✅ Support same environment variables (GIT_WORKTREE_MANAGER_HOME)
6. ✅ Create identical git repository structure
7. ✅ Support same git version requirements (2.0+)

**Contract Test Suite**: Run identical commands through both implementations and compare:
- stdout/stderr output (normalized for timestamps/paths)
- Exit codes
- Resulting git repository state
- Created files and directories

---

## Binary Name Update (2025-10-04)

**Change**: The Go CLI binary is being renamed from `git-worktree-manager` to `gwtm`.

**Impact on Contracts**:
- All command syntax examples above reference `git-worktree-manager` for documentation clarity
- In practice, the Go implementation will use binary name `gwtm`
- The Bash script remains `git-worktree-manager.sh` (unchanged)
- CLI interface (commands, flags, behavior) is identical between both names

**Updated Command Examples**:
```bash
# Go CLI (new binary name):
gwtm <org>/<repo>
gwtm new-branch <branch> [base]
gwtm remove <branch> [--remote]
gwtm list
gwtm prune
gwtm version
gwtm upgrade
gwtm --help

# Bash script (unchanged):
git-worktree-manager.sh <org>/<repo>
git-worktree-manager.sh --new-branch <branch> [base]
# ... same flags and syntax
```

**Contract Validation**:
- All tests remain valid; simply replace binary name `git-worktree-manager` with `gwtm` in test execution
- CLI behavior, flags, output format, and exit codes are identical
- Migration guide will provide syml link instructions for backward compatibility
