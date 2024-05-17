package api

import (
	"bytes"
	"net/http"
)

// MockResponseWriter captures the HTTP response data
type MockResponseWriter struct {
	Headers    http.Header
	Body       bytes.Buffer
	StatusCode int
}

// NewMockResponseWriter creates and initializes a new MockResponseWriter
func NewMockResponseWriter() *MockResponseWriter {
	return &MockResponseWriter{
		Headers:    http.Header{},
		StatusCode: http.StatusOK, // Default to 200 OK
	}
}

// Header returns the header map that will be sent by WriteHeader
func (m *MockResponseWriter) Header() http.Header {
	return m.Headers
}

// Write writes the data to the buffer and satisfies the http.ResponseWriter interface
func (m *MockResponseWriter) Write(data []byte) (int, error) {
	return m.Body.Write(data)
}

// WriteHeader stores the status code
func (m *MockResponseWriter) WriteHeader(statusCode int) {
	m.StatusCode = statusCode
}
