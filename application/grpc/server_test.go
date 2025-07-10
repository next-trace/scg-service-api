package grpc_test

import (
	"testing"

	appgrpc "github.com/next-trace/scg-service-api/application/grpc"
)

func TestServerDefaultConfig(t *testing.T) {
	cfg := appgrpc.DefaultServerConfig()
	if cfg.MaxConcurrentStreams == 0 {
		t.Fatalf("expected non-zero MaxConcurrentStreams")
	}
	if cfg.MaxRecvMsgSize == 0 || cfg.MaxSendMsgSize == 0 {
		t.Fatalf("expected non-zero message sizes")
	}
	if cfg.ConnectionTimeout <= 0 {
		t.Fatalf("expected positive ConnectionTimeout")
	}
	if !cfg.EnableReflection {
		t.Fatalf("expected EnableReflection=true by default")
	}
	if !cfg.EnableHealthCheck {
		t.Fatalf("expected EnableHealthCheck=true by default")
	}
}
