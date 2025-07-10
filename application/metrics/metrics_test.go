package metrics_test

import (
	"testing"

	appmetrics "github.com/next-trace/scg-service-api/application/metrics"
)

func TestDefaultConfig(t *testing.T) {
	cfg := appmetrics.DefaultConfig()
	if cfg.Namespace != "app" {
		t.Fatalf("unexpected Namespace: %s", cfg.Namespace)
	}
	if cfg.Subsystem != "" {
		t.Fatalf("unexpected Subsystem: %s", cfg.Subsystem)
	}
	if !cfg.EnableGoMetrics || !cfg.EnableProcessMetrics {
		t.Fatalf("expected Go and Process metrics enabled by default")
	}
	if cfg.Labels == nil {
		t.Fatalf("expected Labels map to be initialized")
	}
}
