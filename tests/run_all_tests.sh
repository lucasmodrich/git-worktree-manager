#!/usr/bin/env bash
set -euo pipefail

# Test runner for git-worktree-manager
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo "🧪 Running all tests for git-worktree-manager..."
echo

failed_tests=()
passed_tests=()

# Function to run a test and track results
run_test() {
    local test_name="$1"
    local test_script="$2"

    echo "▶️  Running $test_name..."
    if ./"$test_script"; then
        echo "✅ $test_name passed"
        passed_tests+=("$test_name")
    else
        echo "❌ $test_name failed"
        failed_tests+=("$test_name")
    fi
    echo
}

# Run all tests
run_test "Version Comparison Tests" "version_compare_tests.sh"
run_test "Input Validation Tests" "input_validation_tests.sh"
run_test "Dry-run Tests" "dry_run_tests.sh"

# Summary
echo "📊 Test Results Summary:"
echo "  Passed: ${#passed_tests[@]}"
echo "  Failed: ${#failed_tests[@]}"

if [ ${#passed_tests[@]} -gt 0 ]; then
    echo
    echo "✅ Passed tests:"
    for test in "${passed_tests[@]}"; do
        echo "   - $test"
    done
fi

if [ ${#failed_tests[@]} -gt 0 ]; then
    echo
    echo "❌ Failed tests:"
    for test in "${failed_tests[@]}"; do
        echo "   - $test"
    done
    echo
    echo "❌ Some tests failed!"
    exit 1
fi

echo
echo "🎉 All tests passed!"
exit 0