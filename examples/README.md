# Examples

This directory contains example OpenQASM 3.0 files that demonstrate the formatter's capabilities.

## Files

### bell_state.qasm
A simple Bell state preparation circuit that demonstrates:
- Quantum register declarations
- Gate operations (Hadamard and CNOT)
- Measurement operations

### grover.qasm
Grover's search algorithm implementation that shows:
- More complex quantum circuits
- Multi-qubit operations
- Advanced gate sequences

## Usage

You can use these examples to test the formatter:

```bash
# Format an example file
qasmfmt format examples/bell_state.qasm

# Check if examples are properly formatted
qasmfmt check examples/*.qasm

# Format in-place with verbose output
qasmfmt format -w -v examples/bell_state.qasm
```

## Adding New Examples

When adding new examples:
1. Use descriptive filenames
2. Include comments explaining the quantum algorithm
3. Follow OpenQASM 3.0 syntax
4. Test with the formatter to ensure proper formatting