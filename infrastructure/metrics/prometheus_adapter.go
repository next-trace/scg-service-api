// Package metrics provides metrics collection functionality.
//
// Note: This package requires the following dependencies:
// - github.com/prometheus/client_golang
//
// See docs/dependencies.md for more information.
package metrics

import (
	"context"
	"net/http"
	"sync"
	"time"

	appmetrics "github.com/hbttundar/scg-service-base/application/metrics"
	applogger "github.com/hbttundar/scg-service-base/application/logger"
)

// prometheusAdapter implements the metrics.Metrics interface using Prometheus.
// In a real implementation, this would use the Prometheus client library.
// For now, we'll provide a simple implementation that can be replaced later.
type prometheusAdapter struct {
	config     appmetrics.Config
	log        applogger.Logger
	counters   map[string]float64
	gauges     map[string]float64
	histograms map[string][]float64
	labels     map[string]string
	server     *http.Server
	mu         sync.RWMutex
}

// NewPrometheusAdapter creates a new Prometheus metrics adapter.
func NewPrometheusAdapter(config appmetrics.Config, log applogger.Logger) appmetrics.Metrics {
	return &prometheusAdapter{
		config:     config,
		log:        log,
		counters:   make(map[string]float64),
		gauges:     make(map[string]float64),
		histograms: make(map[string][]float64),
		labels:     config.Labels,
	}
}

// WithLabels returns a new Metrics instance with the given labels.
func (p *prometheusAdapter) WithLabels(labels map[string]string) appmetrics.Metrics {
	// Create a new adapter with the same configuration
	newAdapter := &prometheusAdapter{
		config:     p.config,
		log:        p.log,
		counters:   p.counters,
		gauges:     p.gauges,
		histograms: p.histograms,
		labels:     make(map[string]string),
	}

	// Copy existing labels
	for k, v := range p.labels {
		newAdapter.labels[k] = v
	}

	// Add new labels
	for k, v := range labels {
		newAdapter.labels[k] = v
	}

	return newAdapter
}

// Serve starts the metrics server on the given address.
func (p *prometheusAdapter) Serve(ctx context.Context, addr string) error {
	// In a real implementation, this would start a Prometheus HTTP server
	// that exposes metrics on the /metrics endpoint.
	p.log.InfoKV(ctx, "starting metrics server", map[string]interface{}{
		"address": addr,
	})

	// Create a new HTTP server
	p.server = &http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// In a real implementation, this would use the Prometheus handler
			// to expose metrics in the Prometheus format.
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte("# HELP example_metric Example metric\n"))
			w.Write([]byte("# TYPE example_metric gauge\n"))
			w.Write([]byte("example_metric 42\n"))
		}),
	}

	// Start the server in a goroutine
	go func() {
		if err := p.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			p.log.Error(ctx, err, "metrics server error")
		}
	}()

	return nil
}

// Shutdown gracefully shuts down the metrics server.
func (p *prometheusAdapter) Shutdown(ctx context.Context) error {
	if p.server != nil {
		p.log.Info(ctx, "shutting down metrics server")
		return p.server.Shutdown(ctx)
	}
	return nil
}

// Counter methods

// CounterInc increments the counter by 1.
func (p *prometheusAdapter) CounterInc(name string) {
	p.CounterAdd(name, 1)
}

// CounterAdd adds the given value to the counter.
func (p *prometheusAdapter) CounterAdd(name string, value float64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// In a real implementation, this would use the Prometheus Counter.
	p.counters[name] += value
}

// Gauge methods

// GaugeSet sets the gauge to the given value.
func (p *prometheusAdapter) GaugeSet(name string, value float64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// In a real implementation, this would use the Prometheus Gauge.
	p.gauges[name] = value
}

// GaugeInc increments the gauge by 1.
func (p *prometheusAdapter) GaugeInc(name string) {
	p.GaugeAdd(name, 1)
}

// GaugeDec decrements the gauge by 1.
func (p *prometheusAdapter) GaugeDec(name string) {
	p.GaugeSub(name, 1)
}

// GaugeAdd adds the given value to the gauge.
func (p *prometheusAdapter) GaugeAdd(name string, value float64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// In a real implementation, this would use the Prometheus Gauge.
	p.gauges[name] += value
}

// GaugeSub subtracts the given value from the gauge.
func (p *prometheusAdapter) GaugeSub(name string, value float64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// In a real implementation, this would use the Prometheus Gauge.
	p.gauges[name] -= value
}

// Histogram methods

// HistogramObserve adds a single observation to the histogram.
func (p *prometheusAdapter) HistogramObserve(name string, value float64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// In a real implementation, this would use the Prometheus Histogram.
	p.histograms[name] = append(p.histograms[name], value)
}

// Timer methods

// TimerObserveDuration measures the duration of the given function call.
func (p *prometheusAdapter) TimerObserveDuration(name string, f func()) {
	start := time.Now()
	f()
	duration := time.Since(start)
	p.HistogramObserve(name, duration.Seconds())
}

// TimerStart starts a new timer and returns a function to stop it.
func (p *prometheusAdapter) TimerStart(name string) func() time.Duration {
	start := time.Now()
	return func() time.Duration {
		duration := time.Since(start)
		p.HistogramObserve(name, duration.Seconds())
		return duration
	}
}
