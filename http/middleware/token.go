package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/bobbysciacchitano/pkg/authorization"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

type tokenContextKey string

const TokenKey tokenContextKey = "token"

func Token(auth authorization.DirectorAuthorization) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := auth.ValidateRequest(r)

			if err != nil {
				log.Printf("unauthorized: invalid token: %v", err)

				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, TokenKey, token)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func TokenFromContext(ctx context.Context) (jwt.Token, bool) {
	token, ok := ctx.Value(TokenKey).(jwt.Token)

	return token, ok
}
