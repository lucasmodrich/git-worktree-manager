# Git Worktree Manager Constitution

## Core Principles

### I. Single-File Portability (NON-NEGOTIABLE)
**All functionality MUST remain in a single, self-contained bash script.**

- The entire application lives in `git-worktree-manager.sh` (one file, ~500 lines)
- Zero external dependencies beyond: bash 4.0+, git 2.0+, curl, standard Unix tools (grep, sed, etc.)
- Script MUST be distributable as a single file via curl/wget
- Installation directory MUST be configurable via `GIT_WORKTREE_MANAGER_HOME` environment variable (default: `$HOME/.git-worktree-manager`)
- No libraries, no modules, no external config files required for core functionality
- Self-updating capability MUST be maintained (fetches from GitHub main branch)

### II. Bash Best Practices (MANDATORY)
**Code quality and safety standards are non-negotiable.**

- **Error Handling**: `set -e` required at script start; all functions must handle errors gracefully
- **Variable Quoting**: Always use `"$variable"` to prevent word splitting
- **Conditionals**: Use `[[ ]]` for tests (bash-specific), NOT `[ ]` (POSIX)
- **Command Substitution**: Use `$()` syntax, NOT backticks
- **Function Naming**: Descriptive names with underscores: `create_new_branch_worktree()`
- **Variable Scope**: Use `local` keyword for function variables; ALL_CAPS for globals
- **Arrays**: Declare with `local -a array_name`
- **IFS Safety**: Save and restore IFS when modifying: `local oldIFS=$IFS; IFS=...; IFS=$oldIFS`
- **Comments**: Use `# ---` for section headers, `#` for inline explanations

### III. Test-First Development (NON-NEGOTIABLE)
**All new functionality requires test coverage BEFORE implementation.**

- **TDD Workflow**: Write test → Get user approval → Verify test fails → Implement → Verify test passes
- Test files in `tests/` directory using bash with `set -euo pipefail`
- All tests MUST pass (100% success rate) before merge
- Test naming: `*_tests.sh` pattern
- Unified test runner: `./tests/run_all_tests.sh` runs all suites
- Output format: Emoji indicators (🧪 ▶️ ✅ ❌ 📊) for readability
- Test coverage requirements:
  - Core functionality (version comparison, branch operations, worktree management)
  - Input validation and sanitization
  - Dry-run mode verification
  - Error handling paths

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
  - `SCRIPT_VERSION` at line 5 MUST be updated by semantic-release ONLY
  - `VERSION` file maintained by CI/CD
  - `CHANGELOG.md` auto-generated from commits
- **Version Bumping**:
  - MAJOR: Breaking changes (backward-incompatible CLI changes)
  - MINOR: New features (new flags, commands)
  - PATCH: Bug fixes, documentation, refactoring
- **Release Assets**: GitHub releases MUST include:
  - `git-worktree-manager.sh` (primary asset)
  - `README.md`, `LICENSE`, `VERSION`
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
- **Primary**: Linux (Ubuntu/Debian tested)
- **Secondary**: macOS (should work, best-effort support)
- **Unsupported**: Windows (unless WSL/Git Bash)

### Bash Version Requirements
- **Minimum**: Bash 4.0 (for associative arrays, `[[` conditionals)
- **Tested**: Ubuntu default bash (5.x)
- **Not Compatible**: POSIX sh, dash, zsh (bash-specific features used)

### Git Requirements
- **Minimum**: Git 2.0 (for worktree support)
- **Recommended**: Git 2.5+ (improved worktree commands)

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
- Complexity MUST be justified against Principle I (single-file portability)

---

**Version**: 1.0.0
**Ratified**: 2025-10-03
**Last Amended**: 2025-10-03
