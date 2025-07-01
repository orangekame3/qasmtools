# QASM Lint Rules Documentation

This directory contains documentation for all available QASM linting rules.

Each rule includes detailed descriptions, examples, and links to the relevant OpenQASM 3.0 specification sections.

## Error Rules

These rules catch critical issues that prevent valid OpenQASM execution:

- **[QAS002](QAS002.md)** - Identifier is used without being declared.
- **[QAS004](QAS004.md)** - Index of qubit, bit, or array exceeds the defined size.
- **[QAS006](QAS006.md)** - Quantum register sizes do not match during gate call, preventing broadcasting.
- **[QAS007](QAS007.md)** - Index access on registers passed as gate arguments is not allowed in gate body.
- **[QAS008](QAS008.md)** - Qubits can only be declared in global scope.
- **[QAS009](QAS009.md)** - break and continue can only be used inside loops.
- **[QAS010](QAS010.md)** - Using non-unitary instructions (measurement, reset, etc.) within gate definition.
- **[QAS011](QAS011.md)** - Using reserved prefix (__) in identifier.

## Warning Rules

These rules identify potential issues and style violations:

- **[QAS001](QAS001.md)** - Detects qubits that are declared but never used in gates or measurements.
- **[QAS003](QAS003.md)** - Measuring a qubit that has no gates applied. The result will always be |0⟩.
- **[QAS005](QAS005.md)** - Identifier name violates OpenQASM naming conventions.
- **[QAS012](QAS012.md)** - Identifiers should be named in snake_case (lower_snake_case).

## All Rules Summary

| Rule ID | Name | Severity | Tags | Fixable | Specification |
|---------|------|----------|------|---------|---------------|
| [QAS001](QAS001.md) | unused-qubit | warning | qasm3, readability, unused-variables | false | [Link](https://openqasm.com/versions/3.0/language/types.html#qubits) |
| [QAS002](QAS002.md) | undefined-identifier | error | qasm3, scope, semantic | false | [Link](https://openqasm.com/versions/3.0/language/scope.html) |
| [QAS003](QAS003.md) | constant-measured-bit | warning | qasm3, logic, measurement, constant-results | false | [Link](https://openqasm.com/versions/3.0/language/quantum.html#measurement) |
| [QAS004](QAS004.md) | out-of-bounds-index | error | qasm3, bounds-checking, array-access | false | [Link](https://openqasm.com/versions/3.0/language/types.html#index-sets-and-slicing) |
| [QAS005](QAS005.md) | naming-convention-violation | warning | qasm3, style, naming, readability | false | [Link](https://openqasm.com/versions/3.0/language/lexical.html#identifiers) |
| [QAS006](QAS006.md) | gate-register-size-mismatch | error | qasm3, semantic, gate | false | [Link](https://openqasm.com/versions/3.0/language/gates.html#broadcasting) |
| [QAS007](QAS007.md) | gate-parameter-indexing | error | qasm3, gate, syntax | false | [Link](https://openqasm.com/versions/3.0/language/gates.html#hierarchical-gates-definitions) |
| [QAS008](QAS008.md) | qubit-declared-in-local-scope | error | qasm3, scope, qubit | false | [Link](https://openqasm.com/versions/3.0/language/types.html#qubits) |
| [QAS009](QAS009.md) | illegal-break-continue | error | qasm3, control-flow, syntax | false | [Link](https://openqasm.com/versions/3.0/language/classical.html#breaking-and-continuing-loops) |
| [QAS010](QAS010.md) | invalid-instruction-in-gate | error | qasm3, gate, syntax | false | [Link](https://openqasm.com/versions/3.0/language/gates.html#hierarchical-gates-definitions) |
| [QAS011](QAS011.md) | reserved-prefix-usage | error | qasm3, naming, style | false | [Link](https://openqasm.com/versions/3.0/language/lexical.html#identifiers) |
| [QAS012](QAS012.md) | snake-case-required | warning | qasm3, style, naming | false | [Link](https://openqasm.com/versions/3.0/language/lexical.html#identifiers) |

## Usage

To disable specific rules, use the `--disable` flag:

```bash
qasm lint --disable=QAS005,QAS012 input.qasm
```

To enable only specific rules, use the `--enable-only` flag:

```bash
qasm lint --enable-only=QAS001,QAS002 input.qasm
```

