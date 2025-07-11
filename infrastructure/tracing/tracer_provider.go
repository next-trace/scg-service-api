// Package tracing provides OpenTelemetry-based tracing functionality.
// It implements the application's tracing interface with OpenTelemetry.
package tracing

import (
	"os"

	apptracing "github.com/hbttundar/scg-service-base/application/tracing"
)

// NewTracer creates a new tracer using the application's Tracer interface.
// This is the recommended way to create a tracer.
func NewTracer(cfg apptracing.Config) (apptracing.Tracer, error) {
	return NewOtelAdapter(cfg)
}

// DefaultTracerConfig returns a default configuration for development environment.
// It configures a tracer with:
// - stdout exporter for easy local debugging
// - 100% sampling rate
// - development environment tag
func DefaultTracerConfig(serviceName, serviceVersion string) apptracing.Config {
	return apptracing.Config{
		ServiceName:      serviceName,
		ServiceVersion:   serviceVersion,
		Environment:      "development",
		ExporterType:     "stdout",
		SamplingRate:     1.0,
		Output:           os.Stdout,
	}
}

// ProductionTracerConfig returns a configuration suitable for production use.
// It configures a tracer with:
// - Partial sampling (10% by default)
// - Production environment tag
// Note: You should set the ExporterEndpoint based on your infrastructure
func ProductionTracerConfig(serviceName, serviceVersion string) apptracing.Config {
	return apptracing.Config{
		ServiceName:      serviceName,
		ServiceVersion:   serviceVersion,
		Environment:      "production",
		ExporterType:     "otlp", // Change as needed for your production setup
		SamplingRate:     0.1,    // Sample 10% of traces by default
		ExporterEndpoint: "",     // Set this to your collector endpoint
	}
}
