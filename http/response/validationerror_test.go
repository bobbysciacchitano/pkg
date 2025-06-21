package response

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bobbysciacchitano/pkg/validator"
)

func TestWriteValidationError(t *testing.T) {
	// Capture log output
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(nil)

	// Arrange
	rec := httptest.NewRecorder()
	validationErrors := validator.ValidationErrors{
		"name": "Name is required",
	}

	// Act
	WriteValidationError(rec, validationErrors)

	// Assert status
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}

	// Assert JSON body
	var body map[string]string

	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if got := body["name"]; got != "Name is required" {
		t.Errorf("expected error message for 'name', got: %s", got)
	}

	// Assert no encoding errors were logged
	if logBuf.Len() > 0 {
		t.Errorf("unexpected log output: %s", logBuf.String())
	}
}
