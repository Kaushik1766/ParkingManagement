package customerrors

import (
	"encoding/json"
	"log"
	"net/http"

	errorcodes "github.com/Kaushik1766/ParkingManagement/internal/constants/error_codes"
	"github.com/aws/aws-lambda-go/events"
)

type WebError struct {
	Message    string                `json:"message"`
	Code       errorcodes.ErrorCodes `json:"code"`
	statusCode int
}

func NewWebError(err error, code errorcodes.ErrorCodes) WebError {
	log.Printf("error occurrred %v\n", err)
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
	w.Header().Set("Content-Type", "application/json")
	webErr, ok := err.(WebError)
	if !ok {
		log.Printf("cant type assert %v to weberr\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}
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

func LambdaError(code int, msg string) events.APIGatewayProxyResponse {
	body, _ := json.Marshal(map[string]string{
		"message": msg,
	})
	return events.APIGatewayProxyResponse{
		StatusCode: code,
		Body:       string(body),
	}
}
