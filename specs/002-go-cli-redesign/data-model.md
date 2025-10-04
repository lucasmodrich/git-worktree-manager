# Data Model: Go CLI Redesign

**Feature**: 002-go-cli-redesign
**Date**: 2025-10-03
**Source**: Extracted from spec.md Key Entities section

---

## Core Entities

### 1. Repository

**Purpose**: Represents a GitHub repository configuration and local state.

**Fields**:
```go
type Repository struct {
    Organization string // e.g., "lucasmodrich"
    Name         string // e.g., "git-worktree-manager"
    SSHUrl       string // e.g., "git@github.com:lucasmodrich/git-worktree-manager.git"
    DefaultBranch string // e.g., "main" (detected from remote)
    BareClonePath string // Absolute path to .bare directory
}
```

**Validation Rules**:
- Organization must match pattern `[a-zA-Z0-9._-]+`
- Name must match pattern `[a-zA-Z0-9._-]+`
- SSHUrl must match format `git@github.com:<org>/<repo>.git`
- DefaultBranch must exist in remote repository
- BareClonePath must be absolute and end with `.bare`

**State Transitions**: N/A (immutable after creation)

**Relationships**:
- One Repository has many Worktrees (via shared .bare)

---

### 2. Worktree

**Purpose**: Represents a git worktree linked to a branch.

**Fields**:
```go
type Worktree struct {
    Path         string // Absolute path to worktree directory
    BranchName   string // Associated branch name
    DirectoryExists bool // Whether directory exists on filesystem
    TrackingStatus string // "tracking", "untracked", "diverged"
}
```

**Validation Rules**:
- Path must be absolute
- BranchName must be valid git ref name
- Path must not be inside another worktree
- Path must not be the bare repository path

**State Transitions**:
```
[Created] → [Exists, Tracking Remote]
         ↓
      [Removed Locally]
         ↓
      [Removed from Git] → [Pruned]
```

**Relationships**:
- Each Worktree belongs to one Repository
- Each Worktree is associated with one Branch

---

### 3. Branch

**Purpose**: Represents a git branch with local and remote status.

**Fields**:
```go
type Branch struct {
    Name           string // Branch name (e.g., "feature/new-ui")
    BaseBranch     string // Branch this was created from (e.g., "main")
    ExistsLocal    bool   // Present in local repository
    ExistsRemote   bool   // Present on origin
    TrackingRemote bool   // Local branch tracks remote
    HasWorktree    bool   // Has associated worktree
}
```

**Validation Rules**:
- Name must be valid git ref name (no spaces, special chars)
- BaseBranch must exist if creating new branch
- If ExistsRemote, must be fetchable from origin
- TrackingRemote requires both ExistsLocal and ExistsRemote

**State Transitions**:
```
[New] → [Local Only] → [Pushed] → [Tracking Remote]
                    ↓
                 [Deleted Locally]
                    ↓
                 [Deleted Remotely] → [Pruned]
```

**Relationships**:
- Each Branch may have one Worktree
- Each Branch belongs to one Repository

---

### 4. Version

**Purpose**: Represents a semantic version for comparison and upgrade logic.

**Fields**:
```go
type Version struct {
    Major      int      // Major version number
    Minor      int      // Minor version number
    Patch      int      // Patch version number
    Prerelease []string // Prerelease identifiers (split on '.')
    Build      string   // Build metadata (ignored in comparison)
    Original   string   // Original version string (e.g., "v1.2.3-beta.1+build123")
}
```

**Validation Rules**:
- Major, Minor, Patch must be >= 0
- Prerelease identifiers must not be empty strings
- Prerelease numeric identifiers must parse as valid integers
- Build metadata may be any string (not validated)
- Original must parse according to semver 2.0.0 spec

**State Transitions**: N/A (immutable value object)

**Relationships**: N/A (value object, no entity relationships)

**Comparison Logic**:
1. Compare Major, then Minor, then Patch numerically
2. Release (no prerelease) > Prerelease
3. Compare prerelease identifiers left-to-right:
   - Numeric < Alphanumeric (within same position)
   - Numeric identifiers compared numerically
   - Alphanumeric identifiers compared lexically (ASCII)
4. Shorter prerelease list < longer (if all previous match)
5. Build metadata ignored

---

### 5. Command

**Purpose**: Represents a CLI operation with its configuration.

**Fields**:
```go
type Command struct {
    Name       string            // Command name (e.g., "setup", "new-branch")
    Args       []string          // Positional arguments
    Flags      map[string]string // Flag key-value pairs
    DryRun     bool              // Dry-run mode enabled
    GlobalFlags GlobalFlags      // Global flags (verbosity, etc.)
}

type GlobalFlags struct {
    DryRun  bool
    Verbose bool
    NoEmoji bool // Future: if emoji made configurable
}
```

