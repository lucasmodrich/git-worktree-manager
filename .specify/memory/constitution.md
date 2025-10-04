# Git Worktree Manager Constitution

## Core Principles

### I. Multi-Implementation Strategy (MANDATORY)
**The tool MUST be available in multiple implementation forms to serve different user needs.**

**A. Canonical Bash Implementation** (Legacy support - MAINTAINED)
- `git-worktree-manager.sh` remains available as single-file Bash script
- Zero external dependencies beyond: bash 4.0+, git 2.0+, curl, standard Unix tools
- Distributable as a single file via curl/wget
- Maintained for users who prefer shell scripts or have constrained environments
- Feature parity NOT required - may receive only critical bug fixes
- Self-updating capability maintained

**B. Primary Go Implementation** (Current development focus - ACTIVE)
- Go CLI application using Cobra framework for command-line interface
- Compiled binaries for Linux, macOS, and Windows
- All existing features from Bash version MUST be implemented
- Installation directory MUST be configurable via `GIT_WORKTREE_MANAGER_HOME` environment variable (default: `$HOME/.git-worktree-manager`)
- Self-updating capability MUST be maintained (fetches from GitHub releases)
- Superior performance, error handling, and user experience compared to Bash version

**C. Implementation Compatibility**
- Both implementations MUST support identical CLI interface (command compatibility)
- Configuration (environment variables, directory structure) MUST be shared
- Users MUST be able to switch between implementations without workflow changes

### II. Language-Specific Best Practices (MANDATORY)

**A. Bash Best Practices** (for git-worktree-manager.sh)
- **Error Handling**: `set -e` required at script start; all functions must handle errors gracefully
- **Variable Quoting**: Always use `"$variable"` to prevent word splitting
- **Conditionals**: Use `[[ ]]` for tests (bash-specific), NOT `[ ]` (POSIX)
- **Command Substitution**: Use `$()` syntax, NOT backticks
- **Function Naming**: Descriptive names with underscores: `create_new_branch_worktree()`
- **Variable Scope**: Use `local` keyword for function variables; ALL_CAPS for globals
- **Arrays**: Declare with `local -a array_name`
- **IFS Safety**: Save and restore IFS when modifying: `local oldIFS=$IFS; IFS=...; IFS=$oldIFS`
- **Comments**: Use `# ---` for section headers, `#` for inline explanations

**B. Go Best Practices** (for primary implementation)
- **Error Handling**: Return errors, never panic in production code; use custom error types
- **CLI Framework**: Use Cobra for command structure and flag parsing
- **Project Structure**: Follow standard Go layout (cmd/, internal/, pkg/)
- **Dependency Management**: Use Go modules; minimize external dependencies
- **Code Style**: Follow `gofmt` and `golint` standards
- **Testing**: Use standard `testing` package; table-driven tests preferred
- **Naming**: Exported names use PascalCase; unexported use camelCase
- **Comments**: Godoc-style comments for all exported functions/types
- **Concurrency**: Use goroutines and channels appropriately; avoid race conditions

### III. Test-First Development (NON-NEGOTIABLE)
**All new functionality requires test coverage BEFORE implementation.**

**A. Common Requirements (Both Implementations)**
- **TDD Workflow**: Write test → Get user approval → Verify test fails → Implement → Verify test passes
- All tests MUST pass (100% success rate) before merge
- Test coverage requirements:
  - Core functionality (version comparison, branch operations, worktree management)
  - Input validation and sanitization
  - Dry-run mode verification
  - Error handling paths

**B. Bash Implementation Testing**
- Test files in `tests/bash/` directory using bash with `set -euo pipefail`
- Test naming: `*_tests.sh` pattern
- Unified test runner: `./tests/bash/run_all_tests.sh`
- Output format: Emoji indicators (🧪 ▶️ ✅ ❌ 📊) for readability

**C. Go Implementation Testing**
- Test files colocated with source: `*_test.go` pattern
- Use standard `testing` package and table-driven tests
- Test runner: `go test ./...` for all packages
- Coverage requirement: >80% for core packages
- Integration tests in `tests/integration/`
- Contract tests verify CLI compatibility with Bash version

### IV. User Safety & Transparency (MANDATORY)
**Users must be able to preview and understand all operations.**

