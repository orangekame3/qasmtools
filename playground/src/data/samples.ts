export interface QASMSample {
  id: string;
  name: string;
  description: string;
  difficulty: 'beginner' | 'intermediate' | 'advanced';
  category: 'basic' | 'algorithms' | 'gates' | 'measurement' | 'complex';
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
  }
];

export const sampleCategories = [
  { id: 'all', name: 'All Samples', icon: '📋' },
  { id: 'basic', name: 'Basic Circuits', icon: '🔰' },
  { id: 'algorithms', name: 'Algorithms', icon: '🧮' },
  { id: 'gates', name: 'Gate Demos', icon: '⚡' },
  { id: 'measurement', name: 'Measurement', icon: '📊' },
  { id: 'complex', name: 'Advanced', icon: '🚀' }
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