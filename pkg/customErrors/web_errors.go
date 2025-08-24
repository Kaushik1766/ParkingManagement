package customerrors

import (
	"encoding/json"
	"log"
	"net/http"
)

type WebError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func UnauthorizedError(w http.ResponseWriter, err WebError) {
	w.WriteHeader(http.StatusUnauthorized)

	log.Println("Unauthorized access attempt:", err.Message)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(err)
}

func BadRequestError(w http.ResponseWriter, err WebError) {
	w.WriteHeader(http.StatusBadRequest)

	log.Println("Bad request:", err.Message)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(err)
}

func InternalServerError(w http.ResponseWriter, err WebError) {
	w.WriteHeader(http.StatusInternalServerError)

	log.Println("Internale server error:", err.Message)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(err)
}
