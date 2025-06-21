package response

import (
	"log"
	"net/http"
)

func WriteServerError(w http.ResponseWriter, err error) {
	log.Print(err)

	w.WriteHeader(http.StatusInternalServerError)
}
