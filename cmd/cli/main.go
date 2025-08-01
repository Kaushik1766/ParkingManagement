package main

import (
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

func main() {
}
