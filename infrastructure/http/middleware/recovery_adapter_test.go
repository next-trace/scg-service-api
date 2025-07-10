package middleware_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/next-trace/scg-service-api/infrastructure/http/middleware"
	"github.com/next-trace/scg-service-api/infrastructure/logger"
	"github.com/stretchr/testify/assert"
)

// mockHandler is a test handler that can be configured to panic
type mockHandler struct {
	shouldPanic bool
	panicValue  interface{}
}

func (h *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_ = r
	if h.shouldPanic {
		panic(h.panicValue)
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

func TestRecovery(t *testing.T) {
	// Create a buffer to capture log output
	var logBuffer bytes.Buffer
	log := logger.NewSlogAdapter(&logBuffer, "debug")

	t.Run("No panic", func(t *testing.T) {
		// Create a handler that doesn't panic
		handler := &mockHandler{shouldPanic: false}

		// Wrap it with the recovery middleware
		recoveryMiddleware := middleware.NewRecoveryMiddleware(log).Middleware()
		wrappedHandler := recoveryMiddleware(handler)

		// Create a test request
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		// Call the handler
		wrappedHandler.ServeHTTP(w, req)

		// Verify the response
		resp := w.Result()
		t.Cleanup(func() { _ = resp.Body.Close() })

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Empty(t, logBuffer.String()) // No log entries should be created
	})

	t.Run("With panic - error", func(t *testing.T) {
		// Reset the log buffer
		logBuffer.Reset()

		// Create a handler that panics with an error
		testErr := "test error"
		handler := &mockHandler{shouldPanic: true, panicValue: testErr}

		// Wrap it with the recovery middleware
		recoveryMiddleware := middleware.NewRecoveryMiddleware(log).Middleware()
		wrappedHandler := recoveryMiddleware(handler)

		// Create a test request
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		// Call the handler
		wrappedHandler.ServeHTTP(w, req)

		// Verify the response
		resp := w.Result()
		t.Cleanup(func() { _ = resp.Body.Close() })

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		// Verify the log entry
		logOutput := logBuffer.String()
		assert.Contains(t, logOutput, "panic recovered")
		assert.Contains(t, logOutput, "test error")
		assert.Contains(t, logOutput, "stack")
	})

	t.Run("With panic - non-error", func(t *testing.T) {
		// Reset the log buffer
		logBuffer.Reset()

		// Create a handler that panics with a non-error value
		handler := &mockHandler{shouldPanic: true, panicValue: 42}

		// Wrap it with the recovery middleware
		recoveryMiddleware := middleware.NewRecoveryMiddleware(log).Middleware()
		wrappedHandler := recoveryMiddleware(handler)

		// Create a test request
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		// Call the handler
		wrappedHandler.ServeHTTP(w, req)

		// Verify the response
		resp := w.Result()
		t.Cleanup(func() { _ = resp.Body.Close() })

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		// Verify the log entry
		logOutput := logBuffer.String()
		assert.Contains(t, logOutput, "panic recovered")
		assert.Contains(t, logOutput, "42")
		assert.Contains(t, logOutput, "stack")
	})
}
