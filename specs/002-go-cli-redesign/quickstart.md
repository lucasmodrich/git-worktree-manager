# Quickstart: Go CLI Implementation

**Feature**: 002-go-cli-redesign
**Purpose**: Validate end-to-end implementation through executable user scenarios
**Date**: 2025-10-03

---

## Prerequisites

- Go 1.21+ installed
- Git 2.0+ installed and in PATH
- GitHub account with SSH key configured
- Internet connection for GitHub operations

---

## Scenario 1: Initial Setup and Build

**Goal**: Build the Go CLI from source and verify basic functionality.

### Steps

1. **Navigate to repository**:
   ```bash
   cd /path/to/git-worktree-manager
   ```

2. **Initialize Go module** (if not already done):
   ```bash
   go mod init github.com/lucasmodrich/git-worktree-manager
   go mod tidy
   ```

3. **Build the binary**:
   ```bash
   go build -o git-worktree-manager ./cmd/git-worktree-manager
   ```

4. **Verify build**:
   ```bash
   ./git-worktree-manager --version
   ```

### Expected Output
```
git-worktree-manager version 2.0.0
🔍 Checking for newer version on GitHub...
...
```

### Success Criteria
- ✅ Binary builds without errors
- ✅ `--version` command executes successfully
- ✅ Version number displayed
- ✅ Exit code 0

---

## Scenario 2: Full Repository Setup Workflow

**Goal**: Set up a new repository using git worktree workflow.

### Steps

1. **Create test directory**:
   ```bash
   mkdir -p /tmp/quickstart-test
   cd /tmp/quickstart-test
   ```

2. **Run setup** (using public test repo):
   ```bash
   /path/to/git-worktree-manager lucasmodrich/git-worktree-manager
   ```

3. **Verify directory structure**:
   ```bash
   ls -la git-worktree-manager/
   ```

4. **Verify default branch worktree**:
   ```bash
   cd git-worktree-manager
   git worktree list
   ```

5. **Verify .bare repository**:
   ```bash
   cat .git
   git config --list | grep -E '(push.default|branch.autosetup|remote.origin.fetch)'
   ```

### Expected Output
```
📂 Creating project root: git-worktree-manager
📦 Cloning bare repository into .bare
📝 Creating .git file pointing to .bare
⚙️  Configuring Git for auto remote tracking
🔧 Ensuring all remote branches are fetched
📡 Fetching all remote branches
🌱 Creating initial worktree for branch: main
✅ Setup complete!
/tmp/quickstart-test/git-worktree-manager/.bare   (bare)
/tmp/quickstart-test/git-worktree-manager/main    abc1234 [main]
```

### Success Criteria
- ✅ `.bare` directory created
- ✅ `.git` file contains `gitdir: ./.bare`
- ✅ Default branch worktree created (e.g., `main/`)
- ✅ Git config correctly set
- ✅ Exit code 0

### Failure Scenarios
**Test: Run setup again in same directory (should fail)**:
```bash
/path/to/git-worktree-manager lucasmodrich/git-worktree-manager
```

**Expected**:
```
❌ .bare directory already exists in current directory
💡 Remove existing .bare directory or run setup in a different directory
```
- ✅ Exit code 1
- ✅ No changes to existing .bare

---

## Scenario 3: Branch and Worktree Management

**Goal**: Create, list, and remove worktrees.

### Steps (continuing from Scenario 2)

1. **Create new branch worktree**:
   ```bash
   cd /tmp/quickstart-test/git-worktree-manager
   /path/to/git-worktree-manager --new-branch feature/quickstart-test
   ```

2. **List worktrees**:
   ```bash
   /path/to/git-worktree-manager --list
   ```

3. **Verify worktree directory**:
   ```bash
   ls -la feature/
   cd feature/quickstart-test
   git branch --show-current
   ```

4. **Remove worktree (local only)**:
   ```bash
   cd /tmp/quickstart-test/git-worktree-manager
   /path/to/git-worktree-manager --remove feature/quickstart-test
   ```

