package validator

import (
	"encoding/json"
	"fmt"
	"strings"
)

// FieldError represents a validation error tied to a specific field.
type FieldError struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

func (fe FieldError) Error() string {
	return fe.Message
}

// NewError creates a general validation error not tied to a specific field.
func NewFieldError(message string) FieldError {
	return FieldError{
		Message: message,
	}
}

// ValidationErrors is a collection of field-related validation errors.
type ValidationErrors map[string]FieldError

// New returns an empty validation error map.
func NewValidator() ValidationErrors {
	return make(ValidationErrors)
}

// SetFieldError sets a FieldError explicitly.
func (ve ValidationErrors) FieldError(field string, err error) {
	if fe, ok := err.(FieldError); ok {
		ve[field] = fe
	}
}

// HasErrors reports whether any validation errors have been recorded.
func (ve ValidationErrors) HasErrors() bool {
	return len(ve) > 0
}

// Error returns a human-readable string of all validation errors.
func (ve ValidationErrors) Error() string {
	var sb strings.Builder
	for field, err := range ve {
		sb.WriteString(fmt.Sprintf("%s: %s\n", field, err.Message))
	}
	return sb.String()
}

// MarshalJSON outputs a flat map of field -> message for JSON responses.
func (ve ValidationErrors) MarshalJSON() ([]byte, error) {
	out := make(map[string]string, len(ve))
	for field, err := range ve {
		out[field] = err.Message
	}
	return json.Marshal(out)
}
