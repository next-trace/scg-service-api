package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

// Config holds logger configuration.
// If left zero-valued, Init returns a no-op slog.Logger writing to io.Discard.
type Config struct {
	Service    string // service name to include as attribute
	Level      string // debug, info, warn, error
	Pretty     bool   // use text handler when true
	WithCaller bool   // include source information
}

// Init initializes and returns a *slog.Logger based on the provided config.
// - JSON handler for production (Pretty=false)
// - Text handler for development (Pretty=true)
// - WithCaller via HandlerOptions{AddSource:true}
// - Fallback no-op logger with io.Discard when config is empty
func Init(cfg Config) *slog.Logger {
	// Detect empty config: all fields zero values
	if cfg == (Config{}) {
		return slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{}))
	}

	opts := &slog.HandlerOptions{
		Level:     parseLevel(cfg.Level),
		AddSource: cfg.WithCaller,
	}

	var h slog.Handler
	var w io.Writer = os.Stdout
	if cfg.Pretty {
		h = slog.NewTextHandler(w, opts)
	} else {
		h = slog.NewJSONHandler(w, opts)
	}

	logger := slog.New(h)
	if cfg.Service != "" {
		logger = logger.With(slog.String("service", cfg.Service))
	}
	return logger
}

func parseLevel(level string) slog.Leveler {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	case "info", "":
		fallthrough
	default:
		return slog.LevelInfo
	}
}
