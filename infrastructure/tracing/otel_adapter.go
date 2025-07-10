// Package tracing provides OpenTelemetry-based tracing functionality.
package tracing

import (
	"context"
	"fmt"
	"os"

	apptracing "github.com/next-trace/scg-service-api/application/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// Ensure otelAdapter implements the apptracing.Tracer interface.
var _ apptracing.Tracer = (*otelAdapter)(nil)

// Option customizes the creation of the otelAdapter.
// This enables dependency injection and easier testing (e.g., providing a custom exporter).
// By default, NewOtelAdapter keeps the previous behavior; use NewOtelAdapterWithOptions to customize.
type Option func(*options)

type options struct {
	exporter sdktrace.SpanExporter
	res      *resource.Resource
	sampler  sdktrace.Sampler
}

// WithExporter allows providing a custom SpanExporter (e.g., an in-memory test exporter).
func WithExporter(exp sdktrace.SpanExporter) Option { return func(o *options) { o.exporter = exp } }

// WithResource allows overriding the OpenTelemetry resource.
func WithResource(r *resource.Resource) Option { return func(o *options) { o.res = r } }

// WithSampler allows overriding the sampler (defaults to TraceIDRatioBased based on cfg.SamplingRate).
func WithSampler(s sdktrace.Sampler) Option { return func(o *options) { o.sampler = s } }

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
	return NewOtelAdapterWithOptions(cfg)
}

// NewOtelAdapterWithOptions creates a new OpenTelemetry tracer adapter with custom options.
// If no options are provided, defaults are used (stdout exporter, derived resource and sampler).
func NewOtelAdapterWithOptions(cfg apptracing.Config, opts ...Option) (apptracing.Tracer, error) {
	// Validate required configuration
	if cfg.ServiceName == "" {
		return nil, fmt.Errorf("service name is required")
	}

	// Gather options
	var o options
	for _, fn := range opts {
		if fn != nil {
			fn(&o)
		}
	}

	// Exporter: use injected exporter when provided, otherwise create default
	exporter, err := createExporter(cfg, o.exporter)
	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %w", err)
	}

	// Resource
	res := o.res
	if res == nil {
		var err error
		res, err = createResource(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create resource: %w", err)
		}
	}

	// Sampler
	sampler := o.sampler
	if sampler == nil {
		sampler = configureSampler(cfg.SamplingRate)
	}

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
// If exp is provided, it will be used directly; otherwise a default exporter is created.
func createExporter(cfg apptracing.Config, exp sdktrace.SpanExporter) (sdktrace.SpanExporter, error) {
	if exp != nil {
		return exp, nil
	}
	// Only default stdout exporter is supported by default;
	output := cfg.Output
	if output == nil {
		output = os.Stdout
	}
	return stdouttrace.New(
		stdouttrace.WithWriter(output),
		stdouttrace.WithPrettyPrint(),
	)
}

// createResource creates a resource with service information.
func createResource(cfg apptracing.Config) (*resource.Resource, error) {
	return resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			"", // no schema URL to avoid version conflicts
			attribute.String("service.name", cfg.ServiceName),
			attribute.String("service.version", cfg.ServiceVersion),
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
