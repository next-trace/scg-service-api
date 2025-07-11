package logger

import (
	"context"
	"io"
	"os"
	"strings"

	applogger "github.com/hbttundar/scg-service-base/application/logger"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// zerologAdapter implements the logger.Logger interface using the zerolog library.
type zerologAdapter struct {
	log zerolog.Logger
}

// NewZerologAdapter creates a concrete logger adapter.
func NewZerologAdapter(output io.Writer, logLevel string) applogger.Logger {
	if output == nil {
		output = os.Stdout
	}
	level, err := zerolog.ParseLevel(strings.ToLower(logLevel))
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	return &zerologAdapter{
		log: zerolog.New(output).Level(level).With().Timestamp().Logger(),
	}
}

// Helper to add trace context to logs if available.
func (z *zerologAdapter) getEvent(ctx context.Context, level zerolog.Level) *zerolog.Event {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		return z.log.WithLevel(level).
			Str("trace_id", span.SpanContext().TraceID().String()).
			Str("span_id", span.SpanContext().SpanID().String())
	}
	return z.log.WithLevel(level)
}

// Basic logging methods

func (z *zerologAdapter) Debug(ctx context.Context, msg string) {
	z.getEvent(ctx, zerolog.DebugLevel).Msg(msg)
}

func (z *zerologAdapter) Info(ctx context.Context, msg string) {
	z.getEvent(ctx, zerolog.InfoLevel).Msg(msg)
}

func (z *zerologAdapter) Warn(ctx context.Context, msg string) {
	z.getEvent(ctx, zerolog.WarnLevel).Msg(msg)
}

func (z *zerologAdapter) Error(ctx context.Context, err error, msg string) {
	z.getEvent(ctx, zerolog.ErrorLevel).Err(err).Msg(msg)
}

func (z *zerologAdapter) Fatal(ctx context.Context, err error, msg string) {
	z.getEvent(ctx, zerolog.FatalLevel).Err(err).Msg(msg)
}

// Structured logging methods with key-value pairs

func (z *zerologAdapter) DebugKV(ctx context.Context, msg string, keyValues map[string]interface{}) {
	event := z.getEvent(ctx, zerolog.DebugLevel)
	for k, v := range keyValues {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

func (z *zerologAdapter) InfoKV(ctx context.Context, msg string, keyValues map[string]interface{}) {
	event := z.getEvent(ctx, zerolog.InfoLevel)
	for k, v := range keyValues {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

func (z *zerologAdapter) WarnKV(ctx context.Context, msg string, keyValues map[string]interface{}) {
	event := z.getEvent(ctx, zerolog.WarnLevel)
	for k, v := range keyValues {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

func (z *zerologAdapter) ErrorKV(ctx context.Context, err error, msg string, keyValues map[string]interface{}) {
	event := z.getEvent(ctx, zerolog.ErrorLevel).Err(err)
	for k, v := range keyValues {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

func (z *zerologAdapter) FatalKV(ctx context.Context, err error, msg string, keyValues map[string]interface{}) {
	event := z.getEvent(ctx, zerolog.FatalLevel).Err(err)
	for k, v := range keyValues {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

// WithField returns a new logger with the field added to the logger's context
func (z *zerologAdapter) WithField(key string, value interface{}) applogger.Logger {
	return &zerologAdapter{
		log: z.log.With().Interface(key, value).Logger(),
	}
}
