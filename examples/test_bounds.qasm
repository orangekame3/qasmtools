OPENQASM 3.0;
qubit[2] q;
h q[0];
cx q[0], q[2];
measure q[0] -> c[0];