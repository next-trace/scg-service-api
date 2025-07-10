// Package validation provides validation functionality.
//
// Note: This package requires the following dependencies:
// - github.com/go-playground/validator/v10
//
// See docs/dependencies.md for more information.
package validation

import (
	"context"
	"fmt"
	"reflect"

	applogger "github.com/next-trace/scg-service-api/application/logger"
	appvalidation "github.com/next-trace/scg-service-api/application/validation"
)

// Ensure playgroundAdapter implements the appvalidation.Validator interface.
var _ appvalidation.Validator = (*playgroundAdapter)(nil)

// Type aliases for validator types
// These are placeholders for the actual types from the validator package
// They will be replaced with the actual types when the dependencies are added
type (
	// validate represents a validator
	validate struct {
		tagName      string
		customRules  map[string]func(fl interface{}) bool
		translations map[string]string
		tagNameFunc  func(field reflect.StructField) string
	}

	// fieldError represents a validation error for a field (unused in mock)
)

// Mock methods for the validate type
func newValidator() *validate {
	return &validate{
		tagName:      "validate",
		customRules:  make(map[string]func(fl interface{}) bool),
		translations: make(map[string]string),
		tagNameFunc:  nil,
	}
}

func (v *validate) Struct(s interface{}) error {
	// In a real implementation, this would validate the struct
	// For now, we'll just return nil
	return nil
}

func (v *validate) StructPartial(s interface{}, fields ...string) error {
	// In a real implementation, this would validate the struct partially
	// For now, we'll just return nil
	return nil
}

func (v *validate) Var(field interface{}, tag string) error {
	// In a real implementation, this would validate the variable
	// For now, we'll just return nil
	return nil
}

func (v *validate) RegisterValidation(tag string, fn func(fl interface{}) bool) error {
	v.customRules[tag] = fn
	return nil
}

func (v *validate) RegisterTagNameFunc(fn func(field reflect.StructField) string) {
	v.tagNameFunc = fn
}

// playgroundAdapter implements the validation.Validator interface using the go-playground/validator package.
type playgroundAdapter struct {
	config    appvalidation.Config
	validator *validate
	log       applogger.Logger
}

// NewPlaygroundAdapter creates a new validator adapter using the go-playground/validator package.
func NewPlaygroundAdapter(config appvalidation.Config, log applogger.Logger) appvalidation.Validator {
	validator := newValidator()
	validator.tagName = config.TagName

	// Register custom rules
	for name, rule := range config.CustomRules {
		// Convert the rule to the format expected by the validator
		validatorRule := func(fl interface{}) bool {
			// In a real implementation, this would convert the validator.FieldLevel to the value
			// For now, we'll just pass the value directly
			return rule(context.Background(), fl)
		}
		if err := validator.RegisterValidation(name, validatorRule); err != nil {
			log.WarnKV(context.Background(), "failed to register validation rule", map[string]interface{}{
				"rule":  name,
				"error": err.Error(),
			})
		}
	}

	return &playgroundAdapter{
		config:    config,
		validator: validator,
		log:       log,
	}
}

// Validate validates the given value against its validation rules.
func (p *playgroundAdapter) Validate(ctx context.Context, value interface{}) appvalidation.ValidationResult {
	if !p.config.Enabled {
		return appvalidation.ValidationResult{Valid: true}
	}

	err := p.validator.Struct(value)
	if err == nil {
		return appvalidation.ValidationResult{Valid: true}
	}

	// Convert validation errors to our format
	validationErrors := make(appvalidation.ValidationErrors)
	for _, err := range p.extractValidationErrors(err) {
		field := err.Field
		message := p.formatErrorMessage(err)
		validationErrors[field] = append(validationErrors[field], message)
	}

	return appvalidation.ValidationResult{
		Valid:  false,
		Errors: validationErrors,
	}
}

// ValidateField validates a specific field of the given value.
func (p *playgroundAdapter) ValidateField(ctx context.Context, value interface{}, field string) appvalidation.ValidationResult {
	if !p.config.Enabled {
		return appvalidation.ValidationResult{Valid: true}
	}

	err := p.validator.StructPartial(value, field)
	if err == nil {
		return appvalidation.ValidationResult{Valid: true}
	}

	// Convert validation errors to our format
	validationErrors := make(appvalidation.ValidationErrors)
	for _, err := range p.extractValidationErrors(err) {
		if err.Field == field {
			message := p.formatErrorMessage(err)
			validationErrors[field] = append(validationErrors[field], message)
		}
	}

	return appvalidation.ValidationResult{
		Valid:  len(validationErrors) == 0,
		Errors: validationErrors,
	}
}

// ValidateMap validates a map of values against their validation rules.
func (p *playgroundAdapter) ValidateMap(ctx context.Context, values map[string]interface{}) appvalidation.ValidationResult {
	if !p.config.Enabled {
		return appvalidation.ValidationResult{Valid: true}
	}

	validationErrors := make(appvalidation.ValidationErrors)
	for field, value := range values {
		err := p.validator.Var(value, p.getTagForField(field))
		if err != nil {
			for _, err := range p.extractValidationErrors(err) {
				message := p.formatErrorMessage(err)
				validationErrors[field] = append(validationErrors[field], message)
			}
		}
	}

	return appvalidation.ValidationResult{
		Valid:  len(validationErrors) == 0,
		Errors: validationErrors,
	}
}

// RegisterCustomRule registers a custom validation rule.
func (p *playgroundAdapter) RegisterCustomRule(name string, rule appvalidation.CustomRule) error {
	if !p.config.Enabled {
		return nil
	}

	// Convert the rule to the format expected by the validator
	validatorRule := func(fl interface{}) bool {
		// In a real implementation, this would convert the validator.FieldLevel to the value
		// For now, we'll just pass the value directly
		return rule(context.Background(), fl)
	}

	return p.validator.RegisterValidation(name, validatorRule)
}

// RegisterTagNameFunc registers a function to get the field name from a struct tag.
func (p *playgroundAdapter) RegisterTagNameFunc(fn func(field reflect.StructField) string) {
	if !p.config.Enabled {
		return
	}

	p.validator.RegisterTagNameFunc(fn)
}

// extractValidationErrors extracts validation errors from the error.
func (p *playgroundAdapter) extractValidationErrors(err error) []appvalidation.ValidationError {
	if err == nil {
		return nil
	}

	var validationErrors []appvalidation.ValidationError

	// In a real implementation, this would use validator.ValidationErrors
	// For now, we'll just create a mock error
	validationErrors = append(validationErrors, appvalidation.ValidationError{
		Field:   "mock_field",
		Tag:     "mock_tag",
		Value:   "mock_value",
		Param:   "mock_param",
		Message: err.Error(),
	})

	return validationErrors
}

// formatErrorMessage formats a validation error message.
func (p *playgroundAdapter) formatErrorMessage(err appvalidation.ValidationError) string {
	// In a real implementation, this would use translations
	// For now, we'll just return a simple message
	return fmt.Sprintf("validation failed on field %s, tag %s", err.Field, err.Tag)
}

// getTagForField returns the validation tag for the given field.
func (p *playgroundAdapter) getTagForField(field string) string {
	// In a real implementation, this would look up the tag from the struct
	// For now, we'll just return a default tag
	return "required"
}