- **Dry-Run Mode**: ALL destructive operations MUST support `--dry-run` flag
  - Preview actions with `🔍 [DRY-RUN]` prefix
  - No actual changes when dry-run active
  - Show exactly what would happen
- **Input Validation**: MUST prevent command injection and invalid inputs
  - Validate GitHub repo format (`org/repo` pattern)
  - Sanitize all user-provided strings
  - Clear error messages for invalid inputs
- **User Feedback**:
  - Progress indicators for long operations (fetching, cloning)
  - Emoji-prefixed messages for status (✅ ❌ 📡 🌱 ☁️ 🗑)
  - Actionable error messages (tell user HOW to fix)
- **Safe Defaults**: Conservative behavior unless explicitly overridden
  - Local-only deletion by default (`--remote` flag required for remote branch deletion)
  - Confirmation for destructive operations

### V. Semantic Release Compliance (MANDATORY)
**Version management is fully automated and MUST NOT be manual.**

- **Conventional Commits**: Required format: `type(scope): description`
  - Types: `feat`, `fix`, `docs`, `chore`, `test`, `refactor`
  - Breaking changes: Include `BREAKING CHANGE:` in commit body OR `!` after type
- **Automated Versioning**:
  - Bash: `SCRIPT_VERSION` at line 5 updated by semantic-release
  - Go: Version embedded at build time via `-ldflags`
  - `VERSION` file maintained by CI/CD
  - `CHANGELOG.md` auto-generated from commits
- **Version Bumping**:
  - MAJOR: Breaking changes (backward-incompatible CLI changes)
  - MINOR: New features (new flags, commands)
  - PATCH: Bug fixes, documentation, refactoring
- **Release Assets**: GitHub releases MUST include:
  - `git-worktree-manager.sh` (Bash implementation)
  - Go binaries: `git-worktree-manager-linux-amd64`, `git-worktree-manager-darwin-amd64`, `git-worktree-manager-windows-amd64.exe`
  - Go binaries with ARM support: `git-worktree-manager-linux-arm64`, `git-worktree-manager-darwin-arm64`
  - `README.md`, `LICENSE`, `VERSION`
  - Checksums file: `checksums.txt` (SHA256 for all binaries)
  - `release-package.tar.gz` (full package)
- **Branch Strategy**:
  - `main`: Stable releases
  - `dev`: Beta prereleases (`x.y.z-beta.n`)

### VI. Backward Compatibility (MANDATORY)
**Existing users MUST NOT break on upgrades.**

- **CLI Interface Stability**:
  - Existing flags/commands MUST continue to work
  - New flags are additive, not replacements
  - Deprecation warnings required 1 major version before removal
- **Breaking Changes Protocol**:
  - Requires MAJOR version bump
  - Must be documented in CHANGELOG.md
  - Migration guide required in release notes
  - Consider providing compatibility shims
- **Self-Update Safety**:
  - Upgrade MUST preserve user's installation directory
  - Version check before replacing script
  - Atomic download (temp file → move)

## Development Workflow

### Feature Development Process
1. **Branch Naming**: Use semantic branch names: `feature/description`, `fix/issue`, `docs/topic`
2. **Specify Integration**: For complex features:
   - Create feature branch: `###-feature-name` format (e.g., `001-dry-run-support`)
   - Run `/specify "feature description"` to create spec
   - Run `/plan` to generate implementation plan
   - Run `/tasks` to break down work
   - Run `/analyze` to validate consistency
3. **Implementation**: Follow Test-First principle (Principle III)
4. **Testing**: Run `./tests/run_all_tests.sh` before commit
5. **Commit**: Use conventional commit format (Principle V)

### Code Review Gates
- [ ] All tests passing (100% success rate)
- [ ] Dry-run support for destructive operations
- [ ] Input validation for user-provided data
- [ ] Error handling with actionable messages
- [ ] Backward compatibility maintained
- [ ] Documentation updated (README.md, --help text)
- [ ] Conventional commit message format

### Quality Standards
- **Performance**: Operations should complete in <5 seconds for typical repos
- **Error Recovery**: Script should handle network failures gracefully
- **Logging**: Important operations logged to stdout; errors to stderr
- **Exit Codes**:
  - `0` for success
  - `1` for user errors (invalid input, missing requirements)
  - `2` for system errors (network failure, git errors)

## Technical Constraints

