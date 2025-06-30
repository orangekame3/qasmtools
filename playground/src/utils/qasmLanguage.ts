import { TokenInfo } from '@/hooks/useWasm';
import type * as monaco from 'monaco-editor';

// QASM token types mapping to Monaco token types
const TOKEN_TYPE_MAP: Record<string, string> = {
  'keyword': 'keyword',
  'operator': 'operator',
  'identifier': 'identifier',
  'number': 'number',
  'string': 'string',
  'comment': 'comment',
  'gate': 'keyword.gate',
  'measurement': 'keyword.measurement',
  'register': 'type',
  'punctuation': 'delimiter',
  'builtin_gate': 'keyword.builtin',
  'builtin_quantum': 'keyword.builtin',
  'builtin_classical': 'function',
  'builtin_constant': 'constant',
  'access_control': 'keyword',
  'extern': 'keyword',
  'hardware_qubit': 'variable.hardware'
};

// Define QASM language configuration for Monaco Editor
export const qasmLanguageConfig: monaco.languages.LanguageConfiguration = {
  comments: {
    lineComment: '//',
    blockComment: ['/*', '*/']
  },
  brackets: [
    ['{', '}'],
    ['[', ']'],
    ['(', ')']
  ],
  autoClosingPairs: [
    { open: '{', close: '}' },
    { open: '[', close: ']' },
    { open: '(', close: ')' },
    { open: '"', close: '"' }
  ],
  surroundingPairs: [
    { open: '{', close: '}' },
    { open: '[', close: ']' },
    { open: '(', close: ')' },
    { open: '"', close: '"' }
  ]
};