5. **Verify removal**:
   ```bash
   ls -la | grep feature
   /path/to/git-worktree-manager --list
   ```

### Expected Output

**Create**:
```
📡 Fetching latest from origin
🌱 Creating new branch 'feature/quickstart-test' from 'main'
☁️  Pushing new branch 'feature/quickstart-test' to origin
✅ Worktree for 'feature/quickstart-test' is ready
```

**List**:
```
📋 Active Git worktrees:
/tmp/quickstart-test/git-worktree-manager/.bare                   (bare)
/tmp/quickstart-test/git-worktree-manager/main                    abc1234 [main]
/tmp/quickstart-test/git-worktree-manager/feature/quickstart-test def5678 [feature/quickstart-test]
```

**Remove**:
```
🗑 Removing worktree 'feature/quickstart-test'
🧨 Deleting local branch 'feature/quickstart-test'
✅ Removal complete.
```

### Success Criteria
- ✅ New worktree directory created at `feature/quickstart-test/`
- ✅ Branch pushed to remote
- ✅ `git branch --show-current` shows correct branch
- ✅ Worktree appears in `--list` output
- ✅ After removal, directory gone and not in `--list`
- ✅ Remote branch still exists (not deleted without `--remote`)
- ✅ All exit codes 0

---

## Scenario 4: Dry-Run Mode Validation

**Goal**: Verify dry-run mode prevents actual changes.

### Steps

1. **Dry-run new branch creation**:
   ```bash
   cd /tmp/quickstart-test/git-worktree-manager
   /path/to/git-worktree-manager --dry-run --new-branch feature/dry-run-test
   ```

2. **Verify no changes**:
   ```bash
   ls -la | grep dry-run
   git branch --list | grep dry-run
   /path/to/git-worktree-manager --list | grep dry-run
   ```

3. **Dry-run removal**:
   ```bash
   /path/to/git-worktree-manager --dry-run --remove main
   ```

4. **Verify main still exists**:
   ```bash
   ls -la main/
   ```

### Expected Output

**Dry-run create**:
```
🔍 [DRY-RUN] Would fetch latest from origin
🔍 [DRY-RUN] Would create new branch 'feature/dry-run-test' from 'main'
🔍 [DRY-RUN] Would push new branch 'feature/dry-run-test' to origin
🔍 [DRY-RUN] Would list all worktrees
```

**Dry-run remove**:
```
🔍 [DRY-RUN] Would remove worktree 'main'
🔍 [DRY-RUN] Would delete local branch 'main'
```

### Success Criteria
- ✅ No `feature/dry-run-test` directory created
- ✅ No branch created locally or remotely
- ✅ `main` directory still exists after dry-run removal
- ✅ All dry-run output prefixed with `🔍 [DRY-RUN]`
- ✅ Exit codes 0

---

## Scenario 5: Interactive Prompts (Existing Branch)

**Goal**: Test confirmation prompts when branch already exists.

### Steps

1. **Create branch locally** (without worktree):
   ```bash
   cd /tmp/quickstart-test/git-worktree-manager/main
   git branch feature/prompt-test
   ```

2. **Attempt to create worktree for existing branch**:
   ```bash
   cd /tmp/quickstart-test/git-worktree-manager
   /path/to/git-worktree-manager --new-branch feature/prompt-test
   ```
   **When prompted "Branch 'feature/prompt-test' exists locally. Use it? [y/N]"**: Type `y`

   **When prompted "Branch not found on remote. Push to origin? [y/N]"**: Type `y`

3. **Verify worktree created and branch pushed**:
   ```bash
   ls -la feature/prompt-test
   git ls-remote --heads origin feature/prompt-test
   ```

4. **Clean up**:
   ```bash
   /path/to/git-worktree-manager --remove feature/prompt-test --remote
   ```

