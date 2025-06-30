'use client';

import { useState, useCallback } from 'react';
import dynamic from 'next/dynamic';
import { useWasm } from '@/hooks/useWasm';
import Header from './Header';
import SampleSelector from './SampleSelector';
import Footer from './Footer';
import { qasmSamples, type QASMSample } from '@/data/samples';

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

  const { isLoading: wasmLoading, isReady: wasmReady, error: wasmError, formatQASM } = useWasm();

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

  const handleCopyOutput = useCallback(async (e: React.MouseEvent) => {
    e.preventDefault();
    if (!outputCode) return;

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
        }

        document.body.removeChild(textarea);
      }

      // Show success feedback
      const button = e.currentTarget as HTMLButtonElement;
      const originalText = button.textContent;
      button.textContent = 'Copied!';
      setTimeout(() => {
        button.textContent = originalText;
      }, 2000);
    } catch (error) {
      console.error('Failed to copy to clipboard:', error);
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
    <div className="min-h-screen bg-base-100 flex flex-col">
      {/* Header */}
      <Header
        onLoadExample={handleLoadExample}
        onClear={handleClearInput}
        onFormat={handleFormat}
        isFormatting={isFormatting}
        canFormat={wasmReady && !!inputCode.trim()}
      />

      {/* Main Content */}
      <div className="flex-1 flex flex-col md:flex-row mx-0 mb-1">
        {/* Input Panel */}
        <div className="flex-1 flex flex-col border-r-0 md:border-r border-base-300 min-h-0 bg-base-100 rounded-t-lg md:rounded-l-lg md:rounded-tr-none shadow-sm">
          <div className="bg-base-200 px-1 md:px-2 py-2 border-b border-base-300 rounded-t-lg md:rounded-tl-lg md:rounded-tr-none flex flex-col sm:flex-row justify-between items-start sm:items-center gap-2">
            <div className="flex-1 min-w-0">
              <h2 className="font-semibold text-sm md:text-base">Input QASM Code</h2>
              <p className="text-xs opacity-70 hidden sm:block">Write or paste your OpenQASM 3.0 code here</p>
            </div>
            <div className="flex gap-1 md:gap-2">
              <button
                className="btn btn-xs md:btn-sm btn-primary whitespace-nowrap"
                onClick={handleLoadExample}
              >
                <svg xmlns="http://www.w3.org/2000/svg" className="h-3 w-3 md:h-4 md:w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
                <span className="hidden sm:inline">Examples</span>
                <span className="sm:hidden">Ex</span>
              </button>
              <button
                className="btn btn-xs md:btn-sm btn-warning whitespace-nowrap"
                onClick={handleClearInput}
              >
                <svg xmlns="http://www.w3.org/2000/svg" className="h-3 w-3 md:h-4 md:w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
                <span className="hidden sm:inline">Clear</span>
                <span className="sm:hidden">✕</span>
              </button>
              <button
                className={`btn btn-xs md:btn-sm btn-accent whitespace-nowrap ${isFormatting ? 'loading' : ''}`}
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
          <div className="flex-1 p-0 min-h-0 editor-container">
            <MonacoEditor
              height="100%"
              language="text"
              theme="vs-light"
              value={inputCode}
              onChange={(value) => setInputCode(value || '')}
              options={{
                minimap: { enabled: false },
                scrollBeyondLastLine: false,
                fontSize: 13,
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
                folding: false,
                lineDecorationsWidth: 5,
                lineNumbersMinChars: 3,
              }}
            />
          </div>
        </div>

        {/* Output Panel */}
        <div className="flex-1 flex flex-col border-t md:border-t-0 border-base-300 min-h-0 bg-base-100 rounded-b-lg md:rounded-r-lg md:rounded-bl-none shadow-sm">
          <div className="bg-base-200 px-1 md:px-2 py-2 border-b border-base-300 md:rounded-tr-lg flex flex-col sm:flex-row justify-between items-start sm:items-center gap-2">
            <div className="flex-1 min-w-0">
              <h2 className="font-semibold text-sm md:text-base">Formatted Output</h2>
              <p className="text-xs opacity-70 hidden sm:block">
                {outputCode ? 'Formatted code ready to copy' : 'Click Format to see results'}
              </p>
            </div>
            {outputCode && (
              <button
                className="btn btn-xs md:btn-sm btn-success whitespace-nowrap"
                onClick={(e) => handleCopyOutput(e)}
              >
                <svg xmlns="http://www.w3.org/2000/svg" className="h-3 w-3 md:h-4 md:w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                </svg>
                Copy
              </button>
            )}
          </div>
          <div className="flex-1 p-0 min-h-0 editor-container">
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
                language="text"
                theme="vs-light"
                value={outputCode}
                options={{
                  readOnly: true,
                  minimap: { enabled: false },
                  scrollBeyondLastLine: false,
                  fontSize: 13,
                  lineNumbers: 'on',
                  roundedSelection: false,
                  scrollbar: { useShadows: false },
                  automaticLayout: true,
                  wordWrap: 'on',
                  contextmenu: true,
                  folding: false,
                  lineDecorationsWidth: 5,
                  lineNumbersMinChars: 3,
                }}
              />
            )}
          </div>
        </div>
      </div>

      {/* Simple Status Bar */}
      {(isFormatting || formatError || outputCode) && (
        <div className="bg-base-200 border-t border-base-300 px-2 py-2">
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
        </div>
      )}

      {/* Footer */}
      <Footer />

      {/* Sample Selector Modal */}
      <SampleSelector
        isOpen={showSampleSelector}
        onClose={() => setShowSampleSelector(false)}
        onSelectSample={handleSelectSample}
      />
    </div>
  );
}
