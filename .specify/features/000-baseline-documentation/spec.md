# Feature Specification: Git Worktree Manager (Baseline)

**Feature Branch**: `000-baseline-documentation`
**Created**: 2025-10-03
**Status**: Baseline (Documenting Existing Implementation)
**Purpose**: Retrospective specification of the existing git-worktree-manager.sh tool

---

## Overview & Context

Git Worktree Manager is a **self-contained, self-updating Bash CLI tool** that simplifies Git worktree management using a bare clone + worktree workflow. It eliminates the complexity of managing multiple feature branches by creating isolated working directories for each branch, all sharing a single Git metadata store.

### Problem Statement
Developers working on multiple branches simultaneously face several challenges:
- **Context switching overhead**: Stashing/unstashing changes, checking out different branches
- **Detached HEAD states**: Confusion when working directly on commits
- **Disk inefficiency**: Multiple full clones waste space
- **Branch management complexity**: Manually tracking which directories contain which branches

### Solution
A single bash script that automates the bare repository + worktree pattern, providing:
- One-command repository setup from GitHub
- Automatic branch creation with worktree isolation
- Safe cleanup of branches and worktrees
- Self-updating capability for maintenance-free usage
- Dry-run mode for preview before execution

---

## User Scenarios & Testing

### Primary User Stories

#### Story 1: Initial Repository Setup
**As a developer**, I want to quickly set up a repository using the worktree workflow, so I can immediately start working on multiple branches without complex Git commands.

**Given** I have a GitHub repository `org/repo`
**When** I run `./git-worktree-manager.sh org/repo`
**Then** the system:
- Creates a root folder named `repo/`
- Clones the repository as a bare clone in `.bare/`
- Creates a `.git` file pointing to `.bare/`
- Configures Git for automatic remote tracking
- Detects the default branch (e.g., `main`)
- Creates the first worktree for the default branch
- Shows me a list of all worktrees

#### Story 2: Creating a New Feature Branch
**As a developer**, I want to create a new feature branch in its own worktree, so I can work on it without affecting other branches.

**Given** I'm in a repository root (where `.git` points to `.bare`)
**When** I run `./git-worktree-manager.sh --new-branch feature/login-page`
**Then** the system:
- Fetches the latest changes from origin
- Creates a new branch `feature/login-page` from the default branch
- Creates a worktree directory `feature/login-page/`
- Pushes the new branch to origin
- Shows me the updated worktree list

#### Story 3: Safe Cleanup with Preview
**As a developer**, I want to preview what will happen before removing a branch, so I can avoid mistakes.

**Given** I have a worktree `feature/old-feature` I want to remove
**When** I run `./git-worktree-manager.sh --dry-run --remove feature/old-feature --remote`
**Then** the system:
- Shows me what WOULD be removed (worktree, local branch, remote branch)
- Does NOT actually delete anything
- Gives me confidence about what will happen when I run without `--dry-run`

#### Story 4: Self-Updating
**As a user**, I want the script to update itself when a new version is available, so I always have the latest features and fixes.

**Given** a new version of the script is available on GitHub
**When** I run `./git-worktree-manager.sh --version`
**Then** the system:
- Shows my current version
- Checks GitHub for the latest version
- Compares versions using semantic versioning rules
- Tells me if an upgrade is available and how to upgrade

**When** I run `./git-worktree-manager.sh --upgrade`
**Then** the system:
- Downloads the new script from GitHub
- Replaces the old script with the new one
- Preserves my installation directory (respects `GIT_WORKTREE_MANAGER_HOME`)
- Downloads accompanying files (README, LICENSE, VERSION)
- Confirms the upgrade succeeded

### Edge Cases & Error Scenarios

#### Input Validation
- **Invalid repository format**: `./script.sh invalid-format` → Error with usage examples
- **Malicious input**: Command injection attempts are sanitized
- **Missing repository**: Non-existent GitHub repo → Git clone fails with clear message

#### Network & Git Failures
- **Network unavailable during upgrade**: Curl fails → Script shows error, doesn't corrupt existing installation
- **Git operation fails**: Any git command failure exits with error code 1
- **Remote branch doesn't exist**: When creating worktree from non-existent base → Error message

#### Worktree Management
- **Worktree already exists**: Attempting to create duplicate → Git error (existing behavior)
- **Branch doesn't exist for removal**: Clear error message, no crash
- **Removing worktree while inside it**: Git handles this (user must exit first)

#### Version Comparison
- **Semantic version edge cases**:
  - `1.2.3 > 1.2.3-beta.1` (release > prerelease)
  - `1.2.3-beta.2 > 1.2.3-beta.1` (prerelease precedence)
  - `1.2.3 <= 1.2.3` (equal versions, no upgrade)

