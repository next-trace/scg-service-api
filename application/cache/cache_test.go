package cache_test

import (
	"testing"
	"time"

	appcache "github.com/next-trace/scg-service-api/application/cache"
)

func TestDefaultConfig(t *testing.T) {
	cfg := appcache.DefaultConfig()
	if !cfg.Enabled {
		t.Fatalf("expected Enabled=true by default")
	}
	if cfg.StoreType != appcache.StoreTypeMemory {
		t.Fatalf("unexpected StoreType: %s", cfg.StoreType)
	}
	if cfg.DefaultTTL != 5*time.Minute {
		t.Fatalf("unexpected DefaultTTL: %v", cfg.DefaultTTL)
	}
	if cfg.CleanupInterval != time.Minute {
		t.Fatalf("unexpected CleanupInterval: %v", cfg.CleanupInterval)
	}
	if cfg.MaxEntries <= 0 {
		t.Fatalf("expected positive MaxEntries")
	}
	if cfg.Redis.Address == "" {
		t.Fatalf("expected default Redis address")
	}
}