// Define QASM syntax highlighting rules
export const qasmMonarchLanguage: monaco.languages.IMonarchLanguage = {
  defaultToken: 'invalid',
  tokenPostfix: '.qasm',

  keywords: [
    'OPENQASM', 'include', 'def', 'cal', 'defcal', 'gate', 'extern',
    'let', 'break', 'continue', 'if', 'else', 'end', 'return', 'for', 'while', 'in',
    'const', 'readonly', 'mutable', 'reset', 'measure'
  ],

  types: [
    'qreg', 'qubit', 'creg', 'bool', 'bit', 'int', 'uint', 'float', 'angle', 'complex', 'void'
  ],

  builtinGates: [
    'x', 'y', 'z', 'h', 's', 't', 'sx', 'sy',
    'rx', 'ry', 'rz', 'p', 'u1', 'u2', 'u3',
    'cx', 'cy', 'cz', 'cnot', 'ch', 'crx', 'cry', 'crz', 'cp', 'cu1', 'cu3',
    'swap', 'iswap', 'dcx', 'ccx', 'cswap', 'toffoli', 'fredkin',
    'U', 'CX', 'CCX', 'id', 'iden'
  ],

  builtinFunctions: [
    'sin', 'cos', 'tan', 'exp', 'ln', 'sqrt',
    'arcsin', 'arccos', 'arctan', 'abs', 'mod', 'pow'
  ],

  builtinConstants: [
    'pi', 'tau', 'euler'
  ],

  operators: [
    '=', '->', '+', '-', '*', '/', '%', '**',
    '==', '!=', '<', '<=', '>', '>=',
    '&', '|', '^', '~', '<<', '>>', '&&', '||', '!',
    '+=', '-=', '*=', '/=', '%=', '&=', '|=', '^=', '<<=', '>>='
  ],

  symbols: /[=><!~?:&|+\-*\/\^%]+/,
  escapes: /\\(?:[abfnrtv\\"']|x[0-9A-Fa-f]{1,4}|u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8})/,

  tokenizer: {
    root: [
      // Identifiers and keywords
      [/[a-zA-Z_$][\w$]*/, {
        cases: {
          '@keywords': 'keyword',
          '@types': 'type',
          '@builtinGates': 'keyword.builtin',
          '@builtinFunctions': 'function',
          '@builtinConstants': 'constant',
          '@default': 'identifier'
        }
      }],

      // Hardware qubits
      [/\$[0-9]+/, 'variable.hardware'],

      // Whitespace
      { include: '@whitespace' },

      // Delimiters and operators
      [/[{}()\[\]]/, '@brackets'],
      [/[<>](?!@symbols)/, '@brackets'],
      [/@symbols/, {
        cases: {
          '@operators': 'operator',
          '@default': ''
        }
      }],

      // Numbers
      [/\d*\.\d+([eE][\-+]?\d+)?/, 'number.float'],
      [/0[xX][0-9a-fA-F]+/, 'number.hex'],
      [/\d+/, 'number'],

      // Delimiter: after number because of .\d floats
      [/[;,.]/, 'delimiter'],

      // Strings
      [/"([^"\\]|\\.)*$/, 'string.invalid'], // non-terminated string
      [/"/, { token: 'string.quote', bracket: '@open', next: '@string' }],
    ],

    comment: [
      [/[^\/*]+/, 'comment'],
      [/\/\*/, 'comment', '@push'], // nested comment
      ["\\*/", 'comment', '@pop'],
      [/[\/*]/, 'comment']
    ],

    string: [
      [/[^\\"]+/, 'string'],
      [/@escapes/, 'string.escape'],
      [/\\./, 'string.escape.invalid'],
      [/"/, { token: 'string.quote', bracket: '@close', next: '@pop' }]
    ],

    whitespace: [
      [/[ \t\r\n]+/, 'white'],
      [/\/\*/, 'comment', '@comment'],
      [/\/\/.*$/, 'comment'],
    ],
  },
};

// Custom theme for QASM highlighting
export const qasmTheme: monaco.editor.IStandaloneThemeData = {
  base: 'vs-dark',
  inherit: true,
  rules: [
    { token: 'keyword', foreground: 'C586C0' }, // Purple
    { token: 'keyword.builtin', foreground: 'C586C0' }, // Purple
    { token: 'keyword.gate', foreground: 'C586C0' }, // Purple
    { token: 'keyword.measurement', foreground: 'C586C0' }, // Purple
    { token: 'type', foreground: '4EC9B0' }, // Teal
    { token: 'function', foreground: 'DCDCAA' }, // Light Yellow
    { token: 'constant', foreground: '569CD6' }, // Blue
    { token: 'number', foreground: 'B5CEA8' }, // Light Green
    { token: 'string', foreground: 'CE9178' }, // Orange
    { token: 'comment', foreground: '6A9955' }, // Green
    { token: 'operator', foreground: 'D4D4D4' }, // Light Gray
    { token: 'delimiter', foreground: 'D4D4D4' }, // Light Gray
    { token: 'variable.hardware', foreground: '9CDCFE' }, // Light Blue
  ],
  colors: {}
};

// Function to register QASM language with Monaco Editor
export function registerQasmLanguage(monacoInstance: typeof monaco) {
  // Register the language
  monacoInstance.languages.register({ id: 'qasm' });

  // Set the language configuration
  monacoInstance.languages.setLanguageConfiguration('qasm', qasmLanguageConfig);

  // Set the monarch tokenizer
  monacoInstance.languages.setMonarchTokensProvider('qasm', qasmMonarchLanguage);

  // Define and set the theme
  monacoInstance.editor.defineTheme('qasm-theme', qasmTheme);
}

// Function to convert highlighter tokens to Monaco decorations (for enhanced highlighting)
export function tokensToDecorations(tokens: TokenInfo[], monacoInstance: typeof monaco): monaco.editor.IModelDeltaDecoration[] {
  return tokens.map(token => ({
    range: new monacoInstance.Range(
      token.line,
      token.column + 1,
      token.line,
      token.column + 1 + token.length
    ),
    options: {
      inlineClassName: `qasm-token-${token.type.replace('_', '-')}`,
      hoverMessage: { value: `Token: ${token.type}` }
    }
  }));
}
