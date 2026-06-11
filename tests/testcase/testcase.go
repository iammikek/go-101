// Package testcase provides a Laravel-style TestCase base for HTTP and unit tests.
package testcase

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/iammikek/go-101/internal/app"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const DefaultAPIKey = "dev-key-123"

// Case is the base test case, similar to Laravel's TestCase.
type Case struct {
	T      *testing.T
	App    *app.Application
	APIKey string
}

// NewFeature creates a feature test case with a fresh database, like RefreshDatabase.
func NewFeature(t *testing.T, application *app.Application) *Case {
	t.Helper()

	tc := &Case{
		T:      t,
		App:    application,
		APIKey: application.Config.APIKey,
	}
	tc.ResetDatabase()
	return tc
}

// ResetDatabase clears all items so each test starts with an empty table.
func (tc *Case) ResetDatabase() {
	tc.T.Helper()
	require.NoError(tc.T, tc.App.ResetDatabase())
}

// Get issues a GET request against the application router.
func (tc *Case) Get(path string) *Response {
	return tc.request(http.MethodGet, path, "", nil)
}

// Post issues a POST request with a JSON body.
func (tc *Case) Post(path, body string) *Response {
	return tc.request(http.MethodPost, path, body, nil)
}

// Patch issues a PATCH request with a JSON body.
func (tc *Case) Patch(path, body string) *Response {
	return tc.request(http.MethodPatch, path, body, nil)
}

// Delete issues a DELETE request with optional headers.
func (tc *Case) Delete(path string, headers map[string]string) *Response {
	return tc.request(http.MethodDelete, path, "", headers)
}

func (tc *Case) request(method, path, body string, headers map[string]string) *Response {
	tc.T.Helper()

	var bodyReader io.Reader
	if body != "" {
		bodyReader = bytes.NewBufferString(body)
	}

	rec := httptest.NewRecorder()
	req, err := http.NewRequest(method, path, bodyReader)
	require.NoError(tc.T, err)

	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	tc.App.Router.ServeHTTP(rec, req)
	return &Response{T: tc.T, Recorder: rec}
}

// WithAPIKey returns headers containing the test API key.
func (tc *Case) WithAPIKey() map[string]string {
	return map[string]string{"X-API-Key": tc.APIKey}
}

// Response wraps httptest.ResponseRecorder with assertion helpers.
type Response struct {
	T        *testing.T
	Recorder *httptest.ResponseRecorder
}

// AssertStatus asserts the HTTP status code.
func (r *Response) AssertStatus(code int) {
	r.T.Helper()
	assert.Equal(r.T, code, r.Recorder.Code)
}

// AssertJSON asserts the response body matches expected JSON.
func (r *Response) AssertJSON(expected string) {
	r.T.Helper()
	assert.JSONEq(r.T, expected, r.Recorder.Body.String())
}

// JSON unmarshals the response body into a map.
func (r *Response) JSON() map[string]interface{} {
	r.T.Helper()
	var data map[string]interface{}
	require.NoError(r.T, json.Unmarshal(r.Recorder.Body.Bytes(), &data))
	return data
}

// JSONArray unmarshals the response body into a slice of maps.
func (r *Response) JSONArray() []map[string]interface{} {
	r.T.Helper()
	var data []map[string]interface{}
	require.NoError(r.T, json.Unmarshal(r.Recorder.Body.Bytes(), &data))
	return data
}