### Expected Output
```
📡 Fetching latest from origin
📂 Branch 'feature/prompt-test' exists locally — creating worktree from it
⚠️  Branch 'feature/prompt-test' not found on remote
☁️  Push branch to remote? [y/N]: y
☁️  Pushing branch 'feature/prompt-test' to origin
✅ Worktree for 'feature/prompt-test' is ready
```

### Success Criteria
- ✅ Prompt displayed for existing branch
- ✅ Prompt displayed for missing remote branch
- ✅ User input accepted (y/n)
- ✅ Branch pushed after confirmation
- ✅ Worktree created successfully
- ✅ Exit code 0

---

## Scenario 6: Version and Upgrade Simulation

**Goal**: Test version checking and upgrade workflow (simulated).

### Steps

1. **Check current version**:
   ```bash
   /path/to/git-worktree-manager --version
   ```

2. **Simulate upgrade check** (mocked for testing):
   - For real testing: Modify binary version to be older than latest release
   - Run upgrade command:
   ```bash
   /path/to/git-worktree-manager --upgrade
   ```

3. **Verify version after upgrade**:
   ```bash
   /path/to/git-worktree-manager --version
   ```

### Expected Output

**Version check (upgrade available)**:
```
git-worktree-manager version 1.9.0
🔍 Checking for newer version on GitHub...
🔢 Local version: 1.9.0
🌐 Remote version: 2.0.0
2.0.0 > 1.9.0
⬇️  Run 'git-worktree-manager --upgrade' to upgrade to version 2.0.0.
```

**Upgrade**:
```
🔍 Checking for newer version on GitHub...
⬇️  Upgrading to version 2.0.0...
✓ Binary downloaded
✓ Checksum verified
✓ README.md downloaded
✓ VERSION downloaded
✓ LICENSE downloaded
✅ Upgrade complete. Now running version 2.0.0.
```

### Success Criteria
- ✅ Version comparison correct (semver 2.0.0 rules)
- ✅ Upgrade downloads correct binary for OS/arch
- ✅ Checksum verified before replacing binary
- ✅ Binary replaced atomically (no partial state)
- ✅ Executable permissions preserved
- ✅ Exit codes 0

---

## Scenario 7: Error Handling and Recovery

**Goal**: Verify actionable error messages and safe failure modes.

### Tests

**Test 1: Invalid repository format**:
```bash
/path/to/git-worktree-manager invalid-repo-name
```

**Expected**:
```
❌ Invalid repository format. Expected: org/repo or git@github.com:org/repo.git
💡 Examples: acme/webapp, user123/my-project
```
- ✅ Exit code 1
- ✅ No directories created

**Test 2: Not in worktree-managed repo**:
```bash
cd /tmp
/path/to/git-worktree-manager --new-branch test
```

**Expected**:
```
❌ Not in a worktree-managed repository
💡 Run this command from a directory where .git points to .bare
```
- ✅ Exit code 1

**Test 3: Remove non-existent worktree**:
```bash
cd /tmp/quickstart-test/git-worktree-manager
/path/to/git-worktree-manager --remove nonexistent-branch
```

**Expected**:
```
❌ Worktree for branch 'nonexistent-branch' not found.
💡 Use --list to see available worktrees and branches
```
- ✅ Exit code 1

### Success Criteria
- ✅ All errors include ❌ emoji
- ✅ All errors include 💡 actionable guidance
- ✅ Correct exit codes (1 for user errors)
- ✅ No partial state left on errors

---

## Scenario 8: CLI Compatibility with Bash Version

**Goal**: Verify Go implementation produces identical behavior to Bash version.

### Steps

1. **Set up fresh test environment**:
   ```bash
   mkdir -p /tmp/compat-test-go /tmp/compat-test-bash
   ```

2. **Run identical setup with both implementations**:
   ```bash
   cd /tmp/compat-test-bash
   /path/to/git-worktree-manager.sh lucasmodrich/git-worktree-manager

   cd /tmp/compat-test-go
   /path/to/git-worktree-manager lucasmodrich/git-worktree-manager
   ```

