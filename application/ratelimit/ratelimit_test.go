package ratelimit_test

import (
	"testing"
	"time"

	ratelimit "github.com/next-trace/scg-service-api/application/ratelimit"
)

func TestDefaultConfig(t *testing.T) {
	cfg := ratelimit.DefaultConfig()
	if cfg.Rate <= 0 {
		t.Fatalf("expected positive default Rate")
	}
	if cfg.Burst <= 0 {
		t.Fatalf("expected positive default Burst")
	}
	if cfg.WaitTimeout != time.Second {
		t.Fatalf("unexpected WaitTimeout: %v", cfg.WaitTimeout)
	}
}
