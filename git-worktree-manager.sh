#!/usr/bin/env bash
#
# git-worktree-manager.sh

SCRIPT_VERSION="1.0.1-beta.4"
SCRIPT_FOLDER="$HOME/.git-worktree-manager"
SCRIPT_NAME="git-worktree-manager.sh"
GITHUB_REPO="lucasmodrich/git-worktree-manager"
RAW_URL="https://raw.githubusercontent.com/$GITHUB_REPO/main/$SCRIPT_NAME"

set -e

# --- Helper: Compare semantic versions ---
version_gt() {
    [ "$(printf '%s\n' "$1" "$2" | sort -V | head -n1)" != "$1" ]
}

# --- Helper: Detect default branch from remote ---
detect_default_branch() {
    git symbolic-ref refs/remotes/origin/HEAD 2>/dev/null | sed 's@^refs/remotes/origin/@@'
}

# --- Helper: Create and push new branch worktree ---
create_new_branch_worktree() {
    local new_branch="$1"
    local base_branch="$2"

    if [ -z "$base_branch" ]; then
        base_branch=$(detect_default_branch)
        if [ -z "$base_branch" ]; then
            echo "❌ Could not detect default branch. Please specify base branch."
            exit 1
        fi
    fi

    echo "📡 Fetching latest from origin"
    git fetch --all --prune

    if git show-ref --verify --quiet "refs/heads/$new_branch"; then
        echo "📂 Branch '$new_branch' exists locally — creating worktree from it"
        git worktree add "$new_branch" "$new_branch"
    else
        echo "🌱 Creating new branch '$new_branch' from '$base_branch'"
        git worktree add "$new_branch" -b "$new_branch" --track "origin/$base_branch"
        echo "☁️  Pushing new branch '$new_branch' to origin"
        (cd "$new_branch" && git push -u origin "$new_branch")
    fi

    echo "✅ Worktree for '$new_branch' is ready"
    git worktree list
}

# --- Helper: List all worktrees ---
list_worktrees() {
    echo "📋 Active Git worktrees:"
    git worktree list
}

# --- Helper: Prune stale worktrees ---
prune_worktrees() {
    echo "🧹 Pruning stale worktrees..."
    git worktree prune
    echo "✅ Prune complete."
}

# --- Helper: Remove worktree and local branch ---
remove_worktree_and_branch() {
    local branch="$1"

    if [ -z "$branch" ]; then
        echo "Usage: $0 --remove <branch-name>"
        exit 1
    fi

    if ! git worktree list | grep -q "/$branch "; then
        echo "❌ Worktree for branch '$branch' not found."
        exit 1
    fi

    echo "🗑 Removing worktree '$branch'"
    git worktree remove "$branch"

    if git show-ref --verify --quiet "refs/heads/$branch"; then
        echo "🧨 Deleting local branch '$branch'"
        git branch -D "$branch"
    else
        echo "⚠️ Local branch '$branch' not found — nothing to delete."
    fi

    echo "✅ Removal complete."
}

# --- Helper: Show version ---
show_version() {
    echo "git-worktree-manager.sh version $SCRIPT_VERSION"
}

# --- Helper: Upgrade script ---
upgrade_script() {
    echo "🔍 Checking for newer version on GitHub..."
    echo "Script Folder: $SCRIPT_FOLDER"
    mkdir -p "$SCRIPT_FOLDER"

    remote_version=$(curl -s "$RAW_URL" | grep '^SCRIPT_VERSION=' | cut -d'"' -f2)

    if [ -z "$remote_version" ]; then
        echo "❌ Could not retrieve remote version."
        exit 1
    fi

    echo "🔢 Local version: $SCRIPT_VERSION"
    echo "🌐 Remote version: $remote_version"

    if version_gt "$SCRIPT_VERSION" "$remote_version"; then
        echo "✅ You already have the latest version."
    else
        echo "⬇️ Upgrading to version $remote_version..."
        curl -s -o "$SCRIPT_FOLDER/$SCRIPT_NAME.tmp" "$RAW_URL"
        mv "$SCRIPT_FOLDER/$SCRIPT_NAME.tmp" "$SCRIPT_FOLDER/$SCRIPT_NAME"
        chmod +x "$SCRIPT_FOLDER/$SCRIPT_NAME"
        echo "✅ Upgrade complete. Now running version $remote_version."
    fi
    exit 0
}

