const vscode = require('vscode');
const { LanguageClient } = require('vscode-languageclient/node');
const fs = require('fs');
const path = require('path');

let client;
let outputChannel;

function activate(context) {
    try {
        console.log('QASM Extension: Starting activation...');
        
        // Create output channel for extension logs
        outputChannel = vscode.window.createOutputChannel('QASM Extension');
        outputChannel.appendLine('=== QASM EXTENSION ACTIVATION START ===');
        outputChannel.appendLine(`QASM Extension: Activation triggered at ${new Date().toISOString()}`);
        outputChannel.appendLine(`QASM Extension: Extension path: ${context.extensionPath}`);
        outputChannel.appendLine(`QASM Extension: VSCode version: ${vscode.version}`);
        outputChannel.show(); // Immediately show the output channel
        
        // Show a notification to confirm activation
        vscode.window.showInformationMessage('QASM Extension activated successfully!');
    
    const serverPath = context.asAbsolutePath('bin/qasmlsp');
    outputChannel.appendLine(`QASM Extension: Server path: ${serverPath}`);
    
    // Check if server binary exists
    if (!fs.existsSync(serverPath)) {
        outputChannel.appendLine(`QASM Extension: ERROR - Server binary not found at: ${serverPath}`);
        vscode.window.showErrorMessage(`QASM Language Server binary not found at: ${serverPath}`);
        return;
    }
    
    outputChannel.appendLine('QASM Extension: Server binary found, checking permissions...');
    
    const serverOptions = {
        run: {
            command: serverPath,
            args: [],
            options: {
                env: process.env
            }
        },
        debug: {
            command: serverPath,
            args: [],
            options: {
                env: process.env
            }
        }
    };

    const clientOptions = {
        documentSelector: [{ scheme: 'file', language: 'qasm' }],
        synchronize: {
            fileEvents: vscode.workspace.createFileSystemWatcher('**/*.qasm')
        },
        outputChannel: outputChannel,
        traceOutputChannel: outputChannel,
        revealOutputChannelOn: 4 // Never
    };

    outputChannel.appendLine('QASM Extension: Creating Language Client...');
    
    client = new LanguageClient(
        'qasm-lsp',
        'QASM Language Server',
        serverOptions,
        clientOptions
    );

    outputChannel.appendLine('QASM Extension: Starting Language Client...');
    
    client.start().then(() => {
        outputChannel.appendLine('QASM Extension: Language Client started successfully');
        // Enable trace for debugging
        client.trace = 2; // Verbose
        outputChannel.appendLine('QASM Extension: Trace level set to verbose');
    }).catch(error => {
        outputChannel.appendLine(`QASM Extension: Failed to start Language Client: ${error.message}`);
        outputChannel.appendLine(`QASM Extension: Error stack: ${error.stack}`);
        console.error('Failed to start QASM Language Client:', error);
    });
    
    // Register the format command
    const formatCommand = vscode.commands.registerCommand('qasm.format', async () => {
        try {
            outputChannel.appendLine('QASM Extension: Format command invoked');
            outputChannel.show(); // Show the output panel
            
            const editor = vscode.window.activeTextEditor;
            outputChannel.appendLine(`QASM Extension: Active editor: ${editor ? 'found' : 'not found'}`);
            
            if (!editor) {
                const errorMsg = 'No active editor found';
                outputChannel.appendLine(`QASM Extension: ERROR - ${errorMsg}`);
                vscode.window.showErrorMessage(errorMsg);
                return;
            }
            
            outputChannel.appendLine(`QASM Extension: Document language: ${editor.document.languageId}`);
            outputChannel.appendLine(`QASM Extension: Document URI: ${editor.document.uri.toString()}`);
            
            if (editor.document.languageId !== 'qasm') {
                const errorMsg = `Active file is not a QASM file (detected: ${editor.document.languageId})`;
                outputChannel.appendLine(`QASM Extension: ERROR - ${errorMsg}`);
                vscode.window.showErrorMessage(errorMsg);
                return;
            }
            
            outputChannel.appendLine('QASM Extension: Checking LSP client status...');
            if (!client || !client.isRunning()) {
                const errorMsg = 'QASM Language Server is not running';
                outputChannel.appendLine(`QASM Extension: ERROR - ${errorMsg}`);
                vscode.window.showErrorMessage(errorMsg);
                return;
            }
            
            outputChannel.appendLine('QASM Extension: LSP client is running, executing format command...');
            await vscode.commands.executeCommand('editor.action.formatDocument');
            outputChannel.appendLine('QASM Extension: Format command completed successfully');
            
        } catch (error) {
            outputChannel.appendLine(`QASM Extension: EXCEPTION in format command: ${error.message}`);
            outputChannel.appendLine(`QASM Extension: Error stack: ${error.stack}`);
            console.error('QASM Format Command Error:', error);
            vscode.window.showErrorMessage(`Failed to format QASM document: ${error.message}`);
        }
    });
    
    context.subscriptions.push(client);
    context.subscriptions.push(outputChannel);
    context.subscriptions.push(formatCommand);
    
    } catch (error) {
        console.error('QASM Extension: Activation failed:', error);
        if (outputChannel) {
            outputChannel.appendLine(`QASM Extension: ACTIVATION ERROR: ${error.message}`);
            outputChannel.appendLine(`QASM Extension: Error stack: ${error.stack}`);
        }
        vscode.window.showErrorMessage(`QASM Extension failed to activate: ${error.message}`);
        throw error;
    }
}

function deactivate() {
    if (outputChannel) {
        outputChannel.appendLine('QASM Extension: Deactivating...');
    }
    if (!client) {
        return undefined;
    }
    return client.stop();
}

module.exports = {
    activate,
    deactivate
};
