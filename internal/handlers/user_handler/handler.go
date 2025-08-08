package userhandler

import (
	"context"
	"fmt"
	"strings"

	"github.com/Kaushik1766/ParkingManagement/internal/constants/menuconstants"
	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	userservice "github.com/Kaushik1766/ParkingManagement/internal/service/user_service"
	customerrors "github.com/Kaushik1766/ParkingManagement/pkg/customErrors"
	"github.com/Kaushik1766/ParkingManagement/utils"
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
	var name, email, password, office string
	color.Cyan("Enter your new details to update profile:")
	color.Cyan("Name (leave blank to skip):")
	fmt.Scanln(&name)
	color.Yellow("Email (leave blank to skip):")
	fmt.Scanln(&email)
	color.Green("Password (leave blank to skip):")
	fmt.Scanln(&password)
	color.Magenta("Office (leave blank to skip):")
	fmt.Scanln(&office)

	err := handler.userService.UpdateProfile(userCtx, name, email, password, office)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Failed to update profile: %v", err))
		return
	}
	color.Green("Profile updated successfully")
	color.Green("Press Enter to continue...")
	fmt.Scanln()
}

func (handler *CliUserHandler) RegisterVehicle(userCtx context.Context) {
	var numberPlate string
	var vehicleType vehicletypes.VehicleType
	color.Cyan("Enter vehicle details:")
	color.Cyan("Number Plate:")
	fmt.Scanln(&numberPlate)
	color.Yellow("Vehicle Type (0 for Two Wheeler , 1 for Four Wheeler):")
	fmt.Scanln(&vehicleType)

	err := handler.userService.RegisterVehicle(userCtx, numberPlate, vehicleType)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Failed to register vehicle: %v", err))
		return
	}
	color.Green("Vehicle registered successfully")
	color.Green("Press Enter to continue...")
	fmt.Scanln()
}

func (handler *CliUserHandler) UnregisterVehicle(userCtx context.Context) {
	registeredVehicles := handler.userService.GetRegisteredVehicles(userCtx)

	color.Cyan("Enter number of the vehicle to unregister")
	var vehiclesStr []string
	for _, val := range registeredVehicles {
		vehiclesStr = append(vehiclesStr, val.String())
	}

	utils.PrintListInRows(vehiclesStr)

	var vehicleNumber int
	fmt.Scanf("%d", &vehicleNumber)
	numberPlate := registeredVehicles[vehicleNumber-1].NumberPlate
	err := handler.userService.UnregisterVehicle(userCtx, numberPlate)
	if err != nil {
		customerrors.DisplayError("Failed to unregister vehicle")
		return
	}
	color.Green("Unregistered successfully")
	color.Green(menuconstants.PressEnterToContinue)
	fmt.Scanln()
}

func (handler *CliUserHandler) GetUserProfile(userCtx context.Context) {
	userProfile, err := handler.userService.GetUserProfile(userCtx)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Failed to fetch user profile: %v", err))
		return
	}

	var userName string

	userNameList := strings.Split(userProfile.Name, " ")

	for _, name := range userNameList {
		name = strings.ToUpper(name[:1]) + name[1:]
		userName += string(name)
	}

	color.Cyan("User Profile:")
	color.Yellow("User ID: %s", userProfile.UserId)
	color.Yellow("Name: %s", userName)
	color.Yellow("Email: %s", userProfile.Email)
	color.Yellow("Role: %s", userProfile.Role)
	color.Yellow("Office: %s", userProfile.Office)
	color.Green(menuconstants.PressEnterToContinue)
	fmt.Scanln()
}

func (handler *CliUserHandler) GetRegisteredVehicles(userCtx context.Context) {
	vehicles := handler.userService.GetRegisteredVehicles(userCtx)
	if len(vehicles) == 0 {
		customerrors.DisplayError("no vehicles registered")
		return
	}
	color.Cyan("Registered Vehicles:")
	for _, val := range vehicles {
		color.Yellow("Number Plate: %s", val.NumberPlate)
		color.Yellow("Vehicle Type: %s", val.VehicleType)
		color.Yellow("Assigned Slot: %s", val.AssignedSlot)
		fmt.Println()
	}
	color.Green(menuconstants.PressEnterToContinue)
	fmt.Scanln()
}
