// Package grpc defines ports (interfaces) for gRPC client and server so services
// can depend on abstractions and swap underlying gRPC libraries or settings without
// breaking consumers.
//
// This package contains:
//   - Server and ServerConfig: to start/stop servers and register services.
//   - Client: minimal client abstraction defined alongside (see client.go).
//
// See infrastructure/grpc for adapters and health/reflection integration.
package grpc
