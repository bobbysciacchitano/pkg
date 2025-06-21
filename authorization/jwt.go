package authorization

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

type DirectorAuthorization interface {
	CreateToken(string, map[string]string, time.Duration) (*string, error)
	ValidateRequest(*http.Request) (jwt.Token, error)
}

type JWT struct {
	privateKey *rsa.PrivateKey
	publicKey  jwk.Key
}

func NewJWT(path string) (*JWT, error) {
	if path == "" {
		path = "./private.pem"
	}

	var privateKey *rsa.PrivateKey

	// Try to load existing key
	if data, err := os.ReadFile(path); err == nil {
		block, _ := pem.Decode(data)

		if block != nil && block.Type == "RSA PRIVATE KEY" {
			key, err := x509.ParsePKCS1PrivateKey(block.Bytes)

			if err == nil {
				privateKey = key
			} else {
				log.Printf("failed to parse existing private key, generating new one: %v", err)
			}
		}
	}

	// If no valid key loaded, generate one
	if privateKey == nil {
		log.Printf("generating new RSA private key at %s", path)

		var err error

		privateKey, err = rsa.GenerateKey(rand.Reader, 2048)

		if err != nil {
			return nil, fmt.Errorf("failed to generate private key %s: %w", path, err)
		}

		// Write to file
		pemBytes := pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		})

		err = os.WriteFile(path, pemBytes, 0600)

		if err != nil {
			return nil, fmt.Errorf("failed to write private key to %s: %w", path, err)
		}
	}

	publicKey, err := jwk.PublicKeyOf(privateKey)

	if err != nil {
		return nil, fmt.Errorf("failed get get public key: %w", err)
	}

	return &JWT{
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}
