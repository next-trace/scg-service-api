// Package middleware provides HTTP middleware components.
package middleware

import (
	"net/http"
	"strings"

	applogger "github.com/next-trace/scg-service-api/application/logger"
	appratelimit "github.com/next-trace/scg-service-api/application/ratelimit"
)

// RateLimitMiddleware provides middleware to limit the rate of requests.
type RateLimitMiddleware struct {
	limiter appratelimit.Limiter
	config  appratelimit.Config
	log     applogger.Logger
}

// NewRateLimitMiddleware creates a new rate limit middleware.
func NewRateLimitMiddleware(limiter appratelimit.Limiter, config appratelimit.Config, log applogger.Logger) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		limiter: limiter,
		config:  config,
		log:     log,
	}
}

// Middleware returns an http.Handler middleware function.
func (rl *RateLimitMiddleware) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !rl.config.Enabled {
				next.ServeHTTP(w, r)
				return
			}

			// Get the key for rate limiting
			key := rl.getKey(r)

			// Check if the request is allowed
			if !rl.limiter.Allow(r.Context(), key) {
				rl.log.WarnKV(r.Context(), "rate limit exceeded", map[string]interface{}{
					"key":    key,
					"path":   r.URL.Path,
					"method": r.Method,
				})
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}

// getKey returns a key for rate limiting based on the request.
// If the config has a KeyFunc, it will be used. Otherwise, a default key will be generated.
func (rl *RateLimitMiddleware) getKey(r *http.Request) string {
	if rl.config.KeyFunc != nil {
		return rl.config.KeyFunc(r.Context())
	}

	// Default key is based on the client's IP address
	ip := getClientIP(r)
	return "ip:" + ip
}

// getClientIP returns the client's IP address from the request.
func getClientIP(r *http.Request) string {
	// Check for X-Forwarded-For header
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		// X-Forwarded-For can contain multiple IPs, use the first one
		ips := strings.Split(xForwardedFor, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check for X-Real-IP header
	xRealIP := r.Header.Get("X-Real-IP")
	if xRealIP != "" {
		return xRealIP
	}

	// Use RemoteAddr as a fallback
	ip := r.RemoteAddr
	// Remove port if present
	if i := strings.LastIndex(ip, ":"); i != -1 {
		ip = ip[:i]
	}
	return ip
}

// RateLimit provides backward compatibility with the old API.
// Deprecated: Use NewRateLimitMiddleware instead.
func RateLimit(limiter appratelimit.Limiter, config appratelimit.Config, log applogger.Logger) func(http.Handler) http.Handler {
	return NewRateLimitMiddleware(limiter, config, log).Middleware()
}

// WaitRateLimitMiddleware provides middleware that waits for a token instead of rejecting the request.
// This is useful for internal services where you want to throttle but not reject requests.
type WaitRateLimitMiddleware struct {
	limiter appratelimit.Limiter
	config  appratelimit.Config
	log     applogger.Logger
}

// NewWaitRateLimitMiddleware creates a new wait rate limit middleware.
func NewWaitRateLimitMiddleware(limiter appratelimit.Limiter, config appratelimit.Config, log applogger.Logger) *WaitRateLimitMiddleware {
	return &WaitRateLimitMiddleware{
		limiter: limiter,
		config:  config,
		log:     log,
	}
}

// Middleware returns an http.Handler middleware function.
func (wrl *WaitRateLimitMiddleware) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !wrl.config.Enabled {
				next.ServeHTTP(w, r)
				return
			}

			// Get the key for rate limiting
			key := wrl.getKey(r)

			// Wait for a token
			if err := wrl.limiter.Wait(r.Context(), key); err != nil {
				wrl.log.WarnKV(r.Context(), "rate limit wait failed", map[string]interface{}{
					"key":    key,
					"path":   r.URL.Path,
					"method": r.Method,
					"error":  err.Error(),
				})
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}

// getKey returns a key for rate limiting based on the request.
// If the config has a KeyFunc, it will be used. Otherwise, a default key will be generated.
func (wrl *WaitRateLimitMiddleware) getKey(r *http.Request) string {
	if wrl.config.KeyFunc != nil {
		return wrl.config.KeyFunc(r.Context())
	}

	// Default key is based on the client's IP address
	ip := getClientIP(r)
	return "ip:" + ip
}

// WaitRateLimit provides backward compatibility with the old API.
// Deprecated: Use NewWaitRateLimitMiddleware instead.
func WaitRateLimit(limiter appratelimit.Limiter, config appratelimit.Config, log applogger.Logger) func(http.Handler) http.Handler {
	return NewWaitRateLimitMiddleware(limiter, config, log).Middleware()
}
