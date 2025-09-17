const fs = require('fs');

module.exports = {
    branches: [
        { name: 'main', prerelease: false }, // normal releases
        { name: 'dev', prerelease: 'beta' }  // pre-release channel
    ],
    plugins: [
        '@semantic-release/commit-analyzer',
        '@semantic-release/release-notes-generator',
        ['@semantic-release/changelog', {
            changelogFile: 'CHANGELOG.md',
        }],
        {
            verifyRelease: (pluginConfig, context) => {
                fs.writeFileSync('VERSION', context.nextRelease.version);
            }
        },
        ['@semantic-release/exec', {
            prepareCmd: "sed -i '/^SCRIPT_VERSION=/s/.*/SCRIPT_VERSION=\"${nextRelease.version}\"/' git-worktree-manager.sh"
        }],
        ['@semantic-release/git', {
            assets: ['CHANGELOG.md', 'VERSION', 'git-worktree-manager.sh'],
            message: 'chore(release): ${nextRelease.version} [skip ci]\n\n${nextRelease.notes}'
        }],
        ['@semantic-release/github', {
            assets: [
                { path: 'git-worktree-manager.sh', label: 'git-worktree-manager.sh' },
                { path: 'README.md', label: 'README.md' },
                { path: 'LICENSE', label: 'License' },
                { path: 'release-package.tar.gz', label: 'Full Package' }
            ]
        }],
    ]
};