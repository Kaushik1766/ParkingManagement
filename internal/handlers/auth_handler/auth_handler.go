package authhandler

import (
	"context"
	"fmt"
	"os"

	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	"github.com/Kaushik1766/ParkingManagement/internal/constants/menuconstants"
	authenticationmiddleware "github.com/Kaushik1766/ParkingManagement/internal/middleware/authentication_middleware"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	authservice "github.com/Kaushik1766/ParkingManagement/internal/service/auth_service"
	officeservice "github.com/Kaushik1766/ParkingManagement/internal/service/office_service"
	customerrors "github.com/Kaushik1766/ParkingManagement/pkg/customErrors"
	"github.com/Kaushik1766/ParkingManagement/utils"
	"github.com/common-nighthawk/go-figure"
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
	(figure.NewColorFigure("Login", "", "green", true)).Print()
	color.Cyan("Enter your credentials to login:")
	fmt.Print(color.GreenString("Email: "))
	var email, password string
	fmt.Scanln(&email)
	fmt.Print(color.GreenString("Password: "))
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
	(figure.NewColorFigure("Signup", "", "green", true)).Print()
	var name, email, password string
	color.Cyan("Enter your details to signup:")
	fmt.Print(color.GreenString("Name: "))
	fmt.Scanln(&name)
	fmt.Print(color.GreenString("Email: "))
	fmt.Scanln(&email)
	fmt.Print(color.GreenString("Password: "))
	passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		customerrors.DisplayError("Failed to read password: " + err.Error())
		return
	}
	password = string(passwordBytes)

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
	color.Green(menuconstants.PressEnterToContinue)
	fmt.Scanln()
}

func (auth *CliAuthHandler) AdminSignup() error {
	(figure.NewColorFigure("sudo Signup", "", "green", true)).Print()
	var name, email, password string
	color.Cyan("Enter your details to signup as an admin:")
	fmt.Print(color.GreenString("Name: "))
	fmt.Scanln(&name)
	fmt.Print(color.GreenString("Email: "))
	fmt.Scanln(&email)
	fmt.Print(color.GreenString("Password: "))
	fmt.Scanln(&password)
	authErr := auth.authMgr.Signup(name, email, password, constants.AdminOffice, roles.Admin)
	return authErr
}

func (auth *CliAuthHandler) Logout() context.Context {
	os.Remove("token.txt")
	return nil
}
