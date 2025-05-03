package authorization

import (
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

func (j *JWT) CreateToken(subject string) (*string, error) {
	token, err := jwt.NewBuilder().
		Subject(subject).
		IssuedAt(time.Now()).
		Build()

	if err != nil {
		return nil, fmt.Errorf("could not create token: %w", err)
	}

	signed, err := jwt.Sign(token, jwt.WithKey(jwa.RS256(), j.privateKey))

	if err != nil {
		return nil, err
	}

	asString := string(signed)

	return &asString, nil
}
