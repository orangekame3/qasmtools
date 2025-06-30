'use client';

import { useState, useCallback, useEffect, useRef } from 'react';
import dynamic from 'next/dynamic';
import { useWasm, TokenInfo, Violation } from '@/hooks/useWasm';
import Header from './Header';
import SampleSelector from './SampleSelector';
import Footer from './Footer';
import { qasmSamples, type QASMSample } from '@/data/samples';
import { registerQasmLanguage, tokensToDecorations } from '@/utils/qasmLanguage';

// Dynamic import for Monaco Editor to avoid SSR issues
const MonacoEditor = dynamic(
  () => import('@monaco-editor/react').then(mod => mod.default),
  {
    ssr: false,
    loading: () => (
      <div className="flex items-center justify-center h-full">
        <div className="loading loading-spinner loading-lg text-primary"></div>
      </div>
    )
  }
);

const DEFAULT_QASM = qasmSamples[0].code; // Use first sample as default

export default function Playground() {
  const [inputCode, setInputCode] = useState(DEFAULT_QASM);
  const [outputCode, setOutputCode] = useState('');
  const [isFormatting, setIsFormatting] = useState(false);
  const [formatError, setFormatError] = useState<string | null>(null);
  const [showSampleSelector, setShowSampleSelector] = useState(false);
  const [tokens, setTokens] = useState<TokenInfo[]>([]);
  const [violations, setViolations] = useState<Violation[]>([]);
  const [monacoInstance, setMonacoInstance] = useState<typeof import('monaco-editor') | null>(null);
  const editorRef = useRef<import('monaco-editor').editor.IStandaloneCodeEditor | null>(null);
  const decorationsRef = useRef<string[]>([]);

  const { isLoading: wasmLoading, isReady: wasmReady, error: wasmError, formatQASM, highlightQASM, lintQASM } = useWasm();

  const handleFormat = useCallback(async () => {
    if (!wasmReady || !inputCode.trim()) return;

    setIsFormatting(true);
    setFormatError(null);

    try {
      const result = await formatQASM(inputCode);

      if (result.success && result.formatted) {
        setOutputCode(result.formatted);
      } else {
        setFormatError(result.error || 'Unknown formatting error');
        setOutputCode('');
      }
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Formatting failed';
      setFormatError(errorMessage);
      setOutputCode('');
    } finally {
      setIsFormatting(false);
    }
  }, [inputCode, wasmReady, formatQASM]);

  const updateSyntaxHighlighting = useCallback(async (code: string) => {
    if (!wasmReady || !code.trim() || !editorRef.current || !monacoInstance) return;

    try {
      const result = await highlightQASM(code);

      if (!result) {
        console.warn('No result from highlightQASM');
        return;
      }

      if (result.success && result.tokens) {
        setTokens(result.tokens);

        // Clear previous decorations
        if (decorationsRef.current.length > 0) {
          editorRef.current.deltaDecorations(decorationsRef.current, []);
        }

        // Apply new decorations
        const decorations = tokensToDecorations(result.tokens, monacoInstance);
        decorationsRef.current = editorRef.current.deltaDecorations([], decorations);
      } else if (!result.success) {
        console.warn('Syntax highlighting failed:', result.error);
      }
    } catch (error) {
      console.warn('Syntax highlighting failed:', error);
    }
  }, [wasmReady, highlightQASM, monacoInstance]);

  const updateLinting = useCallback(async (code: string) => {
    if (!wasmReady || !code.trim() || !editorRef.current || !monacoInstance) return;

    try {
      const result = await lintQASM(code);

      if (!result) {
        console.warn('No result from lintQASM');
        return;
      }

      if (result.success && result.violations) {
        setViolations(result.violations);

        // Convert violations to Monaco markers
        const markers = result.violations.map(violation => ({
          startLineNumber: violation.line,
          startColumn: violation.column,
          endLineNumber: violation.line,
          endColumn: violation.column + 10, // Estimate end column
          message: `${violation.message} (${violation.rule.id})`,
          severity: violation.severity === 'error' ?
            monacoInstance.MarkerSeverity.Error :
            violation.severity === 'warning' ?
              monacoInstance.MarkerSeverity.Warning :
              monacoInstance.MarkerSeverity.Info,
          code: violation.rule.id,
          source: 'qasm-lint'
        }));

        // Set markers on the model
        const model = editorRef.current.getModel();
        if (model) {
          monacoInstance.editor.setModelMarkers(model, 'qasm-lint', markers);
        }
      } else if (!result.success) {
        console.warn('Linting failed:', result.error);
        setViolations([]);

        // Clear markers
        const model = editorRef.current.getModel();
        if (model) {
          monacoInstance.editor.setModelMarkers(model, 'qasm-lint', []);
        }
      }
    } catch (error) {
      console.warn('Linting failed:', error);
      setViolations([]);
    }
  }, [wasmReady, lintQASM, monacoInstance]);

  const handleEditorDidMount = useCallback((editor: import('monaco-editor').editor.IStandaloneCodeEditor, monaco: typeof import('monaco-editor')) => {
    editorRef.current = editor;
    setMonacoInstance(monaco);

    // Register QASM language
    registerQasmLanguage(monaco);

    // Initial syntax highlighting
    // Temporarily disabled to debug WASM issues
    // if (inputCode) {
    //   updateSyntaxHighlighting(inputCode);
    // }
  }, [inputCode, updateSyntaxHighlighting]);

  const handleCodeChange = useCallback((value: string | undefined) => {
    const newCode = value || '';
    setInputCode(newCode);

    // Update syntax highlighting and linting with debouncing
    if (wasmReady && newCode.trim()) {
      setTimeout(() => {
        // updateSyntaxHighlighting(newCode); // Temporarily disabled
        // updateLinting(newCode); // Temporarily disabled to debug WASM stability
      }, 500); // 500ms debounce for linting (a bit slower than highlighting)
    } else if (!newCode.trim()) {
      // Clear violations when code is empty
      setViolations([]);
      if (editorRef.current && monacoInstance) {
        const model = editorRef.current.getModel();
        if (model) {
          monacoInstance.editor.setModelMarkers(model, 'qasm-lint', []);
        }
      }
    }
  }, [wasmReady, updateSyntaxHighlighting, updateLinting, monacoInstance]);

  // Update linting when WASM becomes ready
  // Temporarily disabled to debug WASM stability
  // useEffect(() => {
  //   if (wasmReady && inputCode.trim()) {
  //     updateLinting(inputCode);
  //   }
  // }, [wasmReady, inputCode, updateLinting]);

  const handleCopyOutput = useCallback(async (e: React.MouseEvent) => {
    e.preventDefault();
    if (!outputCode) return;

    const button = e.currentTarget as HTMLButtonElement;
    const spanElement = button.querySelector('span');

    try {
      // Check if clipboard API is available
      if (navigator?.clipboard?.writeText) {
        await navigator.clipboard.writeText(outputCode);
      } else {
        // Fallback: Create temporary textarea element
        const textarea = document.createElement('textarea');
        textarea.value = outputCode;
        textarea.style.position = 'fixed';
        textarea.style.left = '-9999px';
        textarea.style.top = '-9999px';
        document.body.appendChild(textarea);
        textarea.focus();
        textarea.select();

        try {
          document.execCommand('copy');
        } catch (err) {
          console.error('Fallback: Oops, unable to copy', err);
          throw new Error('Copy operation failed');
        }

        document.body.removeChild(textarea);
      }

      // Show success feedback
      if (spanElement) {
        const originalText = spanElement.textContent || 'Copy';
        spanElement.textContent = 'Copied!';

        setTimeout(() => {
          if (spanElement) {
            spanElement.textContent = originalText;
          }
        }, 2000);
      }
    } catch (error) {
      console.error('Failed to copy to clipboard:', error);

      // Show error feedback
      if (spanElement) {
        const originalText = spanElement.textContent || 'Copy';
        spanElement.textContent = 'Failed!';

        setTimeout(() => {
          if (spanElement) {
            spanElement.textContent = originalText;
          }
        }, 2000);
      }
    }
  }, [outputCode]);

  const handleClearInput = useCallback(() => {
    setInputCode('');
    setOutputCode('');
    setFormatError(null);
  }, []);

  const handleLoadExample = useCallback(() => {
    setShowSampleSelector(true);
  }, []);

  const handleSelectSample = useCallback((sample: QASMSample) => {
    setInputCode(sample.code);
    setOutputCode('');
    setFormatError(null);
  }, []);

  if (wasmLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <div className="loading loading-spinner loading-lg text-primary mb-4"></div>
          <p className="text-lg">Loading QASM Tools...</p>
        </div>
      </div>
    );
  }

  if (wasmError) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="alert alert-error max-w-md">
          <svg xmlns="http://www.w3.org/2000/svg" className="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <div>
            <h3 className="font-bold">Failed to load WASM module</h3>
            <div className="text-xs">{wasmError}</div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="h-screen bg-base-100 flex flex-col overflow-hidden">
      {/* Header */}
      <Header
        onLoadExample={handleLoadExample}
        onClear={handleClearInput}
        onFormat={handleFormat}
        isFormatting={isFormatting}
        canFormat={wasmReady && !!inputCode.trim()}
      />

      {/* Main Content */}
      <div className="flex-1 flex flex-col md:flex-row mx-0">
        {/* Input Panel */}
        <div className="flex-1 flex flex-col border-r-0 md:border-r border-base-300 min-h-0 bg-base-100 rounded-t-lg md:rounded-l-lg md:rounded-tr-none shadow-sm">
          <div className="bg-base-200 px-2 md:px-4 py-3 border-b border-base-300 rounded-t-lg md:rounded-tl-lg md:rounded-tr-none flex flex-col sm:flex-row justify-between items-start sm:items-center gap-2">
            <div className="flex-1 min-w-0">
              <h2 className="font-semibold text-sm md:text-base">Input QASM Code</h2>
              <p className="text-xs opacity-70 hidden sm:block">Write or paste your OpenQASM 3.0 code here</p>
            </div>
            <div className="flex gap-1 md:gap-2">
              <button
                className="btn btn-sm md:btn-md btn-primary whitespace-nowrap"
                onClick={handleLoadExample}
              >
                <svg xmlns="http://www.w3.org/2000/svg" className="h-3 w-3 md:h-4 md:w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
                <span className="hidden sm:inline">Examples</span>
                <span className="sm:hidden">Ex</span>
              </button>
              <button
                className="btn btn-sm md:btn-md btn-warning whitespace-nowrap"
                onClick={handleClearInput}
              >
                <svg xmlns="http://www.w3.org/2000/svg" className="h-3 w-3 md:h-4 md:w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
                <span className="hidden sm:inline">Clear</span>
                <span className="sm:hidden">✕</span>
              </button>
              <button
                className={`btn btn-sm md:btn-md btn-accent whitespace-nowrap ${isFormatting ? 'loading' : ''}`}
                onClick={handleFormat}
                disabled={!wasmReady || isFormatting || !inputCode.trim()}
              >
                {!isFormatting && (
                  <svg xmlns="http://www.w3.org/2000/svg" className="h-3 w-3 md:h-4 md:w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                  </svg>
                )}
                <span className="hidden sm:inline">{isFormatting ? 'Formatting...' : 'Format'}</span>
                <span className="sm:hidden">{isFormatting ? '...' : 'Go'}</span>
              </button>
            </div>
          </div>
          <div className="flex-1 p-0 min-h-0 editor-container bg-[#1e1e1e]">
            <MonacoEditor
              height="100%"
              language="qasm"
              theme="vs-dark"
              value={inputCode}
              onChange={handleCodeChange}
              onMount={handleEditorDidMount}
              options={{
                minimap: { enabled: false },
                scrollBeyondLastLine: true,
                fontSize: 20,
                lineNumbers: 'on',
                roundedSelection: false,
                scrollbar: { useShadows: false },
                automaticLayout: true,
                tabSize: 2,
                insertSpaces: true,
                wordWrap: 'on',
                contextmenu: true,
                selectOnLineNumbers: true,
                glyphMargin: true,
                folding: true,
                lineDecorationsWidth: 5,
                lineNumbersMinChars: 3,
                suggest: {
                  showKeywords: true,
                  showSnippets: true,
                  showFunctions: true,
                  showConstants: true,
                },
                quickSuggestions: {
                  other: true,
                  comments: false,
                  strings: false,
                },
              }}
            />
          </div>
        </div>

        {/* Output Panel */}
        <div className="flex-1 flex flex-col border-t md:border-t-0 border-base-300 min-h-0 bg-base-100 rounded-b-lg md:rounded-r-lg md:rounded-bl-none shadow-sm">
          <div className="bg-base-200 px-2 md:px-4 py-3 border-b border-base-300 md:rounded-tr-lg flex flex-col sm:flex-row justify-between items-start sm:items-center gap-2">
            <div className="flex-1 min-w-0">
              <h2 className="font-semibold text-sm md:text-base">Formatted Output</h2>
              <p className="text-xs opacity-70 hidden sm:block">
                {outputCode ? 'Formatted code ready to copy' : 'Click Format to see results'}
              </p>
            </div>
            {outputCode && (
              <button
                className="btn btn-sm md:btn-md btn-success whitespace-nowrap"
                onClick={(e) => handleCopyOutput(e)}
              >
                <svg xmlns="http://www.w3.org/2000/svg" className="h-3 w-3 md:h-4 md:w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                </svg>
                <span>Copy</span>
              </button>
            )}
          </div>
          <div className="flex-1 p-0 min-h-0 editor-container bg-[#1e1e1e]">
            {!outputCode && !formatError ? (
              <div className="flex items-center justify-center h-full text-center">
                <div className="px-4">
                  <div className="text-4xl md:text-6xl mb-2 md:mb-4">⚡</div>
                  <h3 className="text-base md:text-lg font-semibold mb-1 md:mb-2">Ready to Format</h3>
                  <p className="text-xs md:text-sm opacity-70 mb-3 md:mb-4">
                    Your formatted QASM code will appear here
                  </p>
                  <button
                    className={`btn btn-sm md:btn-md btn-primary ${isFormatting ? 'loading' : ''}`}
                    onClick={handleFormat}
                    disabled={!wasmReady || isFormatting || !inputCode.trim()}
                  >
                    {isFormatting ? 'Formatting...' : 'Format Code'}
                  </button>
                </div>
              </div>
            ) : (
              <MonacoEditor
                height="100%"
                language="qasm"
                theme="vs-dark"
                value={outputCode}
                options={{
                  readOnly: true,
                  minimap: { enabled: false },
                  scrollBeyondLastLine: true,
                fontSize: 20,
                lineNumbers: 'on',
                  roundedSelection: false,
                  scrollbar: { useShadows: false },
                  automaticLayout: true,
                  wordWrap: 'on',
                  contextmenu: true,
                  folding: true,
                  lineDecorationsWidth: 5,
                  lineNumbersMinChars: 3,
                }}
              />
            )}
          </div>
        </div>
      </div>

      {/* Simple Status Bar */}
      {(isFormatting || formatError || outputCode || violations.length > 0) && (
        <div className="bg-base-200 border-t border-base-300 px-2 py-1 space-y-1">
          {isFormatting && (
            <div className="flex items-center gap-2 text-info">
              <div className="loading loading-spinner loading-sm"></div>
              <span className="text-sm">Formatting your QASM code...</span>
            </div>
          )}

          {formatError && (
            <div className="alert alert-error">
              <svg xmlns="http://www.w3.org/2000/svg" className="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <div>
                <h3 className="font-bold text-sm">Formatting Error</h3>
                <div className="text-xs">{formatError}</div>
              </div>
            </div>
          )}

          {outputCode && !formatError && !isFormatting && (
            <div className="flex items-center gap-2 text-success">
              <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <span className="text-sm font-medium">✨ Code formatted successfully!</span>
            </div>
          )}

          {/* Linting Results */}
          {violations.length > 0 && (
            <div className="space-y-1">
              <div className="flex items-center gap-2">
                <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4 text-warning" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.728-.833-2.498 0L4.316 16.5c-.77.833.192 2.5 1.732 2.5z" />
                </svg>
                <span className="text-sm font-medium">
                  Found {violations.length} issue{violations.length !== 1 ? 's' : ''}:
                  {' '}
                  {violations.filter(v => v.severity === 'error').length} error{violations.filter(v => v.severity === 'error').length !== 1 ? 's' : ''},
                  {' '}
                  {violations.filter(v => v.severity === 'warning').length} warning{violations.filter(v => v.severity === 'warning').length !== 1 ? 's' : ''}
                </span>
              </div>
              <div className="max-h-32 overflow-y-auto space-y-1">
                {violations.map((violation, index) => (
                  <div key={index} className={`text-xs p-2 rounded ${
                    violation.severity === 'error' ? 'bg-error/20 text-error-content' :
                    violation.severity === 'warning' ? 'bg-warning/20 text-warning-content' :
                    'bg-info/20 text-info-content'
                  }`}>
                    <div className="font-medium">
                      Line {violation.line}, Column {violation.column}: {violation.message}
                    </div>
                    <div className="opacity-70">
                      Rule: {violation.rule.id} ({violation.rule.name})
                      {violation.rule.documentationUrl && (
                        <a
                          href={violation.rule.documentationUrl}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="ml-2 underline hover:no-underline"
                        >
                          📖 Learn more
                        </a>
                      )}
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      )}

      {/* Footer is hidden to save space */}

      {/* Sample Selector Modal */}
      <SampleSelector
        isOpen={showSampleSelector}
        onClose={() => setShowSampleSelector(false)}
        onSelectSample={handleSelectSample}
      />
    </div>
  );
}
