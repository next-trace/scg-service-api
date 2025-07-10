// Package logger defines a minimal, structured, context-aware logging port.
// Implementations live in infrastructure/logger. The port keeps consumers stable
// while allowing library swaps (e.g., slog, zap, zerolog).
package logger
