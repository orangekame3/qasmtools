'use client';

import { useState, useEffect, useCallback } from 'react';
import { useRouter } from 'next/router';

interface WasmResult {
  success: boolean;
  formatted?: string;
  error?: string;
}

interface WasmModule {
  formatQASM: (code: string) => WasmResult;
}

export const useWasm = () => {
  const [isLoading, setIsLoading] = useState(false);
  const [isReady, setIsReady] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [wasmModule, setWasmModule] = useState<WasmModule | null>(null);
  const router = useRouter();

  const loadWasm = useCallback(async () => {
    if (isReady || isLoading) return;

    setIsLoading(true);
    setError(null);

    try {
      const basePath = router.basePath || '';
      // Load wasm_exec.js
      const wasmExecScript = document.createElement('script');
      wasmExecScript.src = `${basePath}/wasm/wasm_exec.js`;

      await new Promise<void>((resolve, reject) => {
        wasmExecScript.onload = () => resolve();
        wasmExecScript.onerror = () => reject(new Error('Failed to load wasm_exec.js'));
        document.head.appendChild(wasmExecScript);
      });

      // Initialize Go WebAssembly
      const go = new (window as unknown as { Go: new () => { importObject: WebAssembly.Imports; run: (instance: WebAssembly.Instance) => void } }).Go();
      const wasmResponse = await fetch(`${basePath}/wasm/qasmtools.wasm`);
      const wasmBytes = await wasmResponse.arrayBuffer();
      const wasmModule = await WebAssembly.instantiate(wasmBytes, go.importObject);

      // Run the Go program
      go.run(wasmModule.instance);

      // Wait for the formatQASM function to be available
      await new Promise<void>((resolve) => {
        const checkFunction = () => {
          if ((window as any).formatQASM) {
            resolve();
          } else {
            setTimeout(checkFunction, 100);
          }
        };
        checkFunction();
      });

      setWasmModule({
        formatQASM: (window as any).formatQASM
      });

      setIsReady(true);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to load WASM module';
      setError(errorMessage);
      console.error('WASM loading error:', err);
    } finally {
      setIsLoading(false);
    }
  }, [isReady, isLoading, router.basePath]);

  const formatQASM = useCallback((code: string): Promise<WasmResult> => {
    return new Promise((resolve) => {
      if (!wasmModule) {
        resolve({ success: false, error: 'WASM module not loaded' });
        return;
      }

      try {
        const result = wasmModule.formatQASM(code);
        resolve(result);
      } catch (err) {
        const errorMessage = err instanceof Error ? err.message : 'Unknown error occurred';
        resolve({ success: false, error: errorMessage });
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
    reload: loadWasm
  };
};
