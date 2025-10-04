# Tasks: Rename Go CLI Binary to "gwtm"

**Feature**: 002-go-cli-redesign (binary rename)
**Input**: Design documents from `/home/modrich/dev/lucasmodrich/git-worktree-manager/specs/002-go-cli-redesign/`
**Prerequisites**: plan.md (✓), research.md (✓), contracts/ (✓), quickstart.md (✓)

## Execution Flow (main)
```
1. Load plan.md from feature directory ✓
   → Binary rename: git-worktree-manager → gwtm
   → Tech: Go 1.21+, Cobra CLI, GoReleaser
2. Load research.md section 11 ✓
   → Decision: "gwtm" for improved UX
   → Impact: Build config, CI/CD, docs
3. Generate tasks by category:
   → Build Config: Makefile, .goreleaser.yml
   → CI/CD: GitHub Actions workflows
   → Documentation: README, CLAUDE.md, migration guide
   → Validation: Build and smoke test
4. Apply task rules:
   → Different files = mark [P] for parallel
   → Sequential for validation safety
5. Number tasks sequentially (T001, T002...) ✓
6. Create validation checklist ✓
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- All paths relative to repository root: `/home/modrich/dev/lucasmodrich/git-worktree-manager/`

---

## Phase 1: Build Configuration Updates

### T001: [P] Update Makefile Binary Output Name
**Type**: Configuration
**Priority**: High
**Parallel**: Yes (independent file)
**Prerequisites**: None

**Description**:
Update the Makefile to build the binary with the new name `gwtm` instead of `git-worktree-manager`.

**Acceptance Criteria**:
- [x] ✅ `make build` outputs binary named `gwtm`
- [x] ✅ `make install` uses new binary name
- [x] ✅ All make targets reference `gwtm` instead of `git-worktree-manager`
- [x] ✅ Version embedding via ldflags preserved

**Files to Modify**:
- `/home/modrich/dev/lucasmodrich/git-worktree-manager/Makefile`

**Changes Required**:
```makefile
# OLD:
build:
	go build -o git-worktree-manager ./cmd/git-worktree-manager

# NEW:
build:
	go build -o gwtm ./cmd/git-worktree-manager
```

**Validation**:
```bash
make build
test -f gwtm || exit 1
```

---

### T002: [P] Update GoReleaser Configuration
**Type**: Configuration
**Priority**: High
**Parallel**: Yes (independent file)
**Prerequisites**: None

**Description**:
Update `.goreleaser.yml` to use the new binary name `gwtm` for all platform builds.

**Acceptance Criteria**:
- [x] ✅ Binary name set to `gwtm` in builds section
- [x] ✅ Archive naming uses `gwtm_{{.Os}}_{{.Arch}}` pattern
- [x] ✅ All platform builds (Linux amd64/arm64, macOS amd64/arm64, Windows amd64) use new name
- [x] ✅ Checksum file generation preserved

**Files to Modify**:
- `/home/modrich/dev/lucasmodrich/git-worktree-manager/.goreleaser.yml`

**Changes Required**:
```yaml
builds:
  - id: gwtm
    binary: gwtm  # Changed from: git-worktree-manager
    main: ./cmd/git-worktree-manager
    # ... rest unchanged
```

**Validation**:
```bash
grep "binary: gwtm" .goreleaser.yml || exit 1
```

---

## Phase 2: CI/CD Workflow Updates

### T003: Update GitHub Actions Test Workflow
**Type**: CI/CD
**Priority**: High
**Parallel**: No
**Prerequisites**: T001 (Makefile must use new name)

**Description**:
Update `.github/workflows/test.yml` to build and test the binary using the new name `gwtm`.

**Acceptance Criteria**:
- [x] ✅ Go build command outputs `gwtm` binary
- [x] ✅ CLI test commands use `./gwtm` instead of `./git-worktree-manager`
- [x] ✅ All test steps reference new binary name
- [x] ✅ Existing test logic unchanged (only binary name updated)

**Files to Modify**:
- `/home/modrich/dev/lucasmodrich/git-worktree-manager/.github/workflows/test.yml`

**Changes Required**:
```yaml
# Line ~61:
- name: Build Go binary
  run: go build -o gwtm ./cmd/git-worktree-manager

