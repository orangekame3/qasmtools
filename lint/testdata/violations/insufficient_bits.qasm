OPENQASM 3.0;
include "stdgates.qasm";

qubit[3] q;
bit c;  // Only 1 classical bit for 3 measurements

h q[0];
h q[1];
h q[2];

measure q[0] -> c;
measure q[1] -> c;  // This will cause insufficient bits violation
measure q[2] -> c;