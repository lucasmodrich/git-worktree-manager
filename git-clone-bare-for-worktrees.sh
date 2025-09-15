#!/usr/bin/env bash

# Strict error handling
set -euo pipefail

# Script metadata
readonly SCRIPT_NAME="$(basename "$0")"
readonly SCRIPT_VERSION="2.0.0"

# Configuration variables
VERBOSE=false
DRY_RUN=false
NO_COLOR=false
AUTO_YES=false
CLEANUP_ON_ERROR=true

# Color definitions (will be disabled if NO_COLOR=true)
declare -r RED='\033[0;31m'
declare -r GREEN='\033[0;32m'
declare -r YELLOW='\033[1;33m'
declare -r BLUE='\033[0;34m'
declare -r PURPLE='\033[0;35m'
declare -r CYAN='\033[0;36m'
declare -r WHITE='\033[1;37m'
declare -r BOLD='\033[1m'
declare -r NC='\033[0m' # No Color

# Global variables
PROJECT_DIR=""
REPO_URL=""
MAIN_BRANCH=""
WORKTREE_NAME=""
CREATED_DIRS=()

# Cleanup function
cleanup_on_error() {
    local exit_code=$?
    if [[ "$CLEANUP_ON_ERROR" == true && ${#CREATED_DIRS[@]} -gt 0 ]]; then
        print_warning "Cleaning up created directories due to error..."
        for dir in "${CREATED_DIRS[@]}"; do
            if [[ -d "$dir" ]]; then
                print_verbose "Removing directory: $dir"
                rm -rf "$dir" || true
            fi
        done
    fi
    exit $exit_code
}

# Set up error trap
trap cleanup_on_error ERR

# Utility functions
print_color() {
    local color=$1
    local text=$2
    if [[ "$NO_COLOR" == false ]]; then
        echo -e "${color}${text}${NC}"
    else
        echo "$text"
    fi
}

print_success() {
    print_color "$GREEN" "✅ $1"
}

print_error() {
    print_color "$RED" "❌ Error: $1" >&2
}

print_warning() {
    print_color "$YELLOW" "⚠️  $1" >&2
}

print_info() {
    print_color "$CYAN" "ℹ️  $1"
}

print_header() {
    print_color "$BOLD$BLUE" "$1"
}

print_verbose() {
    if [[ "$VERBOSE" == true ]]; then
        print_color "$PURPLE" "🔍 $1"
    fi
}

# Validation functions
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

validate_dependencies() {
    local missing_deps=()
    
    if ! command_exists git; then
        missing_deps+=("git")
    fi
    
    if [[ ${#missing_deps[@]} -gt 0 ]]; then
        print_error "Missing required dependencies: ${missing_deps[*]}"
        print_info "Please install the missing dependencies and try again."
        exit 1
    fi
}

sanitize_dirname() {
    local name="$1"
    # Remove/replace problematic characters, keep alphanumeric, dots, hyphens, underscores
    name=$(echo "$name" | tr -cd '[:alnum:]._-' | sed 's/^[.-]*//' | sed 's/[.-]*$//')
    
    # Ensure it's not empty and not a git internal name
    if [[ -z "$name" || "$name" =~ ^\.git || "$name" == "." || "$name" == ".." ]]; then
        name="repo-$(date +%s)"
    fi
    
    echo "$name"
}

validate_git_url() {
    local url="$1"
    
    print_verbose "Validating Git URL: $url"
    
    # Basic format check
    if [[ ! "$url" =~ ^(https?://|git@|ssh://|file://) ]]; then
        print_error "Invalid Git URL format: $url"
        print_info "Expected formats: https://..., git@..., ssh://..., or file://..."
        return 1
    fi
    
    # Test if URL is accessible (with timeout)
    if ! timeout 10 git ls-remote "$url" >/dev/null 2>&1; then
        print_error "Cannot access Git repository: $url"
        print_info "Please check the URL and your network connection."
        return 1
    fi
    
    print_verbose "Git URL validation successful"
    return 0
}

# User input function (moved outside conditional blocks)
get_user_choice() {
    local prompt="$1"
    local default="$2"
    local choice
    
    if [[ "$AUTO_YES" == true ]]; then
        echo "$default"
        return 0
    fi
    
    while true; do
        if [[ "$NO_COLOR" == false ]]; then
            echo -e -n "${CYAN}$prompt ${WHITE}[$default]${NC}: "
        else
            echo -n "$prompt [$default]: "
        fi
        
        read -r choice
        choice="${choice:-$default}"
        
        # Validate y/n responses for yes/no questions
        if [[ "$prompt" =~ \(y/n\) ]]; then
            if [[ "$choice" =~ ^[YyNn]([Ee][Ss])?$ ]]; then
                echo "$choice"
                return 0
            else
                print_warning "Please enter 'y' for yes or 'n' for no."
                continue
            fi
        fi
        
        # For other inputs, just return the choice
        echo "$choice"
        return 0
    done
}

# Help function
show_help() {
    cat << EOF
$(print_header "USAGE:")
$SCRIPT_NAME [OPTIONS] <git-url> [directory-name]

$(print_header "DESCRIPTION:")
This script clones a Git repository as a bare repository and sets it up for
Git worktrees. The bare repository is stored in a .bare subdirectory, and
worktrees are created as subdirectories within the main project folder.

$(print_header "PARAMETERS:")
  git-url         The Git repository URL to clone (required)
  directory-name  The name of the directory to create (optional)
                  If not provided, uses the repository name from the URL

$(print_header "OPTIONS:")
  -h, --help      Show this help message
  -v, --verbose   Enable verbose output
  -n, --dry-run   Show what would be done without executing
  -y, --yes       Automatically answer 'yes' to prompts
  --no-color      Disable colored output
  --no-cleanup    Don't cleanup on error
  --version       Show version information

$(print_header "EXAMPLES:")
  $SCRIPT_NAME git@github.com:user/repo.git
  => Creates structure:
     repo/
       .bare/     # Bare repository
       main/      # Main branch worktree (if created)

  $SCRIPT_NAME --verbose git@github.com:user/repo.git my-project
  => Creates structure with verbose output:
     my-project/
       .bare/     # Bare repository
       main/      # Main branch worktree (if created)

  $SCRIPT_NAME --dry-run --yes git@github.com:user/repo.git
  => Shows what would be done without executing

$(print_header "AFTER SETUP:")
You can create additional worktrees like:
  cd repo  # or your chosen directory name
  git worktree add ./feature-branch origin/feature-branch
  git worktree add ./hotfix -b hotfix

EOF
}

# Parse command line arguments
parse_arguments() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -v|--verbose)
                VERBOSE=true
                shift
                ;;
            -n|--dry-run)
                DRY_RUN=true
                shift
                ;;
            -y|--yes)
                AUTO_YES=true
                shift
                ;;
            --no-color)
                NO_COLOR=true
                shift
                ;;
            --no-cleanup)
                CLEANUP_ON_ERROR=false
                shift
                ;;
            --version)
                echo "$SCRIPT_NAME version $SCRIPT_VERSION"
                exit 0
                ;;
            -*)
                print_error "Unknown option: $1"
                print_info "Use --help for usage information."
                exit 1
                ;;
            *)
                if [[ -z "$REPO_URL" ]]; then
                    REPO_URL="$1"
                elif [[ -z "$PROJECT_DIR" ]]; then
                    PROJECT_DIR="$1"
                else
                    print_error "Too many arguments: $1"
                    print_info "Use --help for usage information."
                    exit 1
                fi
                shift
                ;;
        esac
    done
    
    # Validate required arguments
    if [[ -z "$REPO_URL" ]]; then
        print_error "Git URL is required."
        echo ""
        show_help
        exit 1
    fi
}

