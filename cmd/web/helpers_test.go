package main

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/berberapan/my-stuff/internal/assert"
)

func TestRender(t *testing.T) {
	app := newTestApplication(t)
	app.templateCache = map[string]*template.Template{
		"test.tmpl": template.Must(template.New("base").Parse("{{define \"base\"}}Testing succesful{{end}}")),
	}

	tests := []struct {
		name           string
		page           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Existing template render",
			page:           "test.tmpl",
			expectedStatus: http.StatusOK,
			expectedBody:   "Testing succesful",
		},
		{
			name:           "Non-exisiting template render",
			page:           "incorrectfile.tmpl",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Internal Server Error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			app.render(rr, r, tt.expectedStatus, tt.page)

			assert.Equal(t, rr.Code, tt.expectedStatus)
			assert.Equal(t, strings.TrimSpace(rr.Body.String()), strings.TrimSpace(tt.expectedBody))
		})
	}
}
