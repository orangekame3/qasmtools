OPENQASM 3.0; // OpenQASM version declaration
include "stdgates.qasm"; // Include standard gates

// Qubit and classical bit declarations
qubit[3] q;
bit[3] c;

// Initialization and setup
reset q[0]; // Reset first qubit
barrier q; // Synchronization barrier

// Single-qubit gates
h q[0]; // Hadamard gate
rz(pi/4) q[1]; // Rotation around Z-axis

// Multi-qubit gates
cphase(pi/2) q[0], q[1]; // Controlled phase gate
cx q[1], q[2]; // CNOT gate

/* 
 * Measurement section
 * This measures all qubits
 */
measure q -> c; // Final measurement