**Validation Rules**:
- Name must match one of the defined Cobra commands
- Args length must match command requirements
- Required flags must be present in Flags map
- DryRun affects execution but not validation

**State Transitions**:
```
[Parsed] → [Validated] → [Executed] → [Success|Error]
```

**Relationships**: N/A (ephemeral, exists only during command execution)

---

### 6. Installation

**Purpose**: Represents the tool installation metadata and file locations.

**Fields**:
```go
type Installation struct {
    BaseDir          string // Installation base directory
    BinaryPath       string // Path to executable
    Version          Version // Installed version
    ComponentFiles   map[string]string // "README": path, "LICENSE": path, etc.
}
```

**Validation Rules**:
- BaseDir must be absolute path
- BaseDir defaults to `$HOME/.git-worktree-manager` if GIT_WORKTREE_MANAGER_HOME unset
- BinaryPath must be executable (Unix: mode 0755, Windows: .exe extension)
- Version must be parseable semantic version

**State Transitions**:
```
[Fresh Install] → [Installed]
               ↓
            [Upgrade Available]
               ↓
            [Downloading]
               ↓
            [Verifying Checksum]
               ↓
            [Upgraded]
```

**Relationships**:
- Installation is global singleton (one per system/user)
- Installation manages Version lifecycle

---

## Entity Relationships Diagram

```
Repository (1) ──────< (N) Worktree
    │                       │
    │                       │
    └────< (N) Branch (1) ──┘
                │
                │ (uses)
                ↓
             Version

Installation (1) ── (has) → Version

Command (ephemeral, no persistent relationships)
```

---

## Derived Data & Computed Fields

### Repository
- `ProjectRoot()`: Directory containing .bare (parent of BareClonePath)
- `IsSetup()`: Whether .bare directory exists and is valid

### Worktree
- `IsStale()`: Directory missing but git worktree list shows entry (needs pruning)
- `BranchRef()`: Full ref path ("refs/heads/<BranchName>")

### Branch
- `NeedsPush()`: ExistsLocal && !ExistsRemote && user wants to push
- `Status()`: "local-only", "remote-only", "tracking", "diverged"

### Version
- `GreaterThan(other *Version) bool`: Implements semver 2.0.0 comparison
- `String()`: Returns Original for display

### Installation
- `NeedsUpgrade(latestVersion *Version) bool`: Compare installed vs latest
- `ExecutableName()`: "git-worktree-manager" (Unix) or "git-worktree-manager.exe" (Windows)

---

## Validation Summary

| Entity | Required Fields | Unique Constraints | Foreign Keys |
|--------|----------------|-------------------|--------------|
| Repository | Organization, Name, SSHUrl | SSHUrl must be unique per setup | N/A |
| Worktree | Path, BranchName | Path must be unique | BranchName → Branch.Name |
| Branch | Name | Name must be unique per repo | BaseBranch → Branch.Name (optional) |
| Version | Major, Minor, Patch | N/A (value object) | N/A |
| Command | Name, Args | N/A (ephemeral) | N/A |
| Installation | BaseDir, BinaryPath | BaseDir per user/system | N/A |

---

## Implementation Notes

### Persistence
- **None Required**: All state persists in git itself (.git/worktrees/, refs/, config)
- Installation.ComponentFiles tracked in filesystem only
- No database, no JSON files, no custom config beyond environment variables

### Concurrency
- Git operations are process-isolated (safe for concurrent runs in different repos)
- Upgrade operations should use file locking (prevent concurrent self-replacement)
- Worktree creation should check existence atomically (race between check and create acceptable - git will error)

### Error Recovery
- Failed clone: Remove partial .bare directory
- Failed worktree create: Remove directory, don't leave broken state
- Failed upgrade: Keep old binary if new binary verification fails

---

## Testing Implications

### Unit Tests
- Version comparison: Table-driven with all semver edge cases
- Path validation: Test Windows vs Unix path handling
- Flag parsing: Verify all command/flag combinations

### Integration Tests
- Repository setup: Verify .bare creation, default branch detection
- Worktree lifecycle: Create, list, remove, prune full workflow
- Upgrade: Mock GitHub API, verify atomic replacement

### Contract Tests
- CLI output format: Compare Go output to Bash output for same inputs
- Exit codes: Verify 0/1/2 codes match Bash version
- Error messages: Ensure actionable guidance present
