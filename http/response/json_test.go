package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteJSON_Success(t *testing.T) {
	rr := httptest.NewRecorder()

	body := map[string]string{"message": "hello"}

	WriteJSON(rr, body)

	result := rr.Result()

	defer result.Body.Close()

	if result.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", result.StatusCode)
	}

	var got map[string]string

	if err := json.NewDecoder(result.Body).Decode(&got); err != nil {
		t.Fatalf("could not decode response body: %v", err)
	}

	if got["message"] != "hello" {
		t.Errorf("expected 'hello', got %q", got["message"])
	}
}

func TestWriteJSON_EncodeFailure(t *testing.T) {
	// response.WriteJSON writes 500 after logging an encode error
	// simulate it by using a value that can't be marshaled
	rr := httptest.NewRecorder()

	// Functions can't be JSON encoded
	type Bad struct {
		Fn func()
	}

	WriteJSON(rr, Bad{Fn: func() {}})

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rr.Code)
	}
}
