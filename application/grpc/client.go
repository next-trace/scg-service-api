// Package grpc defines the abstract interface (PORT) for gRPC client functionality.
package grpc

import (
	"context"
	"time"
)

// Client defines the abstract interface for a gRPC client.
type Client interface {
	// Connect establishes a connection to the gRPC server.
	Connect(ctx context.Context) error

	// Close closes the connection to the gRPC server.
	Close() error

	// GetConnection returns the underlying connection for use by service clients.
	// The actual type will depend on the implementation.
	GetConnection() interface{}
}

// ClientConfig holds configuration for gRPC clients.
type ClientConfig struct {
	// Target is the server address in the format "host:port".
	Target string

	// Timeout is the timeout for connection establishment.
	Timeout time.Duration

	// MaxRecvMsgSize is the maximum message size the client can receive.
	MaxRecvMsgSize int

	// MaxSendMsgSize is the maximum message size the client can send.
	MaxSendMsgSize int

	// EnableRetry enables automatic retry of failed RPCs.
	EnableRetry bool

	// MaxRetryAttempts is the maximum number of retry attempts.
	MaxRetryAttempts int

	// RetryBackoff is the backoff duration between retry attempts.
	RetryBackoff time.Duration

	// EnableTLS enables TLS for secure connections.
	EnableTLS bool

	// TLSCertPath is the path to the TLS certificate file.
	TLSCertPath string
}

// DefaultClientConfig returns the default configuration for a gRPC client.
func DefaultClientConfig() ClientConfig {
	return ClientConfig{
		Target:           "localhost:50051",
		Timeout:          time.Second * 10,
		MaxRecvMsgSize:   4 * 1024 * 1024, // 4MB
		MaxSendMsgSize:   4 * 1024 * 1024, // 4MB
		EnableRetry:      true,
		MaxRetryAttempts: 3,
		RetryBackoff:     time.Second * 1,
		EnableTLS:        false,
		TLSCertPath:      "",
	}
}