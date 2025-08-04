package userhandler

import (
	"context"
	"fmt"

	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	userservice "github.com/Kaushik1766/ParkingManagement/internal/service/user_service"
	"github.com/fatih/color"
)

type CliUserHandler struct {
	userService userservice.UserManager
}

func NewCliUserHandler(userService userservice.UserManager) *CliUserHandler {
	return &CliUserHandler{
		userService: userService,
	}
}

func (handler *CliUserHandler) UpdateProfile(userCtx context.Context) {
	var name, email, password string
	color.Cyan("Enter your new details to update profile:")
	color.Cyan("Name (leave blank to skip):")
	fmt.Scanln(&name)
	color.Yellow("Email (leave blank to skip):")
	fmt.Scanln(&email)
	color.Green("Password (leave blank to skip):")
	fmt.Scanln(&password)

	err := handler.userService.UpdateProfile(userCtx, name, email, password)
	if err != nil {
		color.Red("Failed to update profile: %v", err)
	}
}

func (handler *CliUserHandler) RegisterVehicle(userCtx context.Context) {
	var numberPlate string
	var vehicleType vehicletypes.VehicleType
	color.Cyan("Enter vehicle details:")
	color.Cyan("Number Plate:")
	fmt.Scanln(&numberPlate)
	color.Yellow("Vehicle Type (0 for Two Wheeler , 2 for Four Wheeler):")
	fmt.Scanln(&vehicleType)

	err := handler.userService.RegisterVehicle(userCtx, numberPlate, vehicleType)
	if err != nil {
		color.Red("Failed to register vehicle: %v", err)
	}
}
