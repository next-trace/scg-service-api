// Package health provides health check functionality.
package health

import (
	"context"
	"sync"
	"time"

	apphealth "github.com/next-trace/scg-service-api/application/health"
)

// registry implements the health.Registry interface.
type registry struct {
	checks map[apphealth.CheckType]map[string]apphealth.Check
	mu     sync.RWMutex
}

// NewRegistry creates a new health check registry.
func NewRegistry() apphealth.Registry {
	return &registry{
		checks: map[apphealth.CheckType]map[string]apphealth.Check{
			apphealth.CheckTypeLiveness:  make(map[string]apphealth.Check),
			apphealth.CheckTypeReadiness: make(map[string]apphealth.Check),
		},
	}
}

// RegisterCheck registers a health check with the given name and type.
func (r *registry) RegisterCheck(name string, checkType apphealth.CheckType, check apphealth.Check) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.checks[checkType]; !ok {
		r.checks[checkType] = make(map[string]apphealth.Check)
	}

	r.checks[checkType][name] = check
}

// UnregisterCheck removes a health check with the given name and type.
func (r *registry) UnregisterCheck(name string, checkType apphealth.CheckType) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if checks, ok := r.checks[checkType]; ok {
		delete(checks, name)
	}
}

// GetChecks returns all registered health checks of the given type.
func (r *registry) GetChecks(checkType apphealth.CheckType) map[string]apphealth.Check {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if checks, ok := r.checks[checkType]; ok {
		// Create a copy to avoid concurrent map access
		result := make(map[string]apphealth.Check, len(checks))
		for k, v := range checks {
			result[k] = v
		}
		return result
	}

	return make(map[string]apphealth.Check)
}

// RegisterCommonChecks registers common health checks that are useful for most applications.
func RegisterCommonChecks(registry apphealth.Registry) {
	// Register a basic liveness check that always returns UP
	registry.RegisterCheck("service", apphealth.CheckTypeLiveness, BasicLivenessCheck())
}

// BasicLivenessCheck returns a simple liveness check that always returns UP.
// This is useful as a default liveness check to indicate the service is running.
func BasicLivenessCheck() apphealth.Check {
	return func(ctx context.Context) apphealth.Result {
		return apphealth.Result{
			Status:    apphealth.StatusUp,
			Component: "service",
			Timestamp: time.Now(),
		}
	}
}
