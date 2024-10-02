package main

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/berberapan/my-stuff/internal/assert"
)

func TestLogRequestMiddleware(t *testing.T) {
	expectedValues := []string{
		"received request",
		"ip=192.0.15.25:1123",
		"proto=HTTP/1.1",
		"method=GET",
		"uri=/test",
	}

	var logBuffer bytes.Buffer
	app := newTestApplication(t)
	app.logger = slog.New(slog.NewTextHandler(&logBuffer, nil))

	rr := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://localhost:8080/test", nil)
	r.RemoteAddr = "192.0.15.25:1123"

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	app.logRequest(next).ServeHTTP(rr, r)

	assert.Equal(t, rr.Code, http.StatusOK)

	logOutput := logBuffer.String()
	for _, entry := range expectedValues {
		assert.StringContains(t, logOutput, entry)
	}
}

func TestRecoverPanic(t *testing.T) {
	tests := []struct {
		name           string
		handler        http.Handler
		expectedStatus int
		expectedHeader string
	}{
		{
			name: "No Panic",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
			expectedStatus: http.StatusOK,
			expectedHeader: "",
		},
		{
			name: "Panic",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic(fmt.Errorf("Something gone wrong!"))
			}),
			expectedStatus: http.StatusInternalServerError,
			expectedHeader: "close",
		},
	}

	app := newTestApplication(t)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/", nil)

			app.recoverPanic(test.handler).ServeHTTP(rr, r)
			assert.Equal(t, rr.Code, test.expectedStatus)
			assert.Equal(t, rr.Header().Get("Connection"), test.expectedHeader)
		})
	}
}

func TestCommonHeaders(t *testing.T) {
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	commonHeaders(next).ServeHTTP(rr, r)

	rs := rr.Result()

	expectedValue := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
	assert.Equal(t, rs.Header.Get("Content-Security-Policy"), expectedValue)

	expectedValue = "origin-when-cross-origin"
	assert.Equal(t, rs.Header.Get("Referrer-Policy"), expectedValue)

	expectedValue = "nosniff"
	assert.Equal(t, rs.Header.Get("X-Content-Type-Options"), expectedValue)

	expectedValue = "deny"
	assert.Equal(t, rs.Header.Get("X-Frame-Options"), expectedValue)

	expectedValue = "0"
	assert.Equal(t, rs.Header.Get("X-XSS-Protection"), expectedValue)

	assert.Equal(t, rs.StatusCode, http.StatusOK)

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}
