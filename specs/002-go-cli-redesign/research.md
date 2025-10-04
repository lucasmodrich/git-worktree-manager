# Research: Go CLI Redesign

**Feature**: 002-go-cli-redesign
**Date**: 2025-10-03
**Purpose**: Resolve technical unknowns and establish best practices for Go CLI implementation

---

## 1. Cobra CLI Framework

### Decision
Use `github.com/spf13/cobra` v1.8+ for CLI structure and command parsing.

### Rationale
- **Industry Standard**: Most popular Go CLI framework (used by kubectl, hugo, gh)
- **Rich Feature Set**: Subcommands, flags, aliases, help generation, shell completion
- **Minimal Dependencies**: Only requires pflag (flags) and optionally viper (config)
- **Compatible with Bash CLI**: Can exactly replicate existing flag structure
- **Well Documented**: Extensive documentation and examples

### Alternatives Considered
- **urfave/cli**: Simpler but less powerful flag handling, harder to match Bash interface exactly
- **Standard flag package**: Too low-level, would require significant boilerplate
- **kingpin**: Less actively maintained, smaller ecosystem

### Implementation Notes
- Use `cobra-cli` generator to bootstrap command structure
- One command file per CLI command (setup, branch, list, remove, prune, version, upgrade)
- Global `--dry-run` flag defined in root command
- Persistent flags for common options across commands

---

## 2. Semantic Version Comparison

### Decision
Port existing Bash `version_gt()` function logic to Go, implementing full semver 2.0.0 specification.

### Rationale
- **Proven Logic**: Existing Bash implementation handles all semver edge cases correctly
- **Test Coverage**: Existing tests in `tests/version_compare_tests.sh` can validate Go port
- **No External Dependency**: Simple enough to implement without third-party library
- **Exact Compatibility**: Ensures version checks produce identical results to Bash version

### Alternatives Considered
- **github.com/Masterminds/semver**: Popular library but adds dependency for simple use case
- **golang.org/x/mod/semver**: Limited to Go module versions, doesn't handle all semver 2.0.0 features
- **github.com/blang/semver**: Full semver 2.0.0 but larger dependency

### Implementation Notes
```go
type Version struct {
    Major      int
    Minor      int
    Patch      int
    Prerelease []string  // Split on '.'
    Build      string    // Ignored in comparison
}

func ParseVersion(v string) (*Version, error)
func (v *Version) GreaterThan(other *Version) bool
```

- Strip leading 'v' and build metadata ('+...')
- Split into main version and prerelease ('-...')
- Compare major.minor.patch numerically
- Handle prerelease precedence: release > prerelease, numeric < alphanumeric
- Port unit tests from Bash to Go table-driven tests

---

## 3. Git Command Execution

### Decision
Use `os/exec` package to execute git commands as subprocesses, wrapping in typed functions.

### Rationale
- **No Git Library Needed**: Git itself is more reliable than Go git libraries (go-git, git2go)
- **Identical Behavior**: Produces exact same results as Bash script's git calls
- **Error Handling**: Can capture stdout/stderr separately for better error messages
- **Progress Integration**: Can stream output for progress indicators

### Alternatives Considered
- **go-git/go-git**: Pure Go git implementation, but incomplete worktree support, risk of behavioral differences
- **libgit2/git2go**: CGo dependency, compilation complexity, cross-platform challenges

### Implementation Notes
```go
package git

type Client struct {
    workDir string
}

func NewClient(workDir string) *Client
func (c *Client) Clone(url, target string, bare bool) error
func (c *Client) WorktreeAdd(path, branch string, track bool) error
func (c *Client) WorktreeList() ([]Worktree, error)
func (c *Client) WorktreeRemove(path string) error
func (c *Client) BranchExists(name string, remote bool) (bool, error)
func (c *Client) DetectDefaultBranch() (string, error)
```

- All git commands return structured errors with context
- Capture both stdout and stderr for debugging
- Timeout support for long-running operations
- Dry-run mode: log command without executing

---

## 4. Cross-Platform File Operations

### Decision
Use `io/fs`, `os`, and `path/filepath` from standard library for all file operations.

### Rationale
- **Built-in Cross-Platform**: Works on Linux, macOS, Windows without external dependencies
- **Path Handling**: `filepath` automatically handles OS-specific separators
- **Atomic Operations**: Can use temp files + rename for atomic writes (upgrade safety)

### Implementation Notes
- **Installation Path**: Use `filepath.Join(os.Getenv("GIT_WORKTREE_MANAGER_HOME"), "bin")`
- **Default**: Fall back to `filepath.Join(os.UserHomeDir(), ".git-worktree-manager")`
- **.git File Creation**: Write `gitdir: ./.bare\n` using `os.WriteFile` with 0644 permissions
- **Executable Permissions**: Use `os.Chmod` to set 0755 on downloaded binary

---

## 5. Self-Upgrade Mechanism

