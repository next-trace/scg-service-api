# SCG Service Base

[![Go CI](https://github.com/hbttundar/scg-service-base/actions/workflows/ci.yml/badge.svg)](https://github.com/hbttundar/scg-service-base/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/hbttundar/scg-service-base/branch/main/graph/badge.svg)](https://codecov.io/gh/hbttundar/scg-service-base)
[![Go Report Card](https://goreportcard.com/badge/github.com/hbttundar/scg-service-base)](https://goreportcard.com/report/github.com/hbttundar/scg-service-base)
[![Go Version](https://img.shields.io/github/go-mod/go-version/hbttundar/scg-service-base)](https://github.com/hbttundar/scg-service-base)

A foundational library for Go microservices providing standardized configuration and logging capabilities.

## Features

- **Configuration Management**: 
  - Flexible configuration loading from YAML, JSON, TOML, and other formats
  - Environment variable overrides with optional prefixing
  - Runtime configuration reloading
  - Validation for input parameters
  - Support for environment-only configuration

- **Structured Logging**: 
  - JSON-formatted logging with zerolog for better observability in cloud environments
  - Context-aware logging for request tracing
  - Helper methods for common log levels
  - Service name tagging for multi-service environments
  - Pretty-print option for development environments
  - Structured logging with key-value pairs for better data organization

- **Distributed Tracing**:
  - OpenTelemetry integration for distributed tracing
  - Configurable trace exporters (stdout, Jaeger, OTLP, etc.)
  - Middleware for HTTP request tracing
  - Context propagation for end-to-end tracing
  - Sampling rate configuration for production environments

## Installation

```bash
go get github.com/SupplyChainGuard/scg-service-base
```

## Usage

### Configuration

The `config` package provides a simple way to load configuration from files and environment variables:

```go
package main

import (
    "fmt"
    "log"

    "github.com/SupplyChainGuard/scg-service-base/config"
)

// AppConfig represents your application's configuration structure
type AppConfig struct {
    Server struct {
        Port    int    `yaml:"port"`
        Host    string `yaml:"host"`
        Timeout int    `yaml:"timeout"`
    } `yaml:"server"`
    Database struct {
        Host     string `yaml:"host"`
        Port     int    `yaml:"port"`
        Username string `yaml:"username"`
        Password string `yaml:"password"`
        Name     string `yaml:"name"`
    } `yaml:"database"`
}

func main() {
    // Create a config instance
    cfg := &AppConfig{}

    // Load configuration from the "config" directory and "app-config" file
    err := config.Load("config", "app-config", cfg)
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    fmt.Printf("Server will start on %s:%d\n", cfg.Server.Host, cfg.Server.Port)

    // For more advanced usage with custom options:
    options := config.ConfigOptions{
        ConfigType:        "json",  // Use JSON instead of YAML
        EnvPrefix:         "APP",   // Environment variables will be prefixed with APP_
        AllowEnvOverride:  true,    // Allow environment variables to override file settings
        RequireConfigFile: true,    // Require config file to exist (error if not found)
    }

    err = config.LoadWithOptions("config", "app-config", cfg, options)
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    // Check if a config file exists before attempting to load it
    if config.FileExists("config", "app-config", "yaml") {
        fmt.Println("Configuration file exists")
    }

    // Reload configuration at runtime (e.g., after receiving SIGHUP)
    err = config.Reload("config", "app-config", cfg, options)
    if err != nil {
        log.Fatalf("Failed to reload configuration: %v", err)
    }
}
```

Environment variables can override configuration values. For example, to override `database.host` in the configuration, you can set the environment variable `DATABASE_HOST`. If using a prefix like `APP`, the environment variable would be `APP_DATABASE_HOST`.

### Logging

The `logger` package provides structured JSON logging using zerolog:

```go
package main

import (
    "context"
    "github.com/SupplyChainGuard/scg-service-base/logger"
)

func main() {
    // Initialize the logger with default settings
    logger.Init()

    // Log messages at different levels using the fluent API
    logger.GetLogger("api").Info().Str("user", "john").Msg("User logged in")
    logger.GetLogger("db").Error().Err(err).Msg("Database connection failed")

    // Or use the simplified helper methods
    logger.Info("api", "User logged in")
    logger.Error("db", "Database connection failed", err)

    // For more advanced usage with custom options:
    options := logger.LoggerOptions{
        LogLevel:          "debug",    // Set minimum log level
        PrettyPrint:       true,       // Use human-readable format instead of JSON
        IncludeCallerInfo: true,       // Include file and line number in logs
        ServiceName:       "auth-api", // Add service name to all log entries
    }

    logger.InitWithOptions(options)

    // Context-aware logging for request tracing
    ctx := context.Background()
    reqLogger := logger.GetLogger("http").With().Str("request_id", "12345").Logger()

    // Add logger to context
    ctx = logger.ContextWithLogger(ctx, reqLogger)

    // Later in the request handling chain, retrieve and use the logger
    handlerLogger := logger.GetLoggerFromContext(ctx)
    handlerLogger.Info().Msg("Processing request")
}
```

The log level can be set via the `LOG_LEVEL` environment variable (e.g., "debug", "info", "warn", "error"), defaulting to "info".

## Testing

Both packages include comprehensive test suites. Run the tests with:

```bash
go test ./...
```

## Development Tools

This project includes Go-based development tools that replace traditional Makefile functionality:

### Building the Tools

To build the development tools, run:

```bash
go run cmd/tools/build.go
```

This will create a binary called `scg-tools` in the project root directory.

### Using the Tools

The tools provide the following commands:

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

For more details, see the [tools documentation](cmd/tools/README.md).

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Development Workflow

This project uses GitHub Actions for continuous integration and delivery:

- **CI Workflow**: Runs on every push and pull request to the `main` branch
  - Linting with golangci-lint
  - Verifying Go modules are tidy
  - Testing with multiple Go versions (1.20, 1.21) across multiple operating systems (Linux, macOS, Windows)
  - Building the package on multiple operating systems
  - Security scanning with Gosec
  - Dependency vulnerability checking with govulncheck
  - Code coverage reporting to Codecov
  - Documentation generation verification

- **Release Workflow**: Triggered when a tag with the prefix `v` is pushed
  - Runs tests and builds the package
  - Creates a GitHub release with automatically generated release notes

To create a new release:

```bash
git tag v1.0.0
git push origin v1.0.0
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
