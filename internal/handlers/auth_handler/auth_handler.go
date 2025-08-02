package authhandler

import (
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums"
	authservice "github.com/Kaushik1766/ParkingManagement/internal/service/auth_service"
)

type CliAuthHandler struct {
	authMgr authservice.AuthenticationManager
}

func NewCliAuthHandler(authMgr authservice.AuthenticationManager) *CliAuthHandler {
	return &CliAuthHandler{
		authMgr: authMgr,
	}
}

func (auth *CliAuthHandler) Login(email string, password string) (string, error) {
	token, err := auth.authMgr.Login(email, password)
	if err != nil {
		return token, err
	}

	return token, nil
}

func (auth *CliAuthHandler) CustomerSignup(name string, email string, password string) error {
	authErr := auth.authMgr.Signup(name, email, password, enums.Customer)
	return authErr
}

func (auth *CliAuthHandler) AdminSignup(name string, email string, password string) error {
	authErr := auth.authMgr.Signup(name, email, password, enums.Admin)
	return authErr
}
