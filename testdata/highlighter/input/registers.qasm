OPENQASM 3.0;

// Register declarations
qubit[5] q;
bit[5] c;
qreg r[10];
creg m[10];

// Quantum operations with registers
h q[0];
cx q[0], q[1];
reset q;

// Measurements with registers
measure q[0] -> c[0];
measure q -> c;

// Array indexing
let x = q[2];
let y = c[3];
