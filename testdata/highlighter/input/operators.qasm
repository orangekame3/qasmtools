OPENQASM 3.0;

// Arithmetic operators
let a = 1 + 2;
let b = 3 - 4;
let c = 5 * 6;
let d = 7 / 8;

// Comparison operators
if (a == b) {
    let x = 1;
}

// Quantum operators
cx q[0], q[1];
h q[0];
measure q -> c;
