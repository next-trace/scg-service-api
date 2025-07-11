// Package middleware provides HTTP middleware components.
package middleware

import (
	"bufio"
	"net"
	"net/http"
	"strconv"
	"time"

	appmetrics "github.com/hbttundar/scg-service-base/application/metrics"
)

// MetricsMiddleware provides middleware to collect metrics for HTTP requests.
type MetricsMiddleware struct {
	metrics appmetrics.Metrics
}

// NewMetricsMiddleware creates a new metrics middleware.
func NewMetricsMiddleware(metrics appmetrics.Metrics) *MetricsMiddleware {
	return &MetricsMiddleware{
		metrics: metrics,
	}
}

// Middleware returns an http.Handler middleware function.
func (mm *MetricsMiddleware) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create a response writer wrapper to capture the status code
			rw := newResponseWriterWrapper(w)

			// Start a timer for the request duration
			stopTimer := mm.metrics.TimerStart("http_request_duration_seconds")

			// Increment the request counter
			mm.metrics.CounterInc("http_requests_total")

			// Call the next handler
			next.ServeHTTP(rw, r)

			// Stop the timer and get the duration
			duration := stopTimer()

			// Record metrics with labels
			mm.metrics.WithLabels(map[string]string{
				"method":      r.Method,
				"path":        r.URL.Path,
				"status_code": strconv.Itoa(rw.statusCode),
			}).CounterInc("http_requests_total")

			// Record request size
			if r.ContentLength > 0 {
				mm.metrics.HistogramObserve("http_request_size_bytes", float64(r.ContentLength))
			}

			// Record response size
			mm.metrics.HistogramObserve("http_response_size_bytes", float64(rw.bytesWritten))

			// Log slow requests (more than 1 second)
			if duration > time.Second {
				mm.metrics.CounterInc("http_slow_requests_total")
			}
		})
	}
}

// Metrics provides backward compatibility with the old API.
// Deprecated: Use NewMetricsMiddleware instead.
func Metrics(metrics appmetrics.Metrics) func(http.Handler) http.Handler {
	return NewMetricsMiddleware(metrics).Middleware()
}

// responseWriterWrapper wraps an http.ResponseWriter to capture the status code and bytes written.
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int64
}

// newResponseWriterWrapper creates a new response writer wrapper.
func newResponseWriterWrapper(w http.ResponseWriter) *responseWriterWrapper {
	return &responseWriterWrapper{
		ResponseWriter: w,
		statusCode:     http.StatusOK, // Default status code
	}
}

// WriteHeader captures the status code and calls the underlying WriteHeader.
func (rw *responseWriterWrapper) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

// Write captures the bytes written and calls the underlying Write.
func (rw *responseWriterWrapper) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.bytesWritten += int64(n)
	return n, err
}

// Flush implements the http.Flusher interface if the underlying response writer supports it.
func (rw *responseWriterWrapper) Flush() {
	if f, ok := rw.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// Hijack implements the http.Hijacker interface if the underlying response writer supports it.
func (rw *responseWriterWrapper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := rw.ResponseWriter.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, http.ErrNotSupported
}
