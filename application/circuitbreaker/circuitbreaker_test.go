package circuitbreaker_test

import (
	"testing"
	"time"

	appcb "github.com/next-trace/scg-service-api/application/circuitbreaker"
)

func TestDefaultConfig(t *testing.T) {
	cfg := appcb.DefaultConfig()

	if !cfg.Enabled {
		t.Fatalf("expected Enabled=true by default")
	}
	if cfg.Timeout != time.Second*1 {
		t.Fatalf("unexpected Timeout: %v", cfg.Timeout)
	}
	if cfg.MaxConcurrentRequests != 100 {
		t.Fatalf("unexpected MaxConcurrentRequests: %d", cfg.MaxConcurrentRequests)
	}
	if cfg.RequestVolumeThreshold != 20 {
		t.Fatalf("unexpected RequestVolumeThreshold: %d", cfg.RequestVolumeThreshold)
	}
	if cfg.ErrorThresholdPercentage != 50 {
		t.Fatalf("unexpected ErrorThresholdPercentage: %d", cfg.ErrorThresholdPercentage)
	}
	if cfg.SleepWindow != time.Second*5 {
		t.Fatalf("unexpected SleepWindow: %v", cfg.SleepWindow)
	}
	if cfg.HealthCheckInterval != time.Second*10 {
		t.Fatalf("unexpected HealthCheckInterval: %v", cfg.HealthCheckInterval)
	}
}
