# SCG Service Base Tools

This directory contains Go tools to replace the functionality previously provided by the Makefile.

## Building the Tool

To build the tool, run:

```bash
go run cmd/tools/build.go
```

This will create a binary called `scg-tools` in the project root directory.

## Using the Tool

The tool provides the following commands:

```bash
# Show help message
./scg-tools -help

# Generate Go code from protobuf definitions
./scg-tools -proto

# Clean generated files
./scg-tools -clean

# Install required tools
./scg-tools -install-tools
```

## Why Go Instead of Bash?

Using Go instead of Bash scripts (via Makefile) provides several advantages:

1. **Cross-platform compatibility**: Go code works the same on Windows, macOS, and Linux, while Bash scripts may have compatibility issues.

2. **Type safety and error handling**: Go provides better error handling and type safety compared to Bash scripts.

3. **Maintainability**: Go code is generally easier to maintain, test, and extend than Bash scripts.

4. **Consistency**: Using Go for both the application and build tools provides a consistent development experience.

5. **IDE support**: Go has excellent IDE support for code completion, refactoring, and debugging.

## Implementation Details

The tool is implemented in two main files:

- `main.go`: Contains the main functionality for generating code, cleaning files, and installing tools.
- `build.go`: A simple script to build the tool and place it in the project root.

The tool uses standard Go libraries to implement the functionality:

- `os/exec` for executing external commands
- `filepath` for file path manipulation
- `flag` for command-line argument parsing