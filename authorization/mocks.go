package authorization

import (
	"net/http"
	"time"

	"github.com/lestrrat-go/jwx/v3/jwt"
)

type MockAuthorization struct {
	CreateTokenFunc     func(string, map[string]string, time.Duration) (*string, error)
	ValidateRequestFunc func(*http.Request) (jwt.Token, error)
}

func (m *MockAuthorization) CreateToken(subject string, claims map[string]string, expires time.Duration) (*string, error) {
	return m.CreateTokenFunc(subject, claims, expires)
}

func (m *MockAuthorization) ValidateRequest(r *http.Request) (jwt.Token, error) {
	return m.ValidateRequestFunc(r)
}
