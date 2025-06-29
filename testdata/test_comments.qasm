OPENQASM 3.0;
// This is a single line comment
include "stdgates.qasm";

/* This is a 
   multiline comment */
qubit[2] q;  // Inline comment

h q[0];  // Apply Hadamard gate
cx q[0], q[1];  // CNOT gate