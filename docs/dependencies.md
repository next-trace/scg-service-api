# Dependencies

This document lists the dependencies required for the SCG Service Base library.

## Core Dependencies

These dependencies are already included in the go.mod file:

- Go standard library slog (Go 1.25+) - Structured logging
- github.com/spf13/viper v1.20.1 - Configuration management
- github.com/stretchr/testify v1.10.0 - Testing utilities
- go.opentelemetry.io/otel v1.37.0 - OpenTelemetry API
- go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.37.0 - OpenTelemetry stdout exporter
- go.opentelemetry.io/otel/sdk v1.37.0 - OpenTelemetry SDK
- go.opentelemetry.io/otel/trace v1.37.0 - OpenTelemetry tracing API

## gRPC Dependencies

These dependencies need to be added to support gRPC:

- google.golang.org/grpc v1.75.0 - gRPC framework
- google.golang.org/protobuf v1.36.8 - Protocol Buffers support
- github.com/grpc-ecosystem/go-grpc-middleware v1.4.0 - gRPC middleware
- github.com/grpc-ecosystem/grpc-gateway/v2 v2.19.1 - gRPC Gateway for REST APIs

To add these dependencies, run:

```bash
go get google.golang.org/grpc@v1.75.0
go get google.golang.org/protobuf@v1.36.8
go get github.com/grpc-ecosystem/go-grpc-middleware@v1.4.0
go get github.com/grpc-ecosystem/grpc-gateway/v2@v2.19.1
```

Then run:

```bash
go mod tidy
```

## Dependency Injection

For dependency injection, we recommend:

- github.com/uber-go/dig v1.17.1 - Lightweight dependency injection framework

```bash
go get github.com/uber-go/dig@v1.17.1
```

## Metrics

For metrics collection, we recommend:

- github.com/prometheus/client_golang v1.19.0 - Prometheus client library

```bash
go get github.com/prometheus/client_golang@v1.19.0
```

## Validation

For request validation, we recommend:

- github.com/go-playground/validator/v10 v10.19.0 - Validation library

```bash
go get github.com/go-playground/validator/v10@v10.19.0
```

## Circuit Breaking

For circuit breaking, we recommend:

- github.com/sony/gobreaker v0.5.0 - Circuit breaker implementation

```bash
go get github.com/sony/gobreaker@v0.5.0
```

## Rate Limiting

For rate limiting, we recommend:

- golang.org/x/time/rate v0.5.0 - Rate limiting implementation

```bash
go get golang.org/x/time/rate@v0.5.0
```

## Caching

For caching, we recommend:

- github.com/patrickmn/go-cache v2.1.0 - In-memory caching
- github.com/go-redis/redis/v8 v8.11.5 - Redis client

```bash
go get github.com/patrickmn/go-cache@v2.1.0
go get github.com/go-redis/redis/v8@v8.11.5
```