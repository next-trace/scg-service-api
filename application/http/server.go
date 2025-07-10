package http

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	applogger "github.com/next-trace/scg-service-api/application/logger"
)

// Run starts the given http.Server and performs a graceful shutdown on SIGINT/SIGTERM.
//
// Behavior:
// - Starts srv.ListenAndServe() in a goroutine.
// - Listens for OS signals (os.Interrupt, syscall.SIGTERM) and context cancellation.
// - When a shutdown trigger occurs, logs a message and calls srv.Shutdown with a 30s timeout.
// - Returns the first non-nil error from ListenAndServe (other than http.ErrServerClosed) or from Shutdown.
func Run(ctx context.Context, srv *http.Server, log applogger.Logger) error {
	if srv == nil {
		return nil
	}

	// Log server start if address is known
	if srv.Addr != "" {
		log.InfoKV(ctx, "starting HTTP server", map[string]interface{}{"address": srv.Addr})
	} else {
		log.Info(ctx, "starting HTTP server")
	}

	errCh := make(chan error, 1)

	// Start the HTTP server
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		close(errCh)
	}()

	// Set up signal handling for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(stop)

	select {
	case <-ctx.Done():
		// External context canceled; proceed to graceful shutdown
		log.Info(ctx, "shutdown signal received (context canceled), shutting down gracefully")
	case <-stop:
		// OS termination signal received
		log.Info(ctx, "shutdown signal received, shutting down gracefully")
	case err := <-errCh:
		// Server failed to start or crashed
		if err != nil {
			return err
		}
		// Channel closed without error (server closed), just return
		return nil
	}

	// Perform graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		// Log and return shutdown error
		log.Error(ctx, err, "http server shutdown error")
		return err
	}

	log.Info(ctx, "http server shutdown complete")
	return nil
}
