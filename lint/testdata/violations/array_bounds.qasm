OPENQASM 3.0;
include "stdgates.qasm";

qubit[2] q;    // q[0], q[1] are valid
bit[3] c;      // c[0], c[1], c[2] are valid

h q[0];        // OK
h q[1];        // OK
h q[2];        // ERROR: out of bounds (should be 0-1)

measure q[0] -> c[0];   // OK
measure q[1] -> c[2];   // OK
measure q[0] -> c[3];   // ERROR: out of bounds (should be 0-2)