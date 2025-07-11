package http

import "net/http"

// ResponseWriter defines the abstract interface (PORT) for encoding (serializing)
// data and writing it as a standardized HTTP response.
type ResponseWriter interface {
	// Respond sends a standard structured success response.
	Respond(w http.ResponseWriter, r *http.Request, statusCode int, data interface{})

	// Error sends a standard structured error response.
	Error(w http.ResponseWriter, r *http.Request, err error)
}