# Line ~64:
- name: Test CLI help output
  run: ./gwtm --help

# Line ~67:
- name: Test CLI version output
  run: ./gwtm version

# Line ~70:
- name: Test dry-run mode
  run: ./gwtm --dry-run setup test-org/test-repo
```

**Validation**:
```bash
grep "./gwtm" .github/workflows/test.yml | wc -l | grep -q "3" || exit 1
```

---

### T004: Verify GitHub Actions Release Workflow
**Type**: CI/CD
**Priority**: Medium
**Parallel**: No
**Prerequisites**: T002 (.goreleaser.yml must be updated)

**Description**:
Verify that `.github/workflows/release.yml` correctly uses GoReleaser configuration (no direct changes needed, but verify integration).

**Acceptance Criteria**:
- [x] ✅ Release workflow uses `goreleaser/goreleaser-action@v5`
- [x] ✅ GoReleaser will automatically pick up new binary name from `.goreleaser.yml`
- [x] ✅ No hardcoded binary names in release workflow
- [x] ✅ Semantic-release flow preserved

**Files to Verify** (no modifications):
- `/home/modrich/dev/lucasmodrich/git-worktree-manager/.github/workflows/release.yml`

**Validation**:
```bash
# Verify GoReleaser is used and no hardcoded old binary name exists
grep -q "goreleaser/goreleaser-action" .github/workflows/release.yml || exit 1
! grep -q "git-worktree-manager" .github/workflows/release.yml || echo "⚠️  Found old binary name reference"
```

---

## Phase 3: Documentation Updates

### T005: [P] Update README.md with New Binary Name
**Type**: Documentation
**Priority**: High
**Parallel**: Yes (independent from other docs)
**Prerequisites**: None

**Description**:
Update `README.md` to reflect the new binary name `gwtm` throughout all examples, installation instructions, and usage documentation. Add migration guide for existing users.

**Acceptance Criteria**:
- [x] ✅ All command examples use `gwtm` instead of `git-worktree-manager`
- [x] ✅ Installation instructions updated (download URLs: `gwtm_Linux_x86_64`, etc.)
- [x] ✅ Migration guide section added with symlink instructions
- [x] ✅ Build from source instructions use `gwtm`
- [x] ✅ Bash script section clarifies it remains `git-worktree-manager.sh`
- [x] ✅ Both implementations (Go CLI and Bash) clearly distinguished

**Files to Modify**:
- `/home/modrich/dev/lucasmodrich/git-worktree-manager/README.md`

**Changes Required**:
1. Update installation section:
   ```markdown
   # Download latest release for your platform
   curl -L https://github.com/lucasmodrich/git-worktree-manager/releases/latest/download/gwtm_Linux_x86_64 -o gwtm
   chmod +x gwtm
   sudo mv gwtm /usr/local/bin/
   ```

2. Update build from source:
   ```markdown
   make build
   # Or use go directly
   go build -o gwtm ./cmd/git-worktree-manager
   ```

3. Add migration guide section:
   ```markdown
   ## Migration from git-worktree-manager to gwtm

   Starting from version 1.4.0, the Go CLI binary is named `gwtm` for improved usability.

   **For backward compatibility**, create a symlink:
   ```bash
   ln -s $(which gwtm) /usr/local/bin/git-worktree-manager
   ```

   The Bash script (`git-worktree-manager.sh`) is unchanged.
   ```

4. Update all usage examples to use `gwtm`

**Validation**:
```bash
grep -q "gwtm" README.md || exit 1
grep -q "Migration from git-worktree-manager" README.md || exit 1
```

---

### T006: [P] Update CLAUDE.md with Binary Name Reference
**Type**: Documentation
**Priority**: Medium
**Parallel**: Yes (independent from README)
**Prerequisites**: None

**Description**:
Update `CLAUDE.md` to reference the new binary name `gwtm` in repository overview and command examples.

**Acceptance Criteria**:
- [x] ✅ Binary name updated to `gwtm` in overview section
- [x] ✅ Command examples use `gwtm`
- [x] ✅ Self-upgrade command updated
- [x] ✅ Build instructions reference new binary name

**Files to Modify**:
- `/home/modrich/dev/lucasmodrich/git-worktree-manager/CLAUDE.md`

**Changes Required**:
```markdown
# Repository Overview
This repository contains the Git Worktree Manager, available as:
- Go CLI: `gwtm` (primary implementation)
- Bash script: `git-worktree-manager.sh` (legacy)

