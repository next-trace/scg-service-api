// Package ratelimit provides rate limiting functionality.
//
// Note: This package requires the following dependencies:
// - golang.org/x/time/rate
//
// See docs/dependencies.md for more information.
package ratelimit

import (
	"context"
	"sync"
	"time"

	appratelimit "github.com/hbttundar/scg-service-base/application/ratelimit"
	applogger "github.com/hbttundar/scg-service-base/application/logger"
)

// Type aliases for rate limiter types
// These are placeholders for the actual types from the rate package
// They will be replaced with the actual types when the dependencies are added
type (
	// rateLimiter represents a token bucket rate limiter
	rateLimiter struct {
		limit  float64
		burst  int
		tokens float64
		last   time.Time
		mu     sync.Mutex
	}
)

// Mock methods for the rateLimiter type
func newRateLimiter(limit float64, burst int) *rateLimiter {
	return &rateLimiter{
		limit:  limit,
		burst:  burst,
		tokens: float64(burst),
		last:   time.Now(),
	}
}

func (r *rateLimiter) Allow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	elapsed := now.Sub(r.last).Seconds()
	r.last = now
	r.tokens += elapsed * r.limit
	if r.tokens > float64(r.burst) {
		r.tokens = float64(r.burst)
	}
	if r.tokens < 1 {
		return false
	}
	r.tokens--
	return true
}

func (r *rateLimiter) AllowN(n int) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	elapsed := now.Sub(r.last).Seconds()
	r.last = now
	r.tokens += elapsed * r.limit
	if r.tokens > float64(r.burst) {
		r.tokens = float64(r.burst)
	}
	if r.tokens < float64(n) {
		return false
	}
	r.tokens -= float64(n)
	return true
}

func (r *rateLimiter) Reserve() time.Duration {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	elapsed := now.Sub(r.last).Seconds()
	r.last = now
	r.tokens += elapsed * r.limit
	if r.tokens > float64(r.burst) {
		r.tokens = float64(r.burst)
	}
	if r.tokens >= 1 {
		r.tokens--
		return 0
	}
	tokensNeeded := 1 - r.tokens
	waitTime := tokensNeeded / r.limit
	return time.Duration(waitTime * float64(time.Second))
}

func (r *rateLimiter) ReserveN(n int) time.Duration {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	elapsed := now.Sub(r.last).Seconds()
	r.last = now
	r.tokens += elapsed * r.limit
	if r.tokens > float64(r.burst) {
		r.tokens = float64(r.burst)
	}
	if r.tokens >= float64(n) {
		r.tokens -= float64(n)
		return 0
	}
	tokensNeeded := float64(n) - r.tokens
	waitTime := tokensNeeded / r.limit
	return time.Duration(waitTime * float64(time.Second))
}

func (r *rateLimiter) Wait(ctx context.Context) error {
	waitTime := r.Reserve()
	if waitTime == 0 {
		return nil
	}
	timer := time.NewTimer(waitTime)
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		timer.Stop()
		return ctx.Err()
	}
}

func (r *rateLimiter) WaitN(ctx context.Context, n int) error {
	waitTime := r.ReserveN(n)
	if waitTime == 0 {
		return nil
	}
	timer := time.NewTimer(waitTime)
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		timer.Stop()
		return ctx.Err()
	}
}

// tokenBucketLimiter implements the ratelimit.Limiter interface using a token bucket algorithm.
type tokenBucketLimiter struct {
	config    appratelimit.Config
	limiters  map[string]*rateLimiter
	mu        sync.RWMutex
	log       applogger.Logger
}

// NewTokenBucketLimiter creates a new token bucket rate limiter.
func NewTokenBucketLimiter(config appratelimit.Config, log applogger.Logger) appratelimit.Limiter {
	return &tokenBucketLimiter{
		config:   config,
		limiters: make(map[string]*rateLimiter),
		log:      log,
	}
}

// getLimiter returns a rate limiter for the given key, creating one if it doesn't exist.
func (t *tokenBucketLimiter) getLimiter(key string) *rateLimiter {
	t.mu.RLock()
	limiter, exists := t.limiters[key]
	t.mu.RUnlock()

	if exists {
		return limiter
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	// Check again in case another goroutine created the limiter while we were waiting for the lock
	limiter, exists = t.limiters[key]
	if exists {
		return limiter
	}

	// Calculate rate as tokens per second
	rate := float64(t.config.Rate) / t.config.Period.Seconds()
	limiter = newRateLimiter(rate, t.config.Burst)
	t.limiters[key] = limiter
	return limiter
}

// Allow checks if a request is allowed based on the key.
func (t *tokenBucketLimiter) Allow(ctx context.Context, key string) bool {
	if !t.config.Enabled {
		return true
	}

	limiter := t.getLimiter(key)
	return limiter.Allow()
}

// AllowN checks if n requests are allowed based on the key.
func (t *tokenBucketLimiter) AllowN(ctx context.Context, key string, n int) bool {
	if !t.config.Enabled {
		return true
	}

	limiter := t.getLimiter(key)
	return limiter.AllowN(n)
}

// Wait waits until a request is allowed based on the key.
func (t *tokenBucketLimiter) Wait(ctx context.Context, key string) error {
	if !t.config.Enabled {
		return nil
	}

	limiter := t.getLimiter(key)
	return limiter.Wait(ctx)
}

// WaitN waits until n requests are allowed based on the key.
func (t *tokenBucketLimiter) WaitN(ctx context.Context, key string, n int) error {
	if !t.config.Enabled {
		return nil
	}

	limiter := t.getLimiter(key)
	return limiter.WaitN(ctx, n)
}

// Reserve reserves a token and returns the time to wait before the token is available.
func (t *tokenBucketLimiter) Reserve(ctx context.Context, key string) time.Duration {
	if !t.config.Enabled {
		return 0
	}

	limiter := t.getLimiter(key)
	return limiter.Reserve()
}

// ReserveN reserves n tokens and returns the time to wait before the tokens are available.
func (t *tokenBucketLimiter) ReserveN(ctx context.Context, key string, n int) time.Duration {
	if !t.config.Enabled {
		return 0
	}

	limiter := t.getLimiter(key)
	return limiter.ReserveN(n)
}