# Git operations with error handling
safe_git_clone() {
    local url="$1"
    local target="$2"
    
    print_verbose "Cloning bare repository from $url to $target"
    
    if [[ "$DRY_RUN" == true ]]; then
        print_info "[DRY RUN] Would clone: git clone --bare '$url' '$target'"
        return 0
    fi
    
    if ! git clone --bare "$url" "$target"; then
        print_error "Failed to clone repository from $url"
        return 1
    fi
    
    print_verbose "Successfully cloned bare repository"
    return 0
}

safe_git_config() {
    local config_key="$1"
    local config_value="$2"
    
    print_verbose "Setting git config: $config_key = $config_value"
    
    if [[ "$DRY_RUN" == true ]]; then
        print_info "[DRY RUN] Would set config: git config '$config_key' '$config_value'"
        return 0
    fi
    
    if ! git config "$config_key" "$config_value"; then
        print_error "Failed to set git config: $config_key"
        return 1
    fi
    
    return 0
}

safe_git_fetch() {
    print_verbose "Fetching all branches from origin"
    
    if [[ "$DRY_RUN" == true ]]; then
        print_info "[DRY RUN] Would fetch: git fetch origin"
        return 0
    fi
    
    if ! git fetch origin; then
        print_error "Failed to fetch branches from origin"
        return 1
    fi
    
    print_verbose "Successfully fetched branches"
    return 0
}

