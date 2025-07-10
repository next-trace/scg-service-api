package health_test

import (
	"testing"
	"time"

	health "github.com/next-trace/scg-service-api/application/health"
)

func TestDefaultConfig(t *testing.T) {
	cfg := health.DefaultConfig()
	if !cfg.Enabled {
		t.Fatalf("expected Enabled=true by default")
	}
	if cfg.Path == "" || cfg.LivenessPath == "" || cfg.ReadinessPath == "" {
		t.Fatalf("expected default paths to be set")
	}
	if cfg.Timeout != 5*time.Second {
		t.Fatalf("unexpected default Timeout: %v", cfg.Timeout)
	}
}
