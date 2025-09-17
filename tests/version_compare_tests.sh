#!/usr/bin/env bash
set -euo pipefail

# Simple test runner for version_gt() from git-worktree-manager.sh
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
source "$SCRIPT_DIR/git-worktree-manager.sh"

pass=0
fail=0

assert_gt() {
  if version_gt "$1" "$2"; then
    echo "PASS: '$1' > '$2'"
    pass=$((pass+1))
  else
    echo "FAIL: expected '$1' > '$2'"
    fail=$((fail+1))
  fi
}

assert_not_gt() {
  if version_gt "$1" "$2"; then
    echo "FAIL: expected NOT ('$1' > '$2')"
    fail=$((fail+1))
  else
    echo "PASS: '$1' <= '$2'"
    pass=$((pass+1))
  fi
}

echo "Running semver comparison tests..."

assert_not_gt "1.0.0" "1.0.0"
assert_gt "1.0.1" "1.0.0"
assert_gt "1.1.0" "1.0.99"
assert_not_gt "2.0.0" "10.0.0"
assert_not_gt "v1.2.3" "1.2.3"
assert_not_gt "1.0.0-alpha" "1.0.0"
assert_gt "1.0.0" "1.0.0-alpha"
assert_gt "1.0.0-alpha.1" "1.0.0-alpha"
assert_gt "1.0.0-alpha.beta" "1.0.0-alpha.1"
assert_gt "1.0.0-beta.11" "1.0.0-beta.2"
assert_gt "1.0.0" "1.0.0-rc.1"
assert_not_gt "1.0.0+build.1" "1.0.0+other"
assert_gt "1.0.0-alpha" "1.0.0-1"

echo
echo "Passed: $pass"
echo "Failed: $fail"

if [ $fail -gt 0 ]; then
  echo "Some tests failed"
  exit 2
fi

echo "All tests passed"
exit 0