#!/usr/bin/env bash
set -euo pipefail

# Test input validation functions
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$SCRIPT_DIR"

pass=0
fail=0

# Helper function to test repo validation
test_repo_validation() {
    local repo="$1"
    local should_pass="$2"
    local description="$3"

    # Test org/repo format validation regex
    if [[ "$repo" =~ ^[a-zA-Z0-9._-]+/[a-zA-Z0-9._-]+$ ]]; then
        result="pass"
    else
        result="fail"
    fi

    if [ "$should_pass" = "true" ] && [ "$result" = "pass" ]; then
        echo "PASS: $description - '$repo'"
        pass=$((pass+1))
    elif [ "$should_pass" = "false" ] && [ "$result" = "fail" ]; then
        echo "PASS: $description - '$repo'"
        pass=$((pass+1))
    else
        echo "FAIL: $description - '$repo' (expected $should_pass, got $result)"
        fail=$((fail+1))
    fi
}

echo "Running input validation tests..."

# Valid repository formats
test_repo_validation "user/repo" true "basic org/repo"
test_repo_validation "my-org/my-repo" true "hyphens in org/repo"
test_repo_validation "org123/repo456" true "numbers in org/repo"
test_repo_validation "user.name/repo.name" true "dots in org/repo"
test_repo_validation "org_name/repo_name" true "underscores in org/repo"
test_repo_validation "a/b" true "single letter org/repo"

# Invalid repository formats
test_repo_validation "repo" false "missing org"
test_repo_validation "/repo" false "empty org"
test_repo_validation "org/" false "empty repo"
test_repo_validation "org/repo/extra" false "too many parts"
test_repo_validation "org repo/name" false "space in org"
test_repo_validation "org/repo name" false "space in repo"
test_repo_validation "org@/repo" false "invalid character @"
test_repo_validation "org/repo#tag" false "invalid character #"
test_repo_validation "" false "empty string"

echo
echo "Passed: $pass"
echo "Failed: $fail"

if [ $fail -gt 0 ]; then
    echo "Some input validation tests failed"
    exit 2
fi

echo "All input validation tests passed"
exit 0