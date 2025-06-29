OPENQASM 3.0;
include "stdgates.inc";

// Initialize qubits and classical bits
qubit[2] q;
bit[2] c;

// Apply Hadamard gate to first qubit
h q[0];

// Apply CNOT gate
cx q[0], q[1];

// Measure both qubits
measure q[0] -> c[0];
measure q[1] -> c[1];