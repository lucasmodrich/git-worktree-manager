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
go build -o gwtm ./cmd/git-worktree-manager
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
│       ├── release.yml      # Runs semantic-release on merge to main (versioning + changelog)
│       └── goreleaser.yml   # Triggered by a new tag; builds binaries and creates GitHub release
├── .githooks/
│   └── commit-msg           # Enforces Conventional Commits format
├── .goreleaser.yml          # Cross-platform binary build and GitHub release config
└── release.config.js        # semantic-release: version bump, CHANGELOG.md, VERSION file, git tag
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
go build -o gwtm ./cmd/git-worktree-manager
```

### Test

```bash
go test ./...

# With race detection and coverage (matches CI)
go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
```

Tests use `t.TempDir()` and local bare repositories — no network access required.

### Format

```bash
gofmt -s -w .
```

### Lint

```bash
golangci-lint run   # if installed
```

### Manual smoke test with dry-run

The `--dry-run` flag lets you exercise any command without touching the filesystem or network:

```bash
./gwtm --dry-run setup your-org/your-repo
./gwtm --dry-run new-branch feature-test
```

### Local multi-platform build (snapshot)

To build binaries for all target platforms locally with version info injected, use GoReleaser in snapshot mode:

```bash
goreleaser release --snapshot --clean
# Outputs to ./dist/
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
   gofmt -s -w .
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

### How it works

A single workflow (`.github/workflows/goreleaser.yml`) runs on every push to `main`:

```
Merge to main
     │
     ▼
github-tag-action
  • Reads conventional commits since the last tag
  • Determines semver bump (patch / minor / major)
  • If no releasable commits: stops here (no release)
  • Creates and pushes the new tag (e.g. v2.1.0)
     │
     ▼ (same workflow, next step)
GoReleaser
  • Runs go mod tidy and go test ./...
  • Compiles binaries for Linux, macOS (Intel + Apple Silicon), Windows
  • Creates the GitHub release with generated release notes
  • Attaches binaries and checksums.txt
```

Everything runs in a single workflow — no cross-workflow token issues, no Node.js, no separate release pipeline.

### Manual publishing

If the automated pipeline fails to publish a release (e.g. a transient CI error), go to **Actions → Release → Run workflow**, enter the existing tag (e.g. `v2.1.0`), and click **Run workflow**. The tag step is skipped and GoReleaser runs directly against that tag.

### Version bump rules

Commit type determines the version bump — this is why the Conventional Commits format is required:

| Commit type | Effect |
|---|---|
| `feat` | Minor bump (`1.2.0` → `1.3.0`) |
| `fix`, `perf` | Patch bump (`1.2.0` → `1.2.1`) |
| `feat!` / `BREAKING CHANGE` | Major bump (`1.2.0` → `2.0.0`) |
| `docs`, `chore`, `test`, `refactor` | No release triggered |

### Release history

Full release notes for every version are on the [GitHub Releases page](https://github.com/lucasmodrich/git-worktree-manager/releases). `CHANGELOG.md` in the repository is a historical record up to v2.0.0.
