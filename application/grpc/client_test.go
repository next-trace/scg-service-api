package grpc_test

import (
	"testing"

	appgrpc "github.com/next-trace/scg-service-api/application/grpc"
)

func TestClientDefaultConfig(t *testing.T) {
	cfg := appgrpc.DefaultClientConfig()
	if cfg.Target == "" {
		t.Fatalf("expected non-empty default Target")
	}
	if cfg.Timeout == 0 {
		t.Fatalf("expected non-zero default Timeout")
	}
	// TLS should be disabled by default in the app layer defaults
	if cfg.EnableTLS {
		t.Fatalf("expected EnableTLS=false by default")
	}
}
