package tracing

import (
	"context"
	"io"
)

// Tracer defines the abstract tracing interface (PORT) for all services.
// It provides a context-aware tracing contract.
type Tracer interface {
	// Start begins a new span and returns the updated context and a function to end the span.
	Start(ctx context.Context, spanName string) (context.Context, func())

	// AddEvent adds an event to the current span in the context.
	AddEvent(ctx context.Context, name string, attributes map[string]string)

	// SetAttributes sets attributes on the current span in the context.
	SetAttributes(ctx context.Context, attributes map[string]string)

	// RecordError records an error in the current span.
	RecordError(ctx context.Context, err error)

	// Shutdown gracefully shuts down the tracer, flushing any remaining spans.
	Shutdown(ctx context.Context) error
}

// Config holds configuration for tracers.
type Config struct {
	ServiceName      string
	ServiceVersion   string
	Environment      string
	ExporterType     string // "stdout", "jaeger", "otlp", etc.
	ExporterEndpoint string
	SamplingRate     float64
	Output           io.Writer // For stdout exporter
}
