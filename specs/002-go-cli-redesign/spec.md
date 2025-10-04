# Feature Specification: Go CLI Redesign

**Feature Branch**: `002-go-cli-redesign`
**Created**: 2025-10-03
**Status**: Draft
**Input**: User description: "Go CLI redesign. I want to redesign the shell script as a GO CLI app that uses Cobra to manage the command line interface. All existing features provided by the shell script should be implemented in the initial version."

## Execution Flow (main)
```
1. Parse user description from Input
   → If empty: ERROR "No feature description provided"
2. Extract key concepts from description
   → Identify: actors, actions, data, constraints
3. For each unclear aspect:
   → Mark with [NEEDS CLARIFICATION: specific question]
4. Fill User Scenarios & Testing section
   → If no clear user flow: ERROR "Cannot determine user scenarios"
5. Generate Functional Requirements
   → Each requirement must be testable
   → Mark ambiguous requirements
6. Identify Key Entities (if data involved)
7. Run Review Checklist
   → If any [NEEDS CLARIFICATION]: WARN "Spec has uncertainties"
   → If implementation details found: ERROR "Remove tech details"
8. Return: SUCCESS (spec ready for planning)
```

---

## ⚡ Quick Guidelines
- ✅ Focus on WHAT users need and WHY
- ❌ Avoid HOW to implement (no tech stack, APIs, code structure)
- 👥 Written for business stakeholders, not developers

### Section Requirements
- **Mandatory sections**: Must be completed for every feature
- **Optional sections**: Include only when relevant to the feature
- When a section doesn't apply, remove it entirely (don't leave as "N/A")

---

## Clarifications

### Session 2025-10-03
- Q: Should visual status indicators (emojis) be configurable or always-on? → A: Always display emojis (same as current shell script)
- Q: Which environment variable should be used for installation directory override? → A: GIT_WORKTREE_MANAGER_HOME (matches current shell script)
- Q: Should all commands require being run from worktree-managed repo, or only some? → A: Only worktree operations (list, new-branch, remove, prune) require worktree-managed repo context
- Q: How should the system handle attempting to create a branch that already exists locally? → A: Prompt user for confirmation to use existing branch; if user agrees to continue and use the branch, check if the branch exists on the remote repo and if not ask the user if they want to continue and push the branch to the remote
- Q: What happens when a user runs setup in a directory that already contains a .bare folder? → A: Detect .bare folder existence early and abort immediately, preventing any actions from taking place
- Q: What should happen when a user tries to create a worktree for a branch that exists remotely but NOT locally? → A: Prompt user for confirmation, then fetch the remote branch and create worktree from it if confirmed
- Q: How should the system handle network failures during clone or fetch operations? → A: Prompt user whether to retry or abort when network failure detected
- Q: What happens when GitHub API rate limits are hit during version checks? → A: Fail with error message and suggest trying again later
- Q: How should the system handle repositories that have renamed their default branch (e.g., from 'master' to 'main')? → A: Prompt user to confirm the detected default branch before proceeding
- Q: How should the system handle interrupted operations (e.g., user presses Ctrl+C during partial clone or worktree creation)? → A: Ask user whether to clean up or preserve partial state before exiting

---

## User Scenarios & Testing *(mandatory)*

### Primary User Story
A developer using git-worktree-manager needs to perform all existing worktree management operations (setup, branch creation, listing, removal, version checking, upgrades) through a command-line interface that maintains complete functional compatibility with the current shell script while providing improved performance, better error handling, and enhanced user experience.

### Acceptance Scenarios
1. **Given** a user currently using the shell script, **When** they switch to the new CLI tool, **Then** all existing commands work with identical syntax and produce equivalent results
2. **Given** a user runs `<tool> <org>/<repo>`, **When** the command executes, **Then** a bare clone is created in .bare and the default branch worktree is set up
3. **Given** a user runs the new branch creation command, **When** specifying a branch name, **Then** the worktree is created and the branch is pushed to the remote
4. **Given** a user runs the list command, **When** executed, **Then** all active worktrees are displayed
5. **Given** a user runs the remove command with --remote flag, **When** executed, **Then** both local worktree/branch and remote branch are deleted
6. **Given** a user runs the prune command, **When** executed, **Then** stale worktree references are cleaned up
7. **Given** a user runs the version command, **When** executed, **Then** current version and upgrade availability status are shown
8. **Given** a user runs the upgrade command when a newer version exists, **When** executed, **Then** the tool downloads and installs the latest version
9. **Given** a user runs any destructive command with --dry-run, **When** executed, **Then** the tool shows what would happen without making changes
10. **Given** a user runs help command or uses -h flag, **When** executed, **Then** comprehensive usage documentation is displayed

