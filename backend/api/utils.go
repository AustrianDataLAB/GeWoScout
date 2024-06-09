package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/AustrianDataLAB/GeWoScout/backend/models"
)

// Checks that the user is marked as authenticated by the runtime and that
// the headers that are required are present. These headers include the client
// id and an email address of the user.
func GetClientPrincipalData[Q any](ir *models.InvokeRequest[Q]) (string, string, error) {
	principalIds, ok1 := ir.Data.Req.Headers[models.CLIENT_PRINCIPAL_ID_KEY]
	principals, ok2 := ir.Data.Req.Headers[models.CLIENT_PRINCIPAL_KEY]
	if !ok1 || !ok2 || !ir.Data.Req.Identities[0].IsAuthenticated {
		return "", "", errors.New("user not authenticated")
	}

	clientId := principalIds[0]
	principalB64 := principals[0]

	pDec, _ := base64.StdEncoding.DecodeString(principalB64)

	up := models.UserPrincipal{}
	if err := json.Unmarshal(pDec, &up); err != nil {
		return "", "", errors.New("failed to read user principal")
	}

	var email string
	for _, c := range up.Claims {
		if c.Typ == "preferred_username" {
			email = c.Val
		}
	}
	if email == "" {
		return "", "", errors.New("failed to read user email")
	}

	return clientId, email, nil
}

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