# --- Helper: Show help card ---
show_help() {
    cat <<EOF

🛠 Git Worktree Manager — Help Card

Usage:
  $0 <org>/<repo>                     # Full setup from GitHub
  $0 --new-branch <branch> [base]     # Create new branch worktree
  $0 --remove <branch>                # Remove worktree and local branch
  $0 --list                           # List active worktrees
  $0 --prune                          # Prune stale worktrees
  $0 --version                        # Show script version
  $0 --upgrade                        # Upgrade to latest version
  $0 --help (-h)                      # Show this help card

Examples:
  $0 acme/webapp
  $0 --new-branch feature/login-page
  $0 --remove feature/login-page

Notes:
  - Run from repo root (where .git points to .bare)
  - New branches are pushed to GitHub automatically
  - Remote branch is not deleted by --remove

EOF
}

# --- Mode: Help ---
if [ "$1" == "--help" ] || [ "$1" == "-h" ]; then
    show_help
    exit 0
fi

# --- Mode: Version ---
if [ "$1" == "--version" ]; then
    show_version
    exit 0
fi

# --- Mode: Upgrade ---
if [ "$1" == "--upgrade" ]; then
    upgrade_script
fi

# --- Mode: List worktrees ---
if [ "$1" == "--list" ]; then
    list_worktrees
    exit 0
fi

# --- Mode: Prune worktrees ---
if [ "$1" == "--prune" ]; then
    prune_worktrees
    exit 0
fi

# --- Mode: Remove worktree and branch ---
if [ "$1" == "--remove" ]; then
    remove_worktree_and_branch "$2"
    exit 0
fi

# --- Mode: Create new branch worktree ---
if [ "$1" == "--new-branch" ]; then
    if [ -z "$2" ]; then
        echo "Usage: $0 --new-branch <branch-name> [base-branch]"
        exit 1
    fi
    create_new_branch_worktree "$2" "$3"
    exit 0
fi

# --- Mode: Full setup ---
if [ -z "$1" ]; then
    show_help
    exit 1
fi

# Convert org/repo to full GitHub SSH URL
if [[ "$1" != git@github.com:* ]]; then
    REPO_PATH="$1"
    REPO_URL="git@github.com:$REPO_PATH.git"
else
    REPO_URL="$1"
    REPO_PATH=$(basename -s .git "$REPO_URL")
fi

REPO_NAME=$(basename -s .git "$REPO_PATH")

echo "📂 Creating project root: $REPO_NAME"
mkdir -p "$REPO_NAME"
cd "$REPO_NAME"

echo "📦 Cloning bare repository into .bare"
git clone --bare "$REPO_URL" .bare

echo "📝 Creating .git file pointing to .bare"
echo "gitdir: ./.bare" > .git

echo "⚙️ Configuring Git for auto remote tracking"
git config push.default current
git config branch.autosetupmerge always
git config branch.autosetuprebase always

echo "🔧 Ensuring all remote branches are fetched"
git config remote.origin.fetch "+refs/heads/*:refs/remotes/origin/*"

echo "📡 Fetching all remote branches"
git fetch --all --prune

DEFAULT_BRANCH=$(detect_default_branch)
if [ -z "$DEFAULT_BRANCH" ]; then
    echo "❌ Could not detect default branch. Please specify manually."
    exit 1
fi

echo "🌱 Creating initial worktree for branch: $DEFAULT_BRANCH"
if git show-ref --verify --quiet "refs/heads/$DEFAULT_BRANCH"; then
    git worktree add "$DEFAULT_BRANCH" "$DEFAULT_BRANCH"
else
    git worktree add "$DEFAULT_BRANCH" -b "$DEFAULT_BRANCH" --track "origin/$DEFAULT_BRANCH"
    echo "☁️  Pushing new branch '$DEFAULT_BRANCH' to origin"
    (cd "$DEFAULT_BRANCH" && git push -u origin "$DEFAULT_BRANCH")
fi

echo "✅ Setup complete!"
git worktree list