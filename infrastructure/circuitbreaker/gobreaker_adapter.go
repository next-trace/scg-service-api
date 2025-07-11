// Package circuitbreaker provides circuit breaking functionality.
//
// Note: This package requires the following dependencies:
// - github.com/sony/gobreaker
//
// See docs/dependencies.md for more information.
package circuitbreaker

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	appcircuitbreaker "github.com/hbttundar/scg-service-base/application/circuitbreaker"
	applogger "github.com/hbttundar/scg-service-base/application/logger"
)

// Type aliases for circuit breaker types
// These are placeholders for the actual types from the gobreaker package
// They will be replaced with the actual types when the dependencies are added
type (
	// circuitBreaker represents a circuit breaker
	circuitBreaker struct {
		name                   string
		maxRequests            uint32
		interval               time.Duration
		timeout                time.Duration
		readyToTrip            func(counts interface{}) bool
		onStateChange          func(name string, from, to string)
		counts                 interface{}
		state                  string
		generation             uint64
		lastStateChangeTime    time.Time
		mutex                  sync.Mutex
	}

	// counts represents statistics for the circuit breaker
	counts struct {
		requests             uint32
		totalSuccesses       uint32
		totalFailures        uint32
		consecutiveSuccesses uint32
		consecutiveFailures  uint32
	}

	// settings represents settings for the circuit breaker
	settings struct {
		name          string
		maxRequests   uint32
		interval      time.Duration
		timeout       time.Duration
		readyToTrip   func(counts interface{}) bool
		onStateChange func(name string, from, to string)
	}
)

// Mock constants for circuit breaker states
const (
	stateOpen     = "open"
	stateClosed   = "closed"
	stateHalfOpen = "half-open"
)

// Mock functions for the circuitBreaker type
func newCircuitBreaker(st settings) *circuitBreaker {
	return &circuitBreaker{
		name:                st.name,
		maxRequests:         st.maxRequests,
		interval:            st.interval,
		timeout:             st.timeout,
		readyToTrip:         st.readyToTrip,
		onStateChange:       st.onStateChange,
		counts:              &counts{},
		state:               stateClosed,
		lastStateChangeTime: time.Now(),
	}
}

func (cb *circuitBreaker) Execute(req func() (interface{}, error)) (interface{}, error) {
	cb.mutex.Lock()
	state := cb.state
	cb.mutex.Unlock()

	if state == stateOpen {
		return nil, errors.New("circuit breaker is open")
	}

	result, err := req()

	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	c := cb.counts.(*counts)
	c.requests++

	if err != nil {
		c.totalFailures++
		c.consecutiveFailures++
		c.consecutiveSuccesses = 0

		if cb.state == stateClosed && cb.readyToTrip(c) {
			cb.setState(stateOpen)
		}
	} else {
		c.totalSuccesses++
		c.consecutiveSuccesses++
		c.consecutiveFailures = 0

		if cb.state == stateHalfOpen && c.consecutiveSuccesses >= cb.maxRequests {
			cb.setState(stateClosed)
		}
	}

	return result, err
}

func (cb *circuitBreaker) setState(state string) {
	if cb.state == state {
		return
	}

	oldState := cb.state
	cb.state = state
	cb.lastStateChangeTime = time.Now()
	cb.counts = &counts{}

	if cb.onStateChange != nil {
		cb.onStateChange(cb.name, oldState, state)
	}
}

func (cb *circuitBreaker) State() string {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	if cb.state == stateOpen && time.Since(cb.lastStateChangeTime) > cb.timeout {
		cb.setState(stateHalfOpen)
	}

	return cb.state
}

// gobreakerAdapter implements the circuitbreaker.CircuitBreaker interface using the gobreaker package.
type gobreakerAdapter struct {
	config    appcircuitbreaker.Config
	breakers  map[string]*circuitBreaker
	mu        sync.RWMutex
	log       applogger.Logger
}

// NewGoBreakerAdapter creates a new circuit breaker adapter using the gobreaker package.
func NewGoBreakerAdapter(config appcircuitbreaker.Config, log applogger.Logger) appcircuitbreaker.CircuitBreaker {
	return &gobreakerAdapter{
		config:   config,
		breakers: make(map[string]*circuitBreaker),
		log:      log,
	}
}

// getBreaker returns a circuit breaker for the given name, creating one if it doesn't exist.
func (g *gobreakerAdapter) getBreaker(name string) *circuitBreaker {
	g.mu.RLock()
	breaker, exists := g.breakers[name]
	g.mu.RUnlock()

	if exists {
		return breaker
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	// Check again in case another goroutine created the breaker while we were waiting for the lock
	breaker, exists = g.breakers[name]
	if exists {
		return breaker
	}

	// Create a new circuit breaker
	st := settings{
		name:        name,
		maxRequests: uint32(g.config.MaxConcurrentRequests),
		interval:    g.config.HealthCheckInterval,
		timeout:     g.config.SleepWindow,
		readyToTrip: func(c interface{}) bool {
			counts := c.(*counts)
			return counts.requests >= uint32(g.config.RequestVolumeThreshold) &&
				float64(counts.totalFailures)/float64(counts.requests)*100 >= float64(g.config.ErrorThresholdPercentage)
		},
		onStateChange: func(name string, from, to string) {
			g.log.InfoKV(context.Background(), "circuit breaker state changed", map[string]interface{}{
				"name":      name,
				"from":      from,
				"to":        to,
				"timestamp": time.Now().Format(time.RFC3339),
			})
		},
	}

	breaker = newCircuitBreaker(st)
	g.breakers[name] = breaker
	return breaker
}

// Execute executes the given function with circuit breaking.
func (g *gobreakerAdapter) Execute(ctx context.Context, name string, fn func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	if !g.config.Enabled {
		return fn(ctx)
	}

	breaker := g.getBreaker(name)

	// Create a context with timeout
	execCtx := ctx
	if g.config.Timeout > 0 {
		var cancel context.CancelFunc
		execCtx, cancel = context.WithTimeout(ctx, g.config.Timeout)
		defer cancel()
	}

	// Execute the function with the circuit breaker
	result, err := breaker.Execute(func() (interface{}, error) {
		return fn(execCtx)
	})

	if err != nil {
		// If the circuit is open, return a specific error
		if breaker.State() == stateOpen {
			return nil, fmt.Errorf("circuit breaker '%s' is open", name)
		}
		return nil, err
	}

	return result, nil
}

// ExecuteWithFallback executes the given function with circuit breaking and a fallback.
func (g *gobreakerAdapter) ExecuteWithFallback(ctx context.Context, name string, fn func(ctx context.Context) (interface{}, error), fallback func(ctx context.Context, err error) (interface{}, error)) (interface{}, error) {
	result, err := g.Execute(ctx, name, fn)
	if err != nil && fallback != nil {
		return fallback(ctx, err)
	}
	return result, err
}

// GetState returns the current state of the circuit breaker for the given name.
func (g *gobreakerAdapter) GetState(name string) appcircuitbreaker.State {
	g.mu.RLock()
	breaker, exists := g.breakers[name]
	g.mu.RUnlock()

	if !exists {
		return appcircuitbreaker.StateClosed
	}

	state := breaker.State()
	switch state {
	case stateOpen:
		return appcircuitbreaker.StateOpen
	case stateHalfOpen:
		return appcircuitbreaker.StateHalfOpen
	default:
		return appcircuitbreaker.StateClosed
	}
}

// Reset resets the circuit breaker for the given name to its initial state.
func (g *gobreakerAdapter) Reset(name string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	delete(g.breakers, name)
}