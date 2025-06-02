package authorization

import (
	"net/http"

	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

func (j *JWT) ValidateRequest(r *http.Request) (jwt.Token, error) {
	return jwt.ParseRequest(r, jwt.WithKey(jwa.RS256(), j.publicKey))
}
