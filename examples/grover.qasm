OPENQASM 3.0;
include "stdgates.qasm";

qubit[3] q;
bit[3] c;

// Initialize superposition
h q[0];
h q[1];
h q[2];

// Oracle
cz q[0], q[2];
cz q[1], q[2];

// Diffuser
h q[0];
h q[1];
h q[2];
x q[0];
x q[1];
x q[2];
ccx q[0], q[1], q[2];
x q[0];
x q[1];
x q[2];
h q[0];
h q[1];
h q[2];

measure q -> c;