safe_git_worktree_add() {
    local worktree_path="$1"
    local branch="$2"
    
    print_verbose "Creating worktree: $worktree_path -> $branch"
    
    if [[ "$DRY_RUN" == true ]]; then
        print_info "[DRY RUN] Would create worktree: git worktree add '$worktree_path' '$branch'"
        return 0
    fi
    
    if ! git worktree add "$worktree_path" "$branch"; then
        print_error "Failed to create worktree: $worktree_path"
        return 1
    fi
    
    print_verbose "Successfully created worktree"
    return 0
}

# Main branch detection
detect_main_branch() {
    print_verbose "Detecting main branch"
    
    if [[ "$DRY_RUN" == true ]]; then
        print_info "[DRY RUN] Would detect main branch"
        echo "main"  # Assume main for dry run
        return 0
    fi
    
    local detected_branch=""
    
    if git show-ref --verify --quiet refs/remotes/origin/main; then
        detected_branch="main"
    elif git show-ref --verify --quiet refs/remotes/origin/master; then
        detected_branch="master"
    fi
    
    print_verbose "Detected main branch: ${detected_branch:-none}"
    echo "$detected_branch"
}

# Directory operations
create_project_directory() {
    local dir_name="$1"
    
    print_verbose "Creating project directory: $dir_name"
    
    if [[ -d "$dir_name" ]]; then
        print_error "Directory '$dir_name' already exists."
        print_info "Please choose a different directory name or remove the existing directory."
        exit 1
    fi
    
    if [[ "$DRY_RUN" == true ]]; then
        print_info "[DRY RUN] Would create directory: $dir_name"
        return 0
    fi
    
    if ! mkdir "$dir_name"; then
        print_error "Failed to create directory: $dir_name"
        exit 1
    fi
    
    CREATED_DIRS+=("$dir_name")
    print_verbose "Successfully created directory: $dir_name"
}

