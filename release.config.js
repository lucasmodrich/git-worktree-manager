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
            prepare: (pluginConfig, context) => {
                fs.writeFileSync('VERSION', context.nextRelease.version);
            }
        },
        ['@semantic-release/git', {
            assets: ['CHANGELOG.md', 'VERSION'],
            message: 'chore(release): ${nextRelease.version} [skip ci]\n\n${nextRelease.notes}'
        }],
    ]
};
