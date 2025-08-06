package authhandler

import (
	"context"
	"fmt"
	"os"

	authenticationmiddleware "github.com/Kaushik1766/ParkingManagement/internal/middleware/authentication_middleware"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	authservice "github.com/Kaushik1766/ParkingManagement/internal/service/auth_service"
	"github.com/fatih/color"
	"golang.org/x/term"
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
	fmt.Print(color.YellowString("Email:"))
	var email, password string
	fmt.Scanln(&email)
	fmt.Print(color.GreenString("Password:"))
	passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return nil, fmt.Errorf("failed to read password: %w", err)
	}
	password = string(passwordBytes)
	token, err := auth.authMgr.Login(email, password)
	if err != nil {
		return nil, err
	}

	_ = os.WriteFile("token.txt", []byte(token), 0666)
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
	os.Remove("token.txt")
	return nil
}
