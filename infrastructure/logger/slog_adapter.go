package logger

import (
	"context"
	"io"
	"log/slog"
	"os"

	applogger "github.com/next-trace/scg-service-api/application/logger"
	"go.opentelemetry.io/otel/trace"
)

// slogAdapter implements the logger.Logger interface using Go's slog.
type slogAdapter struct {
	log *slog.Logger
}

// NewSlogAdapter creates a concrete logger adapter.
// If output is nil, it defaults to os.Stdout. Level is one of: debug, info, warn, error.
func NewSlogAdapter(output io.Writer, level string) applogger.Logger {
	if output == nil {
		output = os.Stdout
	}
	// Use internal logger with provided writer; Pretty=false by default for JSON output
	h := slog.NewJSONHandler(output, &slog.HandlerOptions{Level: internallogLevel(level)})
	l := slog.New(h)
	return &slogAdapter{log: l}
}

func internallogLevel(level string) slog.Leveler { // helper to avoid import cycle with internal/logger
	switch level {
	case "debug", "DEBUG":
		return slog.LevelDebug
	case "warn", "warning", "WARN", "WARNING":
		return slog.LevelWarn
	case "error", "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// withTrace returns logger with OTEL trace/span IDs if present in ctx
func (s *slogAdapter) withTrace(ctx context.Context) *slog.Logger {
	span := trace.SpanFromContext(ctx)
	sc := span.SpanContext()
	if sc.IsValid() {
		return s.log.With(
			slog.String("trace_id", sc.TraceID().String()),
			slog.String("span_id", sc.SpanID().String()),
		)
	}
	return s.log
}

// Basic logging methods
func (s *slogAdapter) Debug(ctx context.Context, msg string) {
	s.withTrace(ctx).DebugContext(ctx, msg)
}

func (s *slogAdapter) Info(ctx context.Context, msg string) {
	s.withTrace(ctx).InfoContext(ctx, msg)
}

func (s *slogAdapter) Warn(ctx context.Context, msg string) {
	s.withTrace(ctx).WarnContext(ctx, msg)
}

func (s *slogAdapter) Error(ctx context.Context, err error, msg string) {
	s.withTrace(ctx).ErrorContext(ctx, msg, slog.Any("error", err))
}

func (s *slogAdapter) Fatal(ctx context.Context, err error, msg string) {
	// slog has no Fatal; we log at Error level and then exit with non-zero code for compatibility
	s.withTrace(ctx).ErrorContext(ctx, msg, slog.Any("error", err), slog.String("severity", "FATAL"))
	os.Exit(1)
}

// Structured logging methods with key-value pairs
func (s *slogAdapter) DebugKV(ctx context.Context, msg string, keyValues map[string]interface{}) {
	attrs := make([]any, 0, len(keyValues))
	for k, v := range keyValues {
		attrs = append(attrs, slog.Any(k, v))
	}
	s.withTrace(ctx).DebugContext(ctx, msg, attrs...)
}

func (s *slogAdapter) InfoKV(ctx context.Context, msg string, keyValues map[string]interface{}) {
	attrs := make([]any, 0, len(keyValues))
	for k, v := range keyValues {
		attrs = append(attrs, slog.Any(k, v))
	}
	s.withTrace(ctx).InfoContext(ctx, msg, attrs...)
}

func (s *slogAdapter) WarnKV(ctx context.Context, msg string, keyValues map[string]interface{}) {
	attrs := make([]any, 0, len(keyValues))
	for k, v := range keyValues {
		attrs = append(attrs, slog.Any(k, v))
	}
	s.withTrace(ctx).WarnContext(ctx, msg, attrs...)
}

func (s *slogAdapter) ErrorKV(ctx context.Context, err error, msg string, keyValues map[string]interface{}) {
	attrs := make([]any, 0, len(keyValues)+1)
	attrs = append(attrs, slog.Any("error", err))
	for k, v := range keyValues {
		attrs = append(attrs, slog.Any(k, v))
	}
	s.withTrace(ctx).ErrorContext(ctx, msg, attrs...)
}

func (s *slogAdapter) FatalKV(ctx context.Context, err error, msg string, keyValues map[string]interface{}) {
	attrs := make([]any, 0, len(keyValues)+2)
	attrs = append(attrs, slog.Any("error", err), slog.String("severity", "FATAL"))
	for k, v := range keyValues {
		attrs = append(attrs, slog.Any(k, v))
	}
	s.withTrace(ctx).ErrorContext(ctx, msg, attrs...)
	os.Exit(1)
}

// WithField returns a new logger with the field added to the logger's context
func (s *slogAdapter) WithField(key string, value interface{}) applogger.Logger {
	return &slogAdapter{log: s.log.With(slog.Any(key, value))}
}
