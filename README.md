# Git Worktree Manager

## 📌 Overview
`git-worktree-manager.sh` is a flexible shell script for managing Git repositories using a **bare clone + worktree** workflow.

It supports two modes:

1. **Full Setup Mode** — Clones a GitHub repo into a `.bare` folder, configures it for worktrees, and creates the initial worktree for the default branch.
2. **Branch‑Only Mode** — Skips the bare clone and simply creates a new branch worktree (from a base branch), automatically pushing it to the remote with upstream tracking.

This workflow:
- Keeps Git metadata isolated in `.bare`
- Allows multiple branches to be worked on in parallel without switching
- Saves disk space by sharing the same object database across worktrees
- Ensures all remote branches are fetched and tracked

---

## 📂 Folder Structure

After **Full Setup Mode**, your project will look like:

```
<repo-name>/
├── .bare/             # Bare repository clone (Git metadata only)
├── .git               # Points to .bare
└── <default-branch>/  # Worktree for the default branch
```

All additional worktrees will also be created under the root folder.

---

## 🚀 Usage

### 1. Make the script executable
```bash
chmod +x git-worktree-manager.sh
```

---

### 2. **Full Setup Mode** (Initial Repo Setup)
```bash
./git-worktree-manager.sh <repo-url>
```
Example:
```bash
./git-worktree-manager.sh git@github.com:org/repo.git
```

**What happens:**
1. Creates a root folder named after the repo.
2. Bare-clones the repo into `.bare`.
3. Creates `.git` file pointing to `.bare`.
4. Configures Git to fetch **all** remote branches.
5. Fetches all branches from the remote.
6. Detects the default branch from `origin/HEAD`.
7. Creates a local tracking branch for the default branch in a worktree.
8. Pushes it to GitHub if it’s new.

---

### 3. **Branch‑Only Mode** (Skip Setup, Create New Branch)
From inside an already set‑up repo root:
```bash
./git-worktree-manager.sh --new-branch <branch-name> [base-branch]
```

Examples:
```bash
# Create from default branch
./git-worktree-manager.sh --new-branch feature/login-page

# Create from a specific base branch
./git-worktree-manager.sh --new-branch hotfix/payment-bug develop
```

**What happens:**
1. Fetches all remote branches.
2. If the branch exists locally → creates a worktree from it.
3. If it’s new → creates it from the base branch, sets it to track the remote, and pushes it to GitHub.

---

## 🛠 Common Worktree Commands

| Task | Command |
|------|---------|
| List all worktrees | `git worktree list` |
| Remove a worktree | `git worktree remove <dir>` |
| Prune stale worktrees | `git worktree prune` |
| Fetch all updates | `git fetch --all --prune` |

---

## 💡 Best Practices
- Keep **all worktrees** under the root folder for clarity.
- Use **descriptive branch names** for worktree directories.
- Regularly run `git fetch --all --prune` to keep remotes in sync.
- Remove unused worktrees to avoid clutter.
- Never edit files directly in `.bare`.

---

## 📖 Example Workflow

```bash
# 1. Initial setup
./git-worktree-manager.sh git@github.com:org/repo.git

# 2. Create a new feature branch from default branch
./git-worktree-manager.sh --new-branch feature/login-page

# 3. Work in the new branch
cd feature/login-page
git add .
git commit -m "Implement login page"
git push
```

---

## ✅ Benefits of This Workflow
- **Disk space efficiency** — all worktrees share the same `.bare` repo data.
- **Fast branch switching** — no need to stash or re-checkout.
- **Parallel development** — work on multiple branches at once in separate directories.
- **Clean separation** — Git metadata is isolated from working directories.
- **Automatic remote sync** — new branches are pushed to GitHub immediately.
```

