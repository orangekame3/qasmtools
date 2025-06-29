# Development Guide

This document describes the development workflow and available tasks for the qasmtools project.

## Prerequisites

- Go 1.16 or higher
- [Task](https://taskfile.dev/) - A task runner / build tool
- [ANTLR4](https://www.antlr.org/) - Parser generator
- [golangci-lint](https://golangci-lint.run/) - Go linter
- [gofumpt](https://github.com/mvdan/gofumpt) - A stricter gofmt

## Development Workflow

1. Clone the repository:

   ```bash
   git clone https://github.com/orangekame3/qasmtools.git
   cd qasmtools
   ```

2. Set up the development environment:

   ```bash
   task setup
   ```

   This will:
   - Download dependencies
   - Install required tools (golangci-lint, gofumpt)
   - Create example files
   - Set up the project structure

3. Generate parser code:

   ```bash
   task gen
   ```

   This will:
   - Download the latest OpenQASM grammar files
   - Generate Go parser code using ANTLR4
   - Fix import paths for compatibility

4. Build the project:

   ```bash
   task build
   ```

## Available Tasks

### Build Tasks

- `task build` - Build the application
- `task build:all` - Build for all platforms (Linux, macOS, Windows)
- `task install` - Install the binary to $GOPATH/bin

### Development Tasks

- `task dev` - Run in development mode with live reload
- `task gen` - Generate parser code from ANTLR grammar files
- `task fmt` - Format code
- `task lint` - Run linters

### Testing Tasks

- `task test` - Run tests
- `task test:coverage` - Run tests with coverage
- `task test:race` - Run tests with race detection
- `task bench` - Run benchmarks

### Code Quality Tasks

- `task lint` - Run linters
- `task lint:golangci` - Run golangci-lint
- `task fmt` - Format code
- `task fmt:install` - Install gofumpt

### Dependency Management

- `task deps` - Download dependencies
- `task deps:tidy` - Tidy dependencies
- `task deps:verify` - Verify dependencies
- `task deps:update` - Update dependencies

### Example Tasks

- `task example:create` - Create example QASM files

### Security Tasks

- `task security` - Run security checks (gosec, govulncheck)

### CI/CD Tasks

- `task ci` - Run CI pipeline locally
- `task pre-commit` - Run pre-commit checks

### Cleanup Tasks

- `task clean` - Clean build artifacts
- `task clean:all` - Clean all generated files

### Version Info

- `task version` - Show version information

## Project Structure

```bash
qasmtools/
├── cmd/qasm/          # Main CLI entry point
├── formatter/         # Formatter implementation
├── parser/           # Parser implementation
│   ├── grammar/     # ANTLR grammar files
│   └── gen/        # Generated parser code
├── examples/         # Example QASM files
└── testdata/        # Test files
```

## Code Style

The project uses `.editorconfig` to maintain consistent coding styles. Make sure your editor supports EditorConfig or install the appropriate plugin.

Key style guidelines:

- Go files: tabs for indentation
- QASM files: 2 spaces for indentation
- YAML files: 2 spaces for indentation
- Grammar files: 4 spaces for indentation

## Pre-commit Checks

Before committing changes, run:

```bash
task pre-commit
```

This will:

1. Format the code
2. Run linters
3. Run tests

## Continuous Integration

The project uses GitHub Actions for CI/CD. The following workflows are available:

### CI Workflow (`ci.yml`)

Triggered on push to main and pull requests:

1. Test job:
   - Runs tests with coverage
   - Uploads coverage to Codecov
2. Lint job:
   - Runs golangci-lint
3. Build job:
   - Builds for all platforms
   - Uploads artifacts
4. Security job:
   - Runs security checks (gosec, govulncheck)

Run the CI pipeline locally:

```bash
task ci
```

### Release Workflow (`release.yml`)

Triggered on pushing tags (v*):

1. Runs tests
2. Builds binaries using GoReleaser
3. Creates GitHub release
4. Publishes to:
   - GitHub Releases
   - Homebrew tap
   - Container registry (ghcr.io)

## Release Process

1. Update version information:

   ```bash
   task version
   ```

2. Run the full test suite:

   ```bash
   task test:coverage
   ```

3. Create and push a new tag:

   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

This will trigger the release workflow, which:

- Creates a GitHub release
- Builds and uploads binaries
- Updates the Homebrew formula
- Pushes Docker images

## Troubleshooting

### Parser Generation Issues

If you encounter issues with parser generation:

1. Clean generated files:

   ```bash
   task clean
   ```

2. Regenerate parser code:

   ```bash
   task gen
   ```

### Build Issues

If you encounter build issues:

1. Verify dependencies:

   ```bash
   task deps:verify
   ```

2. Clean and rebuild:

   ```bash
   task clean:all
   task build
   ```

### CI/CD Issues

If CI/CD workflows fail:

1. Check the workflow logs in GitHub Actions
2. Run the CI pipeline locally:

   ```bash
   task ci
   ```

3. Verify GoReleaser configuration:

   ```bash
   goreleaser check
    ```
