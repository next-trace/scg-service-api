package tracing_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	apptracing "github.com/next-trace/scg-service-api/application/tracing"
	impl "github.com/next-trace/scg-service-api/infrastructure/tracing"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type fakeExporter struct{ exported int64 }

func (f *fakeExporter) ExportSpans(_ context.Context, _ []sdktrace.ReadOnlySpan) error {
	atomic.AddInt64(&f.exported, 1)
	return nil
}
func (f *fakeExporter) Shutdown(_ context.Context) error { return nil }

func TestOtelAdapter_Defaults_StartAddEventRecordErrorShutdown(t *testing.T) {
	cfg := apptracing.Config{
		ServiceName:    "svc",
		ServiceVersion: "v1",
		Environment:    "test",
		ExporterType:   "stdout",
		SamplingRate:   1.0,
	}
	tr, err := impl.NewOtelAdapter(cfg)
	if err != nil {
		t.Fatalf("new tracer: %v", err)
	}
	ctx, end := tr.Start(context.Background(), "op")
	tr.AddEvent(ctx, "evt", map[string]string{"k": "v"})
	tr.SetAttributes(ctx, map[string]string{"a": "b"})
	tr.RecordError(ctx, nil) // no-op
	end()
	if err := tr.Shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown: %v", err)
	}
}

func TestOtelAdapter_WithExporter_ExporterReceivesSpansOnShutdown(t *testing.T) {
	fe := &fakeExporter{}
	cfg := apptracing.Config{ServiceName: "svc2", SamplingRate: 1.0}
	// Provide custom exporter via option; keep resource/sampler defaults
	tr, err := impl.NewOtelAdapterWithOptions(cfg, impl.WithExporter(fe), impl.WithResource(resource.Empty()))
	if err != nil {
		t.Fatalf("new tracer with options: %v", err)
	}
	// create and end a span so the batcher has something to flush
	ctx, end := tr.Start(context.Background(), "op2")
	_ = ctx
	end()
	// Flush and trigger export
	if err := tr.Shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown: %v", err)
	}
	// ExportSpans should have been called at least once after shutdown
	if atomic.LoadInt64(&fe.exported) == 0 {
		// Allow some extra time in case of scheduling delays (should be flushed on Shutdown)
		time.Sleep(10 * time.Millisecond)
		if atomic.LoadInt64(&fe.exported) == 0 {
			t.Fatalf("expected exporter to be invoked at least once")
		}
	}
}
