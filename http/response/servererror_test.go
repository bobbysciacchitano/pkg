package response

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriteServerError(t *testing.T) {
	// Setup a buffer to capture log output
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(nil) // Reset log output after the test

	// Setup a ResponseRecorder to capture HTTP response
	rec := httptest.NewRecorder()
	testErr := errors.New("something went wrong")

	// Call the function under test
	WriteServerError(rec, testErr)

	// Verify the status code
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}

	// Verify the error was logged
	logOutput := logBuf.String()
	if !strings.Contains(logOutput, testErr.Error()) {
		t.Errorf("expected log output to contain error message, got: %s", logOutput)
	}
}
