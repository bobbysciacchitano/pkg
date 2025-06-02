package validator

import (
	"encoding/json"
	"fmt"
	"strings"
)

// FieldError is a typed string that implements error
type FieldError string

func NewFieldError(message string) FieldError {
	return FieldError(message)
}

func (e FieldError) Error() string {
	return string(e)
}

// Validator stores validation errors by field and implements error
type Validator map[string]FieldError

// NewValidator creates a new, empty Validator
func NewValidator() Validator {
	return Validator{}
}

// AddError adds a field + message
func (v Validator) AddError(field string, err error) {
	if err != nil {
		v[field] = FieldError(err.Error())
	}
}

// HasErrors returns true if there are any field errors
func (v Validator) HasErrors() bool {
	return len(v) > 0
}

// Error returns a joined summary of all validation messages
func (v Validator) Error() string {
	if len(v) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("validation failed:")

	for field, msg := range v {
		sb.WriteString(fmt.Sprintf("\n- %s: %s", field, msg))
	}

	return sb.String()
}

// MarshalJSON outputs the flat { "field": "message" } structure
func (v Validator) MarshalJSON() ([]byte, error) {
	flat := make(map[string]string, len(v))

	for field, err := range v {
		flat[field] = err.Error()
	}

	return json.Marshal(flat)
}
