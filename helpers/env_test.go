package helpers

import "testing"

func TestGetEnv(t *testing.T) {
	t.Setenv("MY_VALUE", "test")

	if v := Getenv("MY_VALUE", "fallback"); v != "test" {
		t.Errorf("expected value test got %s", v)
	}
}

func TestGetEnvFallback(t *testing.T) {
	if v := Getenv("MY_VALUE", "fallback"); v != "fallback" {
		t.Errorf("expected value test got %s", v)
	}
}
