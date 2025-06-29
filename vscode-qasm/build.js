const { execSync } = require('child_process');
const path = require('path');
const fs = require('fs');

function buildLSPServer() {
    try {
        const rootDir = path.resolve(__dirname, '..');
        const binDir = path.join(__dirname, 'bin');
        const mainFile = path.join(rootDir, 'cmd', 'qasmlsp', 'main.go');
        const outputFile = path.join(binDir, 'qasmlsp');

        // Create bin directory if it doesn't exist
        if (!fs.existsSync(binDir)) {
            fs.mkdirSync(binDir, { recursive: true });
        }

        console.log('Building LSP server...');
        execSync(`go build -o "${outputFile}" "${mainFile}"`, {
            cwd: rootDir,
            stdio: 'inherit'
        });
        console.log('LSP server built successfully');

        // Make the binary executable
        try {
            fs.chmodSync(outputFile, '755');
        } catch (error) {
            console.log('Note: Unable to set executable permissions. This is expected on some systems.');
        }
    } catch (error) {
        console.error('Failed to build LSP server:', error.message);
        process.exit(1);
    }
}

buildLSPServer();
