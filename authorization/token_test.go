package authorization

import (
	"path/filepath"
	"testing"
)

func TestCreateToken(t *testing.T) {
	dir := t.TempDir()

	keyPath := filepath.Join(dir, "testkey.pem")

	j, err := NewJWT(keyPath)

	if err != nil {
		t.Fatalf("failed to init jwt: %v", err)
	}

	token, err := j.CreateToken("user123", map[string]string{
		"role": "admin",
	}, 1)

	if err != nil {
		t.Fatalf("CreateToken failed: %v", err)
	}

	if token == nil || *token == "" {
		t.Error("expected non-empty token")
	}
}
