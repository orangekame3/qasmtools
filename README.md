# qasmtools

`qasmtools` is a command-line interface (CLI) tool written in Go for parsing and formatting OpenQASM 3.0 files. It provides functionalities to process QASM files, making them more readable and consistent.

## Features

* **QASM 3.0 Parsing**: Parses OpenQASM 3.0 files into an Abstract Syntax Tree (AST).
* **QASM 3.0 Formatting**: Formats QASM 3.0 files to adhere to a consistent style.

## Installation

To build `qasmtools` from source, ensure you have Go installed (version 1.16 or higher recommended).

1. Clone the repository:

    ```bash
    git clone https://github.com/orangekame3/qasmtools.git
    cd qasmtools
    ```

2. Build the executable:

    ```bash
    go build -o qasm ./cmd/qasm
    ```

    This will create an executable named `qasm` in the current directory.

## Usage

The `qasm` executable can be used to format QASM files.

### Formatting a QASM file

To format a QASM file and print the output to standard output:

```bash
./qasm format <input_file.qasm>
```

Example:

```bash
./qasm format test.qasm
```

### Overwriting a QASM file (in-place formatting)

To format a QASM file and overwrite the original file:

```bash
./qasm format -w <input_file.qasm>
```

Example:

```bash

./qasm format -w test.qasm
```

## Project Structure

* `cmd/qasm/`: Contains the main entry point for the CLI tool.
* `parser/`: Handles the parsing of QASM 3.0 files and AST generation.
* `formatter/`: Implements the QASM 3.0 formatting logic.
* `test.qasm`: A sample QASM file for testing and demonstration.
