OPENQASM 3.0;

// For loop
for i in [0:3] {
    h q[i];
}

// While loop
while (i > 0) {
    i = i - 1;
}

// If-else statement
if (x == 1) {
    h q[0];
} else {
    x q[0];
}

// Break and continue
for i in [0:10] {
    if (i == 5) {
        break;
    }
    if (i < 3) {
        continue;
    }
    h q[i];
}
