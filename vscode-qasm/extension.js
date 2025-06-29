const vscode = require('vscode');
const { LanguageClient } = require('vscode-languageclient/node');
const path = require('path');

let client;

function activate(context) {
    // LSPサーバーの実行ファイルのパス
    const serverPath = context.asAbsolutePath(path.join('bin', 'qasmlsp'));

    // LSPサーバーの起動オプション
    const serverOptions = {
        command: serverPath,
        transport: 'stdio',
    };

    // LSPクライアントの設定
    const clientOptions = {
        documentSelector: [{ scheme: 'file', language: 'qasm' }],
        synchronize: {
            fileEvents: vscode.workspace.createFileSystemWatcher('**/*.qasm')
        }
    };

    // LSPクライアントの作成と起動
    client = new LanguageClient(
        'qasm',
        'QASM Language Server',
        serverOptions,
        clientOptions
    );

    // クライアントの起動
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
