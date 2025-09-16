# Changelog

All notable changes to this project will be documented in this file.

This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

GitHub tags follow the format: `vX.Y.Z`

---

## [Unreleased] - YYYY-MM-DD

### Added
- Initial release of `git-worktree-manager.sh`
- Full setup mode using GitHub shorthand (`org/repo`)
- Create and push new branch worktrees
- List active worktrees (`--list`)
- Prune stale worktrees (`--prune`)
- Remove worktree and local branch (`--remove`)
- Semantic versioning (`--version`)
- Self-upgrade from GitHub (`--upgrade`)
- Markdown-style help card (`--help`)

### Changed
- Placeholder for improvements

### Fixed
- Placeholder for bug fixes

---

## How to Tag a Release on GitHub

1. Commit your changes:
   ```bash
   git commit -am "Release v1.0.0"
   ```

2. Tag the version:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   ```

3. Push the tag:
   ```bash
   git push origin v1.0.0
   ```

4. Create a release on GitHub:
   - Go to your repo → Releases → “Draft a new release”
   - Choose tag `v1.0.0`
   - Title: `v1.0.0`
   - Paste the changelog section for `v1.0.0`

---

## Future Version Template

## [v1.0.0] - YYYY-MM-DD

### Added
- New features

### Changed
- Enhancements or refactors

### Fixed
- Bug fixes
