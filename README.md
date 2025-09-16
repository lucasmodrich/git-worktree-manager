# Git Worktree Manager

## 📌 Overview
`git-worktree-manager.sh` is a flexible shell script for managing Git repositories using a **bare clone + worktree** workflow.

It supports five modes:

1. **Full Setup Mode** — Clones a GitHub repo using `org/repo` shorthand, configures it for worktrees, and creates the initial worktree for the default branch.
2. **Branch‑Only Mode** — Creates a new branch worktree (from a base branch), automatically pushing it to GitHub.
3. **List Mode** — Displays all active worktrees.
4. **Prune Mode** — Cleans up stale worktree references.
5. **Remove Mode** — Deletes a worktree and its local branch.

This workflow:
- Keeps Git metadata isolated in `.bare`
- Allows multiple branches to be worked on in parallel
- Saves disk space by sharing the same object database
- Ensures remote branches are fetched and tracked
- Automatically pushes new branches to GitHub

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
│ Start: Provide <org>/<repo>  │
└───────────────┬──────────────┘
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

### 🌐 Local ↔ Remote Branch Relationship Diagram

```
          GitHub Remote (origin)
        ┌────────────────────────┐
        │  origin/main            │
        │  origin/feature-x       │
        │  origin/bugfix-y        │
        └───────────┬─────────────┘
                    │ fetch/push
                    ▼
           ┌─────────────────┐
           │   .bare repo     │
           │ (local metadata) │
           └───────┬─────────┘
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

---

## 🚀 Usage

### 1. Make the script executable
```bash
chmod +x git-worktree-manager.sh
```

---

### 2. Full Setup Mode
```bash
./git-worktree-manager.sh your-org/your-repo
```

Creates:
- Root folder named after the repo
- Bare clone in `.bare`
- `.git` pointer
- Fetch config
- Initial worktree for default branch
- Pushes it to GitHub if new

---

### 3. Create New Branch Worktree
```bash
./git-worktree-manager.sh --new-branch <branch-name> [base-branch]
```

Examples:
```bash
./git-worktree-manager.sh --new-branch feature/login-page
./git-worktree-manager.sh --new-branch hotfix/payment-bug develop
```

---

### 4. List Worktrees
```bash
./git-worktree-manager.sh --list
```

---

### 5. Prune Stale Worktrees
```bash
./git-worktree-manager.sh --prune
```

---

### 6. Remove Worktree and Local Branch
```bash
./git-worktree-manager.sh --remove <branch-name>
```

Example:
```bash
./git-worktree-manager.sh --remove feature/login-page
```

---

## 🛠 Common Git Commands

| Task | Command |
|------|---------|
| List all worktrees | `git worktree list` |
| Remove a worktree | `git worktree remove <dir>` |
| Delete a local branch | `git branch -D <branch>` |
| Prune stale worktrees | `git worktree prune` |
| Fetch all updates | `git fetch --all --prune` |

---

## 💡 Best Practices
- Keep all worktrees under the root folder for clarity.
- Use descriptive branch names for worktree directories.
- Regularly run `--prune` to keep `.bare` clean.
- Remove unused worktrees with `--remove`.
- Never edit files directly in `.bare`.

---

## 📖 Example Workflow

```bash
# Initial setup
./git-worktree-manager.sh your-org/your-repo

# Create a new feature branch
./git-worktree-manager.sh --new-branch feature/login-page

# Work in the new branch
cd feature/login-page
git add .
git commit -m "Implement login page"
git push

# Clean up when done
./git-worktree-manager.sh --remove feature/login-page
```

---

## ✅ Benefits of This Workflow
- Disk space efficiency — all worktrees share the same `.bare` repo data.
- Fast branch switching — no need to stash or re-checkout.
- Parallel development — work on multiple branches at once.
- Clean separation — Git metadata is isolated from working directories.
- Automatic remote sync — new branches are pushed to GitHub immediately.

