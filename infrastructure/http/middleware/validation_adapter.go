// Package middleware provides HTTP middleware components.
package middleware

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"

	appvalidation "github.com/hbttundar/scg-service-base/application/validation"
	applogger "github.com/hbttundar/scg-service-base/application/logger"
)

// ValidationMiddleware provides middleware to validate request data.
type ValidationMiddleware struct {
	validator appvalidation.Validator
	config    appvalidation.Config
	log       applogger.Logger
}

// NewValidationMiddleware creates a new validation middleware.
func NewValidationMiddleware(validator appvalidation.Validator, config appvalidation.Config, log applogger.Logger) *ValidationMiddleware {
	return &ValidationMiddleware{
		validator: validator,
		config:    config,
		log:       log,
	}
}

// Middleware returns an http.Handler middleware function.
func (vm *ValidationMiddleware) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !vm.config.Enabled {
				next.ServeHTTP(w, r)
				return
			}

			// Skip validation for certain methods
			if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}

			// Get the validation model from the request context
			model := r.Context().Value("validation_model")
			if model == nil {
				// No validation model, skip validation
				next.ServeHTTP(w, r)
				return
			}

			// Parse the request body
			body, err := io.ReadAll(r.Body)
			if err != nil {
				vm.log.Error(r.Context(), err, "failed to read request body")
				http.Error(w, "Bad Request: failed to read request body", http.StatusBadRequest)
				return
			}

			// Restore the request body for later use
			r.Body = io.NopCloser(strings.NewReader(string(body)))

			// Create a new instance of the model
			modelType := reflect.TypeOf(model)
			if modelType.Kind() == reflect.Ptr {
				modelType = modelType.Elem()
			}
			modelValue := reflect.New(modelType).Interface()

			// Unmarshal the request body into the model
			if err := json.Unmarshal(body, modelValue); err != nil {
				vm.log.Error(r.Context(), err, "failed to unmarshal request body")
				http.Error(w, "Bad Request: invalid JSON", http.StatusBadRequest)
				return
			}

			// Validate the model
			result := vm.validator.Validate(r.Context(), modelValue)
			if !result.Valid {
				// Return validation errors
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"error":  "Validation failed",
					"errors": result.Errors,
				})
				return
			}

			// Store the validated model in the request context
			ctx := context.WithValue(r.Context(), "validated_model", modelValue)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Validate provides a middleware that validates a specific model.
func (vm *ValidationMiddleware) Validate(model interface{}) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !vm.config.Enabled {
				next.ServeHTTP(w, r)
				return
			}

			// Skip validation for certain methods
			if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}

			// Parse the request body
			body, err := io.ReadAll(r.Body)
			if err != nil {
				vm.log.Error(r.Context(), err, "failed to read request body")
				http.Error(w, "Bad Request: failed to read request body", http.StatusBadRequest)
				return
			}

			// Restore the request body for later use
			r.Body = io.NopCloser(strings.NewReader(string(body)))

			// Create a new instance of the model
			modelType := reflect.TypeOf(model)
			if modelType.Kind() == reflect.Ptr {
				modelType = modelType.Elem()
			}
			modelValue := reflect.New(modelType).Interface()

			// Unmarshal the request body into the model
			if err := json.Unmarshal(body, modelValue); err != nil {
				vm.log.Error(r.Context(), err, "failed to unmarshal request body")
				http.Error(w, "Bad Request: invalid JSON", http.StatusBadRequest)
				return
			}

			// Validate the model
			result := vm.validator.Validate(r.Context(), modelValue)
			if !result.Valid {
				// Return validation errors
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"error":  "Validation failed",
					"errors": result.Errors,
				})
				return
			}

			// Store the validated model in the request context
			ctx := context.WithValue(r.Context(), "validated_model", modelValue)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Validation provides backward compatibility with the old API.
// Deprecated: Use NewValidationMiddleware instead.
func Validation(validator appvalidation.Validator, config appvalidation.Config, log applogger.Logger) func(http.Handler) http.Handler {
	return NewValidationMiddleware(validator, config, log).Middleware()
}
