package main

import (
	authhandler "github.com/Kaushik1766/ParkingManagement/internal/handlers/auth_handler"
	userhandler "github.com/Kaushik1766/ParkingManagement/internal/handlers/user_handler"
	userrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/user_repository"
	vehiclerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/vehicle_repository"
	authservice "github.com/Kaushik1766/ParkingManagement/internal/service/auth_service"
	userservice "github.com/Kaushik1766/ParkingManagement/internal/service/user_service"
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
	// err := authController.CustomerSignup("kaushik", "kaushik@gmail.com", "123")
	// err := authController.CustomerSignup()
	// if err != nil {
	// 	fmt.Println("Error during signup:", err)
	// }

	// token, err := authController.Login("kaushik@gmail.com", "123")
	// token, err := authController.Login()
	// if err != nil {
	// 	fmt.Println("Error during login:", err)
	// }
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTQzMzA4MjMsImlhdCI6MTc1NDI0NDQyMywiaWQiOiIyY2Q2YzE0Mi0wODcxLTQ1Y2YtYTQwNy1hMGQ0OGViMWNiZDMiLCJlbWFpbCI6ImthdXNoaWsxQGdtYWlsLmNvbSIsInJvbGUiOjB9.Ep7xG2GgdzNCaD5KO_w0HCsedxPOxvW67CUpbiIaQPw"
	// fmt.Println(token)
	// userHandler.UpdateProfile(token)
	userHandler.RegisterVehicle(token)
}
