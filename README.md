# Git Worktree Manager

## 📌 Overview
`git-worktree-manager.sh` is a self-updating shell script for managing Git repositories using a **bare clone + worktree** workflow.
I created this project following frustration using the standard tooling.

It supports:

- **Full setup** from GitHub using `org/repo` shorthand
- **Branch creation** with automatic remote push
- **Worktree listing**, pruning, and removal
- **Version tracking** and **self-upgrade**
- **Markdown-style help card** for onboarding

---

## 🚀 Installation

By default `git-worktree-manager` will install and update itself into the `$HOME\.git-worktree-manager\` (`~\.git-worktree-manager`) folder. Its hard-wired for this presently. If you want to install it otherwise, you will need to clone or fork the repo and update it to your needs. If your changes are broadly useful, please consider raising a PR to submit them for inclusion here.


To install directly from this GitHub repo, use the following command:
```bash
curl -sSL https://raw.githubusercontent.com/lucasmodrich/git-worktree-manager/refs/heads/main/git-worktree-manager.sh | bash -s -- --upgrade
```

## 🧠 Versioning & Upgrade

- Check version:
  ```bash
  ./git-worktree-manager.sh --version
  ```
- Upgrade to latest from GitHub:
  ```bash
  ./git-worktree-manager.sh --upgrade
  ```

## Configure Git Alias
If you want to utilise the `git-worktree-manager` via a Git alias, add the following to you `~/.gitconfig` file.
```
[alias]
        wtm = "!bash $HOME/.git-worktree-manager/git-worktree-manager.sh"
```

Once the alias has been added to your `~/.gitconfig` file, you can invoke the `git-worktree-manager` using the following:
```bash
git wtm -h
```
> Note: the `--help` flag does not work when called via th `git` cli as it invokes the `git` cli Help.

---

## 📂 Folder Structure

After setup:
```
<repo-name>/
├── .bare/             # Bare repository clone
├── .git               # Points to .bare
└── <default-branch>/  # Initial worktree
```

---

## 🖼 Architecture Diagram

```
                ┌───────────────────────────┐
                │        .bare repo         │
                │  (Git metadata & objects) │
                └─────────────┬─────────────┘
                              │
     ┌────────────────────────┼────────────────────────┐
     │                        │                        │
┌──────────────┐       ┌──────────────┐        ┌──────────────┐
│ main/        │       │ feature-x/   │        │ bugfix-y/    │
│ (worktree)   │       │ (worktree)   │        │ (worktree)   │
└──────────────┘       └──────────────┘        └──────────────┘
```

---

## 🔄 Flow Diagrams

### Full Setup
```
Input: <org>/<repo>
→ Create root folder
→ Clone into .bare
→ Point .git to .bare
→ Configure fetch
→ Fetch branches
→ Detect default branch
→ Create worktree
→ Push if new
```

### Branch Creation
```
Input: --new-branch <branch> [base]
→ Fetch branches
→ Create worktree
→ Push if new
```

---

## 🌐 Local ↔ Remote Relationship

```
GitHub Remote (origin)
┌────────────────────┐
│ origin/main        │
│ origin/feature-x   │
└──────────┬─────────┘
           ▼
       .bare repo
           ▼
┌──────────────┐
│ feature-x/   │
└──────────────┘
```

---

## 🚀 Usage

### Make executable
```bash
chmod +x git-worktree-manager.sh
```

---

### Full Setup
```bash
./git-worktree-manager.sh your-org/your-repo
```

---

### Create New Branch
```bash
./git-worktree-manager.sh --new-branch <branch> [base]
```

---

### Remove Worktree + Branch
```bash
./git-worktree-manager.sh --remove <branch>
```

---

### List Worktrees
```bash
./git-worktree-manager.sh --list
```

---

### Prune Stale Worktrees
```bash
./git-worktree-manager.sh --prune
```

---

### Show Version
```bash
./git-worktree-manager.sh --version
```

---

### Upgrade Script
```bash
./git-worktree-manager.sh --upgrade
```

---

### Help Card
```bash
./git-worktree-manager.sh --help
./git-worktree-manager.sh -h
```

---

## 📖 Example Workflow

```bash
# Setup
./git-worktree-manager.sh acme/webapp

# Create feature branch
./git-worktree-manager.sh --new-branch feature/login-page

# Work
cd feature/login-page
git add .
git commit -m "Add login page"
git push

# Clean up
cd ..
./git-worktree-manager.sh --remove feature/login-page
```

---

## ✅ Benefits

- Disk-efficient multi-branch development
- No detached HEADs
- Easy onboarding with help card
- Self-updating and version-aware
- GitHub-native workflow
