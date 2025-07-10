package grpc_test

import (
	"bytes"
	"context"
	"net"
	"testing"

	appgrpc "github.com/next-trace/scg-service-api/application/grpc"
	grpcimpl "github.com/next-trace/scg-service-api/infrastructure/grpc"
	infraLogger "github.com/next-trace/scg-service-api/infrastructure/logger"
)

type dummyService struct{ registered bool }

func (d *dummyService) Register(_ interface{}) { d.registered = true }

func TestServerAdapter_RegisterStartStop(t *testing.T) {
	var buf bytes.Buffer
	log := infraLogger.NewSlogAdapter(&buf, "info")
	cfg := appgrpc.DefaultServerConfig()
	cfg.EnableHealthCheck = true

	srv := grpcimpl.NewServerAdapter(cfg, log)
	d := &dummyService{}
	if err := srv.RegisterService(d); err != nil {
		t.Fatalf("register service: %v", err)
	}
	if !d.registered {
		t.Fatalf("expected dummy service to be registered")
	}

	lc := &net.ListenConfig{}
	ln, err := lc.Listen(context.Background(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	t.Cleanup(func() { _ = ln.Close() })

	if err := srv.Start(context.Background(), ln); err != nil {
		t.Fatalf("start: %v", err)
	}
	if err := srv.Stop(context.Background()); err != nil {
		t.Fatalf("stop: %v", err)
	}
}