---

## Requirements

### Functional Requirements

#### Core Worktree Operations
- **FR-001**: System MUST support full repository setup from GitHub using `org/repo` shorthand
- **FR-002**: System MUST validate repository format and show clear error messages for invalid input
- **FR-003**: System MUST create bare clone in `.bare/` directory
- **FR-004**: System MUST create `.git` file pointing to bare repository
- **FR-005**: System MUST detect the default branch automatically
- **FR-006**: System MUST create initial worktree for the default branch

#### Branch & Worktree Management
- **FR-007**: Users MUST be able to create new branch worktrees with `--new-branch <branch> [base]`
- **FR-008**: System MUST automatically push new branches to origin
- **FR-009**: System MUST fetch latest changes before creating worktrees
- **FR-010**: Users MUST be able to list all active worktrees with `--list`
- **FR-011**: Users MUST be able to prune stale worktrees with `--prune`

#### Branch Removal
- **FR-012**: Users MUST be able to remove worktrees and local branches with `--remove <branch>`
- **FR-013**: System MUST support optional remote branch deletion with `--remove <branch> --remote`
- **FR-014**: System MUST provide clear feedback about what was deleted (local vs remote)

#### Safety & Preview
- **FR-015**: System MUST support `--dry-run` flag for all destructive operations
- **FR-016**: Dry-run mode MUST preview actions with `🔍 [DRY-RUN]` prefix
- **FR-017**: Dry-run mode MUST NOT execute any actual changes
- **FR-018**: System MUST validate all user inputs to prevent command injection

#### Self-Update & Versioning
- **FR-019**: System MUST check for upgrades against GitHub main branch
- **FR-020**: System MUST compare versions using semantic versioning rules (including prerelease precedence)
- **FR-021**: Users MUST be able to see current version with `--version`
- **FR-022**: Users MUST be able to self-upgrade with `--upgrade`
- **FR-023**: Upgrade MUST download script, README, LICENSE, and VERSION files
- **FR-024**: Upgrade MUST use atomic operations (temp file → move)
- **FR-025**: Upgrade MUST preserve user's installation directory

#### Installation & Configuration
- **FR-026**: Script MUST install to `$HOME/.git-worktree-manager/` by default
- **FR-027**: Installation directory MUST be configurable via `GIT_WORKTREE_MANAGER_HOME` environment variable
- **FR-028**: Script MUST be executable as a single file with no external dependencies (except bash, git, curl)

#### User Experience
- **FR-029**: System MUST provide comprehensive help text with `--help` or `-h`
- **FR-030**: System MUST use emoji-prefixed status messages for clarity (✅ ❌ 📡 🌱 ☁️ 🗑 🔍)
- **FR-031**: Error messages MUST be actionable (tell user how to fix)
- **FR-032**: Long operations MUST show progress indicators

### Non-Functional Requirements

#### Performance
- **NFR-001**: Repository operations should complete in <5 seconds for typical repos (excluding network time)
- **NFR-002**: Version comparison must handle all semantic versioning edge cases correctly
- **NFR-003**: Script startup time should be <100ms

#### Portability
- **NFR-004**: Script MUST work on Linux (Ubuntu/Debian primary target)
- **NFR-005**: Script SHOULD work on macOS (best-effort support)
- **NFR-006**: Script MUST require only bash 4.0+, git 2.0+, curl, standard Unix tools
- **NFR-007**: Script MUST NOT require root/sudo privileges

#### Reliability
- **NFR-008**: All curl operations must check for failures
- **NFR-009**: All git operations must exit on error (`set -e`)
- **NFR-010**: Upgrade operations must not corrupt existing installations on failure
- **NFR-011**: Input validation must prevent command injection attacks

#### Maintainability
- **NFR-012**: Script MUST remain a single file (~500 lines)
- **NFR-013**: All functions must use descriptive names with underscores
- **NFR-014**: All variables must be quoted to prevent word splitting
- **NFR-015**: Script version MUST be updatable by semantic-release automation (line 5)

#### Testing
- **NFR-016**: Test coverage must remain at 100% pass rate
- **NFR-017**: All new features must have corresponding test coverage
- **NFR-018**: Tests must use bash with `set -euo pipefail`
- **NFR-019**: Test output must use emoji indicators for readability

---

## Key Entities

### Repository Setup
- **Bare Clone**: The `.bare/` directory containing Git metadata and objects
- **Worktree**: Individual branch working directories (e.g., `main/`, `feature/login/`)
- **Git File**: The `.git` file (not directory) pointing to `.bare/`

