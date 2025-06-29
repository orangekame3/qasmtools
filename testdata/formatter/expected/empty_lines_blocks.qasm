OPENQASM 3.0;
include "stdgates.qasm";

qubit[2] q;
bit c;

gate custom a, b {
  cx a, b;
}

h q[0];
measure q[0] -> c;
