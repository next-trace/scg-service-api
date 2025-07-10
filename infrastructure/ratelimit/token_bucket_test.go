package ratelimit_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	appratelimit "github.com/next-trace/scg-service-api/application/ratelimit"
	infraLogger "github.com/next-trace/scg-service-api/infrastructure/logger"
	limiterimpl "github.com/next-trace/scg-service-api/infrastructure/ratelimit"
)

func TestTokenBucketLimiter_AllowAndWait(t *testing.T) {
	ctx := context.Background()
	var buf bytes.Buffer
	log := infraLogger.NewSlogAdapter(&buf, "debug")

	cfg := appratelimit.DefaultConfig()
	cfg.Rate = 2
	cfg.Period = 50 * time.Millisecond
	cfg.Burst = 1
	cfg.WaitTimeout = 100 * time.Millisecond

	lim := limiterimpl.NewTokenBucketLimiter(cfg, log)

	key := "user:1"
	// Allow a few rapid requests within burst/rate
	if !lim.Allow(ctx, key) {
		t.Fatalf("expected first allow")
	}
	// Depending on timing, the second may or may not pass, but Wait should not error
	if err := lim.Wait(ctx, key); err != nil {
		t.Fatalf("wait error: %v", err)
	}

	// Reserve returns a non-negative duration
	if d := lim.Reserve(ctx, key); d < 0 {
		t.Fatalf("expected non-negative reserve, got %v", d)
	}
}
