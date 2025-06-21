package authorization

import (
	"path/filepath"
	"testing"
)

func TestNewJWT_GeneratesAndLoadsKey(t *testing.T) {
	dir := t.TempDir()

	keyPath := filepath.Join(dir, "testkey.pem")

	// Should generate a new key
	j1, err := NewJWT(keyPath)

	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	if j1 == nil {
		t.Fatal("expected JWT instance, got nil")
	}

	// Should load the existing key this time
	j2, err := NewJWT(keyPath)

	if err != nil {
		t.Fatalf("failed to load existing key: %v", err)
	}

	if j2 == nil {
		t.Fatal("expected JWT instance, got nil")
	}
}
