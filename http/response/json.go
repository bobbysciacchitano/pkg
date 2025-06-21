package response

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, body any) {
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("could not encode response body to JSON: %v", err)

		WriteServerError(w, err)
	}
}
