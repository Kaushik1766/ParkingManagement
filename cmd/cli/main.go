package main

import (
	"fmt"

	authhandler "github.com/Kaushik1766/ParkingManagement/internal/handlers/auth_handler"
	userrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/user_repository"
	authservice "github.com/Kaushik1766/ParkingManagement/internal/service/auth_service"
	userservice "github.com/Kaushik1766/ParkingManagement/internal/service/user_service"
)

var (
	userDb         userrepository.UserStorage        = nil
	userService    userservice.UserManager           = nil
	authService    authservice.AuthenticationManager = nil
	authController *authhandler.CliAuthHandler       = nil
)

func init() {
	userDb = userrepository.NewFileUserRepository()
	userService = userservice.NewUserService(userDb)
	authService = authservice.NewAuthService(userDb)
	authController = authhandler.NewCliAuthHandler(authService)
}

func cleanup() {
	userDb.(*userrepository.FileUserRepository).SerializeData()
}

func main() {
	defer cleanup()
	err := authController.CustomerSignup("kaushik", "kaushik@gmail.com", "123")
	if err != nil {
		fmt.Println("Error during signup:", err)
	}

	token, err := authController.Login("kaushik@gmail.com", "123")
	if err != nil {
		fmt.Println("Error during login:", err)
	}
	fmt.Println(token)
}
