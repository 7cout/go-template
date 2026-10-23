package httpapi_test

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/7cout/go-template/internal/adapters/httpapi"
	"github.com/7cout/go-template/internal/adapters/memory"
	"github.com/7cout/go-template/internal/app"
)

func newServer(t *testing.T) *httptest.Server {
	t.Helper()
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	srv := httptest.NewServer(httpapi.NewHandler(app.NewItems(memory.NewItems()), log))
	t.Cleanup(srv.Close)
	return srv
}

func TestHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
		wantBody   string
	}{
		{name: "health", method: http.MethodGet, path: "/healthz", wantStatus: http.StatusNoContent},
		{name: "create", method: http.MethodPost, path: "/items", body: `{"name":"apple"}`, wantStatus: http.StatusCreated, wantBody: `"name":"apple"`},
		{name: "create invalid", method: http.MethodPost, path: "/items", body: `{"name":""}`, wantStatus: http.StatusBadRequest},
		{name: "create malformed", method: http.MethodPost, path: "/items", body: `{`, wantStatus: http.StatusBadRequest},
		{name: "get missing", method: http.MethodGet, path: "/items/42", wantStatus: http.StatusNotFound},
		{name: "get bad id", method: http.MethodGet, path: "/items/abc", wantStatus: http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			srv := newServer(t)

			req, err := http.NewRequestWithContext(t.Context(), tt.method, srv.URL+tt.path, strings.NewReader(tt.body))
			if err != nil {
				t.Fatalf("NewRequest: %v", err)
			}
			resp, err := srv.Client().Do(req)
			if err != nil {
				t.Fatalf("%s %s: %v", tt.method, tt.path, err)
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("read body: %v", err)
			}

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("%s %s status = %d, want %d (body %s)", tt.method, tt.path, resp.StatusCode, tt.wantStatus, body)
			}
			if !strings.Contains(string(body), tt.wantBody) {
				t.Errorf("%s %s body = %s, want it to contain %s", tt.method, tt.path, body, tt.wantBody)
			}
		})
	}
}
