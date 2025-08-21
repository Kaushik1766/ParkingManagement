package authhandler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	authservice "github.com/Kaushik1766/ParkingManagement/internal/service/auth_service"
)

type WebAuthHandler struct {
	authServ authservice.AuthenticationManager
}

func NewWebAuthHandler(service authservice.AuthenticationManager) *WebAuthHandler {
	return &WebAuthHandler{
		authServ: service,
	}
}

func (handler WebAuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}

	var req models.RegisterRequestDTO

	err = json.Unmarshal(data, &req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Println("Invalid request body", err)
		return
	}

	err = handler.authServ.Signup(req, roles.Customer)
	if err != nil {
		http.Error(w, "Failed to signup: "+err.Error(), http.StatusInternalServerError)
		log.Println("Failed to signup", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (handler WebAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
}
