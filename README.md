# Git Worktree Setup Script

## 📌 Overview
This repository contains a shell script, `setup-git-worktree.sh`, that automates the setup of a **bare Git repository** with a clean, maintainable **worktree-based workflow** for local development.

The script:
- Creates a **root project folder** named after the GitHub repository.
- Bare-clones the repository into a `.bare` folder inside the root.
- Creates a `.git` file in the root that points to `.bare`.
- Configures Git to **fetch and track all remote branches** automatically.
- Auto-detects the **default branch** from the remote (`main`, `master`, or other).
- Creates the **initial worktree** for the default branch as a local tracking branch (no detached HEAD).
- Avoids the need to use `--git-dir` in daily work.

---

## 📂 Folder Structure

After running the script, your project will look like this:

```
<repo-name>/
├── .bare/           # Bare repository clone (Git metadata only)
├── .git             # Points to .bare
└── <default-branch>/ # Worktree for the default branch
```

All additional worktrees will also be created under the root folder.

---

## 🚀 Usage

### 1. Make the script executable
```bash
chmod +x setup-git-worktree.sh
```

### 2. Run the script
```bash
./setup-git-worktree.sh <repo-url>
```
Example:
```bash
./setup-git-worktree.sh git@github.com:org/repo.git
```

The script will:
1. Create a folder named after the repo.
2. Bare-clone into `.bare`.
3. Configure tracking for all remote branches.
4. Fetch all branches.
5. Detect the default branch.
6. Create the first worktree for that branch.

---

## 🛠 Creating Additional Worktrees

### New branch from default branch
```bash
git worktree add feature/my-feature -b feature/my-feature origin/<default-branch>
```

### Worktree for an existing remote branch
```bash
git worktree add bugfix/issue-123 origin/bugfix/issue-123
```

---

## 🔄 Maintenance Commands

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
# Create a new feature branch worktree
git worktree add feature/login-page -b feature/login-page origin/main

# Switch to it
cd feature/login-page

# Work, commit, push
git add .
git commit -m "Implement login page"
git push -u origin feature/login-page
```

---

## ✅ Benefits of This Workflow
- **Disk space efficiency** — all worktrees share the same `.bare` repo data.
- **Fast branch switching** — no need to stash or re-checkout.
- **Parallel development** — work on multiple branches at once in separate directories.
- **Clean separation** — Git metadata is isolated from working directories.
