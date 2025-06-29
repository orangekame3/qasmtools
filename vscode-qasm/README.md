# QASM Language Support for VSCode

This extension provides language support for OpenQASM 3.0 in Visual Studio Code.

## Features

- **Syntax Highlighting**: Full semantic highlighting for QASM keywords, operators, gates, measurements, and more
- **Language Server**: Integrated LSP server for advanced language features
- **Auto-completion**: Bracket and quote auto-closing
- **Comment Support**: Line and block comments
- **Indentation Rules**: Smart indentation for quantum circuits

## Supported Token Types

- Keywords (OPENQASM, include, gate, measure, etc.)
- Operators (+, -, *, /, ==, ->, etc.)  
- Numbers and strings
- Comments (// and /* */)
- Gates and measurements
- Registers (qubit, bit, etc.)
- Builtin functions and constants

## Installation

1. Install the extension from the VSCode marketplace
2. Open any `.qasm` file to activate language support

## Usage

Create or open a `.qasm` file:

```qasm
OPENQASM 3.0;
include "stdgates.inc";

// Initialize qubits
qubit[2] q;
bit[2] c;

// Apply gates
h q[0];
cx q[0], q[1];

// Measure
measure q -> c;
```

## Requirements

- VSCode 1.60.0 or higher

## Contributing

This extension is part of the [qasmtools](https://github.com/orangekame3/qasmtools) project.

## License

Apache License 2.0