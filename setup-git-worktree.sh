#!/usr/bin/env bash
#
# git-worktree-manager.sh
#
# Modes:
#   1. Full setup:
#        ./git-worktree-manager.sh <repo-url>
#   2. Create new branch worktree (skip setup):
#        ./git-worktree-manager.sh --new-branch <branch-name> [base-branch]
#
# Notes:
#   - In branch-only mode, must be run from the repo root (where .git points to .bare)
#   - Automatically pushes new branches to origin with upstream tracking
#   - Never uses --git-dir

set -e

# --- Helper: Detect default branch from remote ---
detect_default_branch() {
    local branch
    branch=$(git symbolic-ref refs/remotes/origin/HEAD 2>/dev/null | sed 's@^refs/remotes/origin/@@')
    echo "$branch"
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

# --- Branch-only mode ---
if [ "$1" == "--new-branch" ]; then
    if [ -z "$2" ]; then
        echo "Usage: $0 --new-branch <branch-name> [base-branch]"
        exit 1
    fi
    create_new_branch_worktree "$2" "$3"
    exit 0
fi

# --- Full setup mode ---
if [ -z "$1" ]; then
    echo "Usage: $0 <repo-url>"
    exit 1
fi

REPO_URL="$1"
REPO_NAME=$(basename -s .git "$REPO_URL")

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