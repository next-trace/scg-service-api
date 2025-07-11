// Package metrics defines the abstract interface (PORT) for metrics collection.
package metrics

import (
	"context"
	"time"
)

// Metrics defines the interface for collecting metrics.
type Metrics interface {
	// Counter operations

	// CounterInc increments the counter by 1.
	CounterInc(name string)

	// CounterAdd adds the given value to the counter.
	CounterAdd(name string, value float64)

	// Gauge operations

	// GaugeSet sets the gauge to the given value.
	GaugeSet(name string, value float64)

	// GaugeInc increments the gauge by 1.
	GaugeInc(name string)

	// GaugeDec decrements the gauge by 1.
	GaugeDec(name string)

	// GaugeAdd adds the given value to the gauge.
	GaugeAdd(name string, value float64)

	// GaugeSub subtracts the given value from the gauge.
	GaugeSub(name string, value float64)

	// Histogram operations

	// HistogramObserve adds a single observation to the histogram.
	HistogramObserve(name string, value float64)

	// Timer operations

	// TimerObserveDuration measures the duration of the given function call.
	TimerObserveDuration(name string, f func())

	// TimerStart starts a new timer and returns a function to stop it.
	TimerStart(name string) func() time.Duration

	// WithLabels returns a new Metrics instance with the given labels.
	WithLabels(labels map[string]string) Metrics

	// Serve starts the metrics server on the given address.
	Serve(ctx context.Context, addr string) error

	// Shutdown gracefully shuts down the metrics server.
	Shutdown(ctx context.Context) error
}

// Since we've simplified the interface, we no longer need the TimerInstance struct.
// Instead, we use the TimerStart method which returns a function to stop the timer.

// Config holds configuration for metrics.
type Config struct {
	// Namespace is the metrics namespace.
	Namespace string

	// Subsystem is the metrics subsystem.
	Subsystem string

	// Labels are the default labels to add to all metrics.
	Labels map[string]string

	// EnableGoMetrics enables Go runtime metrics.
	EnableGoMetrics bool

	// EnableProcessMetrics enables process metrics.
	EnableProcessMetrics bool
}

// DefaultConfig returns the default configuration for metrics.
func DefaultConfig() Config {
	return Config{
		Namespace:           "app",
		Subsystem:           "",
		Labels:              make(map[string]string),
		EnableGoMetrics:     true,
		EnableProcessMetrics: true,
	}
}
