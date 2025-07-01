const { execSync } = require('child_process');
const path = require('path');
const fs = require('fs');

const targets = [
    { goos: 'linux', goarch: 'amd64' },
    { goos: 'linux', goarch: 'arm64' },
    { goos: 'darwin', goarch: 'amd64' },
    { goos: 'darwin', goarch: 'arm64' },
    { goos: 'windows', goarch: 'amd64' },
    { goos: 'windows', goarch: 'arm64' },
];

function buildLSPServer() {
    try {
        const rootDir = path.resolve(__dirname, '..');
        const binDir = path.join(__dirname, 'bin');
        const mainFile = path.join(rootDir, 'cmd', 'qasmlsp', 'main.go');

        // Create bin directory if it doesn't exist
        if (!fs.existsSync(binDir)) {
            fs.mkdirSync(binDir, { recursive: true });
        }

        console.log('Building LSP server for all platforms...');

        for (const target of targets) {
            const { goos, goarch } = target;
            let outputFile = path.join(binDir, `qasmlsp_${goos}_${goarch}`);
            if (goos === 'windows') {
                outputFile += '.exe';
            }

            console.log(`Building for ${goos}/${goarch}...`);
            execSync(`go build -o "${outputFile}" "${mainFile}"`, {
                cwd: rootDir,
                stdio: 'inherit',
                env: {
                    ...process.env,
                    GOOS: goos,
                    GOARCH: goarch,
                },
            });

            // Make the binary executable (not for windows)
            if (goos !== 'windows') {
                try {
                    fs.chmodSync(outputFile, '755');
                } catch (error) {
                    console.log(`Note: Unable to set executable permissions for ${outputFile}. This is expected on some systems.`);
                }
            }
        }

        console.log('LSP server built successfully for all platforms');

    } catch (error) {
        console.error('Failed to build LSP server:', error.message);
        process.exit(1);
    }
}

buildLSPServer();
