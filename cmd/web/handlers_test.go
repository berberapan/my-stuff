package main

import (
	"net/http"
	"testing"

	"github.com/berberapan/my-stuff/internal/assert"
)

func TestHealthzHandle(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	status, _, body := ts.get(t, "/healthz")

	assert.Equal(t, status, http.StatusOK)
	assert.Equal(t, body, "OK")
}
