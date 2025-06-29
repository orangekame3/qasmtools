const vscode = require('vscode');
const { LanguageClient } = require('vscode-languageclient/node');
const net = require('net');

let client;

function activate(context) {
    const serverPath = context.asAbsolutePath('bin/qasmlsp');
    const serverOptions = {
        run: {
            command: serverPath,
            args: []
        },
        debug: {
            command: serverPath,
            args: []
        }
    };

    const clientOptions = {
        documentSelector: [{ scheme: 'file', language: 'qasm' }],
        synchronize: {
            fileEvents: vscode.workspace.createFileSystemWatcher('**/*.qasm')
        }
    };

    client = new LanguageClient(
        'qasm',
        'QASM Language Server (Debug)',
        serverOptions,
        clientOptions
    );

    client.start();
}

function deactivate() {
    if (!client) {
        return undefined;
    }
    return client.stop();
}

module.exports = {
    activate,
    deactivate
};
