OPENQASM 3.0;
include "stdgates.qasm";

qubit[2] q;
bit[2] c;

gate my_gate q {
    h q;
}

h q[0];
cx q[0], q[1];
my_gate q[0];
measure q -> c;