OPENQASM 3.0;
include "stdgates.inc";

// This is a comment
gate h q { }

qubit q;
h q;
measure q -> c;
