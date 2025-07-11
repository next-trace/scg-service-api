package middleware

import (
	"net/http"
	"runtime/debug"

	applogger "github.com/hbttundar/scg-service-base/application/logger"
)

// RecoveryMiddleware provides middleware to recover from panics and prevent server crashes.
type RecoveryMiddleware struct {
	log applogger.Logger
}

// NewRecoveryMiddleware creates a new recovery middleware.
func NewRecoveryMiddleware(log applogger.Logger) *RecoveryMiddleware {
	return &RecoveryMiddleware{
		log: log,
	}
}

// Middleware returns an http.Handler middleware function.
func (rm *RecoveryMiddleware) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					rm.log.ErrorKV(r.Context(), nil, "panic recovered", map[string]interface{}{
						"stack": string(debug.Stack()),
						"error": err,
					})
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// Recovery provides backward compatibility with the old API.
// Deprecated: Use NewRecoveryMiddleware instead.
func Recovery(log applogger.Logger) func(http.Handler) http.Handler {
	return NewRecoveryMiddleware(log).Middleware()
}
