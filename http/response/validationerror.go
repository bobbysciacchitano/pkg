package response

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bobbysciacchitano/pkg/validator"
)

func WriteValidationError(w http.ResponseWriter, body validator.ValidationErrors) {
	w.WriteHeader(http.StatusBadRequest)

	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("could not encode validation error to JSON: %v", err)

		w.WriteHeader(http.StatusInternalServerError)
	}
}
