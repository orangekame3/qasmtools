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

interface WasmModule {
  formatQASM: (code: string) => WasmResult;
  highlightQASM: (code: string) => HighlightResult;
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

      // Load wasm_exec.js
      const wasmExecScript = document.createElement('script');
      wasmExecScript.src = '/qasmtools/wasm/wasm_exec.js';

      await new Promise<void>((resolve, reject) => {
        wasmExecScript.onload = () => resolve();
        wasmExecScript.onerror = () => reject(new Error('Failed to load wasm_exec.js'));
        document.head.appendChild(wasmExecScript);
      });

      // Initialize Go WebAssembly
      const go = new (window as unknown as { Go: new () => { importObject: WebAssembly.Imports; run: (instance: WebAssembly.Instance) => void } }).Go();
      const wasmResponse = await fetch('/qasmtools/wasm/qasmtools.wasm');
      const wasmBytes = await wasmResponse.arrayBuffer();
      const wasmModule = await WebAssembly.instantiate(wasmBytes, go.importObject);

      // Run the Go program
      go.run(wasmModule.instance);

      // Wait for the formatQASM and highlightQASM functions to be available
      await new Promise<void>((resolve) => {
        const checkFunction = () => {
          if ((window as any).qasmToolsReady && (window as any).formatQASM && (window as any).highlightQASM) {
            resolve();
          } else {
            setTimeout(checkFunction, 100);
          }
        };
        checkFunction();
      });

      setWasmModule({
        formatQASM: (window as any).formatQASM,
        highlightQASM: (window as any).highlightQASM
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

  const formatQASM = useCallback((code: string): Promise<WasmResult> => {
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

        const result = wasmModule.formatQASM(code);
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

  // Auto-load WASM on component mount
  useEffect(() => {
    loadWasm();
  }, [loadWasm]);

  return {
    isLoading,
    isReady,
    error,
    formatQASM,
    highlightQASM,
    reload: loadWasm
  };
};
