# AGENTS.md - Development Guidelines

## Repository Context
This repository contains `git-worktree-manager.sh`, a self-updating Bash script for Git worktree management. See `CLAUDE.md` for detailed repository overview and architecture.

## Build/Test Commands
- Run all tests: `./tests/run_all_tests.sh`
- Run single test: `./tests/version_compare_tests.sh` (or other test files)
- No build step - single bash script deployment
- Release: Automated via GitHub Actions on main branch push
- Test script locally: `./git-worktree-manager.sh --help`

## Code Style Guidelines
- **Language**: Bash scripting (requires bash, not sh-compatible)
- **Error handling**: Use `set -e` for strict error handling
- **Variables**: Use `local` for function variables, ALL_CAPS for globals
- **Quoting**: Always quote variables: `"$variable"`, use `$()` for command substitution
- **Functions**: Descriptive names with underscores: `create_new_branch_worktree()`
- **Comments**: Use `# ---` for section headers, `#` for inline comments
- **Conditionals**: Use `[[ ]]` for tests, not `[ ]`
- **Arrays**: Declare with `local -a array_name`

## Architecture
- **Single file**: All functionality in `git-worktree-manager.sh`
- **Version**: Hardcoded at line 5: `SCRIPT_VERSION="x.y.z"`
- **Self-contained**: No external dependencies beyond git and standard Unix tools
- **Installation**: Script installs to `$HOME/.git-worktree-manager/`

## Testing
- Test files in `tests/` directory use bash with `set -euo pipefail`
- Use emoji prefixes for output: 🧪 ▶️ ✅ ❌ 📊
- Test structure: setup, run, verify, cleanup