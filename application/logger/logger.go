package logger

import "context"

// Logger defines the abstract logging interface (PORT) for all services.
// It provides a structured, context-aware logging contract.
type Logger interface {
	// Basic logging methods
	Debug(ctx context.Context, msg string)
	Info(ctx context.Context, msg string)
	Warn(ctx context.Context, msg string)
	Error(ctx context.Context, err error, msg string)
	Fatal(ctx context.Context, err error, msg string)

	// Structured logging methods with key-value pairs
	DebugKV(ctx context.Context, msg string, keyValues map[string]interface{})
	InfoKV(ctx context.Context, msg string, keyValues map[string]interface{})
	WarnKV(ctx context.Context, msg string, keyValues map[string]interface{})
	ErrorKV(ctx context.Context, err error, msg string, keyValues map[string]interface{})
	FatalKV(ctx context.Context, err error, msg string, keyValues map[string]interface{})

	// WithField returns a new logger with the field added to the logger's context
	WithField(key string, value interface{}) Logger
}
