#!/usr/bin/env bash
set -euo pipefail

# Simple test runner for version_gt() from git-worktree-manager.sh
#SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
#source "$SCRIPT_DIR/git-worktree-manager.sh"

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
    IFS=. read -r -a ma <<< "$mainA"
    IFS=. read -r -a mb <<< "$mainB"

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
    IFS=. read -r -a pa <<< "$preA"
    IFS=. read -r -a pb <<< "$preB"
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