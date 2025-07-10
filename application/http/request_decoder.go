package http

import "net/http"

// RequestDecoder defines the abstract interface (PORT) for decoding (unserializing)
// an HTTP request body into a target data structure.
type RequestDecoder interface {
	Decode(r *http.Request, v interface{}) error
}
