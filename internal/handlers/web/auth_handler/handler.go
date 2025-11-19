package authhandler

import (
	"encoding/json"
	"io"
	"net/http"

	errorcodes "github.com/Kaushik1766/ParkingManagement/internal/constants/error_codes"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	authservice "github.com/Kaushik1766/ParkingManagement/internal/service/auth_service"
	customerrors "github.com/Kaushik1766/ParkingManagement/pkg/customErrors"
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
	w.Header().Set("Content-Type", "application/json")
	data, _ := io.ReadAll(r.Body)

	var req models.RegisterRequestDTO

	err := json.Unmarshal(data, &req)
	if err != nil {
		customerrors.SendError(w, customerrors.NewWebError(err, errorcodes.InvalidInput))
		return
	}

	err = handler.authServ.Signup(r.Context(), req, roles.Customer)
	if err != nil {
		customerrors.SendError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (handler WebAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data, _ := io.ReadAll(r.Body)

	var req models.LoginRequestDTO
	err := json.Unmarshal(data, &req)
	if err != nil {
		customerrors.SendError(w, customerrors.NewWebError(err, errorcodes.InvalidInput))
		return
	}

	token, err := handler.authServ.Login(r.Context(), req)
	if err != nil {
		customerrors.SendError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"jwt":"` + token + `"}`))
}
