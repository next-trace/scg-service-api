// Package ratelimit defines the abstract interface (PORT) for rate limiting.
package ratelimit

import (
	"context"
	"time"
)

// Limiter defines the interface for rate limiting.
type Limiter interface {
	// Allow checks if a request is allowed based on the key.
	// It returns true if the request is allowed, false otherwise.
	Allow(ctx context.Context, key string) bool

	// AllowN checks if n requests are allowed based on the key.
	// It returns true if the requests are allowed, false otherwise.
	AllowN(ctx context.Context, key string, n int) bool

	// Wait waits until a request is allowed based on the key.
	// It returns an error if the context is canceled or the wait time exceeds the context deadline.
	Wait(ctx context.Context, key string) error

	// WaitN waits until n requests are allowed based on the key.
	// It returns an error if the context is canceled or the wait time exceeds the context deadline.
	WaitN(ctx context.Context, key string, n int) error

	// Reserve reserves a token and returns the time to wait before the token is available.
	// If the token cannot be reserved, it returns a negative wait time.
	Reserve(ctx context.Context, key string) time.Duration

	// ReserveN reserves n tokens and returns the time to wait before the tokens are available.
	// If the tokens cannot be reserved, it returns a negative wait time.
	ReserveN(ctx context.Context, key string, n int) time.Duration
}

// Strategy defines the rate limiting strategy.
type Strategy string

const (
	// StrategyToken is a token bucket rate limiting strategy.
	StrategyToken Strategy = "token"

	// StrategyLeakyBucket is a leaky bucket rate limiting strategy.
	StrategyLeakyBucket Strategy = "leaky_bucket"

	// StrategyFixedWindow is a fixed window rate limiting strategy.
	StrategyFixedWindow Strategy = "fixed_window"

	// StrategySlidingWindow is a sliding window rate limiting strategy.
	StrategySlidingWindow Strategy = "sliding_window"
)

// Config holds configuration for rate limiting.
type Config struct {
	// Enabled determines if rate limiting is enabled.
	Enabled bool

	// Strategy is the rate limiting strategy to use.
	Strategy Strategy

	// Rate is the number of requests allowed per period.
	Rate int

	// Period is the time period for the rate.
	Period time.Duration

	// Burst is the maximum number of requests allowed to exceed the rate.
	Burst int

	// WaitTimeout is the maximum time to wait for a token.
	WaitTimeout time.Duration

	// KeyFunc is a function that generates a key from a context.
	// If nil, a default key function will be used.
	KeyFunc func(ctx context.Context) string
}

// DefaultConfig returns the default configuration for rate limiting.
func DefaultConfig() Config {
	return Config{
		Enabled:     true,
		Strategy:    StrategyToken,
		Rate:        100,
		Period:      time.Minute,
		Burst:       10,
		WaitTimeout: time.Second,
	}
}