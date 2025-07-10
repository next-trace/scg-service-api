package middleware_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	appmetrics "github.com/next-trace/scg-service-api/application/metrics"
	"github.com/next-trace/scg-service-api/infrastructure/http/middleware"
	"github.com/stretchr/testify/assert"
)

// fakeMetrics is a simple in-memory implementation of appmetrics.Metrics used for testing
// It records calls so we can make assertions without needing a real metrics backend.
type fakeMetrics struct {
	labels             map[string]string
	counters           map[string]float64
	gauges             map[string]float64
	histograms         map[string][]float64
	timersStartedNames []string
}

func newFakeMetrics() *fakeMetrics {
	return &fakeMetrics{
		labels:     map[string]string{},
		counters:   map[string]float64{},
		gauges:     map[string]float64{},
		histograms: map[string][]float64{},
	}
}

func (f *fakeMetrics) CounterInc(name string)                { f.counters[name]++ }
func (f *fakeMetrics) CounterAdd(name string, value float64) { f.counters[name] += value }
func (f *fakeMetrics) GaugeSet(name string, value float64)   { f.gauges[name] = value }
func (f *fakeMetrics) GaugeInc(name string)                  { f.gauges[name]++ }
func (f *fakeMetrics) GaugeDec(name string)                  { f.gauges[name]-- }
func (f *fakeMetrics) GaugeAdd(name string, value float64)   { f.gauges[name] += value }
func (f *fakeMetrics) GaugeSub(name string, value float64)   { f.gauges[name] -= value }
func (f *fakeMetrics) HistogramObserve(name string, value float64) {
	f.histograms[name] = append(f.histograms[name], value)
}

func (f *fakeMetrics) TimerObserveDuration(_ string, fn func()) {
	start := time.Now()
	fn()
	_ = time.Since(start)
}

func (f *fakeMetrics) TimerStart(name string) func() time.Duration {
	f.timersStartedNames = append(f.timersStartedNames, name)
	start := time.Now()
	return func() time.Duration { return time.Since(start) }
}

func (f *fakeMetrics) WithLabels(labels map[string]string) appmetrics.Metrics {
	nm := newFakeMetrics()
	nm.labels = labels
	return nm
}
func (f *fakeMetrics) Serve(_ context.Context, _ string) error { return nil }

func (f *fakeMetrics) Shutdown(_ context.Context) error { return nil }

func TestMetricsMiddleware_NormalRequest(t *testing.T) {
	fm := newFakeMetrics()
	mw := middleware.NewMetricsMiddleware(fm)

	// Handler writes a small JSON body and returns 201
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("{}"))
	})

	wrapped := mw.Middleware()(h)

	req := httptest.NewRequest(http.MethodPost, "/widgets", bytes.NewBufferString("payload"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	wrapped.ServeHTTP(w, req)
	resp := w.Result()
	t.Cleanup(func() { _ = resp.Body.Close() })

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Verify counters incremented
	assert.GreaterOrEqual(t, fm.counters["http_requests_total"], 1.0)
	// Verify response size recorded
	sizes := fm.histograms["http_response_size_bytes"]
	if assert.Len(t, sizes, 1) {
		assert.Greater(t, sizes[0], 0.0)
	}
	// Verify request size recorded
	reqSizes := fm.histograms["http_request_size_bytes"]
	if assert.Len(t, reqSizes, 1) {
		assert.Equal(t, float64(len("payload")), reqSizes[0])
	}
	// Verify we at least started the duration timer
	assert.Contains(t, fm.timersStartedNames, "http_request_duration_seconds")
}

func TestMetricsMiddleware_SlowRequest(t *testing.T) {
	fm := newFakeMetrics()
	mw := middleware.NewMetricsMiddleware(fm)

	// Slow handler intentionally sleeps > 1s to trigger slow request counter
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})

	wrapped := mw.Middleware()(h)

	req := httptest.NewRequest(http.MethodGet, "/slow", nil)
	w := httptest.NewRecorder()

	wrapped.ServeHTTP(w, req)
	resp := w.Result()
	t.Cleanup(func() { resp.Body.Close() })

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	// Should have incremented slow requests
	assert.Equal(t, 1.0, fm.counters["http_slow_requests_total"])
}
