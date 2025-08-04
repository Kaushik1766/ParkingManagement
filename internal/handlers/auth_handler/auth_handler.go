package authhandler

import (
	"context"
	"fmt"

	authenticationmiddleware "github.com/Kaushik1766/ParkingManagement/internal/middleware/authentication_middleware"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	authservice "github.com/Kaushik1766/ParkingManagement/internal/service/auth_service"
	"github.com/fatih/color"
)

type CliAuthHandler struct {
	authMgr authservice.AuthenticationManager
}

func NewCliAuthHandler(authMgr authservice.AuthenticationManager) *CliAuthHandler {
	return &CliAuthHandler{
		authMgr: authMgr,
	}
}

func (auth *CliAuthHandler) Login(baseCtx context.Context) (context.Context, error) {
	color.Cyan("Enter your credentials to login:")
	color.Yellow("Email:")
	var email, password string
	fmt.Scanln(&email)
	color.Green("Password:")
	fmt.Scanln(&password)
	token, err := auth.authMgr.Login(email, password)
	if err != nil {
		return nil, err
	}
	userCtx, err := authenticationmiddleware.CliAuthenticate(baseCtx, token)

	return userCtx, err
}

func (auth *CliAuthHandler) CustomerSignup() error {
	var name, email, password string
	color.Cyan("Enter your details to signup:")
	color.Cyan("Name:")
	fmt.Scanln(&name)
	color.Yellow("Email:")
	fmt.Scanln(&email)
	color.Green("Password:")
	fmt.Scanln(&password)
	authErr := auth.authMgr.Signup(name, email, password, roles.Customer)
	return authErr
}

func (auth *CliAuthHandler) AdminSignup() error {
	var name, email, password string
	color.Cyan("Enter your details to signup as an admin:")
	color.Cyan("Name:")
	fmt.Scanln(&name)
	color.Yellow("Email:")
	fmt.Scanln(&email)
	color.Green("Password:")
	fmt.Scanln(&password)
	authErr := auth.authMgr.Signup(name, email, password, roles.Admin)
	return authErr
}

func (auth *CliAuthHandler) Logout() context.Context {
	return nil
}
