package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bobbysciacchitano/pkg/authorization"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

func TestTokenMiddleware_Success(t *testing.T) {
	expectedSub := "user123"

	mock := &authorization.MockAuthorization{
		ValidateRequestFunc: func(r *http.Request) (jwt.Token, error) {
			token := jwt.New()
			_ = token.Set(jwt.SubjectKey, expectedSub)
			_ = token.Set(jwt.IssuedAtKey, time.Now())
			return token, nil
		},
	}

	// Dummy handler to assert token presence
	var called bool

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true

		token, ok := TokenFromContext(r.Context())

		if !ok {
			t.Fatal("token not found in context")
		}

		subject, _ := token.Subject()

		if subject != expectedSub {
			t.Errorf("expected subject %q, got %q", expectedSub, subject)
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	mw := Token(mock)
	mw(handler).ServeHTTP(rr, req)

	if !called {
		t.Error("handler was not called")
	}
}

func TestTokenMiddleware_Unauthorized(t *testing.T) {
	mock := &authorization.MockAuthorization{
		ValidateRequestFunc: func(r *http.Request) (jwt.Token, error) {
			return nil, errors.New("invalid token")
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	mw := Token(mock)

	mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called on unauthorized request")
	})).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, status)
	}
}
