// Package health defines the abstract interface (PORT) for health checks.
package health

import (
	"context"
	"time"
)

// Status represents the health status of a component.
type Status string

const (
	// StatusUp indicates the component is healthy and operational.
	StatusUp Status = "UP"

	// StatusDown indicates the component is unhealthy or not operational.
	StatusDown Status = "DOWN"

	// StatusDegraded indicates the component is operational but with reduced functionality.
	StatusDegraded Status = "DEGRADED"
)

// CheckType represents the type of health check.
type CheckType string

const (
	// CheckTypeLiveness is used to determine if the application is alive.
	// A failed liveness check typically results in the application being restarted.
	CheckTypeLiveness CheckType = "liveness"

	// CheckTypeReadiness is used to determine if the application is ready to receive traffic.
	// A failed readiness check typically results in the application being removed from service discovery.
	CheckTypeReadiness CheckType = "readiness"
)

// Result represents the result of a health check.
type Result struct {
	// Status is the health status of the component.
	Status Status `json:"status"`

	// Component is the name of the component being checked.
	Component string `json:"component"`

	// Details contains additional information about the health check.
	Details map[string]interface{} `json:"details,omitempty"`

	// Error is the error message if the check failed.
	Error string `json:"error,omitempty"`

	// Timestamp is when the check was performed.
	Timestamp time.Time `json:"timestamp"`
}

// Check defines a health check function.
type Check func(ctx context.Context) Result

// Registry defines the interface for registering health checks.
type Registry interface {
	// RegisterCheck registers a health check with the given name and type.
	RegisterCheck(name string, checkType CheckType, check Check)

	// UnregisterCheck removes a health check with the given name and type.
	UnregisterCheck(name string, checkType CheckType)

	// GetChecks returns all registered health checks of the given type.
	GetChecks(checkType CheckType) map[string]Check
}

// Handler defines the interface for handling health check requests.
type Handler interface {
	// LivenessHandler returns an HTTP handler for liveness checks.
	LivenessHandler() interface{}

	// ReadinessHandler returns an HTTP handler for readiness checks.
	ReadinessHandler() interface{}

	// HealthHandler returns an HTTP handler for all health checks.
	HealthHandler() interface{}
}

// Config holds configuration for health checks.
type Config struct {
	// Enabled determines if health checks are enabled.
	Enabled bool

	// Path is the base path for health check endpoints.
	Path string

	// LivenessPath is the path for liveness checks.
	LivenessPath string

	// ReadinessPath is the path for readiness checks.
	ReadinessPath string

	// Timeout is the maximum time to wait for a health check to complete.
	Timeout time.Duration
}

// DefaultConfig returns the default configuration for health checks.
func DefaultConfig() Config {
	return Config{
		Enabled:       true,
		Path:          "/health",
		LivenessPath:  "/health/liveness",
		ReadinessPath: "/health/readiness",
		Timeout:       time.Second * 5,
	}
}
