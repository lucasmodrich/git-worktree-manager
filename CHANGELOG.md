## [1.1.2](https://github.com/lucasmodrich/git-worktree-manager/compare/v1.1.1...v1.1.2) (2025-09-17)


### Bug Fixes

* Update git-worktree-manager.sh to ensure $RAW_BRANCH_URL variable is used correctly. ([0489853](https://github.com/lucasmodrich/git-worktree-manager/commit/0489853b17e00e8a02f6a7b5f730d79c85c3e518))

## [1.1.1](https://github.com/lucasmodrich/git-worktree-manager/compare/v1.1.0...v1.1.1) (2025-09-17)


### Bug Fixes

* Update git-worktree-manager.sh to remove additional trailing \ from RAW_REPO_URL ([1cfea71](https://github.com/lucasmodrich/git-worktree-manager/commit/1cfea7195a0a6b76d087c6d75e39b24cb4def209))

# [1.1.0](https://github.com/lucasmodrich/git-worktree-manager/compare/v1.0.2...v1.1.0) (2025-09-17)


### Features

* **upgrade:** Improve upgrade to also include download of README, VERSION and LICENCE files ([43aa5ac](https://github.com/lucasmodrich/git-worktree-manager/commit/43aa5ac68aa8e66fff0b4a1c2cf8bcfef6b592c0))

## [1.0.2](https://github.com/lucasmodrich/git-worktree-manager/compare/v1.0.1...v1.0.2) (2025-09-17)


### Bug Fixes

* correct the RAW_URL variable where the upgrade is pulled from. ([0002fe7](https://github.com/lucasmodrich/git-worktree-manager/commit/0002fe769ee789bff87fcd32190f0214c91d8e1c))
* merge conflict ([c080d08](https://github.com/lucasmodrich/git-worktree-manager/commit/c080d08e12f3fe1e8c1cdd696bf013c1b18324ba))

## [1.0.1](https://github.com/lucasmodrich/git-worktree-manager/compare/v1.0.0...v1.0.1) (2025-09-17)

### Bug Fixes

* correcting management of the version in git-worktree-manager.sh' ([b292506](https://github.com/lucasmodrich/git-worktree-manager/commit/b2925060d90878fd1da6b0cd45a22dec3aa9cb86))
* ensure the version number update made to git-worktree-manager.sh is committed as part of the release process.' ([445b899](https://github.com/lucasmodrich/git-worktree-manager/commit/445b8997355bf0279e024e33f3c9649e62d88154))
* futher fixes for release management ([9956406](https://github.com/lucasmodrich/git-worktree-manager/commit/99564063e14d70366f8dede5e4385786d6ea144b))
* Improve release process to trigger pre-release from dev branch. Also ensure version number in git-worktree-manager.sh is updated during the release process. ([4d7de9d](https://github.com/lucasmodrich/git-worktree-manager/commit/4d7de9d35bcf61ad92d363c93b9b5dd4d21cc0ba))
* improve script self-upgrade process. Now deploys to /home/modrich/.git-worktree-manager ([be2850b](https://github.com/lucasmodrich/git-worktree-manager/commit/be2850b87fe54e6e5a5b3606a47082dffa8bf450))

## [1.0.1-beta.4](https://github.com/lucasmodrich/git-worktree-manager/compare/v1.0.1-beta.3...v1.0.1-beta.4) (2025-09-17)


### Bug Fixes

* improve script self-upgrade process. Now deploys to /home/modrich/.git-worktree-manager ([be2850b](https://github.com/lucasmodrich/git-worktree-manager/commit/be2850b87fe54e6e5a5b3606a47082dffa8bf450))

## [1.0.1-beta.3](https://github.com/lucasmodrich/git-worktree-manager/compare/v1.0.1-beta.2...v1.0.1-beta.3) (2025-09-17)


### Bug Fixes

* ensure the version number update made to git-worktree-manager.sh is committed as part of the release process.' ([445b899](https://github.com/lucasmodrich/git-worktree-manager/commit/445b8997355bf0279e024e33f3c9649e62d88154))

## [1.0.1-beta.2](https://github.com/lucasmodrich/git-worktree-manager/compare/v1.0.1-beta.1...v1.0.1-beta.2) (2025-09-17)


### Bug Fixes

* correcting management of the version in git-worktree-manager.sh' ([b292506](https://github.com/lucasmodrich/git-worktree-manager/commit/b2925060d90878fd1da6b0cd45a22dec3aa9cb86))
* futher fixes for release management ([9956406](https://github.com/lucasmodrich/git-worktree-manager/commit/99564063e14d70366f8dede5e4385786d6ea144b))

## [1.0.1-beta.1](https://github.com/lucasmodrich/git-worktree-manager/compare/v1.0.0...v1.0.1-beta.1) (2025-09-17)


### Bug Fixes

* Improve release process to trigger pre-release from dev branch. Also ensure version number in git-worktree-manager.sh is updated during the release process. ([4d7de9d](https://github.com/lucasmodrich/git-worktree-manager/commit/4d7de9d35bcf61ad92d363c93b9b5dd4d21cc0ba))

# 1.0.0 (2025-09-17)


### Bug Fixes

* trigger clean release after tag removal ([aec3ece](https://github.com/lucasmodrich/git-worktree-manager/commit/aec3ecee1d19d184e10e06e86f8ffe3b35ca9375))
* Update commandline flags to support -h as --help triggers Git help when called via usage: git [-v | --version] [-h | --help] [-C <path>] [-c <name>=<value>] ([f22a7a7](https://github.com/lucasmodrich/git-worktree-manager/commit/f22a7a7b085122c50c1c7513f3344d7f081a23f2))


### chore

* **ci-release:** add automated release workflow with semantic-release ([ffb291f](https://github.com/lucasmodrich/git-worktree-manager/commit/ffb291f3e85d4f8494a0a1c240baac53bf5cdac6))


### Features

* Added new script  that helps to automate the cloning of a repo from GitHub to be used in conjunction with Git Worktrees. Additionally added the README.md to provide documentation for use ([9dfc1f1](https://github.com/lucasmodrich/git-worktree-manager/commit/9dfc1f17297f03dce01d8b27f102af1b6d156574))
* enhance git worktree manager with additional modes and update docs ([685743c](https://github.com/lucasmodrich/git-worktree-manager/commit/685743c248cc999156421cdf5f8ae84463790292))
* Enhance setup script to include automatic pushing of new branches to GitHub ([9221b59](https://github.com/lucasmodrich/git-worktree-manager/commit/9221b59d0364588f71f3ff7e0d88756d238f5c02))
* enhance worktree manager with versioning and self-upgrade ([b8d5576](https://github.com/lucasmodrich/git-worktree-manager/commit/b8d557652bd7e8dbcd8f7fe3c2ae72f35f00b85c))
* Refactor setup script to improve branch creation and detection logic ([1201340](https://github.com/lucasmodrich/git-worktree-manager/commit/12013401869bd6bc5a5126f33853fcf9d9c9f811))


### BREAKING CHANGES

* **ci-release:** Introduces mandatory conventional commit format for future commits via hook
* Updated script structure may require re-sourcing or re-executing in some environments
