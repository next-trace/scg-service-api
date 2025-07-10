// Package validation defines the abstract interface (PORT) for validation.
package validation

import (
	"context"
	"reflect"
)

// Validator defines the interface for validating data.
type Validator interface {
	// Validate validates the given value against its validation rules.
	// It returns a ValidationResult containing any validation errors.
	Validate(ctx context.Context, value interface{}) ValidationResult

	// ValidateField validates a specific field of the given value.
	// It returns a ValidationResult containing any validation errors for that field.
	ValidateField(ctx context.Context, value interface{}, field string) ValidationResult

	// ValidateMap validates a map of values against their validation rules.
	// It returns a ValidationResult containing any validation errors.
	ValidateMap(ctx context.Context, values map[string]interface{}) ValidationResult

	// RegisterCustomRule registers a custom validation rule.
	RegisterCustomRule(name string, rule CustomRule) error

	// RegisterTagNameFunc registers a function to get the field name from a struct tag.
	RegisterTagNameFunc(fn func(field reflect.StructField) string)
}

// ValidationResult represents the result of a validation operation.
type ValidationResult struct {
	// Valid indicates whether the validation passed.
	Valid bool

	// Errors contains validation errors, if any.
	Errors ValidationErrors
}

// ValidationErrors is a map of field names to validation error messages.
type ValidationErrors map[string][]string

// CustomRule defines a custom validation rule.
type CustomRule func(ctx context.Context, value interface{}, params ...string) bool

// Config holds configuration for validation.
type Config struct {
	// Enabled determines if validation is enabled.
	Enabled bool

	// TagName is the struct tag name to use for validation rules.
	TagName string

	// CustomRules is a map of custom validation rules.
	CustomRules map[string]CustomRule
}

// DefaultConfig returns the default configuration for validation.
func DefaultConfig() Config {
	return Config{
		Enabled:     true,
		TagName:     "validate",
		CustomRules: map[string]CustomRule{},
	}
}

// ValidationError represents a validation error.
type ValidationError struct {
	// Field is the name of the field that failed validation.
	Field string

	// Tag is the validation tag that failed.
	Tag string

	// Value is the value that was validated.
	Value interface{}

	// Param is the parameter for the validation tag.
	Param string

	// Message is a human-readable error message.
	Message string
}

// Error returns the error message.
func (e ValidationError) Error() string {
	return e.Message
}
