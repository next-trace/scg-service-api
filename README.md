# SCG Service API

[![CI](https://github.com/next-trace/scg-service-api/actions/workflows/ci.yml/badge.svg)](https://github.com/next-trace/scg-service-api/actions/workflows/ci.yml)
[![Coverage](https://codecov.io/gh/next-trace/scg-service-api/branch/main/graph/badge.svg)](https://codecov.io/gh/next-trace/scg-service-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/next-trace/scg-service-api)](https://goreportcard.com/report/github.com/next-trace/scg-service-api)
[![Go Version](https://img.shields.io/github/go-mod/go-version/next-trace/scg-service-api)](https://github.com/next-trace/scg-service-api)

## Overview

SCG Service API is a lightweight base API/service skeleton for building production-grade Go microservices aligned with Domain-Driven Design (DDD). It provides foundational building blocks for HTTP/gRPC services: health endpoints, structured logging, configuration hooks, tracing, metrics, validation, rate limiting, pagination, and more. Use it as a starting point to keep your services consistent and operationally ready.

### Design principles
- Code to interfaces: application layer defines stable ports (logger, metrics, tracing, validation, http, grpc, etc.) and infrastructure implements adapters. Swapping third-party libraries does not affect consumers.
- Idiomatic Go: contexts as first parameter, clear error handling, minimal surface area, and package-level documentation.
- KISS and DRY: small focused interfaces and helpers; avoid duplication; keep defaults sensible.
- SOLID via ports/adapters: single-responsibility packages, dependency inversion (adapters depend on app ports), and open for extension via constructors returning interfaces.

## Install

Requires Go 1.25 or newer.

```bash
go get github.com/next-trace/scg-service-api
```

## Quickstart

Example service layout:

```
my-service/
├─ cmd/
│  └─ my-service/
│     └─ main.go
├─ internal/
│  ├─ app/
│  └─ transport/
├─ configs/
│  └─ app.yaml
├─ go.mod
└─ go.sum
```

Minimal main.go with healthz/readyz and structured logging:

```go
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "time"

    apphealth "github.com/next-trace/scg-service-api/application/health"
    infralog "github.com/next-trace/scg-service-api/infrastructure/logger"
    infrahealth "github.com/next-trace/scg-service-api/infrastructure/health"
    apphttp "github.com/next-trace/scg-service-api/application/http"
)

func main() {
    // Logger (JSON by default). Level can be set via LOG_LEVEL (debug, info, warn, error).
    logLevel := os.Getenv("LOG_LEVEL")
    logger := infralog.NewSlogAdapter(os.Stdout, logLevel)

    // Health registry + default checks
    registry := infrahealth.NewRegistry()
    infrahealth.RegisterCommonChecks(registry)

    // Health config and HTTP handlers
    hc := apphealth.DefaultConfig()
    mux := http.NewServeMux()
    handler := infrahealth.NewHTTPHandler(registry, hc, logger)
    infrahealth.RegisterHTTPHandlers(handler, mux, hc)

    // Basic root handler for demonstration
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        logger.Info(r.Context(), "hello world")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("ok"))
    })

    srv := &http.Server{
        Addr:         ":8080",
        Handler:      mux,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    logger.Info(context.Background(), "server starting on :8080")

    // Graceful shutdown using the helper from application/http
    ctx := context.Background()
    if err := apphttp.Run(ctx, srv, logger); err != nil {
        logger.Error(context.Background(), err, "http server exited with error")
    }
}
```

Available endpoints out of the box:
- GET /health           — JSON payload with liveness and readiness statuses
- GET /health/liveness  — liveness only
- GET /health/readiness — readiness only

## Configuration

Integrate with scg-config to centralize configuration loading (files + env overrides):

```go
package main

import (
    "log"

    scgconfig "github.com/next-trace/scg-config"
)

type AppConfig struct {
    Server struct {
        Port int    `yaml:"port"`
        Host string `yaml:"host"`
    } `yaml:"server"`
}

func loadConfig() AppConfig {
    var cfg AppConfig

    // Load configs/app (e.g., configs/app.yaml). Environment variables can override values.
    if err := scgconfig.Load("configs", "app", &cfg); err != nil {
        log.Fatalf("failed to load config: %v", err)
    }

    // Or with advanced options
    _ = scgconfig.LoadWithOptions("configs", "app", &cfg, scgconfig.ConfigOptions{
        ConfigType:       "yaml",
        EnvPrefix:        "APP",
        AllowEnvOverride: true,
    })

    return cfg
}
```

Tip: set APP_SERVER_PORT=8081 to override server.port when using EnvPrefix: "APP".

## Logging & Errors

Use scg-logger for structured, context-aware logs, and scg-error for rich errors.

```go
package main

import (
    "context"
    "errors"
    "os"

    scglog "github.com/next-trace/scg-logger"
    scgerr "github.com/next-trace/scg-error"
)

func demoLoggingAndErrors() {
    log := scglog.New(os.Stdout, os.Getenv("LOG_LEVEL")) // debug|info|warn|error
    ctx := context.Background()

    // Simple logs
    log.Info(ctx, "service started")

    // Structured logs
    log.InfoKV(ctx, "request processed", map[string]interface{}{"request_id": "abc-123", "user": "alice"})

    // Wrap and log errors with context
    err := errors.New("database connect failed")
    werr := scgerr.Wrap(err, scgerr.WithCode("DB_CONN"), scgerr.WithMsg("cannot connect to primary DB"))
    log.Error(ctx, werr, "operation failed")
}
```

## Testing

Run tests with race detector and coverage:

```bash
go test ./... -race -cover
```

## Documentation

- Package-level documentation is available via doc.go files in each package.
- Generate browsable docs locally:

```bash
./scg docs
```

This uses go doc -all ./... and optionally starts a local godoc server if present.

## CI & Quality

- Continuous Integration: GitHub Actions [CI workflow](https://github.com/next-trace/scg-service-api/actions/workflows/ci.yml)
- Badges:
  - [![CI](https://github.com/next-trace/scg-service-api/actions/workflows/ci.yml/badge.svg)](https://github.com/next-trace/scg-service-api/actions/workflows/ci.yml)
  - [![Coverage](https://codecov.io/gh/next-trace/scg-service-api/branch/main/graph/badge.svg)](https://codecov.io/gh/next-trace/scg-service-api)

The CI runs build, lint, security checks, tests, and uploads coverage.

## Versioning

This project follows [Semantic Versioning](https://semver.org/) (`MAJOR.MINOR.PATCH`).

- **MAJOR**: Breaking API changes
- **MINOR**: New features (backward-compatible)
- **PATCH**: Bug fixes and improvements (backward-compatible)

Consumers should always pin to a specific tag (e.g. `v1.2.3`) to avoid accidental breaking changes.

```bash
go get github.com/next-trace/scg-service-api@v1.0.0
```

Use the latest stable tag published in the repository’s Releases.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
