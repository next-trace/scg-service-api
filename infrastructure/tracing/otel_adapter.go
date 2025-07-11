// Package tracing provides OpenTelemetry-based tracing functionality.
package tracing

import (
	"context"
	"fmt"
	"os"

	apptracing "github.com/hbttundar/scg-service-base/application/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

// otelAdapter implements the apptracing.Tracer interface using OpenTelemetry.
// It encapsulates the OpenTelemetry-specific implementation details.
type otelAdapter struct {
	tracer         trace.Tracer
	tracerProvider *sdktrace.TracerProvider
}

// NewOtelAdapter creates a new OpenTelemetry tracer adapter.
// It configures the tracer based on the provided configuration and
// returns an implementation of the apptracing.Tracer interface.
func NewOtelAdapter(cfg apptracing.Config) (apptracing.Tracer, error) {
	// Validate required configuration
	if cfg.ServiceName == "" {
		return nil, fmt.Errorf("service name is required")
	}

	// Create exporter based on configuration
	exporter, err := createExporter(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %w", err)
	}

	// Create resource with service information
	res, err := createResource(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Configure sampling rate
	sampler := configureSampler(cfg.SamplingRate)

	// Create tracer provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sampler),
	)

	// Set global tracer provider and propagator
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Create tracer
	tracer := tp.Tracer(cfg.ServiceName)

	return &otelAdapter{
		tracer:         tracer,
		tracerProvider: tp,
	}, nil
}

// createExporter creates a span exporter based on the provided configuration.
func createExporter(cfg apptracing.Config) (sdktrace.SpanExporter, error) {
	switch cfg.ExporterType {
	case "stdout":
		output := cfg.Output
		if output == nil {
			output = os.Stdout
		}
		return stdouttrace.New(
			stdouttrace.WithWriter(output),
			stdouttrace.WithPrettyPrint(),
		)
	// Add other exporters as needed (Jaeger, OTLP, etc.)
	// case "jaeger":
	//     exporter, err = jaeger.New(...)
	// case "otlp":
	//     exporter, err = otlptrace.New(...)
	default:
		// Default to stdout for development
		return stdouttrace.New(stdouttrace.WithPrettyPrint())
	}
}

// createResource creates a resource with service information.
func createResource(cfg apptracing.Config) (*resource.Resource, error) {
	return resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion(cfg.ServiceVersion),
			attribute.String("environment", cfg.Environment),
		),
	)
}

// configureSampler returns an appropriate sampler based on the sampling rate.
func configureSampler(samplingRate float64) sdktrace.Sampler {
	if samplingRate >= 1.0 {
		return sdktrace.AlwaysSample()
	} else if samplingRate <= 0.0 {
		return sdktrace.NeverSample()
	}
	return sdktrace.TraceIDRatioBased(samplingRate)
}

// Start begins a new span and returns the updated context and a function to end the span.
// The returned context contains the new span, and the function should be called
// when the operation being traced is complete.
func (o *otelAdapter) Start(ctx context.Context, spanName string) (context.Context, func()) {
	if ctx == nil {
		ctx = context.Background()
	}

	ctx, span := o.tracer.Start(ctx, spanName)
	return ctx, func() { 
		if span != nil {
			span.End() 
		}
	}
}

// AddEvent adds an event to the current span in the context.
// Events are timestamped annotations with optional attributes.
// If the context doesn't contain a valid span, this is a no-op.
func (o *otelAdapter) AddEvent(ctx context.Context, name string, attributes map[string]string) {
	if ctx == nil || name == "" {
		return
	}

	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return
	}

	// Convert string map to OpenTelemetry attributes
	attrs := convertToAttributes(attributes)
	span.AddEvent(name, trace.WithAttributes(attrs...))
}

// SetAttributes sets attributes on the current span in the context.
// Attributes provide additional information about the operation being traced.
// If the context doesn't contain a valid span, this is a no-op.
func (o *otelAdapter) SetAttributes(ctx context.Context, attributes map[string]string) {
	if ctx == nil || len(attributes) == 0 {
		return
	}

	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return
	}

	for k, v := range attributes {
		span.SetAttributes(attribute.String(k, v))
	}
}

// RecordError records an error in the current span.
// This marks the span as failed and adds error information to it.
// If err is nil or the context doesn't contain a valid span, this is a no-op.
func (o *otelAdapter) RecordError(ctx context.Context, err error) {
	if ctx == nil || err == nil {
		return
	}

	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return
	}

	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}

// Shutdown gracefully shuts down the tracer, flushing any remaining spans.
// It should be called when the application is shutting down to ensure all
// spans are properly exported.
func (o *otelAdapter) Shutdown(ctx context.Context) error {
	if o.tracerProvider == nil {
		return nil
	}

	if ctx == nil {
		ctx = context.Background()
	}

	return o.tracerProvider.Shutdown(ctx)
}

// convertToAttributes converts a string map to OpenTelemetry attributes.
func convertToAttributes(attributes map[string]string) []attribute.KeyValue {
	if len(attributes) == 0 {
		return nil
	}

	attrs := make([]attribute.KeyValue, 0, len(attributes))
	for k, v := range attributes {
		attrs = append(attrs, attribute.String(k, v))
	}
	return attrs
}
