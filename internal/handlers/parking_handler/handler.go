package parkinghandler

import (
	"context"
	"fmt"
	"time"

	parkinghistoryservice "github.com/Kaushik1766/ParkingManagement/internal/service/parking_history_service"
	userservice "github.com/Kaushik1766/ParkingManagement/internal/service/user_service"
	vehicleservice "github.com/Kaushik1766/ParkingManagement/internal/service/vehicle_service"
	customerrors "github.com/Kaushik1766/ParkingManagement/pkg/customErrors"
	"github.com/fatih/color"
)

type CliParkingHandler struct {
	vehicleService        vehicleservice.VehicleMgr
	userService           userservice.UserManager
	parkingHistoryService parkinghistoryservice.ParkingHistoryMgr
}

func (handler *CliParkingHandler) Park(ctx context.Context) {
	color.Cyan("Select the vehicle to park:")

	vehicles := handler.userService.GetRegisteredVehicles(ctx)

	for i, vehicle := range vehicles {
		color.Magenta("%d. Number Plate: %s, Vehicle Type: %s\n", i+1, vehicle.NumberPlate, vehicle.VehicleType)
	}

	var choice int
	fmt.Scanln(&choice)

	if choice < 1 || choice > len(vehicles) {
		color.Red("Invalid choice. Please try again.")
		return
	}

	selectedVehicle := vehicles[choice-1]
	ticket, err := handler.vehicleService.Park(ctx, selectedVehicle.NumberPlate)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Failed to park vehicle: %v", err))
		return
	}

	color.Green("Vehicle parked successfully. Ticket ID: %s", ticket)
	color.Green("Press Enter to continue...")
	fmt.Scanln()
}

func (handler *CliParkingHandler) Unpark(ctx context.Context) {
	activeParkings, err := handler.parkingHistoryService.GetActiveUserParkings(ctx)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Failed to fetch active parkings: %v", err))
		return
	}

	color.Cyan("Select the parking to unpark:")

	for i, parking := range activeParkings {
		color.Magenta("%d. Ticket ID: %s, Vehicle Number Plate: %s, StartTime: %s", i+1, parking.TicketId, parking.NumberPlate, parking.StartTime)
	}

	var choice int
	fmt.Scanln(&choice)

	if choice < 1 || choice > len(activeParkings) {
		customerrors.DisplayError("Invalid choice. Please try again.")
		return
	}
	selectedParking := activeParkings[choice-1]

	err = handler.vehicleService.Unpark(ctx, selectedParking.TicketId)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Failed to unpark vehicle: %v", err))
		return
	}
	color.Green("Vehicle unparked successfully. Ticket ID: %s", selectedParking.TicketId)
	color.Green("Press Enter to continue...")
	fmt.Scanln()
}

func (handler *CliParkingHandler) ViewParkingHistory(ctx context.Context) {
	startDate := time.Now().AddDate(0, -1, 0).Format(time.DateOnly)
	endDate := time.Now().Format(time.DateOnly)
	history, err := handler.parkingHistoryService.GetParkingHistory(ctx, startDate, endDate)
	if err != nil {
		customerrors.DisplayError(fmt.Sprintf("Failed to fetch parking history: %v", err))
		return
	}

	color.Cyan("Parking History:")
	for _, record := range history {
		color.Magenta("Ticket ID: %s, Vehicle Number Plate: %s, Start Time: %s, End Time: %s",
			record.TicketId, record.NumberPlate, record.StartTime, record.EndTime)
	}

	color.Green("Press Enter to continue...")
	fmt.Scanln()
}

func NewCliParkingHandler(vehicleService vehicleservice.VehicleMgr, userService userservice.UserManager, parkingHistoryService parkinghistoryservice.ParkingHistoryMgr) *CliParkingHandler {
	return &CliParkingHandler{
		vehicleService:        vehicleService,
		userService:           userService,
		parkingHistoryService: parkingHistoryService,
	}
}
