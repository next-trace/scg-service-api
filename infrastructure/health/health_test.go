package health_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	apphealth "github.com/next-trace/scg-service-api/application/health"
	healthimpl "github.com/next-trace/scg-service-api/infrastructure/health"
	infraLogger "github.com/next-trace/scg-service-api/infrastructure/logger"
)

func TestHealthHandlers_RegistryAndHTTP(t *testing.T) {
	var buf bytes.Buffer
	log := infraLogger.NewSlogAdapter(&buf, "info")

	reg := healthimpl.NewRegistry()
	// Register common liveness
	healthimpl.RegisterCommonChecks(reg)
	// Register a readiness check that returns degraded
	reg.RegisterCheck("db", apphealth.CheckTypeReadiness, func(_ context.Context) apphealth.Result {
		return apphealth.Result{Status: apphealth.StatusDegraded, Component: "db", Timestamp: time.Now()}
	})

	cfg := apphealth.DefaultConfig()
	cfg.Enabled = true
	h := healthimpl.NewHTTPHandler(reg, cfg, log)

	mux := http.NewServeMux()
	healthimpl.RegisterHTTPHandlers(h, mux, cfg)

	// Test /health endpoint
	req := httptest.NewRequest(http.MethodGet, cfg.Path, nil)
	rw := httptest.NewRecorder()
	mux.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(rw.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if body["status"] == nil {
		t.Fatalf("missing status in body")
	}
}
