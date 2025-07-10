// Package circuitbreaker defines the abstract interface (PORT) for circuit breaking.
package circuitbreaker

import (
	"context"
	"time"
)

// State represents the state of a circuit breaker.
type State string

const (
	// StateClosed indicates the circuit is closed and requests are allowed.
	StateClosed State = "CLOSED"

	// StateOpen indicates the circuit is open and requests are not allowed.
	StateOpen State = "OPEN"

	// StateHalfOpen indicates the circuit is half-open and a limited number of requests are allowed.
	StateHalfOpen State = "HALF_OPEN"
)

// Result represents the result of a circuit breaker execution.
type Result struct {
	// Success indicates whether the execution was successful.
	Success bool

	// Error is the error that occurred during execution, if any.
	Error error

	// Duration is the time it took to execute the function.
	Duration time.Duration
}

// CircuitBreaker defines the interface for circuit breaking.
type CircuitBreaker interface {
	// Execute executes the given function with circuit breaking.
	// If the circuit is open, it returns an error without executing the function.
	// If the circuit is closed or half-open, it executes the function and updates the circuit state.
	Execute(ctx context.Context, name string, fn func(ctx context.Context) (interface{}, error)) (interface{}, error)

	// ExecuteWithFallback executes the given function with circuit breaking and a fallback.
	// If the circuit is open or the function fails, it executes the fallback function.
	ExecuteWithFallback(ctx context.Context, name string, fn func(ctx context.Context) (interface{}, error), fallback func(ctx context.Context, err error) (interface{}, error)) (interface{}, error)

	// GetState returns the current state of the circuit breaker for the given name.
	GetState(name string) State

	// Reset resets the circuit breaker for the given name to its initial state.
	Reset(name string)
}

// Config holds configuration for circuit breakers.
type Config struct {
	// Enabled determines if circuit breaking is enabled.
	Enabled bool

	// Timeout is the maximum time to wait for a function to complete.
	Timeout time.Duration

	// MaxConcurrentRequests is the maximum number of concurrent requests allowed.
	MaxConcurrentRequests int

	// RequestVolumeThreshold is the minimum number of requests needed before a circuit can trip.
	RequestVolumeThreshold int

	// ErrorThresholdPercentage is the error percentage at which the circuit should trip.
	ErrorThresholdPercentage int

	// SleepWindow is the time to wait before transitioning from open to half-open.
	SleepWindow time.Duration

	// HealthCheckInterval is the interval at which to check the health of the circuit.
	HealthCheckInterval time.Duration
}

// DefaultConfig returns the default configuration for circuit breakers.
func DefaultConfig() Config {
	return Config{
		Enabled:                  true,
		Timeout:                  time.Second * 1,
		MaxConcurrentRequests:    100,
		RequestVolumeThreshold:   20,
		ErrorThresholdPercentage: 50,
		SleepWindow:              time.Second * 5,
		HealthCheckInterval:      time.Second * 10,
	}
}
