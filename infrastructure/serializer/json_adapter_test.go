package serializer_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/next-trace/scg-service-api/infrastructure/serializer"
	"github.com/stretchr/testify/assert"
)

type TestData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TestJSONAdapter_Decode(t *testing.T) {
	adapter := serializer.NewJSONAdapter()

	t.Run("Valid JSON", func(t *testing.T) {
		// Create a request with valid JSON body
		jsonData := `{"id": 123, "name": "Test Item"}`
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(jsonData))

		// Decode the request
		var data TestData
		err := adapter.Decode(req, &data)

		// Verify the result
		assert.NoError(t, err)
		assert.Equal(t, 123, data.ID)
		assert.Equal(t, "Test Item", data.Name)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		// Create a request with invalid JSON body
		invalidJSON := `{"id": 123, "name": "Missing closing quote}`
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(invalidJSON))

		// Decode the request
		var data TestData
		err := adapter.Decode(req, &data)

		// Verify the result
		assert.Error(t, err)
	})

	t.Run("Empty body", func(t *testing.T) {
		// Create a request with empty body
		req := httptest.NewRequest(http.MethodPost, "/test", nil)

		// Decode the request
		var data TestData
		err := adapter.Decode(req, &data)

		// Verify the result
		assert.Error(t, err)
	})
}

func TestJSONAdapter_Respond(t *testing.T) {
	adapter := serializer.NewJSONAdapter()

	t.Run("Respond with data", func(t *testing.T) {
		// Create a response recorder
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)

		// Test data
		data := TestData{ID: 123, Name: "Test Item"}

		// Call Respond
		adapter.Respond(w, req, http.StatusOK, data)

		// Verify the response
		resp := w.Result()
		defer resp.Body.Close()

		// Check status code
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Check content type
		assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

		// Check body
		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		var responseData TestData
		err = json.Unmarshal(body, &responseData)
		assert.NoError(t, err)
		assert.Equal(t, data.ID, responseData.ID)
		assert.Equal(t, data.Name, responseData.Name)
	})

	t.Run("Respond with nil data", func(t *testing.T) {
		// Create a response recorder
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)

		// Call Respond with nil data
		adapter.Respond(w, req, http.StatusNoContent, nil)

		// Verify the response
		resp := w.Result()
		defer resp.Body.Close()

		// Check status code
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Check content type
		assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

		// Check body (should be empty)
		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Empty(t, string(body))
	})
}

func TestJSONAdapter_Error(t *testing.T) {
	adapter := serializer.NewJSONAdapter()

	t.Run("Error response", func(t *testing.T) {
		// Create a response recorder
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)

		// Create an error
		testErr := errors.New("test error")

		// Call Error
		adapter.Error(w, req, testErr)

		// Verify the response
		resp := w.Result()
		defer resp.Body.Close()

		// Check status code (should be 500 by default)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		// Check content type
		assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))

		// Check body
		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		var errorResp struct {
			Error   string `json:"error"`
			TraceID string `json:"trace_id,omitempty"`
		}
		err = json.Unmarshal(body, &errorResp)
		assert.NoError(t, err)
		assert.Equal(t, "test error", errorResp.Error)
		// Note: TraceID will be empty in tests unless we mock the trace context
	})
}
