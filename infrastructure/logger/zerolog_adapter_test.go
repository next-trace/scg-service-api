package logger_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	applogger "github.com/hbttundar/scg-service-base/application/logger"
	"github.com/hbttundar/scg-service-base/infrastructure/logger"
	"github.com/stretchr/testify/assert"
)

func TestNewZerologAdapter(t *testing.T) {
	t.Run("Default output", func(t *testing.T) {
		// When output is nil, it should default to os.Stdout
		log := logger.NewZerologAdapter(nil, "info")
		assert.NotNil(t, log)
	})

	t.Run("Custom output", func(t *testing.T) {
		var buf bytes.Buffer
		log := logger.NewZerologAdapter(&buf, "info")
		assert.NotNil(t, log)
	})

	t.Run("Invalid log level", func(t *testing.T) {
		// Should default to info level
		var buf bytes.Buffer
		log := logger.NewZerologAdapter(&buf, "invalid")
		assert.NotNil(t, log)
	})
}

func TestLoggingMethods(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewZerologAdapter(&buf, "debug")
	ctx := context.Background()

	t.Run("Debug", func(t *testing.T) {
		buf.Reset()
		log.Debug(ctx, "debug message")
		assert.Contains(t, buf.String(), "debug message")
		assert.Contains(t, buf.String(), `"level":"debug"`)
	})

	t.Run("Info", func(t *testing.T) {
		buf.Reset()
		log.Info(ctx, "info message")
		assert.Contains(t, buf.String(), "info message")
		assert.Contains(t, buf.String(), `"level":"info"`)
	})

	t.Run("Warn", func(t *testing.T) {
		buf.Reset()
		log.Warn(ctx, "warn message")
		assert.Contains(t, buf.String(), "warn message")
		assert.Contains(t, buf.String(), `"level":"warn"`)
	})

	t.Run("Error", func(t *testing.T) {
		buf.Reset()
		err := errors.New("test error")
		log.Error(ctx, err, "error message")
		assert.Contains(t, buf.String(), "error message")
		assert.Contains(t, buf.String(), `"level":"error"`)
		assert.Contains(t, buf.String(), `"error":"test error"`)
	})

	// We don't test Fatal because it would exit the program
}

func TestStructuredLoggingMethods(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewZerologAdapter(&buf, "debug")
	ctx := context.Background()
	keyValues := map[string]interface{}{
		"user_id": 123,
		"action":  "login",
	}

	t.Run("DebugKV", func(t *testing.T) {
		buf.Reset()
		log.DebugKV(ctx, "debug message", keyValues)
		assert.Contains(t, buf.String(), "debug message")
		assert.Contains(t, buf.String(), `"level":"debug"`)
		assert.Contains(t, buf.String(), `"user_id":123`)
		assert.Contains(t, buf.String(), `"action":"login"`)
	})

	t.Run("InfoKV", func(t *testing.T) {
		buf.Reset()
		log.InfoKV(ctx, "info message", keyValues)
		assert.Contains(t, buf.String(), "info message")
		assert.Contains(t, buf.String(), `"level":"info"`)
		assert.Contains(t, buf.String(), `"user_id":123`)
		assert.Contains(t, buf.String(), `"action":"login"`)
	})

	t.Run("WarnKV", func(t *testing.T) {
		buf.Reset()
		log.WarnKV(ctx, "warn message", keyValues)
		assert.Contains(t, buf.String(), "warn message")
		assert.Contains(t, buf.String(), `"level":"warn"`)
		assert.Contains(t, buf.String(), `"user_id":123`)
		assert.Contains(t, buf.String(), `"action":"login"`)
	})

	t.Run("ErrorKV", func(t *testing.T) {
		buf.Reset()
		err := errors.New("test error")
		log.ErrorKV(ctx, err, "error message", keyValues)
		assert.Contains(t, buf.String(), "error message")
		assert.Contains(t, buf.String(), `"level":"error"`)
		assert.Contains(t, buf.String(), `"error":"test error"`)
		assert.Contains(t, buf.String(), `"user_id":123`)
		assert.Contains(t, buf.String(), `"action":"login"`)
	})

	// We don't test FatalKV because it would exit the program
}

func TestWithField(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewZerologAdapter(&buf, "info")
	ctx := context.Background()

	// Create a new logger with a field
	logWithField := log.WithField("request_id", "abc123")
	assert.NotNil(t, logWithField)
	assert.Implements(t, (*applogger.Logger)(nil), logWithField)

	// Test that the field is included in log messages
	buf.Reset()
	logWithField.Info(ctx, "test message")
	assert.Contains(t, buf.String(), "test message")
	assert.Contains(t, buf.String(), `"request_id":"abc123"`)

	// Test that the field is included in structured log messages
	buf.Reset()
	logWithField.InfoKV(ctx, "test message", map[string]interface{}{"user_id": 123})
	assert.Contains(t, buf.String(), "test message")
	assert.Contains(t, buf.String(), `"request_id":"abc123"`)
	assert.Contains(t, buf.String(), `"user_id":123`)

	// Test chaining WithField calls
	buf.Reset()
	logWithTwoFields := logWithField.WithField("session_id", "xyz789")
	logWithTwoFields.Info(ctx, "test message")
	assert.Contains(t, buf.String(), "test message")
	assert.Contains(t, buf.String(), `"request_id":"abc123"`)
	assert.Contains(t, buf.String(), `"session_id":"xyz789"`)
}
