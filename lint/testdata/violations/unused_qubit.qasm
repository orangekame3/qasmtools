OPENQASM 3.0;
include "stdgates.qasm";

qubit unused_q;  // This qubit is never used
qubit[2] q;
bit[2] c;

h q[0];
cx q[0], q[1];
measure q -> c;