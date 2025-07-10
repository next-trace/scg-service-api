package serializer

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	apphttp "github.com/next-trace/scg-service-api/application/http"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// JSONAdapter implements both RequestDecoder and ResponseWriter interfaces.
type JSONAdapter struct{}

// Ensure JSONAdapter implements the apphttp.RequestDecoder interface
var _ apphttp.RequestDecoder = (*JSONAdapter)(nil)

// Ensure JSONAdapter implements the apphttp.ResponseWriter interface
var _ apphttp.ResponseWriter = (*JSONAdapter)(nil)

// NewJSONAdapter creates a new adapter for JSON serialization.
func NewJSONAdapter() *JSONAdapter {
	return &JSONAdapter{}
}

func (a *JSONAdapter) Decode(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func (a *JSONAdapter) Respond(w http.ResponseWriter, _ *http.Request, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if data == nil {
		w.WriteHeader(statusCode)
		return
	}
	// Encode first into a buffer to avoid sending partial/incorrect responses
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	_, _ = w.Write(buf.Bytes())
}

// Error writes a standardized error response with appropriate status code.
// It maps different error types to appropriate HTTP status codes.
func (a *JSONAdapter) Error(w http.ResponseWriter, r *http.Request, err error) {
	type errorResponse struct {
		Error   string `json:"error"`
		TraceID string `json:"trace_id,omitempty"`
		Code    string `json:"code,omitempty"`
	}

	// Extract trace ID if available
	span := trace.SpanFromContext(r.Context())
	traceID := ""
	if span.SpanContext().IsValid() {
		traceID = span.SpanContext().TraceID().String()
	}

	// Default status code and error code
	statusCode := http.StatusInternalServerError
	errorCode := "internal_error"

	// Map error types to appropriate status codes
	// This can be extended with custom error types
	switch {
	case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
		statusCode = http.StatusGatewayTimeout
		errorCode = "timeout"
	case errors.Is(err, http.ErrBodyNotAllowed), errors.Is(err, http.ErrMissingFile),
		errors.Is(err, http.ErrNotMultipart), errors.Is(err, http.ErrNoCookie):
		statusCode = http.StatusBadRequest
		errorCode = "invalid_request"
	case errors.Is(err, http.ErrHandlerTimeout):
		statusCode = http.StatusServiceUnavailable
		errorCode = "service_timeout"
	case errors.Is(err, http.ErrAbortHandler):
		statusCode = http.StatusInternalServerError
		errorCode = "request_aborted"
	}

	// Create the error response
	resp := errorResponse{
		Error:   err.Error(),
		TraceID: traceID,
		Code:    errorCode,
	}

	// Record the error in the span if available
	if span.SpanContext().IsValid() {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	a.Respond(w, r, statusCode, resp)
}
