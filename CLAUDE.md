# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Overview

`gwtm` is a Go CLI tool for managing Git repositories using a **bare clone + worktree** workflow. It provides commands for repository setup, branch creation, worktree management, and self-updating.

## Key Architecture

### Core Components

- **Go module**: `github.com/lucasmodrich/git-worktree-manager` (Go 1.25.1)

1. **Entry point**: `cmd/git-worktree-manager/main.go` — injects version via `-ldflags` at build time
2. **CLI commands**: `internal/commands/` — one file per subcommand, using the Cobra framework
3. **Git client**: `internal/git/` — thin wrapper around `exec.Command("git", ...)`, all methods respect `DryRun`
4. **Config**: `internal/config/` — install directory and binary path resolution
5. **UI**: `internal/ui/` — stdout/stderr formatting, dry-run output, error messages
6. **Version**: `internal/version/` — semver parsing and self-upgrade logic

### Workflow Structure

```
<repo-name>/
├── .bare/             # Bare repository clone
├── .git               # File pointing to .bare
└── <branches>/        # Individual worktrees for each branch
```

### Key Design Decisions

- **No `os.Chdir`**: All paths are computed as absolute values from `os.Getwd()` + `filepath.Join`
- **`findWorktreeRoot()`**: Walks up from CWD checking for `.git` file (not directory) — enables commands to run from any subdirectory
- **Branch names**: Slashes (`/`, `\`) are rejected; hyphens must be used instead to avoid filesystem conflicts
- **Upgrade installs to configured home dir**: Always installs to `GIT_WORKTREE_MANAGER_HOME` (default `$HOME/.git-worktree-manager`), never to the running binary's location
- **Error messages**: Always go to stderr via `ui.PrintError`; every error includes actionable guidance

## Subcommands

| Command | Signature | Notes |
|---------|-----------|-------|
| `setup` | `setup <org>/<repo>` | Clone as bare repo + initial worktree |
| `new-branch` | `new-branch <branch-name> [base-branch]` | Create branch + worktree |
| `list` | `list` | List active worktrees |
| `remove` | `remove <branch>` | Remove worktree + branch; `--remote` also deletes remote branch |
| `prune` | `prune` | Prune stale worktree refs |
| `upgrade` | `upgrade` | Self-update binary |
| `version` | `version` | Print version/build info, check GitHub for updates |

## Common Commands

### Build
```bash
go build -o gwtm ./cmd/git-worktree-manager
```

### Format / Lint
```bash
gofmt -s -w .
golangci-lint run   # if installed
```

### Clean
```bash
rm -f gwtm gwtm.exe
```

### Test
```bash
go test ./...
go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
```

### Local snapshot build (all platforms, with version injection)
```bash
goreleaser release --snapshot --clean
```

### Smoke test
```bash
./gwtm --help
./gwtm version
./gwtm --dry-run setup your-org/your-repo
```

## Release Configuration

Releases are fully automated via a single workflow (`.github/workflows/goreleaser.yml`) that triggers on every push to `main`:

1. `mathieudutour/github-tag-action` analyses commits since the last tag and determines the next semver version from conventional commit types. No tag is created for `docs:`, `chore:`, `test:`, or `refactor:` commits.
2. If a new tag is warranted, GoReleaser compiles binaries for Linux, macOS (Intel + Apple Silicon), and Windows, then creates the GitHub release with release notes and uploads all artifacts.

Canonical release history is on the [GitHub Releases page](https://github.com/lucasmodrich/git-worktree-manager/releases). `CHANGELOG.md` is a historical record up to v2.0.0.
