---
name: using-gwtm
description: Use when working in a git worktree-managed repository (has a .git file, not directory,
  at its root) or when a user asks to set up a new repo with the bare clone + worktree workflow.
  gwtm replaces manual git worktree commands with setup, new-branch, list, remove, and prune subcommands.
---

# Using gwtm

## Overview

`gwtm` is a CLI tool for bare-clone + worktree Git workflows. Each branch lives in its own
sibling directory — switch contexts by changing directory, not `git checkout`. No stashing,
no detached HEADs.

## When to Use gwtm (vs plain git)

Detect a gwtm-managed repo: the repo root has a `.git` **file** (not directory).

| Situation | Use |
|---|---|
| Set up a new repo for worktree workflow | `gwtm setup` (not manual `git clone --bare`) |
| Create a branch and start working on it | `gwtm new-branch` (not `git checkout -b` + `git worktree add`) |
| Remove a branch when done | `gwtm remove` (not manual `git worktree remove` + `git branch -d`) |
| List active worktrees | `gwtm list` |
| Clean stale refs | `gwtm prune` |

## Directory Structure

```
<repo-name>/
├── .bare/          ← bare clone (Git objects; never work here directly)
├── .git            ← FILE pointing to .bare (not a directory)
├── main/           ← worktree for main branch
├── feature-login/  ← worktree for feature-login branch
└── bugfix-crash/   ← worktree for bugfix-crash branch
```

To "switch branches": `cd ../feature-login` — the other directory is already checked out.

## Key Constraints

- **Branch names: hyphens only.** Slashes and backslashes are rejected. `feature-login` ✅  `feature/login` ❌
- **Most commands require a gwtm repo root.** gwtm walks up from CWD looking for a `.git` file.
- **`--dry-run` is available for every command.** Always suggest it before destructive operations.
- **`GIT_WORKTREE_MANAGER_HOME`** env var controls where `gwtm upgrade` installs (default `$HOME/.git-worktree-manager`).

## Command Reference

| Command | Syntax | Notes |
|---|---|---|
| `setup` | `gwtm setup <org>/<repo>` | One-time per repo. Also accepts `git@host:org/repo.git` and `https://...` URLs. Creates `.bare/`, `.git` file, and initial worktree. |
| `new-branch` | `gwtm new-branch <name> [base]` | Creates new branch from base (or default), or checks out existing local/remote branch. Prompts before fetching remote-only branches. Auto-pushes new branches. |
| `list` | `gwtm list` | Lists all active worktrees. |
| `remove` | `gwtm remove <branch> [--remote]` | Removes worktree + local branch. Add `--remote` to also delete the remote branch. |
| `prune` | `gwtm prune` | Cleans stale worktree refs from `.bare`. |
| `version` | `gwtm version` | Shows version and checks GitHub for a newer release. |
| `upgrade` | `gwtm upgrade` | Self-updates the binary from GitHub releases. |
| `--dry-run` | `gwtm --dry-run <cmd>` | Previews any operation without executing. |

## Typical Workflow

```bash
# Once per repo — clone and set up worktree structure
gwtm setup acme/webapp
cd acme/webapp

# Start a feature (hyphens, not slashes)
gwtm new-branch feature-login
cd feature-login

# ... do your work ...
git add .
git commit -m "feat: add login page"
git push

# Done — clean up local and remote
cd ..
gwtm remove feature-login --remote

# Check out a colleague's branch
gwtm new-branch feature-payments   # detects remote branch, prompts to check out
cd feature-payments
```

## Common Mistakes

| Mistake | Correct approach |
|---|---|
| `git checkout -b feature/login` | `gwtm new-branch feature-login` |
| `git worktree add feature-login -b feature/login` | `gwtm new-branch feature-login` (gwtm manages the worktree) |
| Branch name `feature/login` | `feature-login` (slashes are rejected) |
| `git checkout feature-login` to switch | `cd ../feature-login` (already checked out in its own dir) |
| `gwtm remove feature-login` to also clean remote | `gwtm remove feature-login --remote` |
| Deleting files manually instead of removing worktree | `gwtm remove <branch>` cleans up refs too |
