package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/next-trace/scg-service-api/infrastructure/http/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

// MockTracer is a mock implementation of the tracing.Tracer interface
type MockTracer struct {
	mock.Mock
}

func (m *MockTracer) Start(ctx context.Context, spanName string) (context.Context, func()) {
	args := m.Called(ctx, spanName)
	return args.Get(0).(context.Context), args.Get(1).(func())
}

func (m *MockTracer) AddEvent(ctx context.Context, name string, attributes map[string]string) {
	m.Called(ctx, name, attributes)
}

func (m *MockTracer) SetAttributes(ctx context.Context, attributes map[string]string) {
	m.Called(ctx, attributes)
}

func (m *MockTracer) RecordError(ctx context.Context, err error) {
	m.Called(ctx, err)
}

func (m *MockTracer) Shutdown(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestTracingMiddleware(t *testing.T) {
	t.Run("Middleware creates spans", func(t *testing.T) {
		// Create a mock tracer
		mockTracer := new(MockTracer)

		// Define a context key
		testKey := contextKey("test")

		// Set up expectations
		mockCtx := context.WithValue(t.Context(), testKey, "value")
		mockTracer.On("Start", mock.Anything, "/test").Return(mockCtx, func() {})
		mockTracer.On("SetAttributes", mockCtx, mock.MatchedBy(func(attrs map[string]string) bool {
			return attrs["http.method"] == "GET" &&
				attrs["http.url"] == "/test" &&
				attrs["http.host"] == "example.com"
		})).Return()

		// Create the middleware
		tracingMiddleware := middleware.NewTracingMiddleware(mockTracer)

		// Create a test handler
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify that the context was updated
			assert.Equal(t, "value", r.Context().Value(testKey))
			w.WriteHeader(http.StatusOK)
		})

		// Wrap the handler with the middleware
		wrappedHandler := tracingMiddleware.Middleware()(testHandler)

		// Create a test request
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Host = "example.com"
		w := httptest.NewRecorder()

		// Call the handler
		wrappedHandler.ServeHTTP(w, req)

		// Verify the response
		resp := w.Result()
		t.Cleanup(func() { _ = resp.Body.Close() })

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify that the tracer methods were called
		mockTracer.AssertExpectations(t)
	})
}

func TestDeprecatedTracing(t *testing.T) {
	t.Run("Deprecated tracing function works", func(t *testing.T) {
		// Create a test handler
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Just verify that the handler is called
			w.WriteHeader(http.StatusOK)
		})

		// Wrap the handler with the deprecated middleware
		wrappedHandler := middleware.Tracing("test-service")(testHandler)

		// Create a test request
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		// Call the handler
		wrappedHandler.ServeHTTP(w, req)

		// Verify the response
		resp := w.Result()
		t.Cleanup(func() { _ = resp.Body.Close() })

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		// We can't easily verify the tracing behavior since it uses global state
	})
}
