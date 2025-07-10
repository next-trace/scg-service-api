package metrics_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	appmetrics "github.com/next-trace/scg-service-api/application/metrics"
	infraLogger "github.com/next-trace/scg-service-api/infrastructure/logger"
	metricsimpl "github.com/next-trace/scg-service-api/infrastructure/metrics"
)

func TestPrometheusAdapter_BasicCalls(t *testing.T) {
	ctx := context.Background()
	var buf bytes.Buffer
	log := infraLogger.NewSlogAdapter(&buf, "info")

	cfg := appmetrics.DefaultConfig()
	m := metricsimpl.NewPrometheusAdapter(cfg, log)

	m.CounterInc("c1")
	m.CounterAdd("c1", 3)
	m.GaugeSet("g1", 1)
	m.GaugeInc("g1")
	m.GaugeDec("g1")
	m.GaugeAdd("g1", 2)
	m.GaugeSub("g1", 1)
	m.HistogramObserve("h1", 123)
	m.TimerObserveDuration("op", func() { time.Sleep(time.Millisecond) })
	stop := m.TimerStart("op2")
	_ = stop()

	with := m.WithLabels(map[string]string{"a": "b"})
	with.CounterInc("c2")

	// We won't actually Serve a HTTP server in tests to avoid flakes, but Shutdown should be safe.
	if err := m.Shutdown(ctx); err != nil {
		t.Fatalf("shutdown error: %v", err)
	}
}