### Edge Cases
- When a user tries to create a branch that already exists locally, the system prompts for confirmation to use the existing branch. If confirmed, the system checks if the branch exists on the remote repository. If the branch does not exist remotely, the system prompts the user for confirmation to push the branch to the remote.
- When a user tries to create a worktree for a branch that exists remotely but not locally, the system prompts the user for confirmation, then fetches the remote branch and creates the worktree from it if the user confirms.
- When network failures occur during clone or fetch operations, the system prompts the user whether to retry the operation or abort, allowing the user to decide based on their network situation.
- When GitHub API rate limits are hit during version checks, the system fails with a clear error message and suggests the user try again later when the rate limit resets.
- When repositories have renamed their default branch (e.g., from 'master' to 'main'), the system automatically detects the current default branch using git symbolic-ref and prompts the user to confirm the detected default branch before proceeding with setup.
- When a user runs setup in a directory that already contains a .bare folder, the system detects this early and aborts immediately with an error message requiring manual cleanup before proceeding.
- When operations are interrupted (e.g., user presses Ctrl+C during partial clone or worktree creation), the system catches the interruption signal and asks the user whether to clean up the partial state automatically or preserve it for manual inspection.
- What happens when a user tries to remove a worktree that doesn't exist?
- How does the system handle permission errors during installation or upgrade?
- What happens when semantic version comparison encounters malformed version strings?

## Requirements *(mandatory)*

### Functional Requirements

**Core Workflow Compatibility**
- **FR-001**: System MUST support full repository setup using `<org>/<repo>` shorthand notation
- **FR-002**: System MUST support full repository setup using SSH URL format `git@github.com:<org>/<repo>.git`
- **FR-002a**: System MUST validate that `.bare` directory does not already exist before beginning setup, and abort immediately with clear error message if it exists
- **FR-003**: System MUST create bare clone in `.bare` directory within project root
- **FR-004**: System MUST create `.git` file pointing to `.bare` directory
- **FR-005**: System MUST automatically detect default branch from remote repository
- **FR-005a**: System MUST prompt user to confirm the detected default branch before proceeding with initial worktree creation
- **FR-006**: System MUST create initial worktree for the default branch
- **FR-007**: System MUST configure git settings (push.default, branch.autosetupmerge, branch.autosetuprebase, remote fetch refspec)

**Branch and Worktree Management**
- **FR-008**: System MUST allow users to create new branch worktrees with command `--new-branch <branch-name>`
- **FR-009**: System MUST allow users to specify optional base branch when creating new branch worktrees
- **FR-010**: System MUST fetch latest changes from remote before creating worktrees
- **FR-011**: System MUST create new branches from specified base branch when branch doesn't exist locally
- **FR-012**: System MUST create worktrees from existing local branches when branch exists
- **FR-013**: System MUST automatically push new branches to remote with tracking setup
- **FR-013a**: System MUST prompt user for confirmation when attempting to create a branch that already exists locally
- **FR-013b**: System MUST check remote branch existence after user confirms using existing local branch, and prompt for confirmation to push if branch does not exist remotely
- **FR-013c**: System MUST prompt user for confirmation when attempting to create a worktree for a branch that exists remotely but not locally, then fetch the remote branch and create worktree from it if user confirms
- **FR-014**: System MUST allow users to list all active worktrees with `--list` command
- **FR-015**: System MUST allow users to remove worktrees and local branches with `--remove <branch>` command
- **FR-016**: System MUST allow users to delete remote branches when `--remote` flag is provided with remove command
- **FR-017**: System MUST allow users to prune stale worktree references with `--prune` command
- **FR-018**: System MUST validate worktree existence before attempting removal
- **FR-019**: System MUST validate branch existence before attempting deletion

