package tracing_test

import (
	"os"
	"testing"

	apptracing "github.com/next-trace/scg-service-api/application/tracing"
	"github.com/next-trace/scg-service-api/infrastructure/tracing"
	"github.com/stretchr/testify/assert"
)

func TestDefaultTracerConfig(t *testing.T) {
	// Test default configuration
	cfg := tracing.DefaultTracerConfig("test-service", "1.0.0")

	assert.Equal(t, "test-service", cfg.ServiceName)
	assert.Equal(t, "1.0.0", cfg.ServiceVersion)
	assert.Equal(t, "development", cfg.Environment)
	assert.Equal(t, "stdout", cfg.ExporterType)
	assert.Equal(t, 1.0, cfg.SamplingRate)
	assert.Equal(t, os.Stdout, cfg.Output)
}

func TestProductionTracerConfig(t *testing.T) {
	// Test production configuration
	cfg := tracing.ProductionTracerConfig("test-service", "1.0.0")

	assert.Equal(t, "test-service", cfg.ServiceName)
	assert.Equal(t, "1.0.0", cfg.ServiceVersion)
	assert.Equal(t, "production", cfg.Environment)
	assert.Equal(t, "otlp", cfg.ExporterType)
	assert.Equal(t, 0.1, cfg.SamplingRate)
	assert.Empty(t, cfg.ExporterEndpoint) // Should be empty by default
}

func TestNewTracer(t *testing.T) {
	// Test creating a tracer with default config
	cfg := apptracing.Config{
		ServiceName:    "test-service",
		ServiceVersion: "1.0.0",
		Environment:    "test",
		ExporterType:   "stdout",
		SamplingRate:   1.0,
		Output:         os.Stdout,
	}

	tracer, err := tracing.NewTracer(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, tracer)
	assert.Implements(t, (*apptracing.Tracer)(nil), tracer)

	// Shutdown the tracer to clean up resources
	if tracer != nil {
		err = tracer.Shutdown(t.Context())
		assert.NoError(t, err)
	}
}
