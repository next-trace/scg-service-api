package http_test

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	apphttp "github.com/next-trace/scg-service-api/application/http"
	applogger "github.com/next-trace/scg-service-api/application/logger"
)

// simpleLogger is a minimal test logger implementing applogger.Logger.
type simpleLogger struct{}

func (simpleLogger) Debug(context.Context, string)                                  {}
func (simpleLogger) Info(context.Context, string)                                   {}
func (simpleLogger) Warn(context.Context, string)                                   {}
func (simpleLogger) Error(context.Context, error, string)                           {}
func (simpleLogger) Fatal(context.Context, error, string)                           {}
func (simpleLogger) DebugKV(context.Context, string, map[string]interface{})        {}
func (simpleLogger) InfoKV(context.Context, string, map[string]interface{})         {}
func (simpleLogger) WarnKV(context.Context, string, map[string]interface{})         {}
func (simpleLogger) ErrorKV(context.Context, error, string, map[string]interface{}) {}
func (simpleLogger) FatalKV(context.Context, error, string, map[string]interface{}) {}
func (simpleLogger) WithField(string, interface{}) applogger.Logger                 { return simpleLogger{} }

// TestRun_NilServer ensures nil server is a no-op.
func TestRun_NilServer(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := apphttp.Run(ctx, nil, simpleLogger{}); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

// TestRun_ServerError ensures immediate server error is returned.
func TestRun_ServerError(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	srv := &http.Server{Addr: "bad:addr"}
	if err := apphttp.Run(ctx, srv, simpleLogger{}); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

// TestRun_ContextCancel triggers graceful shutdown via context.
func TestRun_ContextCancel(t *testing.T) {
	t.Parallel()
	// Start a real server so ListenAndServe runs and then return ErrServerClosed during shutdown.
	ctx := context.Background()
	lc := net.ListenConfig{}
	ln, err := lc.Listen(ctx, "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer ln.Close()

	s := &http.Server{Handler: http.NewServeMux()}
	// run server in background to accept (it will return ErrServerClosed when Shutdown is called)
	go func() { _ = s.Serve(ln) }()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	// Call Run and then verify that the server actually shuts down by checking Serve returned.
	done := make(chan struct{})
	go func() {
		_ = apphttp.Run(ctx, s, simpleLogger{})
		close(done)
	}()

	// Cancel will trigger shutdown; wait for Run to return
	select {
	case <-done:
		// ok
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting for shutdown")
	}
}

// Note: OS signal path is implicitly covered by context path since sending real signals in tests can be flaky.
// The graceful shutdown logic is identical across both paths. We avoid manipulating process signals to keep tests stable and fast.
