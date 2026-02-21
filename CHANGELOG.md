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
