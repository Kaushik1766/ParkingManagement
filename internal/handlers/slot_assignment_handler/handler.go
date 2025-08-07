package slotassignmenthandler

import (
	"context"
	"fmt"

	slotassignment "github.com/Kaushik1766/ParkingManagement/internal/service/slot_assignment"
	customerrors "github.com/Kaushik1766/ParkingManagement/pkg/customErrors"
	"github.com/fatih/color"
)

type CliSlotAssignmentHandler struct {
	assignmentService slotassignment.SlotAssignmentMgr
}

func (sah *CliSlotAssignmentHandler) ViewVehiclesWithUnassignedSlots(ctx context.Context) {
	unassignedVehicles, err := sah.assignmentService.GetVehiclesWithUnassignedSlots(ctx)
	if err != nil {
		customerrors.DisplayError("error fetching vehicles")
	}

	for i, vehicle := range unassignedVehicles {
		fmt.Printf("%d. UserId: %s, NumberPlate: %s, Vehicle Type: %s\n", i, vehicle.UserId.String(), vehicle.NumberPlate, vehicle.VehicleType)
	}
	color.Green("Press enter to continue...")
	fmt.Scanln()
}

// func (sah *CliSlotAssignmentHandler) AssignSlot(ctx context.Context) {
// 	color.Cyan("Select the number of vehicle to assign a slot to:")
// 	unassignedVehicles, err := sah.assignmentService.GetVehiclesWithUnassignedSlots(ctx)
// 	if err != nil {
// 		customerrors.DisplayError("error fetching vehicles")
// 	}
//
// 	for i, vehicle := range unassignedVehicles {
// 		fmt.Printf("%d. UserId: %s, NumberPlate: %s, Vehicle Type: %s\n", i, vehicle.UserId.String(), vehicle.NumberPlate, vehicle.VehicleType)
// 	}
//
// 	var vehicleNumber int
// 	fmt.Scanf("%d", &vehicleNumber)
//
// 	vehicle := unassignedVehicles[vehicleNumber-1]
//
// 	// err:= sah.assignmentService.AssignSlot(ctx, vehicle.VehicleId.String(), slot slot.Slot)
// }
