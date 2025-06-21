package authorization

import (
	"net/http"
	"path/filepath"
	"testing"
	"time"
)

func TestValidateRequest(t *testing.T) {
	dir := t.TempDir()

	keyPath := filepath.Join(dir, "testkey.pem")

	j, err := NewJWT(keyPath)

	if err != nil {
		t.Fatalf("failed to init jwt: %v", err)
	}

	tokenStr, err := j.CreateToken("test-user", nil, time.Hour*24*30*12)

	if err != nil {
		t.Fatalf("failed to create token: %v", err)
	}

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+*tokenStr)

	token, err := j.ValidateRequest(req)

	if err != nil {
		t.Fatalf("expected token to be valid, got error: %v", err)
	}

	subject, _ := token.Subject()

	if subject != "test-user" {
		t.Errorf("expected subject 'test-user', got: %s", subject)
	}
}
