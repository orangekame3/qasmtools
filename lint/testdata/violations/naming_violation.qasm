OPENQASM 3.0;
include "stdgates.qasm";

qubit BadName;  // Violates naming convention (should start with lowercase)
qubit[2] q;
bit[2] c;

gate MyGate q {  // Also violates naming convention
    h q;
}

h q[0];
measure q -> c;