# QASM Lint Rules Documentation

This directory contains documentation for all available QASM linting rules.

## Error Rules

- **[QAS002](QAS002.md)** - Detects when the number of classical bits is insufficient for measurements.
- **[QAS004](QAS004.md)** - Detects when qubit or classical bit array indices exceed declared bounds.

## Warning Rules

- **[QAS001](QAS001.md)** - Detects qubits that are declared but never used in gates or measurements.
- **[QAS003](QAS003.md)** - Detects measurements of qubits that have not been affected by any gates, resulting in constant measurement outcomes.
- **[QAS005](QAS005.md)** - Detects violations of OpenQASM naming conventions for variables and circuits.

## All Rules

| Rule ID | Name | Severity | Tags | Fixable |
|---------|------|----------|------|---------|
| [QAS001](QAS001.md) | unused-qubit | warning | qasm3, readability, unused-variables | false |
| [QAS002](QAS002.md) | insufficient-classical-bits | error | qasm3, correctness, measurement | false |
| [QAS003](QAS003.md) | constant-measured-bit | warning | qasm3, logic, measurement, constant-results | false |
| [QAS004](QAS004.md) | exceeding-qubit-limits | error | qasm3, bounds-checking, array-access | false |
| [QAS005](QAS005.md) | naming-convention-violation | warning | qasm3, style, naming, readability | false |
