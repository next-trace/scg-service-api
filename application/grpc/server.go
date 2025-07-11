// Package grpc defines the abstract interface (PORT) for gRPC server functionality.
package grpc

import (
	"context"
	"net"
)

// Server defines the abstract interface for a gRPC server.
type Server interface {
	// Start starts the gRPC server on the given listener.
	Start(ctx context.Context, listener net.Listener) error

	// Stop gracefully stops the gRPC server.
	Stop(ctx context.Context) error

	// RegisterService registers a gRPC service with the server.
	// The implementation should handle type assertions and registrations.
	RegisterService(service interface{}) error
}

// ServerConfig holds configuration for gRPC servers.
type ServerConfig struct {
	// MaxConcurrentStreams is the maximum number of concurrent streams to each client.
	MaxConcurrentStreams uint32

	// MaxRecvMsgSize is the maximum message size the server can receive.
	MaxRecvMsgSize int

	// MaxSendMsgSize is the maximum message size the server can send.
	MaxSendMsgSize int

	// ConnectionTimeout is the timeout for connection establishment.
	ConnectionTimeout int

	// EnableReflection enables server reflection for tools like grpcurl.
	EnableReflection bool

	// EnableHealthCheck enables the gRPC health checking service.
	EnableHealthCheck bool
}

// DefaultServerConfig returns the default configuration for a gRPC server.
func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		MaxConcurrentStreams: 1000,
		MaxRecvMsgSize:       4 * 1024 * 1024, // 4MB
		MaxSendMsgSize:       4 * 1024 * 1024, // 4MB
		ConnectionTimeout:    120,             // 120 seconds
		EnableReflection:     true,
		EnableHealthCheck:    true,
	}
}