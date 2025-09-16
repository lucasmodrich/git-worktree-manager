#!/usr/bin/env bash
#
# setup-git-worktree.sh
#
# Usage:
#   ./setup-git-worktree.sh <repo-url>
#
# Example:
#   ./setup-git-worktree.sh git@github.com:org/repo.git
#
# This script:
#   - Creates a root folder named after the repo
#   - Bare clones into .bare
#   - Points .git to .bare
#   - Configures fetch to track ALL remote branches
#   - Auto-detects default branch from remote
#   - Creates initial worktree for default branch on a local tracking branch
#   - Automatically pushes new branches to GitHub when created
#   - Never uses --git-dir

set -e

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

echo "🔍 Detecting default branch from remote"
DEFAULT_BRANCH=$(git symbolic-ref refs/remotes/origin/HEAD | sed 's@^refs/remotes/origin/@@')

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

echo
echo "💡 To create and push a new branch worktree:"
echo "   git worktree add <dir> -b <branch> origin/$DEFAULT_BRANCH && (cd <dir> && git push -u origin <branch>)"