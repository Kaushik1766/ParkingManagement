package customerrors

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Kaushik1766/ParkingManagement/internal/constants/error_codes"
)

type WebError struct {
	Message    string                `json:"message"`
	Code       errorcodes.ErrorCodes `json:"code"`
	statusCode int
}

func NewWebError(code errorcodes.ErrorCodes) WebError {
	return WebError{
		Message:    code.String(),
		Code:       code,
		statusCode: code.Status(),
	}
}

func (err WebError) Error() string {
	return err.Message
}

func SendError(w http.ResponseWriter, err error) {
	webErr, ok := err.(WebError)
	if !ok {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(webErr.statusCode)
	json.NewEncoder(w).Encode(err)
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

	log.Println("Internal server error:", err.Message)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(err)
}
