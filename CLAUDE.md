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
make build          # produces ./gwtm
go build -o gwtm ./cmd/git-worktree-manager
```

### Format / Lint / Clean
```bash
make fmt            # gofmt -s -w .
make lint           # golangci-lint run (if installed)
make clean          # remove ./gwtm and ./gwtm.exe; go clean
```

### Test
```bash
go test ./...
go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
```

### Install
```bash
make install        # copies to $GOPATH/bin or $HOME/.git-worktree-manager/
```

### Smoke test
```bash
./gwtm --help
./gwtm version
./gwtm --dry-run setup your-org/your-repo
```

## Release Configuration

Releases are fully automated:
1. **semantic-release** (`release.config.js`) runs on merge to `main` — analyses commits, bumps version, updates `CHANGELOG.md` and `VERSION`, creates a Git tag
2. **GoReleaser** (`.github/workflows/goreleaser.yml`) triggers on the new tag — compiles binaries for Linux, macOS (Intel + Apple Silicon), and Windows

The `VERSION` file and `CHANGELOG.md` are managed by CI/CD. Do not edit them manually.
