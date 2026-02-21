# Changelog

> **Note:** From v2.0.0 onwards, release notes are published on the
> [GitHub Releases page](https://github.com/lucasmodrich/git-worktree-manager/releases).
> This file is a historical record of changes up to and including v2.0.0.

---

# [2.0.0](https://github.com/lucasmodrich/git-worktree-manager/compare/v1.3.0...v2.0.0) (2026-02-21)


### Features

* Go CLI redesign with GoReleaser-owned releases ([#18](https://github.com/lucasmodrich/git-worktree-manager/issues/18)) ([552d1a2](https://github.com/lucasmodrich/git-worktree-manager/commit/552d1a2c5eeed1cbda354c34777ce6312e2b7aef))


### BREAKING CHANGES

* Binary name changed from git-worktree-manager to gwtm

The shorter binary name (4 chars vs 20 chars) significantly improves
user experience by reducing typing friction. This follows industry
standards where popular CLI tools use memorable short names (gh, kubectl,
helm, hugo).

Changes:
- Update Makefile to build 'gwtm' binary
- Update .goreleaser.yml for multi-platform releases
- Update GitHub Actions workflows to use new binary name
- Add migration guide in README.md with symlink instructions
- Update CLAUDE.md with build instructions
- Add complete Go CLI implementation with Cobra framework
- Update constitution to document multi-implementation strategy

The Bash script (git-worktree-manager.sh) remains unchanged.

Migration path for existing users:
  ln -s $(which gwtm) /usr/local/bin/git-worktree-manager

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>

* fix(tests): resolve test failures in CI environment

- Fix branch_test.go to dynamically detect default branch (main/master)
  instead of hardcoding branch names
- Fix worktree.go WorktreeAdd implementation to correctly pass branch
  argument in both track=true and track=false modes
- Fix worktree_test.go to use proper bare repo + worktree setup that
  matches production workflow
- Update .gitignore to exclude build artifacts and coverage files

These changes address test failures in GitHub Actions CI where git
creates different default branch names depending on version.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>

* fix(tests): detect default branch in remote tests

Update TestFetch, TestPush, and TestDetectDefaultBranch to dynamically
detect the default branch (main/master) instead of hardcoding "main".

This resolves the final CI test failures in GitHub Actions where git
creates different default branch names depending on version.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>

* fix(ci): upgrade Go version to 1.25.1 to fix covdata error

Upgrade from Go 1.21 to 1.25.1 in CI workflow to resolve "no such tool
'covdata'" error when running tests with -race and -coverprofile flags.

Go 1.25.1 is the latest version and properly supports coverage with race
detection.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>

* doc: added coverage.txt

* feat: comprehensive Go CLI review and improvements

Release pipeline:
- Add separate goreleaser.yml workflow triggered on tag push
- Fix release.yml Go version (1.21 -> 1.25.1 to match go.mod)

Binary naming consistency:
- Rename root command Use to gwtm throughout
- Fix GetBinaryPath() and getBinaryName() to use gwtm naming
- Fix go.mod: cobra marked as direct dependency

Self-upgrade robustness:
- Switch version check to GitHub Releases API (authoritative)
- Add HTTP timeouts to all network calls
- Download auxiliary files from release tag, not main branch
- Add cross-filesystem rename fallback (copy+delete)
- Print upgrade progress inline as each step completes
- Fix deprecated strings.Title

Correctness and robustness:
- PrintError now writes to stderr instead of stdout
- Replace verifyWorktreeRepo with findWorktreeRoot that walks up
  from CWD, so all commands work from any subdirectory
- Add branch name slash validation with helpful error message
- Move baseBranch to function scope (was package-level global)
- Remove fragile hardcoded stderr check in WorktreeAdd

Setup command:
- Remove os.Chdir; use absolute paths throughout
- Add cleanup on failure (RemoveAll if any step fails)
- Add --version flag via Cobra native support
- Allow dots in org/repo names in parseRepoSpec regex

CI config:
- Fix release.config.js: VERSION write moved from verifyRelease
  to prepare lifecycle hook

Docs:
- Rewrite README focused on Go CLI; remove all Bash-centric content
- Add CONTRIBUTING.md with setup, commit conventions, and PR guide

Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>

* chore: improve .gitignore, update license year, remove stale todo

- Fix .gitignore bug where node_modules/ and coverage.txt were merged
  into a single path that matched nothing; separate into distinct entries
- Add gwtm.exe, dist/, *.test, .vscode/, .idea/, .DS_Store, Thumbs.db
- Remove stale git-worktree-manager entry (old binary name)
- Update LICENSE copyright year to 2026
- Delete PROJECT_TODO.md (Bash script review artifact, no longer needed)

Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>

* chore: remove Bash script and all associated scaffolding

The Go CLI (gwtm) is now the sole implementation. All Bash-era files
are no longer needed and have been removed.

Deleted:
- git-worktree-manager.sh (deprecated Bash script)
- tests/ (Bash script tests: version comparison, input validation, dry-run)
- AGENTS.md (Bash-centric agent guidelines)
- specs/ (completed Go CLI redesign planning docs + unfilled templates)
- .specify/ (spec-generation framework: templates, constitution, scripts)
- .claude/commands/ (specify/plan/tasks workflow skill implementations)
- coverage.txt (generated test artifact, now in .gitignore)

Updated:
- CLAUDE.md: rewritten to reflect Go CLI only; removed all Bash references
- release.config.js: removed Bash script version update step and
  git-worktree-manager.sh/release-package.tar.gz from release assets

Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>

* feat: inject commit hash and build date into version output

Adds commit and date build vars alongside version, injected via ldflags
in both Makefile (local builds) and .goreleaser.yml (release builds).
Dev builds show commit+date in version output; "dev" builds skip version
comparison and print an informative message instead.

Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>

* doc: improve CLAUDE.md with subcommands, module path, and missing make targets

Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>

* doc: fix CLAUDE.md inaccuracies flagged in code review

- Correct Go version to 1.25.1 (matches go.mod)
- Use exact Cobra Use strings for new-branch and remove signatures
- Move --remote flag from signature to notes column for remove command
- Add go clean to make clean description
- Clarify version command checks GitHub for updates

Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>

* doc: rewrite CHANGELOG.md in Keep a Changelog format

Replaces the semantic-release auto-generated format with the
Keep a Changelog (https://keepachangelog.com/en/1.1.0/) convention.
Adds an [Unreleased] section capturing the Go CLI rewrite, binary
rename, Bash script removal, and CI/test fixes from this branch.
Historical releases v1.0.0–v1.3.0 are preserved with changes mapped
to Added/Changed/Removed/Fixed sections and comparison links at the
bottom.

Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>

* fix: restore executable bit in moveBinary cross-filesystem fallback

os.Create sets mode 0666 (before umask), losing the executable bit
applied to the temp binary. Add os.Chmod(dst, 0755) after the copy
completes. On Windows this is a no-op since executability is determined
by file extension, not permission bits.

Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>

* refactor: second-pass code review improvements

- Extract FetchLatestVersion into version package (version/fetch.go),
  eliminating duplication between version.go and upgrade.go
- Make PromptYesNo accept io.Reader instead of hard-coding os.Stdin,
  enabling unit tests without stdin mocking
- Fix remove.go: use absolute worktree path in WorktreeRemove call and
  look up --remote flag via cmd.Flags() instead of package-level var
- Make UpgradeToLatest accept info/warn callbacks instead of calling
  fmt.Println directly, decoupling progress reporting from UI layer
- Expand parseRepoSpec to accept HTTPS URLs and generic SSH hosts
  (not just github.com), with improved error messages
- Move git.NewClient initialisation before dry-run early-return in setup
  so DryRun is always set before use
- Normalise go.mod directive to go 1.25 (patch omitted per convention)
- Remove stale bash-script CI jobs from test.yml

New tests:
- internal/commands/utils_test.go: TestFindWorktreeRoot (happy path,
  not-found, .git-directory)
- internal/commands/setup_test.go: TestParseRepoSpec (8 cases covering
  shorthand, SSH, HTTPS, invalid inputs)

Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>

* doc: add .claude/ to .gitignore for local Claude settings

* fix: correct glob entry for CHANGELOG.md in release extra_files

* feat: add Taskfile for build, test, install, clean, fmt, and lint tasks

* fix: update checksum name_template quotes and clean up changelog filters

* refactor: remove Makefile and Taskfile for streamlined build process. Using GoReleaser

* refactor: transfer GitHub release ownership to GoReleaser

- Remove @semantic-release/github plugin; semantic-release now only
  handles versioning, CHANGELOG.md, VERSION, and git tag creation
- GoReleaser owns the full GitHub release: notes, binaries, checksums
- Update .goreleaser.yml: drop redundant extra_files and release.github
  config; add changelog.use: github for modern release note generation
- Upgrade goreleaser-action v5 → v6 and pin version to ~> v2
- Upgrade codecov-action v3 → v5
- Remove make references from CLAUDE.md, README.md, and CONTRIBUTING.md
- Rewrite CONTRIBUTING.md release process section with ownership table,
  flow diagram, and version bump rules

Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>

* fix: replace deprecated archives.format with archives.formats

goreleaser check flagged archives.format as deprecated in GoReleaser v2.
Updated to the current archives.formats list syntax.

Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>

# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Inject commit hash and build date into `version` command output
- Comprehensive Go CLI implementation with full test suite
- Spec-kit framework with baseline documentation

### Changed

- Binary renamed from `git-worktree-manager` to `gwtm`
- Improved `.gitignore`, updated license year

### Removed

- Bash script (`git-worktree-manager.sh`) and all associated scaffolding — replaced by the Go CLI

### Fixed

- CI: upgraded Go toolchain to 1.25.1 to resolve `covdata` error
- Tests: detect default branch dynamically in remote integration tests
- Tests: resolved test failures in CI environment

## [1.3.0] - 2025-09-18

### Added

- Progress indicators for long-running operations
- GitHub Actions workflow for PR testing

### Fixed

- Disabled dry-run tests for PRs that were not behaving as intended

## [1.2.0] - 2025-09-18

### Added

- Major improvements to the git-worktree-manager script (additional flags, improved UX)

### Fixed

- `--dry-run` argument parsing bug

## [1.1.7] - 2025-09-17

### Fixed

- Further tweaks to upgrade version-check logic

## [1.1.6] - 2025-09-17

### Fixed

- Test case corrections

## [1.1.5] - 2025-09-17

### Fixed

- Further fixes to the self-upgrade process

## [1.1.4] - 2025-09-17

### Fixed

- Improved upgrade version detection; modularised script for better reusability

## [1.1.3] - 2025-09-17

### Fixed

- Script now echoes progress of downloaded files during upgrade

## [1.1.2] - 2025-09-17

### Fixed

- `$RAW_BRANCH_URL` variable used correctly in upgrade flow

## [1.1.1] - 2025-09-17

### Fixed

- Removed trailing backslash from `RAW_REPO_URL` that broke upgrade downloads

## [1.1.0] - 2025-09-17

### Added

- Upgrade command now also downloads `README`, `VERSION`, and `LICENSE` files

## [1.0.2] - 2025-09-17

### Fixed

- Corrected `RAW_URL` variable used by the self-upgrade command

## [1.0.1] - 2025-09-17

### Fixed

- Version number management in script now handled correctly during release
- Version file committed as part of the automated release process
- Self-upgrade now deploys to `~/.git-worktree-manager`
- Release process: pre-releases triggered correctly from dev branch

## [1.0.0] - 2025-09-17

### Added

- Initial `git-worktree-manager.sh` script for cloning repos as bare clones and managing worktrees
- Support for multiple modes (setup, new-branch, list)
- Automatic pushing of new branches to GitHub remote
- Versioning and self-upgrade capability
- Automated release workflow via semantic-release

### Changed

- Improved branch creation and detection logic
- Command-line flag `-h` used for help (avoids conflict with `git -h` which invokes Git's own help)

### Fixed

- Clean release after stale tag removal

[Unreleased]: https://github.com/lucasmodrich/git-worktree-manager/compare/v1.3.0...HEAD
[1.3.0]: https://github.com/lucasmodrich/git-worktree-manager/compare/v1.2.0...v1.3.0
[1.2.0]: https://github.com/lucasmodrich/git-worktree-manager/compare/v1.1.7...v1.2.0
[1.1.7]: https://github.com/lucasmodrich/git-worktree-manager/compare/v1.1.6...v1.1.7
[1.1.6]: https://github.com/lucasmodrich/git-worktree-manager/compare/v1.1.5...v1.1.6
[1.1.5]: https://github.com/lucasmodrich/git-worktree-manager/compare/v1.1.4...v1.1.5
[1.1.4]: https://github.com/lucasmodrich/git-worktree-manager/compare/v1.1.3...v1.1.4
[1.1.3]: https://github.com/lucasmodrich/git-worktree-manager/compare/v1.1.2...v1.1.3
[1.1.2]: https://github.com/lucasmodrich/git-worktree-manager/compare/v1.1.1...v1.1.2
[1.1.1]: https://github.com/lucasmodrich/git-worktree-manager/compare/v1.1.0...v1.1.1
[1.1.0]: https://github.com/lucasmodrich/git-worktree-manager/compare/v1.0.2...v1.1.0
[1.0.2]: https://github.com/lucasmodrich/git-worktree-manager/compare/v1.0.1...v1.0.2
[1.0.1]: https://github.com/lucasmodrich/git-worktree-manager/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/lucasmodrich/git-worktree-manager/releases/tag/v1.0.0
