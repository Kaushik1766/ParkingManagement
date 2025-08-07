package authhandler

import (
	"context"
	"fmt"
	"os"

	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	authenticationmiddleware "github.com/Kaushik1766/ParkingManagement/internal/middleware/authentication_middleware"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	authservice "github.com/Kaushik1766/ParkingManagement/internal/service/auth_service"
	officeservice "github.com/Kaushik1766/ParkingManagement/internal/service/office_service"
	customerrors "github.com/Kaushik1766/ParkingManagement/pkg/customErrors"
	"github.com/Kaushik1766/ParkingManagement/utils"
	"github.com/fatih/color"
	"golang.org/x/term"
)

type CliAuthHandler struct {
	authMgr   authservice.AuthenticationManager
	officeMgr officeservice.OfficeMgr
}

func NewCliAuthHandler(authMgr authservice.AuthenticationManager, officeMgr officeservice.OfficeMgr) *CliAuthHandler {
	return &CliAuthHandler{
		authMgr:   authMgr,
		officeMgr: officeMgr,
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
		customerrors.DisplayError("Login failed: " + err.Error())
		return nil, err
	}

	_ = os.WriteFile("token.txt", []byte(token), 0666)
	userCtx, err := authenticationmiddleware.CliAuthenticate(baseCtx, token)

	return userCtx, err
}

func (auth *CliAuthHandler) CustomerSignup() {
	var name, email, password string
	color.Cyan("Enter your details to signup:")
	color.Cyan("Name:")
	fmt.Scanln(&name)
	color.Yellow("Email:")
	fmt.Scanln(&email)
	color.Green("Password:")
	fmt.Scanln(&password)

	color.Cyan("Select office number:")
	offices, err := auth.officeMgr.GetAllOfficeNames(context.Background())
	if err != nil {
		customerrors.DisplayError(err.Error())
		return
	}

	utils.PrintListInRows(offices)

	var officeNumber int
	fmt.Scanln(&officeNumber)

	if officeNumber < 1 || officeNumber > len(offices) {
		customerrors.DisplayError("Invalid office number selected.")
		return
	}

	officeName := offices[officeNumber-1]

	authErr := auth.authMgr.Signup(name, email, password, officeName, roles.Customer)
	if authErr != nil {
		customerrors.DisplayError(authErr.Error())
		return
	}
	color.Green("Signup successful")
	color.Green("Press Enter to continue...")
	fmt.Scanln()
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
	authErr := auth.authMgr.Signup(name, email, password, constants.AdminOffice, roles.Admin)
	return authErr
}

func (auth *CliAuthHandler) Logout() context.Context {
	os.Remove("token.txt")
	return nil
}
