// Package middleware provides HTTP middleware components.
package middleware

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"

	applogger "github.com/next-trace/scg-service-api/application/logger"
	appvalidation "github.com/next-trace/scg-service-api/application/validation"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

// Context keys
const (
	validationModelKey contextKey = "validated_model"
	validationKey      contextKey = "validation_model"
)

// ValidationMiddleware provides middleware to validate request data.
type ValidationMiddleware struct {
	validator appvalidation.Validator
	config    appvalidation.Config
	log       applogger.Logger
}

// writeValidationErrorResponse writes validation error response to the response writer
func (vm *ValidationMiddleware) writeValidationErrorResponse(w http.ResponseWriter, r *http.Request, result appvalidation.ValidationResult) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"error":  "Validation failed",
		"errors": result.Errors,
	}); err != nil {
		vm.log.Error(r.Context(), err, "failed to encode validation errors response")
	}
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
			model := r.Context().Value(validationKey)
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
				vm.writeValidationErrorResponse(w, r, result)
				return
			}

			// Store the validated model in the request context
			ctx := context.WithValue(r.Context(), validationModelKey, modelValue)
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
				vm.writeValidationErrorResponse(w, r, result)
				return
			}

			// Store the validated model in the request context
			ctx := context.WithValue(r.Context(), validationModelKey, modelValue)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Validation provides backward compatibility with the old API.
// Deprecated: Use NewValidationMiddleware instead.
func Validation(validator appvalidation.Validator, config appvalidation.Config, log applogger.Logger) func(http.Handler) http.Handler {
	return NewValidationMiddleware(validator, config, log).Middleware()
}
