# Contributing to gwtm

Thank you for your interest in contributing. This document covers everything you need to set up a development environment and submit a pull request.

---

## Table of Contents

- [Prerequisites](#prerequisites)
- [Setting Up Your Environment](#setting-up-your-environment)
- [Project Structure](#project-structure)
- [Development Workflow](#development-workflow)
- [Commit Message Convention](#commit-message-convention)
- [Submitting a Pull Request](#submitting-a-pull-request)
- [Release Process](#release-process)

---

## Prerequisites

| Tool | Minimum Version | Purpose |
|------|----------------|---------|
| [Go](https://go.dev/dl/) | 1.25.1 | Build and test |
| [Git](https://git-scm.com/) | 2.5 | Worktree support required for integration tests |
| [make](https://www.gnu.org/software/make/) | any | Convenience targets (optional) |
| [golangci-lint](https://golangci-lint.run/usage/install/) | any | Linting (optional, but recommended) |

---

## Setting Up Your Environment

### 1. Fork and clone

Fork the repository on GitHub, then clone your fork:

```bash
git clone git@github.com:<your-username>/git-worktree-manager.git
cd git-worktree-manager
```

### 2. Install the commit message hook

The repository enforces [Conventional Commits](#commit-message-convention) via a Git hook. Install it with:

```bash
git config core.hooksPath .githooks
```

This will reject any commit whose message does not follow the required format. See [Commit Message Convention](#commit-message-convention) for details.

### 3. Download dependencies

```bash
go mod download
```

### 4. Build and verify

```bash
make build          # produces ./gwtm in the project root
./gwtm --help       # confirm the binary works
```

---

## Project Structure

```
.
├── cmd/
│   └── git-worktree-manager/
│       └── main.go          # Binary entry point; injects version via ldflags
├── internal/
│   ├── commands/            # One file per CLI subcommand (Cobra)
│   │   ├── root.go          # Root command and global flags (--dry-run, --version)
│   │   ├── setup.go         # gwtm setup
│   │   ├── branch.go        # gwtm new-branch
│   │   ├── list.go          # gwtm list
│   │   ├── remove.go        # gwtm remove
│   │   ├── prune.go         # gwtm prune
│   │   ├── version.go       # gwtm version
│   │   ├── upgrade.go       # gwtm upgrade
│   │   └── utils.go         # Shared helpers (findWorktreeRoot)
│   ├── git/                 # Git client wrapper around exec.Command
│   │   ├── client.go        # ExecGit, dry-run support
│   │   ├── branch.go        # Branch CRUD
│   │   ├── remote.go        # Clone, fetch, push, DetectDefaultBranch
│   │   ├── worktree.go      # Worktree add/list/remove/prune
│   │   └── config.go        # git config helpers
│   ├── config/              # Installation directory and binary path resolution
│   ├── ui/                  # Output formatting (stdout/stderr, dry-run, errors)
│   └── version/             # Semver parsing and self-upgrade logic
├── .github/
│   └── workflows/
│       ├── test.yml         # Runs on every PR and push
│       ├── release.yml      # Runs semantic-release on merge to main
│       └── goreleaser.yml   # Builds binaries when a tag is pushed
├── .githooks/
│   └── commit-msg           # Enforces Conventional Commits format
├── .goreleaser.yml          # Cross-platform binary build config
└── release.config.js        # semantic-release configuration
```

### Key design decisions

- **`internal/git`** is a thin wrapper around `exec.Command("git", ...)`. It does not use any Git library. All methods respect the `DryRun` flag — in dry-run mode they log what would be executed without running anything.
- **`internal/commands`** contains only CLI glue — argument parsing, user prompts, and calling into `internal/git`. Business logic lives in the `git` package.
- **Error messages** always go to stderr via `ui.PrintError`. Actionable guidance is printed alongside every error.
- **`findWorktreeRoot()`** walks up from the current directory to locate the repo root, so all commands work from any subdirectory.

---

## Development Workflow

### Build

```bash
make build
# or
go build -o gwtm ./cmd/git-worktree-manager
```

### Test

```bash
make test
# or
go test ./...

# With race detection and coverage (matches CI)
go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
```

Tests use `t.TempDir()` and local bare repositories — no network access required.

### Format

```bash
make fmt
# or
gofmt -s -w .
```

### Lint

```bash
make lint
# or (if golangci-lint is installed)
golangci-lint run
```

### Manual smoke test with dry-run

The `--dry-run` flag lets you exercise any command without touching the filesystem or network:

```bash
./gwtm --dry-run setup your-org/your-repo
./gwtm --dry-run new-branch feature-test
```

---

## Commit Message Convention

Commit messages are **enforced** by the `.githooks/commit-msg` hook. Every commit must follow the [Conventional Commits](https://www.conventionalcommits.org/) format:

```
<type>(<optional scope>): <description>
```

**Allowed types:**

| Type | Description | Version bump |
|------|-------------|-------------|
| `feat` | A new feature | Minor (`1.2.0` → `1.3.0`) |
| `fix` | A bug fix | Patch (`1.2.0` → `1.2.1`) |
| `perf` | A performance improvement | Patch |
| `refactor` | Code change that is not a feature or fix | None |
| `docs` | Documentation only | None |
| `test` | Adding or fixing tests | None |
| `chore` | Build system, dependencies, config | None |

**Breaking changes** trigger a major version bump. Indicate them with a `!` after the type or a `BREAKING CHANGE:` footer:

```
feat!: remove --flag that was deprecated in v1.2.0

BREAKING CHANGE: the --flag option has been removed. Use --new-flag instead.
```

### Examples

```bash
git commit -m "feat: add shell completion support"
git commit -m "fix: handle missing HOME env var on Windows"
git commit -m "feat(setup): support HTTPS clone URLs"
git commit -m "docs: update installation instructions"
git commit -m "chore: upgrade cobra to v1.11.0"
```

Commits that do not match the pattern will be rejected locally by the hook:

```
❌ Commit message must follow Conventional Commits format
```

---

## Submitting a Pull Request

1. **Create a branch** from `main`:
   ```bash
   git checkout -b fix-my-bug
   ```

2. **Make your changes.** Keep each PR focused on a single concern. If you are fixing a bug and adding a feature, split them into separate PRs.

3. **Add or update tests** for any changed behaviour. All tests must pass:
   ```bash
   go test ./...
   ```

4. **Format your code:**
   ```bash
   make fmt
   ```

5. **Push and open a PR against `main`:**
   ```bash
   git push origin fix-my-bug
   ```
   Then open a pull request on GitHub.

6. **Fill in the PR description.** Include:
   - What the change does and why
   - How you tested it
   - Any follow-up work that is out of scope

7. **CI must pass.** The test workflow runs Go tests with `-race`, builds the binary, and runs a dry-run smoke test. PRs that fail CI will not be merged.

### What reviewers look for

- Tests cover the new or changed behaviour
- Error paths are handled and produce useful messages
- Commits follow the Conventional Commits format (required for correct versioning)
- No unnecessary changes outside the scope of the PR

---

## Release Process

Releases are fully automated — you do not need to do anything manually.

1. A PR is merged into `main`.
2. The **Release** workflow runs [`semantic-release`](https://semantic-release.gitbook.io/), which analyses commit messages since the last release.
3. If there are `feat`, `fix`, or `perf` commits, semantic-release bumps the version, updates `CHANGELOG.md` and `VERSION`, and creates a Git tag.
4. The new tag triggers the **GoReleaser** workflow, which compiles binaries for Linux, macOS (Intel + Apple Silicon), and Windows, and attaches them to the GitHub release.

This means:
- `docs:`, `chore:`, `test:`, and `refactor:` commits will **not** trigger a release on their own.
- Every `fix:` commit produces a patch release.
- Every `feat:` commit produces a minor release.
- A breaking change produces a major release.
