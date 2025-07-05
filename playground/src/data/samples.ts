export interface QASMSample {
  id: string;
  name: string;
  description: string;
  difficulty: 'beginner' | 'intermediate' | 'advanced';
  category: 'basic' | 'algorithms' | 'gates' | 'measurement' | 'complex' | 'string' | 'linting';
  code: string;
}

export const qasmSamples: QASMSample[] = [
  {
    id: 'bell-state-messy',
    name: 'Bell State (Unformatted)',
    description: 'Classic Bell state with poor formatting - perfect for testing the formatter',
    difficulty: 'beginner',
    category: 'basic',
    code: `OPENQASM 3.0;include"stdgates.qasm";qubit[2]q;bit[2]c;h q[0];cx q[0],q[1];measure q->c;`
  },
  {
    id: 'grover-messy',
    name: "Grover's Search (Unformatted)",
    description: 'Grover algorithm with inconsistent spacing and formatting issues',
    difficulty: 'intermediate',
    category: 'algorithms',
    code: `OPENQASM 3.0;include"stdgates.qasm";qubit[3]q;bit[3]c;h q[0];h q[1];h q[2];cz q[0],q[2];cz q[1],q[2];h q[0];h q[1];h q[2];x q[0];x q[1];x q[2];ccx q[0],q[1],q[2];x q[0];x q[1];x q[2];h q[0];h q[1];h q[2];measure q->c;`
  },
  {
    id: 'deutsch-jozsa-messy',
    name: 'Deutsch-Jozsa (Mixed Formatting)',
    description: 'Algorithm with mixed indentation and spacing issues',
    difficulty: 'intermediate',
    category: 'algorithms',
    code: `OPENQASM 3.0;
include"stdgates.qasm";
qubit[4]q;bit[3]c;
x q[3];h q[0];h q[1];h q[2];h q[3];
cx q[0],q[3];cx q[1],q[3];cx q[2],q[3];
h q[0];h q[1];h q[2];
measure q[0]->c[0];measure q[1]->c[1];measure q[2]->c[2];`
  },
  {
    id: 'superposition-messy',
    name: 'Superposition (No Spacing)',
    description: 'Simple superposition circuit with no spaces or proper formatting',
    difficulty: 'beginner',
    category: 'basic',
    code: `OPENQASM 3.0;include"stdgates.qasm";qubit[4]q;bit[4]c;h q[0];h q[1];h q[2];h q[3];measure q->c;`
  },
  {
    id: 'teleportation-messy',
    name: 'Quantum Teleportation (Poor Format)',
    description: 'Teleportation circuit with inconsistent formatting perfect for cleanup',
    difficulty: 'advanced',
    category: 'complex',
    code: `OPENQASM 3.0;include"stdgates.qasm";qubit[3]q;bit[2]c;h q[0];h q[1];cx q[1],q[2];cx q[0],q[1];h q[0];measure q[0]->c[0];measure q[1]->c[1];if(c[1]==1)x q[2];if(c[0]==1)z q[2];`
  },
  {
    id: 'chain-cnot-escaped',
    name: 'Chain CNOT (Escaped String)',
    description: 'JSON-style escaped string with \\n and \\" - enable Unescape to format properly',
    difficulty: 'beginner',
    category: 'string',
    code: `"OPENQASM 3.0;\\ninclude \\"stdgates.inc\\";\\nqubit[5] q;\\nbit[5] c;\\n\\nh q[0];\\ncx q[0], q[1];\\ncx q[1], q[2];\\ncx q[2], q[3];\\ncx q[3], q[4];\\nc = measure q;"`
  },
  {
    id: 'bell-state-escaped',
    name: 'Bell State (Escaped String)',
    description: 'Classic Bell state as an escaped JSON string - try with Unescape enabled',
    difficulty: 'beginner',
    category: 'string',
    code: `"OPENQASM 3.0;\\ninclude \\"stdgates.inc\\";\\nqubit[2] q;\\nbit[2] c;\\nh q[0];\\ncx q[0], q[1];\\nmeasure q -> c;"`
  },
  {
    id: 'complex-circuit-escaped',
    name: 'Complex Circuit (Escaped String)',
    description: 'Multi-gate quantum circuit with conditional logic as escaped string',
    difficulty: 'intermediate',
    category: 'string',
    code: `"OPENQASM 3.0;\\ninclude \\"stdgates.inc\\";\\nqubit[3] q;\\nbit[3] c;\\n\\n// Initialize superposition\\nh q[0];\\nh q[1];\\nh q[2];\\n\\n// Apply gates\\ncx q[0], q[1];\\ncz q[1], q[2];\\nrx(pi/4) q[0];\\n\\n// Measurement and conditional\\nmeasure q[0] -> c[0];\\nif (c[0] == 1) x q[1];\\nmeasure q -> c;"`
  },
  // Linting Examples - designed to trigger specific lint rules
  {
    id: 'lint-unused-qubit',
    name: 'QAS001: Unused Qubit',
    description: 'Demonstrates unused qubit detection - qubits declared but never used',
    difficulty: 'beginner',
    category: 'linting',
    code: `OPENQASM 3.0;
include "stdgates.qasm";

// These qubits are declared but never used
qubit[2] q;
qubit unused_qubit;
qubit[3] another_unused;

// Only q[0] is used - others will trigger QAS001 warnings
h q[0];
bit c;
measure q[0] -> c;`
  },
  {
    id: 'lint-undefined-identifier',
    name: 'QAS002: Undefined Identifier',
    description: 'Shows undefined identifier errors - using variables that were never declared',
    difficulty: 'beginner',
    category: 'linting',
    code: `OPENQASM 3.0;
include "stdgates.qasm";

qubit[2] q;
bit[2] c;

// These will trigger QAS002 errors - undefined identifiers
h undefined_qubit;
cx q[0], unknown_qubit;
measure nonexistent_q -> c[0];

// Valid operations
h q[0];
cx q[0], q[1];
measure q -> c;`
  },
  {
    id: 'lint-constant-measured-bit',
    name: 'QAS003: Constant Measured Bit',
    description: 'Warning for measuring qubits with no gates applied - result will always be |0⟩',
    difficulty: 'beginner',
    category: 'linting',
    code: `OPENQASM 3.0;
include "stdgates.qasm";

qubit[3] q;
bit[3] c;

// Apply gate to only one qubit
h q[0];

// These measurements will trigger QAS003 warnings
// q[1] and q[2] have no gates applied
measure q[0] -> c[0];  // OK - has gate applied
measure q[1] -> c[1];  // WARNING - no gates, always |0⟩
measure q[2] -> c[2];  // WARNING - no gates, always |0⟩`
  },
  {
    id: 'lint-out-of-bounds',
    name: 'QAS004: Out of Bounds Index',
    description: 'Error for accessing array indices that exceed the declared size',
    difficulty: 'beginner',
    category: 'linting',
    code: `OPENQASM 3.0;
include "stdgates.qasm";

qubit[2] q;  // Valid indices: 0, 1
bit[2] c;    // Valid indices: 0, 1

// Valid operations
h q[0];
h q[1];

// These will trigger QAS004 errors - out of bounds
h q[2];      // ERROR - q[2] doesn't exist
cx q[0], q[3];  // ERROR - q[3] doesn't exist
measure q[5] -> c[0];  // ERROR - q[5] doesn't exist
measure q[0] -> c[5];  // ERROR - c[5] doesn't exist`
  },
  {
    id: 'lint-naming-convention',
    name: 'QAS005: Naming Convention',
    description: 'Warning for identifiers that violate OpenQASM naming conventions',
    difficulty: 'beginner',
    category: 'linting',
    code: `OPENQASM 3.0;
include "stdgates.qasm";

// These names violate conventions (start with uppercase)
qubit[2] MyQubit;     // WARNING - should start with lowercase
bit[2] BigBit;        // WARNING - should start with lowercase
qubit SingleQubit;    // WARNING - should start with lowercase

// Valid names (start with lowercase)
qubit[2] myQubit;     // OK
bit[2] result_bits;   // OK
qubit helper_qubit;   // OK

// Operations on the invalid names will still work
h MyQubit[0];
measure MyQubit -> BigBit;`
  },
  {
    id: 'lint-multiple-errors',
    name: 'Multiple Lint Issues',
    description: 'Complex example with multiple different lint violations for testing',
    difficulty: 'intermediate',
    category: 'linting',
    code: `OPENQASM 3.0;
include "stdgates.qasm";

// QAS005: Naming convention violations
qubit[2] MyQubit;     // WARNING - uppercase start
bit[2] BigBit;        // WARNING - uppercase start
qubit UnusedQubit;    // WARNING - uppercase start

// QAS001: Unused qubit
qubit[3] never_used;  // WARNING - declared but never used

// QAS004: Out of bounds access
h MyQubit[0];         // OK
h MyQubit[2];         // ERROR - out of bounds

// QAS002: Undefined identifier
cx MyQubit[0], undefined_qubit;  // ERROR - undefined

// QAS003: Constant measured bit
qubit untouched;
measure untouched -> BigBit[0];  // WARNING - no gates applied

// Some valid operations
h MyQubit[1];
measure MyQubit -> BigBit;`
  }
];

export const sampleCategories = [
  { id: 'all', name: 'All Samples', icon: '📋' },
  { id: 'basic', name: 'Basic Circuits', icon: '🔰' },
  { id: 'algorithms', name: 'Algorithms', icon: '🧮' },
  { id: 'gates', name: 'Gate Demos', icon: '⚡' },
  { id: 'measurement', name: 'Measurement', icon: '📊' },
  { id: 'complex', name: 'Advanced', icon: '🚀' },
  { id: 'string', name: 'Escaped Strings', icon: '📝' },
  { id: 'linting', name: 'Lint Examples', icon: '🔍' }
];

export const difficultyColors = {
  beginner: 'badge-success',
  intermediate: 'badge-warning', 
  advanced: 'badge-error'
};

export const difficultyLabels = {
  beginner: 'Beginner',
  intermediate: 'Intermediate',
  advanced: 'Advanced'
};