**Version Management**
- **FR-020**: System MUST display current version with `--version` command
- **FR-021**: System MUST check for newer versions on GitHub when displaying version
- **FR-021a**: System MUST fail with clear error message and suggest trying again later when GitHub API rate limits are encountered during version checks
- **FR-022**: System MUST compare versions using semantic versioning 2.0.0 rules
- **FR-023**: System MUST handle prerelease identifiers correctly in version comparison (numeric vs alphanumeric precedence)
- **FR-024**: System MUST ignore build metadata when comparing versions
- **FR-025**: System MUST support self-upgrade with `--upgrade` command
- **FR-026**: System MUST download executable, README, VERSION, and LICENSE files during upgrade
- **FR-027**: System MUST verify version is newer before replacing existing installation
- **FR-028**: System MUST preserve executable permissions after upgrade

**User Experience**
- **FR-029**: System MUST provide help documentation with `--help` or `-h` flag
- **FR-030**: System MUST support `--dry-run` flag for all destructive operations
- **FR-031**: System MUST show preview of actions when in dry-run mode without executing them
- **FR-032**: System MUST provide clear error messages when operations fail
- **FR-033**: System MUST provide actionable guidance in error messages
- **FR-034**: System MUST show progress indicators for long-running operations (clone, fetch)
- **FR-035**: System MUST always display emoji visual status indicators (success, error, warning, info) matching current shell script behavior
- **FR-036**: System MUST validate repository format before attempting operations
- **FR-037**: System MUST validate command arguments and provide usage instructions when invalid

**Installation and Distribution**
- **FR-038**: System MUST install to a configurable directory location
- **FR-039**: System MUST respect GIT_WORKTREE_MANAGER_HOME environment variable for installation directory override, defaulting to $HOME/.git-worktree-manager if not set
- **FR-040**: System MUST fetch updates from GitHub repository main branch
- **FR-041**: System MUST handle download failures gracefully with clear error messages

**Safety and Data Integrity**
- **FR-042**: System MUST verify git command availability before executing operations
- **FR-043**: System MUST handle network failures gracefully without corrupting repository state
- **FR-043a**: System MUST prompt user whether to retry or abort when network failure is detected during clone or fetch operations
- **FR-044**: System MUST prevent data loss when operations are interrupted
- **FR-044a**: System MUST catch interruption signals (SIGINT, SIGTERM) and prompt user whether to clean up partial state automatically or preserve it for manual inspection
- **FR-045**: System MUST verify successful download before replacing existing files during upgrade
- **FR-046**: System MUST validate worktree-managed repo context only for worktree operations (list, new-branch, remove, prune), while setup, version, upgrade, and help commands work from any directory

### Key Entities

- **Repository**: Represents a GitHub repository with organization name, repository name, SSH URL, default branch, and local bare clone path
- **Worktree**: Represents a git worktree with associated branch name, directory path, and tracking status
- **Branch**: Represents a git branch with name, base branch reference, existence status (local/remote), and tracking configuration
- **Version**: Represents semantic version with major version number, minor version number, patch version number, optional prerelease identifiers, and optional build metadata
- **Command**: Represents a CLI operation with command name, required arguments, optional flags, and dry-run status
- **Installation**: Represents the tool installation with installation directory path, version, and component files (executable, README, VERSION, LICENSE)

---

## Review & Acceptance Checklist
*GATE: Automated checks run during main() execution*

### Content Quality
- [ ] No implementation details (languages, frameworks, APIs)
- [ ] Focused on user value and business needs
- [ ] Written for non-technical stakeholders
- [ ] All mandatory sections completed

### Requirement Completeness
- [ ] No [NEEDS CLARIFICATION] markers remain
- [ ] Requirements are testable and unambiguous
- [ ] Success criteria are measurable
- [ ] Scope is clearly bounded
- [ ] Dependencies and assumptions identified

---

## Execution Status
*Updated by main() during processing*

- [x] User description parsed
- [x] Key concepts extracted
- [x] Ambiguities marked
- [x] User scenarios defined
- [x] Requirements generated
- [x] Entities identified
- [ ] Review checklist passed

---