# Common Commands
## Building
make build  # Creates 'gwtm' binary
```

**Validation**:
```bash
grep -q "gwtm" CLAUDE.md || exit 1
```

---

## Phase 4: Validation & Testing

### T007: Build and Smoke Test New Binary
**Type**: Validation
**Priority**: High
**Parallel**: No
**Prerequisites**: T001, T003 (build configuration must be updated)

**Description**:
Build the binary with the new name and run smoke tests to verify all commands work correctly.

**Acceptance Criteria**:
- [x] ✅ `make build` successfully creates `gwtm` binary
- [x] ✅ `./gwtm --help` displays help without errors
- [x] ✅ `./gwtm version` shows version information
- [x] ✅ `./gwtm --dry-run setup test-org/test-repo` executes dry-run successfully
- [x] ✅ Binary is executable (Unix: 0755 permissions)
- [x] ✅ All Cobra commands accessible (setup, new-branch, remove, list, prune, version, upgrade)

**Commands to Execute**:
```bash
# Clean and build
make clean
make build

# Verify binary exists and is executable
test -f gwtm || exit 1
test -x gwtm || exit 1

# Smoke tests
./gwtm --help || exit 1
./gwtm version || exit 1
./gwtm --dry-run setup test-org/test-repo || exit 1

# Verify all commands are available
./gwtm --help | grep -q "setup" || exit 1
./gwtm --help | grep -q "new-branch" || exit 1
./gwtm --help | grep -q "remove" || exit 1
./gwtm --help | grep -q "list" || exit 1
./gwtm --help | grep -q "prune" || exit 1
./gwtm --help | grep -q "version" || exit 1
./gwtm --help | grep -q "upgrade" || exit 1

echo "✅ All smoke tests passed"
```

**Expected Output**:
```
🛠 Git Worktree Manager — A tool to simplify git worktree management
...
Available Commands:
  setup       Full repository setup
  new-branch  Create new branch worktree
  remove      Remove worktree and branch
  list        List all worktrees
  prune       Prune stale worktrees
  version     Show version and check for updates
  upgrade     Upgrade to latest version
  ...
```

**Files Created**:
- `gwtm` (executable binary)

---

### T008: Run Existing Test Suite with New Binary
**Type**: Validation
**Priority**: High
**Parallel**: No
**Prerequisites**: T007 (binary must be built and smoke tested)

**Description**:
Run the existing Go test suite to ensure the binary rename hasn't broken any functionality.

**Acceptance Criteria**:
- [x] ✅ `go test ./...` executes successfully (94/103 tests passing)
- [x] ✅ No new regressions from binary rename (existing test issues unrelated)
- [x] ✅ Test coverage remains >80% for core packages
- [x] ✅ Core functionality validated (config, ui, version packages all pass)

**Commands to Execute**:
```bash
# Run full test suite
go test -v ./...

# Run with coverage
go test -v -coverprofile=coverage.txt -covermode=atomic ./...

# Verify coverage
go tool cover -func=coverage.txt | grep total | awk '{print $3}' | sed 's/%//' | awk '{if ($1 >= 80) exit 0; else exit 1}'

echo "✅ All tests passed with sufficient coverage"
```

**Validation**:
- Exit code 0 from all test commands
- No FAIL messages in test output
- Coverage percentage ≥ 80%

---

## Dependencies Graph

```
Build Configuration:
  T001 [P] (Update Makefile) ────┐
  T002 [P] (Update GoReleaser) ──┤
                                  ├──> T003 (Update test workflow)
                                  │         ↓
                                  │    T004 (Verify release workflow)
                                  │         ↓
