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
        ['@semantic-release/git', {
            assets: ['CHANGELOG.md', 'VERSION'],
            message: 'chore(release): ${nextRelease.version} [skip ci]\n\n${nextRelease.notes}'
        }],
        ['@semantic-release/github', {
            assets: [
                { path: 'git-worktree-manager.sh', label: 'Bash Script' },
                { path: 'README.md', label: 'Documentation' },
                { path: 'LICENCE', label: 'Licence' },
                { path: 'release-package.tar.gz', label: 'Full Package' }
            ]
        }],
        '@semantic-release/github'
    ]
};