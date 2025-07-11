// Package errors provides domain-specific error types.
package errors

import (
	"errors"
	"fmt"
)

// Standard error types that can be used for error handling.
var (
	// ErrNotFound indicates that a requested resource was not found.
	ErrNotFound = errors.New("resource not found")

	// ErrAlreadyExists indicates that a resource already exists.
	ErrAlreadyExists = errors.New("resource already exists")

	// ErrInvalidInput indicates that the input is invalid.
	ErrInvalidInput = errors.New("invalid input")

	// ErrUnauthorized indicates that the user is not authorized.
	ErrUnauthorized = errors.New("unauthorized")

	// ErrForbidden indicates that the user is forbidden from accessing a resource.
	ErrForbidden = errors.New("forbidden")

	// ErrInternal indicates an internal server error.
	ErrInternal = errors.New("internal error")

	// ErrTimeout indicates that an operation timed out.
	ErrTimeout = errors.New("operation timed out")

	// ErrUnavailable indicates that a service is unavailable.
	ErrUnavailable = errors.New("service unavailable")
)

// DomainError represents a domain-specific error.
type DomainError struct {
	// Err is the underlying error.
	Err error

	// Code is a machine-readable error code.
	Code string

	// Message is a human-readable error message.
	Message string

	// Details contains additional error details.
	Details map[string]interface{}
}

// Error returns the error message.
func (e *DomainError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

// Unwrap returns the underlying error.
func (e *DomainError) Unwrap() error {
	return e.Err
}

// WithCode sets the error code.
func (e *DomainError) WithCode(code string) *DomainError {
	e.Code = code
	return e
}

// WithMessage sets the error message.
func (e *DomainError) WithMessage(format string, args ...interface{}) *DomainError {
	e.Message = fmt.Sprintf(format, args...)
	return e
}

// WithDetail adds a detail to the error.
func (e *DomainError) WithDetail(key string, value interface{}) *DomainError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// NewNotFound creates a new not found error.
func NewNotFound(entity string, id interface{}) *DomainError {
	err := &DomainError{
		Err:  ErrNotFound,
		Code: "not_found",
		Details: map[string]interface{}{
			"entity": entity,
			"id":     id,
		},
	}
	return err.WithMessage("%s with ID %v not found", entity, id)
}

// NewAlreadyExists creates a new already exists error.
func NewAlreadyExists(entity string, id interface{}) *DomainError {
	err := &DomainError{
		Err:  ErrAlreadyExists,
		Code: "already_exists",
		Details: map[string]interface{}{
			"entity": entity,
			"id":     id,
		},
	}
	return err.WithMessage("%s with ID %v already exists", entity, id)
}

// NewInvalidInput creates a new invalid input error.
func NewInvalidInput(reason string) *DomainError {
	err := &DomainError{
		Err:  ErrInvalidInput,
		Code: "invalid_input",
		Details: map[string]interface{}{
			"reason": reason,
		},
	}
	return err.WithMessage("invalid input: %s", reason)
}

// NewUnauthorized creates a new unauthorized error.
func NewUnauthorized(reason string) *DomainError {
	err := &DomainError{
		Err:  ErrUnauthorized,
		Code: "unauthorized",
		Details: map[string]interface{}{
			"reason": reason,
		},
	}
	return err.WithMessage("unauthorized: %s", reason)
}

// NewForbidden creates a new forbidden error.
func NewForbidden(reason string) *DomainError {
	err := &DomainError{
		Err:  ErrForbidden,
		Code: "forbidden",
		Details: map[string]interface{}{
			"reason": reason,
		},
	}
	return err.WithMessage("forbidden: %s", reason)
}

// NewInternal creates a new internal error.
func NewInternal(err error) *DomainError {
	domainErr := &DomainError{
		Err:  ErrInternal,
		Code: "internal_error",
	}
	return domainErr.WithMessage("internal error: %v", err)
}

// IsNotFound returns true if the error is a not found error.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsAlreadyExists returns true if the error is an already exists error.
func IsAlreadyExists(err error) bool {
	return errors.Is(err, ErrAlreadyExists)
}

// IsInvalidInput returns true if the error is an invalid input error.
func IsInvalidInput(err error) bool {
	return errors.Is(err, ErrInvalidInput)
}

// IsUnauthorized returns true if the error is an unauthorized error.
func IsUnauthorized(err error) bool {
	return errors.Is(err, ErrUnauthorized)
}

// IsForbidden returns true if the error is a forbidden error.
func IsForbidden(err error) bool {
	return errors.Is(err, ErrForbidden)
}

// IsInternal returns true if the error is an internal error.
func IsInternal(err error) bool {
	return errors.Is(err, ErrInternal)
}

// IsTimeout returns true if the error is a timeout error.
func IsTimeout(err error) bool {
	return errors.Is(err, ErrTimeout)
}

// IsUnavailable returns true if the error is an unavailable error.
func IsUnavailable(err error) bool {
	return errors.Is(err, ErrUnavailable)
}
