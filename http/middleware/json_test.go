package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSONMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"ok":true}`))
	})

	// Wrap the handler with the middleware
	wrapped := JSONMiddleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	wrapped.ServeHTTP(rr, req)

	if got := rr.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("expected Content-Type to be 'application/json', got '%s'", got)
	}

	expectedBody := `{"ok":true}`

	if rr.Body.String() != expectedBody {
		t.Errorf("expected body to be '%s', got '%s'", expectedBody, rr.Body.String())
	}
}