### Version Information
- **Script Version**: Hardcoded at line 5 of script (`SCRIPT_VERSION="x.y.z"`)
- **Remote Version**: Fetched from GitHub main branch via curl
- **VERSION File**: Text file containing current release version

### Configuration
- **Installation Directory**: Where script and assets are stored (default: `~/.git-worktree-manager/`)
- **Environment Variable**: `GIT_WORKTREE_MANAGER_HOME` for custom installation location

### Command Operations
- **Full Setup**: Initial repository clone and worktree creation
- **Branch Creation**: New branch with automatic worktree and push
- **Branch Removal**: Worktree removal, local branch deletion, optional remote deletion
- **Dry Run**: Preview mode for destructive operations
- **Upgrade**: Self-update from GitHub

---

## Acceptance Criteria

### Core Functionality
- [ ] Repository setup creates correct directory structure (`.bare/`, `.git` file, default branch worktree)
- [ ] New branch creation fetches, creates worktree, and pushes to origin
- [ ] Branch removal deletes worktree and local branch (remote optional with `--remote`)
- [ ] Dry-run mode previews all actions without executing
- [ ] Worktree list shows all active worktrees
- [ ] Worktree prune removes stale entries

### Version Management
- [ ] Version check compares local vs remote correctly
- [ ] Semantic version comparison handles all edge cases (prerelease, build metadata, etc.)
- [ ] Upgrade downloads all files (script, README, LICENSE, VERSION)
- [ ] Upgrade preserves installation directory setting
- [ ] Upgrade fails safely without corrupting existing installation

### Safety & Quality
- [ ] Input validation prevents command injection
- [ ] Invalid repository formats show clear error messages
- [ ] All curl operations check for failures
- [ ] Error messages are actionable
- [ ] Exit codes are consistent (0=success, 1=user error, 2=system error)

### Testing
- [ ] All 36 tests pass (version comparison, input validation, dry-run)
- [ ] Test runner provides clear pass/fail reporting
- [ ] Emoji indicators improve test output readability

---

## Success Metrics

**Baseline Performance** (as of v1.3.0):
- **Test Coverage**: 36 tests across 3 suites, 100% passing ✅
- **Script Size**: 511 lines, single file
- **Dependencies**: Zero (beyond bash, git, curl)
- **Supported Platforms**: Linux (Ubuntu tested), macOS (best-effort)
- **Install Locations**: Configurable via environment variable

**Quality Gates**:
- All tests must pass before merge
- No command injection vulnerabilities
- Backward compatibility maintained
- Conventional commit messages required
- Semantic release automation functional

---

## Dependencies & Assumptions

### External Dependencies
- **Bash 4.0+**: For associative arrays, `[[ ]]` conditionals
- **Git 2.0+**: For worktree support (2.5+ recommended)
- **curl**: For downloading upgrades and version checks
- **Standard Unix tools**: grep, sed, basename, mkdir, mv, chmod

### Assumptions
- User has GitHub access (SSH or HTTPS)
- User runs script from appropriate directory (repo root for branch operations)
- Network connectivity available for GitHub operations
- User has write permissions for installation directory
- Repository follows standard GitHub conventions (default branch detection)

### GitHub API/Patterns
- Raw file access via `https://raw.githubusercontent.com/org/repo/refs/heads/main/`
- Script version extractable via `grep '^SCRIPT_VERSION='`
- Repository format: `org/repo` or `git@github.com:org/repo.git`

---

## Out of Scope (Current Implementation)

The following are explicitly NOT supported in the current baseline:
- GitLab, Bitbucket, or other Git hosting platforms
- Windows native support (only WSL/Git Bash)
- Interactive mode or guided wizards
- Configuration file support
- Command completion (bash/zsh)
- Rollback/undo functionality
- Multi-repository management
- Branch protection awareness
- CI/CD integration beyond GitHub Actions

---

## Future Considerations

Potential enhancements identified but not yet implemented:
1. **Platform Support**: GitLab/Bitbucket integration
2. **Shell Completion**: Bash/Zsh tab completion
3. **Interactive Mode**: Guided setup for new users
4. **Configuration File**: `.git-worktree-manager.conf` for defaults
5. **Rollback Support**: Undo last operation
6. **Better Error Recovery**: Automatic cleanup on failed operations
7. **Multi-Repo Management**: Manage multiple repositories
8. **Branch Protection**: Check before deleting protected branches

---

**Status**: ✅ BASELINE ESTABLISHED
**Next Steps**: Use this spec as reference for future feature development with Specify workflow
