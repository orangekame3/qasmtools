OPENQASM 3.0;
include "stdgates.qasm";

qubit[2] q;
bit[2] c;

h q[0];  // Only q[0] is affected by a gate

measure q[0] -> c[0];  // This is fine
measure q[1] -> c[1];  // This will always measure |0⟩