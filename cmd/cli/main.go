package main

import (
	"context"
	"fmt"

	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	authhandler "github.com/Kaushik1766/ParkingManagement/internal/handlers/auth_handler"
	userhandler "github.com/Kaushik1766/ParkingManagement/internal/handlers/user_handler"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	userjwt "github.com/Kaushik1766/ParkingManagement/internal/models/user_jwt"
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
	var ctx context.Context = nil
	for {
		if ctx == nil {
			clearScreen()
			color.Cyan("Welcome to Parking Management System")
			color.Yellow("1. Login")
			color.Yellow("2. Signup")
			color.Yellow("3. Exit")
			fmt.Scanf("%d", &choice)
			clearScreen()
			switch choice {
			case 1:
				ctx = context.Background()
				userCtx, err := authController.Login(ctx)
				if err != nil {
					color.Red("Error during login: %v", err)
					color.Cyan("Please try again or signup if you don't have an account.")
					ctx = nil
					fmt.Scanln()
				} else {
					ctx = userCtx
				}
			case 2:
				err := authController.CustomerSignup()
				fmt.Println(err)
			default:
				return
			}
		} else {
			user, ok := ctx.Value(constants.User).(userjwt.UserJwt)
			if !ok {
				continue
			}
			if user.Role == roles.Customer {
				clearScreen()
				color.Cyan("Enter your choice:")
				color.Yellow("1. Update profile")
				color.Yellow("2. Register vehicle")
				color.Yellow("3. View Profile")
				color.Yellow("4. View Registered Vehicles")
				color.Yellow("5. Logout")
				color.Yellow("6. Exit")
				fmt.Scanf("%d", &choice)
				clearScreen()
				switch choice {
				case 1:
					userHandler.UpdateProfile(ctx)
				case 2:
					userHandler.RegisterVehicle(ctx)
				case 3:
					userHandler.GetUserProfile(ctx)
				case 4:
					userHandler.GetRegisteredVehicles(ctx)
				case 5:
					ctx = authController.Logout()
				default:
					return
				}
			} else {
				fmt.Println("Admin functionalities are not implemented in CLI version yet.")
			}
		}
	}
}

func clearScreen() {
	// var cmd *exec.Cmd
	// switch runtime.GOOS {
	// case "windows":
	// 	cmd = exec.Command("cmd", "/c", "cls")
	// default:
	// 	cmd = exec.Command("clear")
	// }
	// cmd.Stdout = os.Stdout
	// cmd.Run()
}
