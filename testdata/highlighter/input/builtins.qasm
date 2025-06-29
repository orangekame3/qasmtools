OPENQASM 3.0;

// Builtin gates
U(pi/2, 0, pi) q[0];
CX q[0], q[1];

// Builtin quantum operations
reset q;
barrier q;
measure q -> c;

// Builtin classical functions
let angle = pi/2;
let x = sin(angle);
let y = cos(angle);
let z = sqrt(x*x + y*y);

// Builtin constants
const real_pi = pi;
const real_tau = tau;
const real_euler = euler;

// Hardware qubits
$0 = U(0, 0, 0) $1;

// Access control
const int x = 1;
mutable float y = 2.0;
