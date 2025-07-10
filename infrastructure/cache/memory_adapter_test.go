package cache_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	appcache "github.com/next-trace/scg-service-api/application/cache"
	cacheimpl "github.com/next-trace/scg-service-api/infrastructure/cache"
	infraLogger "github.com/next-trace/scg-service-api/infrastructure/logger"
)

func TestMemoryAdapter_BasicOpsAndTTL(t *testing.T) {
	ctx := context.Background()
	var buf bytes.Buffer
	log := infraLogger.NewSlogAdapter(&buf, "debug")

	cfg := appcache.DefaultConfig()
	cfg.CleanupInterval = 10 * time.Millisecond
	cfg.DefaultTTL = 0

	c := cacheimpl.NewMemoryAdapter(cfg, log)
	t.Cleanup(func() { _ = c.Close() })

	// Set & Get
	if err := c.Set(ctx, "k", "v", 0); err != nil {
		t.Fatalf("set error: %v", err)
	}
	val, ok := c.Get(ctx, "k")
	if !ok || val.(string) != "v" {
		t.Fatalf("get mismatch: ok=%v val=%v", ok, val)
	}
	if !c.Has(ctx, "k") {
		t.Fatalf("expected key to exist")
	}

	// TTL expiry
	if err := c.Set(ctx, "tmp", 123, 20*time.Millisecond); err != nil {
		t.Fatalf("set ttl error: %v", err)
	}
	time.Sleep(30 * time.Millisecond)
	if c.Has(ctx, "tmp") {
		t.Fatalf("expected tmp to expire")
	}

	// Increment/Decrement
	if v, err := c.Increment(ctx, "cnt", 2); err != nil || v != 2 {
		t.Fatalf("inc got v=%d err=%v", v, err)
	}
	if v, err := c.Decrement(ctx, "cnt", 1); err != nil || v != 1 {
		t.Fatalf("dec got v=%d err=%v", v, err)
	}

	// Multi ops
	items := map[string]interface{}{"a": 1, "b": 2}
	if err := c.SetMulti(ctx, items, 0); err != nil {
		t.Fatalf("setmulti: %v", err)
	}
	vals, missing := c.GetMulti(ctx, []string{"a", "b", "c"})
	if len(missing) != 1 || missing[0] != "c" {
		t.Fatalf("missing unexpected: %v", missing)
	}
	if len(vals) != 2 {
		t.Fatalf("expected 2 values, got %d", len(vals))
	}

	// Delete and Clear
	if err := c.Delete(ctx, "k"); err != nil {
		t.Fatalf("delete error: %v", err)
	}
	if c.Has(ctx, "k") {
		t.Fatalf("expected k to be deleted")
	}
	if err := c.Clear(ctx); err != nil {
		t.Fatalf("clear error: %v", err)
	}
}