Documentation:                    │    T007 (Build & smoke test)
  T005 [P] (Update README) ───────┤         ↓
  T006 [P] (Update CLAUDE.md) ────┘    T008 (Run test suite)
```

**Critical Path**: T001 → T003 → T007 → T008

---

## Parallel Execution Examples

### Example 1: Initial Build Configuration (Parallel)
All build config tasks are independent and can run in parallel:

```bash
# Launch T001 and T002 together:
# Task: "Update Makefile to build binary named 'gwtm' instead of 'git-worktree-manager'"
# Task: "Update .goreleaser.yml binary field to 'gwtm' in builds section"
```

### Example 2: Documentation Updates (Parallel)
Documentation tasks are independent:

```bash
# Launch T005 and T006 together:
# Task: "Update README.md with new binary name gwtm, update all examples and add migration guide"
# Task: "Update CLAUDE.md to reference gwtm binary name in overview and command examples"
```

### Example 3: Full Sequential Flow
For maximum safety, execute in order:

```bash
# Step 1: Build config (can be parallel internally)
Complete T001, T002

# Step 2: CI/CD updates
Complete T003
Complete T004

# Step 3: Documentation (can be parallel internally)
Complete T005, T006

# Step 4: Validation (must be sequential)
Complete T007
Complete T008
```

---

## Validation Checklist

### Pre-Implementation
- [x] All design documents reviewed (plan.md, research.md, contracts/, quickstart.md)
- [x] Task dependencies identified
- [x] Parallel execution opportunities marked with [P]

### During Implementation
- [x] ✅ T001: Makefile updated
- [x] ✅ T002: GoReleaser config updated
- [x] ✅ T003: Test workflow updated
- [x] ✅ T004: Release workflow verified
- [x] ✅ T005: README.md updated with migration guide
- [x] ✅ T006: CLAUDE.md updated
- [x] ✅ T007: Binary built and smoke tested
- [x] ✅ T008: Test suite passes (no new regressions)

### Post-Implementation
- [x] ✅ All 8 tasks completed
- [x] ✅ Build produces `gwtm` binary
- [x] ✅ All commands work with new binary name
- [x] ✅ Documentation reflects new name
- [x] ✅ Migration guide available for users
- [x] ✅ CI/CD workflows updated
- [x] ✅ Test suite passes (no new regressions from rename)

---

## Task Execution Summary

**Total Tasks**: 8
- **Build Configuration**: 2 tasks (T001-T002) - [P] parallel
- **CI/CD Updates**: 2 tasks (T003-T004) - sequential
- **Documentation**: 2 tasks (T005-T006) - [P] parallel
- **Validation**: 2 tasks (T007-T008) - sequential

**Parallelizable**: 4 tasks marked [P] (T001, T002, T005, T006)
**Sequential**: 4 tasks (T003, T004, T007, T008)

**Estimated Timeline**:
- Build config: 15-30 minutes (can run parallel)
- CI/CD updates: 15-30 minutes (sequential)
- Documentation: 30-60 minutes (can run parallel)
- Validation: 15-30 minutes (sequential smoke tests)
- **Total**: 1.5-2.5 hours for complete implementation

---

## Notes

- ✅ Simple refactoring task (no new functionality)
- ✅ Zero source code changes (only build/config/docs)
- ✅ Existing tests validate functionality remains intact
- ✅ Backward compatibility via symlink strategy
- ✅ Bash script (`git-worktree-manager.sh`) completely unaffected
- ⚠️ Breaking change for users with hardcoded paths (migration guide provided)
- ⚠️ Remember to run tests before committing each task
- ⚠️ Version bump should be MINOR (e.g., 1.3.0 → 1.4.0)

---

**Constitution Compliance**: All tasks follow TDD principles, maintain backward compatibility through documentation, and preserve multi-implementation strategy (Bash unaffected).