### Decision
Download binary from GitHub Releases using `net/http`, verify checksum, replace atomically.

### Rationale
- **GitHub Releases API**: Standard approach for distributing Go binaries
- **SHA256 Verification**: Checksums file ensures download integrity
- **Atomic Replace**: Download to temp → verify → rename prevents partial upgrades

### Implementation Notes
```go
func DownloadLatestRelease(arch, os string) (path string, err error)
func VerifyChecksum(binaryPath, checksumURL string) error
func ReplaceExecutable(newPath string) error
```

- Detect current OS/arch: `runtime.GOOS`, `runtime.GOARCH`
- Fetch from `https://github.com/lucasmodrich/git-worktree-manager/releases/latest/download/git-worktree-manager-{os}-{arch}`
- Verify against `checksums.txt` from same release
- Preserve current executable path: use `os.Executable()` to find self
- Atomic replace: write to `.new`, verify, rename over existing
- Windows: May need to write `.old`, rename new, delete old (can't replace running exe)

---

## 6. Progress Indicators

### Decision
Use emoji output matching Bash version, detect terminal capabilities for safe rendering.

### Rationale
- **Clarification Requirement**: Must always display emojis (Session 2025-10-03)
- **Simple Implementation**: Just string constants, no external libraries
- **Cross-Platform**: UTF-8 emoji support on modern Linux, macOS, Windows Terminal

### Implementation Notes
```go
const (
    Success   = "✅"
    Error     = "❌"
    Info      = "📡"
    Warning   = "⚠️"
    DryRun    = "🔍"
    Progress  = "⏳"
)

func PrintStatus(emoji, message string)
func PrintProgress(current, total int, message string)
```

- For long operations (clone, fetch): Use `fmt.Printf("\r")` to update same line
- For multi-step operations: Print each step with emoji prefix
- stderr for errors, stdout for normal output
- No spinner libraries needed - keep output simple and fast

---

## 7. Interactive Prompts

### Decision
Use `bufio.Scanner` on `os.Stdin` for simple yes/no confirmations.

### Rationale
- **Clarification Requirements**: Must prompt when branch exists (FR-013a, FR-013b)
- **No External Library**: Standard library sufficient for simple confirmations
- **Cross-Platform**: Works on all platforms with stdin available

### Implementation Notes
```go
func PromptYesNo(question string) (bool, error) {
    fmt.Printf("%s [y/N]: ", question)
    scanner := bufio.NewScanner(os.Stdin)
    if !scanner.Scan() {
        return false, scanner.Err()
    }
    answer := strings.ToLower(strings.TrimSpace(scanner.Text()))
    return answer == "y" || answer == "yes", nil
}
```

- Default to "No" for safety
- Handle EOF gracefully (non-interactive environments)
- In dry-run mode: skip prompts, log what would be asked

---

## 8. Error Handling Strategy

### Decision
Use custom error types with context, return errors up call stack, format with actionable guidance.

### Rationale
- **Go Best Practice**: Never panic in production code, return errors
- **Better UX**: Structured errors can include suggestions for resolution
- **Testing**: Easier to test error conditions with typed errors

### Implementation Notes
```go
type GitError struct {
    Operation string
    Command   string
    Stderr    string
    Cause     error
}

type ValidationError struct {
    Field   string
    Value   string
    Message string
}

func (e *GitError) Error() string {
    return fmt.Sprintf("git %s failed: %s\nCommand: %s\nOutput: %s",
        e.Operation, e.Cause, e.Command, e.Stderr)
}
```

- Wrap errors with context using `fmt.Errorf("context: %w", err)`
- Top-level commands format errors with emoji and guidance
- Exit codes: 0 success, 1 user error, 2 system error

---

## 9. Testing Strategy

### Decision
- **Unit Tests**: Table-driven tests for version comparison, validation logic
- **Integration Tests**: Real git operations in temp directories
- **Contract Tests**: Compare Go CLI output to Bash CLI output for same commands

### Rationale
- **Constitution Requirement**: >80% coverage for core packages
- **TDD Workflow**: Write tests first, verify failures, implement, verify pass
- **Compatibility Verification**: Contract tests ensure CLI interface matches exactly

### Implementation Notes
```go
// Unit test example
func TestVersionComparison(t *testing.T) {
    tests := []struct{
        name string
        v1, v2 string
        want bool
    }{
        {"major greater", "2.0.0", "1.9.9", true},
        {"prerelease precedence", "1.0.0", "1.0.0-beta", true},
        // ... more cases from Bash tests
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := VersionGreaterThan(tt.v1, tt.v2)
            if got != tt.want {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}

// Integration test example
func TestSetupWorkflow(t *testing.T) {
    tmpDir := t.TempDir()
    // Create test git repo
    // Run setup command
    // Verify .bare exists
    // Verify default branch worktree exists
}

// Contract test example
func TestCLICompatibility(t *testing.T) {
    // Run same command with Bash and Go
    // Compare output format (ignoring timestamps)
}
```

---

## 10. Build and Release

### Decision
Use GoReleaser for automated multi-platform builds and GitHub Releases.

### Rationale
- **Multi-Platform**: Builds for Linux/macOS/Windows amd64/arm64 automatically
- **GitHub Integration**: Creates releases, uploads assets, generates checksums
- **Semantic Release Compatible**: Works with semantic-release workflow

### Implementation Notes
`.goreleaser.yml`:
```yaml
builds:
  - id: git-worktree-manager
    binary: git-worktree-manager
    main: ./cmd/git-worktree-manager
    ldflags:
      - -s -w
      - -X main.version={{.Version}}
    goos: [linux, darwin, windows]
    goarch: [amd64, arm64]

archives:
  - format: binary

checksum:
  name_template: 'checksums.txt'

release:
  github:
    owner: lucasmodrich
    name: git-worktree-manager
```

- Version embedded via ldflags at build time
- Binary-only distribution (no tar.gz needed for single-file tool)
- Checksums automatically generated
- Integrate with existing semantic-release GitHub Action

---

## Summary

All technical unknowns resolved:
- ✅ Cobra for CLI framework
- ✅ Port semver logic from Bash (no external lib)
- ✅ os/exec for git commands (no go-git)
- ✅ Standard library for file operations
- ✅ GitHub Releases API for self-upgrade
- ✅ Simple emoji output (always enabled)
- ✅ bufio.Scanner for prompts
- ✅ Custom error types with context
- ✅ Table-driven + integration + contract tests
- ✅ GoReleaser for multi-platform builds

**No NEEDS CLARIFICATION remain** - ready for Phase 1 design.

---

## 11. Binary Rename: git-worktree-manager → gwtm

**Date Added**: 2025-10-04
**Context**: User request to rename Go CLI binary to "gwtm" for improved usability

### Decision
Rename the Go CLI binary from `git-worktree-manager` to `gwtm` (Git Worktree Manager).

### Rationale
- **User Experience**: Shorter name (4 chars vs 20 chars) significantly reduces typing friction
- **Industry Standard**: Popular CLI tools use short, memorable names (e.g., `gh`, `kubectl`, `helm`, `hugo`)
- **Memorability**: Acronym `gwtm` is pronounceable and clearly relates to the tool's purpose
- **No Conflicts**: `gwtm` is unique and unlikely to collide with other common tools
- **Muscle Memory**: Users type CLI commands frequently; shorter names improve productivity

### Alternatives Considered
- **Keep `git-worktree-manager`**: Too verbose for frequent command-line use
- **Use `wtm`**: Too short, lacks git context
- **Use `gwm`**: Could confuse with generic "git workflow manager" tools
- **Use `worktree`**: Too generic, might conflict with `git worktree` subcommand

### Implementation Impact

**Build Configuration**:
```makefile
# Makefile
build:
    go build -o gwtm ./cmd/git-worktree-manager
```

``yaml
# .goreleaser.yml
builds:
  - binary: gwtm  # Changed from git-worktree-manager
```

**CI/CD Workflows**:
- `.github/workflows/test.yml`: Update build and test commands to use `gwtm`
- `.github/workflows/release.yml`: GoReleaser will automatically use new binary name

**Documentation**:
- `README.md`: Update all examples to use `gwtm` instead of `git-worktree-manager`
- Installation instructions: Update download URLs (e.g., `gwtm_Linux_x86_64`)
- Migration guide: Provide symlink instructions for backward compatibility

**Backward Compatibility**:
- Semantic version: MINOR bump (e.g., 1.3.0 → 1.4.0) - new feature, not breaking change
- Migration strategy: Document symlink creation for users with hardcoded paths
  ```bash
  ln -s $(which gwtm) /usr/local/bin/git-worktree-manager
  ```
- Bash script: Unaffected (`git-worktree-manager.sh` keeps its name)

**No Impact On**:
- Source code structure (`cmd/git-worktree-manager/` can remain as package path)
- Installation directory (`$HOME/.git-worktree-manager/` unchanged)
- Environment variables (`GIT_WORKTREE_MANAGER_HOME` unchanged)
- CLI interface (commands, flags, behavior all identical)

### Testing Requirements
- Build verification: Confirm `gwtm` binary is created
- Smoke tests: Execute `./gwtm --help`, `./gwtm version`, `./gwtm --dry-run setup`
- Existing integration tests: No changes needed (use built binary path, not hardcoded name)

### Release Assets Pattern
```
# Old pattern:
git-worktree-manager_Linux_x86_64
git-worktree-manager_Darwin_arm64

# New pattern:
gwtm_Linux_x86_64
gwtm_Darwin_arm64
```

---

**UPDATED SUMMARY**: All technical unknowns resolved, including binary rename strategy. Ready for Phase 1 design.