# Main execution function
main() {
    print_verbose "Starting $SCRIPT_NAME v$SCRIPT_VERSION"
    
    # Parse arguments
    parse_arguments "$@"
    
    # Validate dependencies
    validate_dependencies
    
    # Validate Git URL
    if ! validate_git_url "$REPO_URL"; then
        exit 1
    fi
    
    # Determine project directory name
    if [[ -z "$PROJECT_DIR" ]]; then
        local repo_basename="${REPO_URL##*/}"
        PROJECT_DIR="${repo_basename%.*}"
    fi
    
    # Sanitize directory name
    PROJECT_DIR=$(sanitize_dirname "$PROJECT_DIR")
    
    print_info "Repository URL: $REPO_URL"
    print_info "Project directory: $PROJECT_DIR"
    
    if [[ "$DRY_RUN" == true ]]; then
        print_header "DRY RUN MODE - No changes will be made"
    fi
    
    # Create project directory
    create_project_directory "$PROJECT_DIR"
    
    if [[ "$DRY_RUN" == false ]]; then
        cd "$PROJECT_DIR" || {
            print_error "Failed to change to directory: $PROJECT_DIR"
            exit 1
        }
    fi
    
    # Clone bare repository
    print_info "Cloning bare repository..."
    if ! safe_git_clone "$REPO_URL" ".bare"; then
        exit 1
    fi
    
    # Create .git file pointing to .bare
    if [[ "$DRY_RUN" == false ]]; then
        echo "gitdir: ./.bare" > .git
    else
        print_info "[DRY_RUN] Would create .git file pointing to .bare"
    fi
    
    # Configure git
    print_info "Configuring remote fetch settings..."
    if ! safe_git_config "remote.origin.fetch" "+refs/heads/*:refs/remotes/origin/*"; then
        exit 1
    fi
    
    # Fetch all branches
    print_info "Fetching all branches from origin..."
    if ! safe_git_fetch; then
        exit 1
    fi
    
    # Detect main branch
    MAIN_BRANCH=$(detect_main_branch)
    
    # Handle main branch worktree creation
    if [[ -n "$MAIN_BRANCH" ]]; then
        echo ""
        print_info "Detected main branch: origin/$MAIN_BRANCH"
        echo ""
        
        local create_worktree
        create_worktree=$(get_user_choice "Would you like to create a worktree for the $MAIN_BRANCH branch? (y/n)" "y")
        
        if [[ "$create_worktree" =~ ^[Yy]([Ee][Ss])?$ ]]; then
            print_info "Creating worktree for $MAIN_BRANCH branch..."
            
            # Ask for worktree directory name
            WORKTREE_NAME=$(get_user_choice "Enter name for the $MAIN_BRANCH worktree directory" "$MAIN_BRANCH")
            WORKTREE_NAME=$(sanitize_dirname "$WORKTREE_NAME")
            
            # Create the worktree
            if ! safe_git_worktree_add "./$WORKTREE_NAME" "origin/$MAIN_BRANCH"; then
                exit 1
            fi
            
            echo ""
            print_success "Successfully created worktree '$WORKTREE_NAME' tracking origin/$MAIN_BRANCH"
            echo ""
            print_header "Your directory structure is now:"
            print_color "$BLUE" "  $PROJECT_DIR/"
            print_color "$PURPLE" "    .bare/           # Bare repository"
            print_color "$PURPLE" "    $WORKTREE_NAME/  # Worktree for $MAIN_BRANCH branch"
            echo ""
            print_header "To create additional worktrees, run:"
            print_color "$GREEN" "  cd $PROJECT_DIR"
            print_color "$GREEN" "  git worktree add ./feature-branch-name origin/branch-name"
            print_color "$YELLOW" "  # or for new branches:"
            print_color "$GREEN" "  git worktree add ./new-feature -b new-feature"
        else
            echo ""
            print_info "Skipped creating worktree. You can create one later with:"
            print_color "$GREEN" "  git worktree add ./$MAIN_BRANCH origin/$MAIN_BRANCH"
        fi
    else
        echo ""
        print_warning "Could not detect a main branch (main or master)."
        print_header "Available remote branches:"
        if [[ "$DRY_RUN" == false ]]; then
            git branch -r
        else
            print_info "[DRY RUN] Would show: git branch -r"
        fi
        echo ""
        print_info "You can create worktrees manually with:"
        print_color "$GREEN" "  git worktree add ./branch-name origin/branch-name"
    fi
    
    # Final success message
    echo ""
    print_success "🎉 Repository setup complete!"
    echo ""
    print_info "Current location: $(pwd)"
    
    if [[ -n "$MAIN_BRANCH" && "$create_worktree" =~ ^[Yy]([Ee][Ss])?$ ]]; then
        print_header "To start working:"
        print_color "$GREEN" "  cd ./$WORKTREE_NAME"
    else
        print_header "Next steps:"
        print_color "$PURPLE" "  Bare repository location: ./.bare"
        print_color "$GREEN" "  Create worktrees with: git worktree add ./worktree-name origin/branch-name"
    fi
    
    print_verbose "Script completed successfully"
}

# Run main function with all arguments
main "$@"
