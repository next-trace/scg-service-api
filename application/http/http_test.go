package http_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	apphttp "github.com/hbttundar/scg-service-base/application/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRequestDecoder is a mock implementation of the RequestDecoder interface
type MockRequestDecoder struct {
	mock.Mock
}

// Ensure MockRequestDecoder implements the apphttp.RequestDecoder interface
var _ apphttp.RequestDecoder = (*MockRequestDecoder)(nil)

func (m *MockRequestDecoder) Decode(r *http.Request, v interface{}) error {
	args := m.Called(r, v)
	return args.Error(0)
}

// MockResponseWriter is a mock implementation of the ResponseWriter interface
type MockResponseWriter struct {
	mock.Mock
}

// Ensure MockResponseWriter implements the apphttp.ResponseWriter interface
var _ apphttp.ResponseWriter = (*MockResponseWriter)(nil)

func (m *MockResponseWriter) Respond(w http.ResponseWriter, r *http.Request, statusCode int, data interface{}) {
	m.Called(w, r, statusCode, data)
}

func (m *MockResponseWriter) Error(w http.ResponseWriter, r *http.Request, err error) {
	m.Called(w, r, err)
}

// TestData is a sample data structure for testing
type TestData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TestRequestDecoder(t *testing.T) {
	t.Run("Successful decode", func(t *testing.T) {
		// Create a mock decoder
		mockDecoder := new(MockRequestDecoder)

		// Create a test request
		body := strings.NewReader(`{"id": 123, "name": "Test Item"}`)
		req := httptest.NewRequest(http.MethodPost, "/test", body)

		// Set up expectations
		var data TestData
		mockDecoder.On("Decode", req, &data).Return(nil).Run(func(args mock.Arguments) {
			// Simulate decoding by setting values on the data struct
			v := args.Get(1).(*TestData)
			v.ID = 123
			v.Name = "Test Item"
		})

		// Call the decoder
		err := mockDecoder.Decode(req, &data)

		// Verify results
		assert.NoError(t, err)
		assert.Equal(t, 123, data.ID)
		assert.Equal(t, "Test Item", data.Name)
		mockDecoder.AssertExpectations(t)
	})

	t.Run("Failed decode", func(t *testing.T) {
		// Create a mock decoder
		mockDecoder := new(MockRequestDecoder)

		// Create a test request with invalid data
		body := strings.NewReader(`{"id": "not a number", "name": "Test Item"}`)
		req := httptest.NewRequest(http.MethodPost, "/test", body)

		// Set up expectations
		var data TestData
		expectedErr := errors.New("invalid type for field id")
		mockDecoder.On("Decode", req, &data).Return(expectedErr)

		// Call the decoder
		err := mockDecoder.Decode(req, &data)

		// Verify results
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		mockDecoder.AssertExpectations(t)
	})
}

func TestResponseWriter(t *testing.T) {
	t.Run("Successful response", func(t *testing.T) {
		// Create a mock response writer
		mockWriter := new(MockResponseWriter)

		// Create a test request and response
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		// Test data
		data := TestData{ID: 123, Name: "Test Item"}

		// Set up expectations
		mockWriter.On("Respond", w, req, http.StatusOK, data).Return()

		// Call the writer
		mockWriter.Respond(w, req, http.StatusOK, data)

		// Verify results
		mockWriter.AssertExpectations(t)
	})

	t.Run("Error response", func(t *testing.T) {
		// Create a mock response writer
		mockWriter := new(MockResponseWriter)

		// Create a test request and response
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		// Test error
		testErr := errors.New("test error")

		// Set up expectations
		mockWriter.On("Error", w, req, testErr).Return()

		// Call the writer
		mockWriter.Error(w, req, testErr)

		// Verify results
		mockWriter.AssertExpectations(t)
	})
}
