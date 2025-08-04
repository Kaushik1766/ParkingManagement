package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	authhandler "github.com/Kaushik1766/ParkingManagement/internal/handlers/auth_handler"
	userhandler "github.com/Kaushik1766/ParkingManagement/internal/handlers/user_handler"
	userrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/user_repository"
	vehiclerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/vehicle_repository"
	authservice "github.com/Kaushik1766/ParkingManagement/internal/service/auth_service"
	userservice "github.com/Kaushik1766/ParkingManagement/internal/service/user_service"
	"github.com/fatih/color"
)

var (
	userDb         userrepository.UserStorage        = nil
	vehicleDb      vehiclerepository.VehicleStorage  = nil
	userService    userservice.UserManager           = nil
	authService    authservice.AuthenticationManager = nil
	authController *authhandler.CliAuthHandler       = nil
	userHandler    *userhandler.CliUserHandler       = nil
)

func init() {
	userDb = userrepository.NewFileUserRepository()
	vehicleDb = vehiclerepository.NewFileVehicleRepository(userDb)
	userService = userservice.NewUserService(userDb, vehicleDb)
	authService = authservice.NewAuthService(userDb)
	authController = authhandler.NewCliAuthHandler(authService)
	userHandler = userhandler.NewCliUserHandler(userService)
}

func cleanup() {
	userDb.(*userrepository.FileUserRepository).SerializeData()
	vehicleDb.(*vehiclerepository.FileVehicleRepository).SerializeData()
}

func main() {
	defer cleanup()
	var choice int
	token := ""
	for {
		clearScreen()
		if token == "" {
			color.Cyan("Welcome to Parking Management System")
			color.Yellow("1. Login")
			color.Yellow("2. Signup")
			color.Yellow("3. Exit")
			fmt.Scanf("%d", &choice)
			switch choice {
			case 1:
				clearScreen()
				jwtToken, err := authController.Login()
				if err != nil {
					color.Red("Error during login: %v", err)
				} else {
					token = jwtToken
				}
			case 2:
				clearScreen()
				err := authController.CustomerSignup()
				fmt.Println(err)
			default:
				break
			}
		} else {
			color.Cyan("Enter your choice:")
			color.Yellow("1. Update profile")
			color.Yellow("2. Register vehicle")
			color.Yellow("3. Exit")
			fmt.Scanf("%d", &choice)
			switch choice {
			case 1:
				clearScreen()
				userHandler.UpdateProfile(token)
			case 2:
				clearScreen()
				userHandler.RegisterVehicle(token)
			default:
				break
			}
		}
	}
}

func clearScreen() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}
