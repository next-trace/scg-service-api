// Package health provides health check functionality.
package health

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	apphealth "github.com/next-trace/scg-service-api/application/health"
	applogger "github.com/next-trace/scg-service-api/application/logger"
)

// httpHandler implements the health.Handler interface.
type httpHandler struct {
	registry apphealth.Registry
	config   apphealth.Config
	log      applogger.Logger
}

// NewHTTPHandler creates a new HTTP handler for health checks.
func NewHTTPHandler(registry apphealth.Registry, config apphealth.Config, log applogger.Logger) apphealth.Handler {
	return &httpHandler{
		registry: registry,
		config:   config,
		log:      log,
	}
}

// LivenessHandler returns an HTTP handler for liveness checks.
func (h *httpHandler) LivenessHandler() interface{} {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.handleChecks(w, r, apphealth.CheckTypeLiveness)
	})
}

// ReadinessHandler returns an HTTP handler for readiness checks.
func (h *httpHandler) ReadinessHandler() interface{} {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.handleChecks(w, r, apphealth.CheckTypeReadiness)
	})
}

// HealthHandler returns an HTTP handler for all health checks.
func (h *httpHandler) HealthHandler() interface{} {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.handleAllChecks(w, r)
	})
}

// handleChecks runs all health checks of the given type and returns the results.
func (h *httpHandler) handleChecks(w http.ResponseWriter, r *http.Request, checkType apphealth.CheckType) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), h.config.Timeout)
	defer cancel()

	// Get all checks of the given type
	checks := h.registry.GetChecks(checkType)

	// Run all checks and collect results
	results := make(map[string]apphealth.Result)
	overallStatus := apphealth.StatusUp

	for name, check := range checks {
		result := check(ctx)
		results[name] = result

		// Update overall status based on the result
		if result.Status == apphealth.StatusDown {
			overallStatus = apphealth.StatusDown
		} else if result.Status == apphealth.StatusDegraded && overallStatus != apphealth.StatusDown {
			overallStatus = apphealth.StatusDegraded
		}
	}

	// Create the response
	response := map[string]interface{}{
		"status":    overallStatus,
		"timestamp": time.Now(),
		"checks":    results,
	}

	// Set the status code based on the overall status
	statusCode := http.StatusOK
	switch overallStatus {
	case apphealth.StatusDown:
		statusCode = http.StatusServiceUnavailable
	case apphealth.StatusDegraded:
		statusCode = http.StatusOK // Still OK, but degraded
	}

	// Write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.log.Error(ctx, err, "failed to encode health check response")
	}
}

// handleAllChecks runs all health checks and returns the results.
func (h *httpHandler) handleAllChecks(w http.ResponseWriter, r *http.Request) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), h.config.Timeout)
	defer cancel()

	// Get all checks
	livenessChecks := h.registry.GetChecks(apphealth.CheckTypeLiveness)
	readinessChecks := h.registry.GetChecks(apphealth.CheckTypeReadiness)

	// Run all checks and collect results
	livenessResults := make(map[string]apphealth.Result)
	readinessResults := make(map[string]apphealth.Result)
	overallStatus := apphealth.StatusUp

	// Run liveness checks
	for name, check := range livenessChecks {
		result := check(ctx)
		livenessResults[name] = result

		// Update overall status based on the result
		if result.Status == apphealth.StatusDown {
			overallStatus = apphealth.StatusDown
		} else if result.Status == apphealth.StatusDegraded && overallStatus != apphealth.StatusDown {
			overallStatus = apphealth.StatusDegraded
		}
	}

	// Run readiness checks
	for name, check := range readinessChecks {
		result := check(ctx)
		readinessResults[name] = result

		// Update overall status based on the result
		if result.Status == apphealth.StatusDown {
			overallStatus = apphealth.StatusDown
		} else if result.Status == apphealth.StatusDegraded && overallStatus != apphealth.StatusDown {
			overallStatus = apphealth.StatusDegraded
		}
	}

	// Create the response
	response := map[string]interface{}{
		"status":    overallStatus,
		"timestamp": time.Now(),
		"checks": map[string]interface{}{
			"liveness":  livenessResults,
			"readiness": readinessResults,
		},
	}

	// Set the status code based on the overall status
	statusCode := http.StatusOK
	switch overallStatus {
	case apphealth.StatusDown:
		statusCode = http.StatusServiceUnavailable
	case apphealth.StatusDegraded:
		statusCode = http.StatusOK // Still OK, but degraded
	}

	// Write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.log.Error(ctx, err, "failed to encode health check response")
	}
}

// RegisterHTTPHandlers registers the health check handlers with the given HTTP server.
func RegisterHTTPHandlers(handler apphealth.Handler, mux *http.ServeMux, config apphealth.Config) {
	if !config.Enabled {
		return
	}

	// Register handlers
	if h, ok := handler.LivenessHandler().(http.Handler); ok {
		mux.Handle(config.LivenessPath, h)
	}
	if h, ok := handler.ReadinessHandler().(http.Handler); ok {
		mux.Handle(config.ReadinessPath, h)
	}
	if h, ok := handler.HealthHandler().(http.Handler); ok {
		mux.Handle(config.Path, h)
	}
}
