package grpc_test

import (
	"context"
	"testing"
	"time"

	appgrpc "github.com/next-trace/scg-service-api/application/grpc"
	applogger "github.com/next-trace/scg-service-api/application/logger"
	infragrpc "github.com/next-trace/scg-service-api/infrastructure/grpc"
)

// stubLogger is a simple no-op logger implementing the application logger interface
// used to satisfy dependencies in tests.
type stubLogger struct{}

func (s stubLogger) Debug(_ context.Context, _ string)                                      {}
func (s stubLogger) Info(_ context.Context, _ string)                                       {}
func (s stubLogger) Warn(_ context.Context, _ string)                                       {}
func (s stubLogger) Error(_ context.Context, _ error, _ string)                             {}
func (s stubLogger) Fatal(_ context.Context, _ error, _ string)                             {}
func (s stubLogger) DebugKV(_ context.Context, _ string, _ map[string]interface{})          {}
func (s stubLogger) InfoKV(_ context.Context, _ string, _ map[string]interface{})           {}
func (s stubLogger) WarnKV(_ context.Context, _ string, _ map[string]interface{})           {}
func (s stubLogger) ErrorKV(_ context.Context, _ error, _ string, _ map[string]interface{}) {}
func (s stubLogger) FatalKV(_ context.Context, _ error, _ string, _ map[string]interface{}) {}
func (s stubLogger) WithField(_ string, _ interface{}) applogger.Logger                     { return s }

func TestClientAdapter_ConnectCloseLifecycle(t *testing.T) {
	cfg := appgrpc.ClientConfig{
		Target:           "localhost:12345",
		Timeout:          50 * time.Millisecond,
		MaxRecvMsgSize:   1024,
		MaxSendMsgSize:   1024,
		EnableRetry:      true,
		MaxRetryAttempts: 2,
		RetryBackoff:     10 * time.Millisecond,
		EnableTLS:        false,
	}
	log := stubLogger{}
	c := infragrpc.NewClientAdapter(cfg, log)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Initially not connected
	if got := c.GetConnection(); got != nil {
		t.Fatalf("expected nil connection initially, got %T", got)
	}

	// Connect should succeed and set connection
	if err := c.Connect(ctx); err != nil {
		t.Fatalf("connect error: %v", err)
	}
	if got := c.GetConnection(); got == nil {
		t.Fatalf("expected non-nil connection after connect")
	}

	// Connect again should be idempotent
	if err := c.Connect(ctx); err != nil {
		t.Fatalf("second connect error: %v", err)
	}

	// Close should reset state
	if err := c.Close(); err != nil {
		t.Fatalf("close error: %v", err)
	}
	if got := c.GetConnection(); got != nil {
		t.Fatalf("expected nil connection after close, got %T", got)
	}

	// Close again should be no-op
	if err := c.Close(); err != nil {
		t.Fatalf("second close error: %v", err)
	}
}
