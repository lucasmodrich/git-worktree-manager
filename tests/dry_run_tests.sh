#!/usr/bin/env bash
set -euo pipefail

# Test dry-run functionality
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$SCRIPT_DIR"

pass=0
fail=0

# Test helper that captures output
test_output() {
    local description="$1"
    shift
    local expected_pattern="$1"
    shift

    echo "Testing: $description"

    # Capture output (allow non-zero exit codes for validation tests)
    if output=$(./git-worktree-manager.sh "$@" 2>&1) || true; then
        if echo "$output" | grep -q "$expected_pattern"; then
            echo "PASS: Found expected pattern '$expected_pattern'"
            pass=$((pass+1))
        else
            echo "FAIL: Expected pattern '$expected_pattern' not found"
            echo "Actual output:"
            echo "$output"
            fail=$((fail+1))
        fi
    else
        echo "FAIL: Command failed or timed out"
        fail=$((fail+1))
    fi
    echo
}

echo "Running dry-run tests..."
echo

# Test dry-run help message
test_output "Help shows dry-run option" "dry-run" --help

# Test dry-run with invalid repo (should show validation without attempting operations)
test_output "Dry-run shows validation error" "Invalid repository format" --dry-run "invalid-repo"

# Test dry-run argument parsing
test_output "Dry-run version check" "git-worktree-manager.sh version" --dry-run --version

# Test dry-run argument order flexibility (regression test for the bug you found)
test_output "Dry-run at end of arguments" "DRY-RUN.*Would create new branch worktree" --new-branch test-branch main --dry-run
test_output "Dry-run at beginning of arguments" "DRY-RUN.*Would create new branch worktree" --dry-run --new-branch test-branch2 main
test_output "Dry-run in middle of arguments" "DRY-RUN.*Would create new branch worktree" --new-branch --dry-run test-branch3 main

echo "Summary:"
echo "Passed: $pass"
echo "Failed: $fail"

if [ $fail -gt 0 ]; then
    echo "Some dry-run tests failed"
    exit 2
fi

echo "All dry-run tests passed"
exit 0