package middleware

import (
	"net/http"

	"github.com/next-trace/scg-service-api/application/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// TracingMiddleware provides middleware to handle trace propagation.
type TracingMiddleware struct {
	tracer     tracing.Tracer
	propagator propagation.TextMapPropagator
}

// NewTracingMiddleware creates a new tracing middleware.
func NewTracingMiddleware(tracer tracing.Tracer) *TracingMiddleware {
	return &TracingMiddleware{
		tracer:     tracer,
		propagator: otel.GetTextMapPropagator(),
	}
}

// Middleware returns an http.Handler middleware function.
func (tm *TracingMiddleware) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract trace context from incoming request
			ctx := tm.propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))

			// Start a new span
			spanCtx, endSpan := tm.tracer.Start(ctx, r.URL.Path)
			defer endSpan()

			// Add HTTP request attributes to the span
			tm.tracer.SetAttributes(spanCtx, map[string]string{
				"http.method": r.Method,
				"http.url":    r.URL.String(),
				"http.host":   r.Host,
			})

			// Pass the new context with the span down to the next handlers
			next.ServeHTTP(w, r.WithContext(spanCtx))
		})
	}
}

// Tracing provides backward compatibility with the old API.
// Deprecated: Use NewTracingMiddleware instead.
func Tracing(serviceName string) func(http.Handler) http.Handler {
	// Use the global tracer for backward compatibility
	tracer := otel.Tracer(serviceName)
	propagator := otel.GetTextMapPropagator()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))
			ctx, span := tracer.Start(ctx, r.URL.Path)
			defer span.End()

			// Pass the new context with the span down to the next handlers.
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