### Supported Platforms
**Bash Implementation:**
- **Primary**: Linux (Ubuntu/Debian tested)
- **Secondary**: macOS (best-effort support)
- **Limited**: Windows (WSL/Git Bash only)

**Go Implementation:**
- **Primary**: Linux (amd64, arm64)
- **Primary**: macOS (amd64/Intel, arm64/Apple Silicon)
- **Primary**: Windows (amd64)

### Language Version Requirements
**Bash:**
- **Minimum**: Bash 4.0 (for associative arrays, `[[` conditionals)
- **Tested**: Ubuntu default bash (5.x)
- **Not Compatible**: POSIX sh, dash, zsh (bash-specific features used)

**Go:**
- **Minimum**: Go 1.21 (for standard library features)
- **Recommended**: Go 1.22+ (latest stable)
- **Build**: Cross-compilation for all supported platforms

### Git Requirements
- **Minimum**: Git 2.0 (for worktree support)
- **Recommended**: Git 2.5+ (improved worktree commands)
- **Required for both implementations**: Git must be available in PATH

### Security Requirements
- **No Secrets**: Never commit credentials, tokens, or sensitive data
- **Input Sanitization**: All user inputs MUST be validated/sanitized
- **Safe Downloads**: HTTPS only; verify downloads complete successfully
- **Command Injection**: Use proper quoting; avoid `eval`
- **Permissions**: Script installed with user permissions only (no sudo required)

## Governance

### Constitutional Authority
- **Non-Negotiable**: Principles I-VI supersede all other preferences or conventions
- **Constitution First**: When in doubt, consult this document
- **Amendments**: Changes to constitution require:
  1. Documented justification
  2. Impact analysis on existing code
  3. Migration plan if breaking changes
  4. Update version and amendment date below

### Development Guidance
- **Primary Reference**: `CLAUDE.md` for repository overview and architecture
- **Agent Guidelines**: `AGENTS.md` for coding standards and patterns
- **Project Status**: `PROJECT_TODO.md` for review findings and progress

### Compliance Verification
- All PRs MUST verify adherence to constitutional principles
- CI/CD MUST enforce test coverage and quality gates
- Breaking changes MUST be explicitly justified against Principle VI
- Multi-implementation compatibility MUST be maintained per Principle I

---

## Amendment History

### Amendment 1 (Version 2.0.0) - 2025-10-03
**Rationale**: Enable Go CLI implementation to provide superior performance, error handling, and user experience while maintaining Bash version for compatibility.

**Impact Analysis**:
- **Existing Code**: Bash implementation (`git-worktree-manager.sh`) continues to work unchanged
- **New Code**: Go implementation provides same CLI interface with improved internals
- **Users**: Can choose implementation; both share configuration and workflow
- **Migration**: Optional - users can continue using Bash version indefinitely

**Changes**:
1. **Principle I**: Changed from "Single-File Portability" to "Multi-Implementation Strategy"
   - Added support for Go implementation as primary development focus
   - Maintained Bash implementation for legacy support
   - Required CLI interface compatibility between implementations

2. **Principle II**: Changed from "Bash Best Practices" to "Language-Specific Best Practices"
   - Separated Bash and Go coding standards
   - Added Go-specific requirements (Cobra, standard layout, etc.)

3. **Principle III**: Updated "Test-First Development" for both languages
   - Added Go testing requirements (>80% coverage, table-driven tests)
   - Separated test directory structure for each implementation

4. **Principle V**: Updated "Semantic Release" for multi-platform binaries
   - Added Go binary assets for Linux/macOS/Windows (amd64/arm64)
   - Added checksum file requirement
   - Version embedding via ldflags for Go

5. **Technical Constraints**: Added Go language requirements
   - Go 1.21+ minimum version
   - Cross-compilation support
   - Expanded platform support (native Windows)

**Justification**: The single-file Bash constraint was appropriate for initial development but limits:
- Performance (compiled Go vs interpreted Bash)
- Error handling (Go's error system vs Bash error handling)
- Cross-platform support (Windows native support)
- Maintainability (Go's type system and testing infrastructure)
- User experience (better progress indicators, interactive prompts)

The multi-implementation approach preserves existing investment while enabling modernization.

---

**Version**: 2.0.0
**Ratified**: 2025-10-03
**Last Amended**: 2025-10-03 (Amendment 1 - Go CLI Implementation)
