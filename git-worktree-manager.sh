#!/usr/bin/env bash
#
# git-worktree-manager.sh

SCRIPT_VERSION="1.1.7"
# Allow override via environment variable, default to $HOME/.git-worktree-manager
SCRIPT_FOLDER="${GIT_WORKTREE_MANAGER_HOME:-$HOME/.git-worktree-manager}"
SCRIPT_NAME="git-worktree-manager.sh"
GITHUB_REPO="lucasmodrich/git-worktree-manager"
RAW_BRANCH_URL="https://raw.githubusercontent.com/$GITHUB_REPO/refs/heads/main"
RAW_URL="https://raw.githubusercontent.com/$GITHUB_REPO/refs/heads/main/$SCRIPT_NAME"

# Global flag for dry-run mode
DRY_RUN=false

set -e

# --- Helper: Compare semantic versions ---
# Portable semver comparison (semver 2.0.0 rules, ignores build metadata)
# Returns 0 (true) if first version > second version, otherwise returns 1.
version_gt() {
    local a b mainA mainB preA preB
    local -a ma mb pa pb
    local i len ia ib na nb

    # strip leading 'v' and build metadata (+...)
    a="${1#v}"; b="${2#v}"
    a="${a%%+*}"; b="${b%%+*}"

    # split main and prerelease parts
    if [[ "$a" == *-* ]]; then
        mainA="${a%%-*}"
        preA="${a#*-}"
    else
        mainA="$a"
        preA=""
    fi

    if [[ "$b" == *-* ]]; then
        mainB="${b%%-*}"
        preB="${b#*-}"
    else
        mainB="$b"
        preB=""
    fi

    # split major.minor.patch (missing parts treated as 0)
    local oldIFS=$IFS
    IFS=. read -r -a ma <<< "$mainA"
    IFS=. read -r -a mb <<< "$mainB"
    IFS=$oldIFS

    for i in 0 1 2; do
        na="${ma[i]:-0}"
        nb="${mb[i]:-0}"
        # numeric comparison (use 10# to avoid octal interpretation)
        if (( 10#${na} > 10#${nb} )); then
            return 0
        fi
        if (( 10#${na} < 10#${nb} )); then
            return 1
        fi
    done

    # main versions equal -> handle prerelease precedence
    # Absence of prerelease (a release) has higher precedence than any prerelease
    if [[ -z "$preA" && -z "$preB" ]]; then
        return 1  # equal
    fi
    if [[ -z "$preA" && -n "$preB" ]]; then
        return 0  # release > prerelease
    fi
    if [[ -n "$preA" && -z "$preB" ]]; then
        return 1  # prerelease < release
    fi

    # both have prerelease -> compare dot-separated identifiers
    local oldIFS=$IFS
    IFS=. read -r -a pa <<< "$preA"
    IFS=. read -r -a pb <<< "$preB"
    IFS=$oldIFS
    len=${#pa[@]}
    [[ ${#pb[@]} -gt $len ]] && len=${#pb[@]}

    for ((i=0;i<len;i++)); do
        ia="${pa[i]:-}"
        ib="${pb[i]:-}"

        # if one identifier list is shorter
        if [[ -z "$ia" && -n "$ib" ]]; then
            return 1
        elif [[ -n "$ia" && -z "$ib" ]]; then
            return 0
        fi

        # both present: numeric vs alphanumeric rules
        if [[ "$ia" =~ ^[0-9]+$ ]] && [[ "$ib" =~ ^[0-9]+$ ]]; then
            # numeric identifiers compared numerically
            if (( 10#${ia} > 10#${ib} )); then
                return 0
            fi
            if (( 10#${ia} < 10#${ib} )); then
                return 1
            fi
        elif [[ "$ia" =~ ^[0-9]+$ ]] && ! [[ "$ib" =~ ^[0-9]+$ ]]; then
            # numeric identifiers have lower precedence than non-numeric
            return 1
        elif ! [[ "$ia" =~ ^[0-9]+$ ]] && [[ "$ib" =~ ^[0-9]+$ ]]; then
            return 0
        else
            # both non-numeric: ASCII lexical comparison
            if [[ "$ia" > "$ib" ]]; then
                return 0
            fi
            if [[ "$ia" < "$ib" ]]; then
                return 1
            fi
        fi
        # otherwise equal -> continue to next identifier
    done

    # all parts equal
    return 1
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

    if [ "$DRY_RUN" = true ]; then
        echo "🔍 [DRY-RUN] Would create new branch worktree:"
        echo "  Branch: $new_branch"
        echo "  Base: $base_branch"
        echo "🔍 [DRY-RUN] Would fetch latest from origin"

        if git show-ref --verify --quiet "refs/heads/$new_branch"; then
            echo "🔍 [DRY-RUN] Branch '$new_branch' exists locally — would create worktree from it"
        else
            echo "🔍 [DRY-RUN] Would create new branch '$new_branch' from '$base_branch'"
            echo "🔍 [DRY-RUN] Would push new branch '$new_branch' to origin"
        fi
        echo "🔍 [DRY-RUN] Would list all worktrees"
        return 0
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
    local remove_remote=false

    if [ -z "$branch" ]; then
        echo "Usage: $0 --remove <branch-name> [--remote]"
        exit 1
    fi

    # Check for --remote flag in remaining arguments
    shift
    while [[ $# -gt 0 ]]; do
        case $1 in
            --remote)
                remove_remote=true
                shift
                ;;
            *)
                echo "❌ Unknown option: $1"
                echo "Usage: $0 --remove <branch-name> [--remote]"
                exit 1
                ;;
        esac
    done

    if ! git worktree list | grep -q "/$branch "; then
        echo "❌ Worktree for branch '$branch' not found."
        exit 1
    fi

    if [ "$DRY_RUN" = true ]; then
        echo "🔍 [DRY-RUN] Would remove worktree '$branch'"
        if git show-ref --verify --quiet "refs/heads/$branch"; then
            echo "🔍 [DRY-RUN] Would delete local branch '$branch'"
        fi
        if [ "$remove_remote" = true ]; then
            if git show-ref --verify --quiet "refs/remotes/origin/$branch"; then
                echo "🔍 [DRY-RUN] Would delete remote branch 'origin/$branch'"
            fi
        fi
        return 0
    fi

    echo "🗑 Removing worktree '$branch'"
    git worktree remove "$branch"

    if git show-ref --verify --quiet "refs/heads/$branch"; then
        echo "🧨 Deleting local branch '$branch'"
        git branch -D "$branch"
    else
        echo "⚠️ Local branch '$branch' not found — nothing to delete."
    fi

    if [ "$remove_remote" = true ]; then
        if git show-ref --verify --quiet "refs/remotes/origin/$branch"; then
            echo "☁️ Deleting remote branch 'origin/$branch'"
            git push origin --delete "$branch"
        else
            echo "⚠️ Remote branch 'origin/$branch' not found — nothing to delete."
        fi
    fi

    echo "✅ Removal complete."
}

# --- Helper: Check Upgrade script ---
check_upgrade_available_script() {
    echo "🔍 Checking for newer version on GitHub..."

    remote_version=$(curl -s "$RAW_URL" | grep '^SCRIPT_VERSION=' | cut -d'"' -f2)

    if [ -z "$remote_version" ]; then
        echo "❌ Could not retrieve remote version."
        return 1  # Return failure instead of exit 0
    fi

    echo "🔢 Local version: $SCRIPT_VERSION"
    echo "🌐 Remote version: $remote_version"

    if version_gt "$remote_version" "$SCRIPT_VERSION"; then
        echo "$remote_version > $SCRIPT_VERSION"
        return 0    # upgrade available
    else
        echo "$remote_version <= $SCRIPT_VERSION"
        return 1    # no upgrade available
    fi
}


# --- Helper: Show version ---
show_version() {
    echo "git-worktree-manager.sh version $SCRIPT_VERSION"
    if  check_upgrade_available_script; then
        echo "⬇️ Run '$0 --upgrade' to upgrade to version $remote_version."
    else        
        echo "✅ You already have the latest version."
    fi 
}

# --- Helper: Upgrade script ---
upgrade_script() { 
    mkdir -p "$SCRIPT_FOLDER"

    if check_upgrade_available_script; then
        echo "⬇️ Upgrading to version $remote_version..."

        # Download with error checking
        if ! curl -fsSL -o "$SCRIPT_FOLDER/$SCRIPT_NAME.tmp" "$RAW_URL"; then
            echo "❌ Failed to download script"
            exit 1
        fi
        mv "$SCRIPT_FOLDER/$SCRIPT_NAME.tmp" "$SCRIPT_FOLDER/$SCRIPT_NAME"
        chmod +x "$SCRIPT_FOLDER/$SCRIPT_NAME"
        echo "✓ Script downloaded"

        if ! curl -fsSL -o "$SCRIPT_FOLDER/README.md" "$RAW_BRANCH_URL/README.md"; then
            echo "⚠️  Failed to download README.md (continuing)"
        else
            echo "✓ README.md downloaded"
        fi

        if ! curl -fsSL -o "$SCRIPT_FOLDER/VERSION" "$RAW_BRANCH_URL/VERSION"; then
            echo "⚠️  Failed to download VERSION (continuing)"
        else
            echo "✓ VERSION downloaded"
        fi

        if ! curl -fsSL -o "$SCRIPT_FOLDER/LICENSE" "$RAW_BRANCH_URL/LICENSE"; then
            echo "⚠️  Failed to download LICENSE (continuing)"
        else
            echo "✓ LICENSE downloaded"
        fi

        echo "✅ Upgrade complete. Now running version $remote_version."
    else
        echo "✅ You already have the latest version."
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
  $0 --remove <branch> [--remote]     # Remove worktree and local branch
  $0 --list                           # List active worktrees
  $0 --prune                          # Prune stale worktrees
  $0 --version                        # Show script version
  $0 --upgrade                        # Upgrade to latest version
  $0 --help (-h)                      # Show this help card

Global Options:
  --dry-run                           # Preview actions without executing

Examples:
  $0 acme/webapp
  $0 --new-branch feature/login-page
  $0 --remove feature/login-page --remote
  $0 --dry-run --new-branch feature/test

Notes:
  - Run from repo root (where .git points to .bare)
  - New branches are pushed to GitHub automatically
  - Use --remote with --remove to also delete the remote branch
  - Installation directory: ${SCRIPT_FOLDER}

EOF
}

# --- Mode: Help ---
# --- Helper: Execute command with dry-run support ---
run_command() {
    if [ "$DRY_RUN" = true ]; then
        echo "🔍 [DRY-RUN] Would execute: $*"
    else
        "$@"
    fi
}

# Only run CLI/top-level logic when executed directly (not when sourced)
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then

    # Parse global options from anywhere in the argument list
    FILTERED_ARGS=()
    while [[ $# -gt 0 ]]; do
        case $1 in
            --dry-run)
                DRY_RUN=true
                shift
                ;;
            --help|-h)
                show_help
                exit 0
                ;;
            *)
                FILTERED_ARGS+=("$1")
                shift
                ;;
        esac
    done

    # Restore filtered arguments
    set -- "${FILTERED_ARGS[@]}"

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
        shift
        remove_worktree_and_branch "$@"
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

    # Validate and convert org/repo to full GitHub SSH URL
    if [[ "$1" != git@github.com:* ]]; then
        REPO_PATH="$1"
        # Validate org/repo format
        if ! [[ "$REPO_PATH" =~ ^[a-zA-Z0-9._-]+/[a-zA-Z0-9._-]+$ ]]; then
            echo "❌ Invalid repository format. Expected: org/repo or git@github.com:org/repo.git"
            echo "   Examples: acme/webapp, user123/my-project"
            exit 1
        fi
        REPO_URL="git@github.com:$REPO_PATH.git"
    else
        REPO_URL="$1"
        # Extract repo path from SSH URL
        if ! [[ "$REPO_URL" =~ ^git@github\.com:([a-zA-Z0-9._-]+/[a-zA-Z0-9._-]+)\.git$ ]]; then
            echo "❌ Invalid GitHub SSH URL format. Expected: git@github.com:org/repo.git"
            exit 1
        fi
        REPO_PATH="${BASH_REMATCH[1]}"
    fi

    REPO_NAME=$(basename "$REPO_PATH")

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
fi
