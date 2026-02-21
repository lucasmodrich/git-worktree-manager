# Design: `using-gwtm` Claude Skill

**Date:** 2026-02-21
**Status:** Approved

---

## Purpose

A personal Claude Code skill (`~/.claude/skills/using-gwtm/SKILL.md`) that teaches Claude to
operate `gwtm` as a user — knowing which commands to run, when to use gwtm instead of plain git,
and the key constraints that prevent mistakes.

## Audience

Claude instances assisting users who work in `gwtm`-managed repositories on any project.

## Location

`~/.claude/skills/using-gwtm/SKILL.md` — personal, available in every Claude Code session.

## Approach

Option B: Rich single-file SKILL.md (~400 words) covering workflow guidance and command reference.
Not loaded every session (tool-specific trigger), so token cost is acceptable.

## Skill Structure

### Frontmatter

```yaml
name: using-gwtm
description: Use when working in a git worktree-managed repository (has a .git file, not directory,
  at its root) or when a user asks to set up a new repo with the bare clone + worktree workflow.
  gwtm replaces manual git worktree commands with setup, new-branch, list, remove, and prune subcommands.
```

### Sections

1. **Overview** — bare clone + worktree mental model; navigate via `cd` not `git checkout`
2. **When to Use** — detect worktree repos by `.git` FILE; use gwtm instead of plain git for branch ops
3. **Directory Structure** — canonical layout with `.bare/`, `.git` file, and sibling worktree dirs
4. **Key Constraints** — hyphens not slashes in branch names; must run from within managed repo; `--dry-run`
5. **Command Reference** — table of all 7 commands with syntax and notes
6. **Typical Workflow** — end-to-end example: setup → new-branch → work → remove
7. **Common Mistakes** — slash branch names, using plain git worktree, forgetting `--remote`

## Test Plan (RED-GREEN-REFACTOR)

- **RED**: Run a subagent scenario without the skill; document where Claude defaults to plain git
  (e.g. `git checkout -b feature/login`, wrong branch naming, wrong navigation pattern)
- **GREEN**: With skill present, verify Claude uses `gwtm new-branch feature-login`, hyphens, `cd`
- **REFACTOR**: Close any loopholes found in green testing

## Success Criteria

- Claude correctly uses `gwtm new-branch` instead of `git checkout -b`
- Claude rejects slash branch names and suggests hyphen equivalents
- Claude knows to `cd <branch>` to switch, not `git checkout`
- Claude suggests `--dry-run` before destructive operations
- Claude understands the directory structure (worktrees are sibling dirs, not nested)
