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

## 🖼 Visual Architecture Diagram

```
                ┌───────────────────────────┐
                │        .bare repo          │
                │  (Git metadata & objects)  │
                └─────────────┬─────────────┘
                              │
     ┌────────────────────────┼────────────────────────┐
     │                        │                        │
┌──────────────┐       ┌──────────────┐        ┌──────────────┐
│ main/        │       │ feature-x/   │        │ bugfix-y/    │
│ (worktree)   │       │ (worktree)   │        │ (worktree)   │
└──────────────┘       └──────────────┘        └──────────────┘
     │                        │                        │
     ▼                        ▼                        ▼
 Files for main branch   Files for feature-x     Files for bugfix-y
 (tracked by .bare)      (tracked by .bare)      (tracked by .bare)
```

---

## 🔄 Flow Diagrams

### **Full Setup Mode**
```
┌──────────────────────────────┐
│ Start: Provide <repo-url>     │
└───────────────┬───────────────┘
                ▼
      Create root folder
                ▼
      Bare clone into .bare
                ▼
  Create .git pointing to .bare
                ▼
 Configure fetch for all branches
                ▼
       Fetch all branches
                ▼
 Detect default branch from remote
                ▼
 Does local branch exist?
     ┌──────────┴──────────┐
     │ Yes                 │ No
     ▼                     ▼
 Create worktree     Create worktree
 from local branch   from remote branch
                     & push to origin
                ▼
        Setup complete
```

---

### **Branch‑Only Mode**
```
┌──────────────────────────────────────┐
│ Start: --new-branch <branch> [base]   │
└───────────────────┬──────────────────┘
                    ▼
         Fetch all branches
                    ▼
   Does local branch exist?
     ┌──────────┴──────────┐
     │ Yes                 │ No
     ▼                     ▼
 Create worktree     Create worktree
 from local branch   from base branch
                     & push to origin
                    ▼
             Worktree ready
```

---

## 🌐 Local ↔ Remote Branch Relationship Diagram

```
          GitHub Remote (origin)
        ┌────────────────────────┐
        │  origin/main           │
        │  origin/feature-x      │
        │  origin/bugfix-y       │
        └───────────┬────────────┘
                    │ fetch/push
                    ▼
           ┌──────────────────┐
           │   .bare repo     │
           │ (local metadata) │
           └───────┬──────────┘
                   │
   ┌───────────────┼────────────────┐
   │               │                │
┌───────┐     ┌───────────┐    ┌───────────┐
│ main/ │     │ feature-x/│    │ bugfix-y/ │
│ local │     │ local     │    │ local     │
└───────┘     └───────────┘    └───────────┘
   │             │                │
   ▼             ▼                ▼
 Files for   Files for        Files for
 main        feature-x        bugfix-y
```

**How it works:**
- **Remote branches** live on GitHub under `origin/*`.
- The `.bare` repo stores **local tracking branches** (`refs/remotes/origin/*`) and local branches (`refs/heads/*`).
- Each worktree is checked out to a **local branch** that tracks its corresponding remote branch.
- Fetch updates `.bare` from GitHub; push sends changes from a worktree’s branch to GitHub.

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
