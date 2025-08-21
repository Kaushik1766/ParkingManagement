package authhandler

import authservice "github.com/Kaushik1766/ParkingManagement/internal/service/auth_service"

type WebAuthHandler struct {
	authServ authservice.AuthenticationManager
}

func NewWebAuthHandler(service authservice.AuthenticationManager) *WebAuthHandler {
	return &WebAuthHandler{
		authServ: service,
	}
}
