'use client';

import { useState, useEffect, useCallback } from 'react';

interface WasmResult {
  success: boolean;
  formatted?: string;
  error?: string;
}

export interface TokenInfo {
  type: string;
  content: string;
  line: number;
  column: number;
  length: number;
}

export interface HighlightResult {
  success: boolean;
  tokens?: TokenInfo[];
  error?: string;
}

export interface Violation {
  file: string;
  line: number;
  column: number;
  severity: string;
  rule_id: string;
  message: string;
  documentation_url: string;
  rule_details?: {
    name: string;
    description: string;
    tags: string[];
    fixable: boolean;
    specification_url: string;
    examples: {
      incorrect: string;
      correct: string;
    };
  };
}

export interface LintResult {
  success: boolean;
  violations?: Violation[];
  error?: string;
  summary?: {
    total: number;
    errors: number;
    warnings: number;
    info: number;
  };
}

interface WasmModule {
  formatQASM: (code: string, unescape?: boolean) => WasmResult;
  highlightQASM: (code: string) => HighlightResult;
  lintQASM: (code: string) => LintResult;
}

export const useWasm = () => {
  const [isLoading, setIsLoading] = useState(false);
  const [isReady, setIsReady] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [wasmModule, setWasmModule] = useState<WasmModule | null>(null);

  const loadWasm = useCallback(async () => {
    if (isReady || isLoading) return;

    setIsLoading(true);
    setError(null);

    try {
      if (typeof window === 'undefined') return; // Skip during SSR

      // Wait for Go to be available
      await new Promise<void>((resolve) => {
        const checkGo = () => {
          if ((window as any).Go) {
            resolve();
          } else {
            setTimeout(checkGo, 100);
          }
        };
        checkGo();
      });

      // Initialize Go WebAssembly
      const go = new (window as any).Go();
      const wasmPath = process.env.NODE_ENV === 'development' ? '/wasm/qasmtools.wasm' : '/qasmtools/wasm/qasmtools.wasm';
      const wasmResponse = await fetch(wasmPath);
      const wasmBytes = await wasmResponse.arrayBuffer();
      const wasmModule = await WebAssembly.instantiate(wasmBytes, go.importObject);

      // Run the Go program
      go.run(wasmModule.instance);

      // Wait for the formatQASM, highlightQASM and lintQASM functions to be available
      await new Promise<void>((resolve) => {
        const checkFunction = () => {
          if ((window as any).qasmToolsReady && (window as any).formatQASM && (window as any).highlightQASM && (window as any).lintQASM) {
            resolve();
          } else {
            setTimeout(checkFunction, 100);
          }
        };
        checkFunction();
      });

      // Store functions
      const formatFunc = (window as any).formatQASM;
      const highlightFunc = (window as any).highlightQASM;
      const lintFunc = (window as any).lintQASM;

      setWasmModule({
        formatQASM: formatFunc,
        highlightQASM: highlightFunc,
        lintQASM: lintFunc
      });

      setIsReady(true);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to load WASM module';
      setError(errorMessage);
      console.error('WASM loading error:', err);
    } finally {
      setIsLoading(false);
    }
  }, [isReady, isLoading]);

  const formatQASM = useCallback((code: string, unescape?: boolean): Promise<WasmResult> => {
    return new Promise((resolve) => {
      if (!wasmModule) {
        resolve({ success: false, error: 'WASM module not loaded' });
        return;
      }

      try {
        // Check if the Go program is still running
        if (!(window as any).qasmToolsReady) {
          resolve({ success: false, error: 'WASM module has exited. Please reload the page.' });
          return;
        }

        const result = wasmModule.formatQASM(code, unescape);
        resolve(result);
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : 'Unknown error occurred';
        // If the error suggests the program has exited, provide helpful message
        if (errorMessage.includes('Go program has already exited')) {
          resolve({ success: false, error: 'WASM module has exited. Please reload the page.' });
        } else {
          resolve({ success: false, error: errorMessage });
        }
      }
    });
  }, [wasmModule]);

  const highlightQASM = useCallback((code: string): Promise<HighlightResult> => {
    return new Promise((resolve) => {
      if (!wasmModule) {
        resolve({ success: false, error: 'WASM module not loaded' });
        return;
      }

      try {
        // Check if the Go program is still running
        if (!(window as any).qasmToolsReady) {
          resolve({ success: false, error: 'WASM module has exited. Please reload the page.' });
          return;
        }

        const result = wasmModule.highlightQASM(code);

        // Check if result is valid
        if (!result || typeof result !== 'object') {
          resolve({ success: false, error: 'Invalid response from WASM module' });
          return;
        }

        resolve(result);
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : 'Unknown error occurred';
        console.error('highlightQASM error:', err);

        // If the error suggests the program has exited, provide helpful message
        if (errorMessage.includes('Go program has already exited')) {
          resolve({ success: false, error: 'WASM module has exited. Please reload the page.' });
        } else {
          resolve({ success: false, error: errorMessage });
        }
      }
    });
  }, [wasmModule]);

  const lintQASM = useCallback((code: string): Promise<LintResult> => {
    return new Promise((resolve) => {
      console.log('lintQASM called with code:', code);
      if (!wasmModule) {
        console.error('WASM module not loaded');
        resolve({ success: false, error: 'WASM module not loaded' });
        return;
      }

      try {
        // Check if the Go program is still running
        if (!(window as any).qasmToolsReady) {
          console.error('WASM module has exited');
          resolve({ success: false, error: 'WASM module has exited. Please reload the page.' });
          return;
        }

        console.log('Calling WASM lintQASM function');
        // Call WASM lintQASM function
        const result = wasmModule.lintQASM(code);
        console.log('WASM lintQASM result:', result);

        // Check if result is valid
        if (!result || typeof result !== 'object') {
          resolve({ success: false, error: 'Invalid response from WASM module' });
          return;
        }

        // Always return violations if they exist, even if there's an error
        resolve({
          success: result.success,
          error: result.error,
          violations: result.violations || [],
          summary: {
            total: result.violations?.length || 0,
            errors: result.violations?.filter(v => v.severity === 'error').length || 0,
            warnings: result.violations?.filter(v => v.severity === 'warning').length || 0,
            info: result.violations?.filter(v => v.severity === 'info').length || 0
          }
        });

      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : 'Unknown error occurred';
        console.error('lintQASM error:', err);

        // If the error suggests the program has exited, provide helpful message
        if (errorMessage.includes('Go program has already exited')) {
          resolve({ success: false, error: 'WASM module has exited. Please reload the page.' });
        } else {
          resolve({ success: false, error: errorMessage });
        }
      }
    });
  }, [wasmModule]);

  // Auto-load WASM on component mount and cleanup on unmount
  useEffect(() => {
    if (!isReady && !isLoading) {
      loadWasm();
    }
  }, [isReady, isLoading, loadWasm]);

  return {
    isLoading,
    isReady,
    error,
    formatQASM,
    highlightQASM,
    lintQASM,
    reload: loadWasm
  };
};
