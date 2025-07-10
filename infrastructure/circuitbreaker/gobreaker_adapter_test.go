package circuitbreaker_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	appcb "github.com/next-trace/scg-service-api/application/circuitbreaker"
	cbimpl "github.com/next-trace/scg-service-api/infrastructure/circuitbreaker"
	infraLogger "github.com/next-trace/scg-service-api/infrastructure/logger"
)

func TestGoBreakerAdapter_ExecuteAndState(t *testing.T) {
	var buf bytes.Buffer
	log := infraLogger.NewSlogAdapter(&buf, "info")
	cfg := appcb.DefaultConfig()
	br := cbimpl.NewGoBreakerAdapter(cfg, log)

	ctx := context.Background()
	// Successful execution keeps state CLOSED
	res, err := br.Execute(ctx, "svc", func(context.Context) (interface{}, error) { return 42, nil })
	if err != nil || res.(int) != 42 {
		t.Fatalf("unexpected exec result: %v, %v", res, err)
	}
	if st := br.GetState("svc"); st != appcb.StateClosed {
		t.Fatalf("expected CLOSED, got %s", st)
	}

	// ExecuteWithFallback should return fallback value on error
	res, err = br.ExecuteWithFallback(ctx, "svc", func(context.Context) (interface{}, error) {
		return nil, errors.New("boom")
	}, func(context.Context, error) (interface{}, error) { return "fallback", nil })
	if err != nil || res.(string) != "fallback" {
		t.Fatalf("unexpected fallback result: %v, %v", res, err)
	}

	// Reset should not panic
	br.Reset("svc")
}
