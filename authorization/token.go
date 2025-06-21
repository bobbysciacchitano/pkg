package authorization

import (
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

func (j *JWT) CreateToken(subject string, claims map[string]string, ttl time.Duration) (*string, error) {
	builder := jwt.NewBuilder().
		Subject(subject).
		IssuedAt(time.Now()).
		Expiration(time.Now().Add(ttl))

	for k, v := range claims {
		builder = builder.Claim(k, v)
	}

	token, err := builder.Build()

	if err != nil {
		return nil, fmt.Errorf("could not create token: %w", err)
	}

	signed, err := jwt.Sign(token, jwt.WithKey(jwa.RS256(), j.privateKey))

	if err != nil {
		return nil, fmt.Errorf("could not sign token: %w", err)
	}

	asString := string(signed)

	return &asString, nil
}
