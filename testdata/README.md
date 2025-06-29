# Test Data

This directory contains test files used for automated testing of the qasmfmt formatter.

## Files

### test_simple.qasm
Basic OpenQASM 3.0 structure for simple formatting tests.

### test_malformed.qasm
Contains intentionally malformed QASM code to test preprocessing and error handling.

### test_unformatted.qasm
Poorly formatted QASM code used to test the formatter's correction capabilities.

### test_gates.qasm
Various gate operations and quantum circuits for comprehensive testing.

### test_comments.qasm
Test file for comment preservation (currently limited by parser).

### debug.qasm & debug_bit.qasm
Debug files used during development for specific formatting scenarios.

## Usage

These files are used by the test suite:

```bash
# Run all tests
go test ./cmd/...

# Run tests with coverage
go test -cover ./cmd/...

# Test specific formatting scenarios
qasmfmt format testdata/test_malformed.qasm
```

## Note

These files are for testing purposes only and may contain intentionally incorrect or malformed QASM code.