3. **Compare directory structures**:
   ```bash
   diff -r /tmp/compat-test-bash/git-worktree-manager/.git \
            /tmp/compat-test-go/git-worktree-manager/.git
   ```

4. **Compare git configs**:
   ```bash
   cd /tmp/compat-test-bash/git-worktree-manager
   git config --list > /tmp/bash-config.txt

   cd /tmp/compat-test-go/git-worktree-manager
   git config --list > /tmp/go-config.txt

   diff /tmp/bash-config.txt /tmp/go-config.txt
   ```

5. **Compare worktree operations**:
   ```bash
   # Create same branch with both
   cd /tmp/compat-test-bash/git-worktree-manager
   /path/to/git-worktree-manager.sh --new-branch feature/compat-test

   cd /tmp/compat-test-go/git-worktree-manager
   /path/to/git-worktree-manager --new-branch feature/compat-test

   # Compare results
   diff <(cd /tmp/compat-test-bash/git-worktree-manager && git worktree list) \
        <(cd /tmp/compat-test-go/git-worktree-manager && git worktree list)
   ```

### Expected Output
```
[No diff output - structures should be identical]
```

### Success Criteria
- ✅ `.git` file content identical
- ✅ Git config values identical (ignoring timestamps)
- ✅ Worktree list output identical (ignoring paths)
- ✅ Branch tracking configuration identical
- ✅ Both implementations produce same exit codes

---

## Cleanup

After all scenarios complete:

```bash
# Remove test repositories
rm -rf /tmp/quickstart-test
rm -rf /tmp/compat-test-go
rm -rf /tmp/compat-test-bash

# Clean up remote branches created during testing
cd /path/to/git-worktree-manager
git push origin --delete feature/quickstart-test
git push origin --delete feature/prompt-test
git push origin --delete feature/compat-test
```

---

## Test Execution Checklist

Run all scenarios in order and verify success criteria:

- [ ] Scenario 1: Build and basic version check ✅
- [ ] Scenario 2: Full repository setup ✅
- [ ] Scenario 3: Branch and worktree management ✅
- [ ] Scenario 4: Dry-run mode validation ✅
- [ ] Scenario 5: Interactive prompts ✅
- [ ] Scenario 6: Version and upgrade ✅
- [ ] Scenario 7: Error handling ✅
- [ ] Scenario 8: Bash compatibility ✅

**Overall Success**: All scenarios pass with no failures

---

## Performance Benchmarks

Time each operation and compare to Bash version:

| Operation | Bash Time | Go Time | Target |
|-----------|-----------|---------|--------|
| Setup (clone) | ~10s | ? | <10s |
| New branch | ~2s | ? | <2s |
| List worktrees | <1s | ? | <1s |
| Remove worktree | ~1s | ? | <1s |
| Version check | ~2s | ? | <2s |

**Success**: Go version performs at least as fast as Bash version

---

## Notes

- This quickstart serves as both **documentation** and **integration test suite**
- Each scenario maps to acceptance scenarios from spec.md
- All commands should be runnable as-is for validation
- Failures indicate implementation bugs or spec violations

---

## Binary Name Update (2025-10-04)

**Change**: The Go CLI binary has been renamed from `git-worktree-manager` to `gwtm`.

**Impact on Quickstart**:
- All command examples above reference `git-worktree-manager` for clarity
- To execute these scenarios with the new binary name, replace `git-worktree-manager` with `gwtm`
- Build command becomes: `go build -o gwtm ./cmd/git-worktree-manager`

**Updated Quick Reference**:
```bash
# Build:
go build -o gwtm ./cmd/git-worktree-manager

# All commands use "gwtm" instead of "git-worktree-manager":
./gwtm --version
./gwtm <org>/<repo>
./gwtm new-branch <branch>
./gwtm list
./gwtm remove <branch>
./gwtm --help
```

**Note**: The Bash script (`git-worktree-manager.sh`) remains unchanged and can coexist with the Go binary.
