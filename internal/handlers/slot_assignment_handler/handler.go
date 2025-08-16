package slotassignmenthandler

import (
	"context"
	"fmt"

	"github.com/Kaushik1766/ParkingManagement/internal/constants/menuconstants"
	officeservice "github.com/Kaushik1766/ParkingManagement/internal/service/office_service"
	slotassignment "github.com/Kaushik1766/ParkingManagement/internal/service/slot_assignment"
	slotservice "github.com/Kaushik1766/ParkingManagement/internal/service/slot_service"
	userservice "github.com/Kaushik1766/ParkingManagement/internal/service/user_service"
	customerrors "github.com/Kaushik1766/ParkingManagement/pkg/customErrors"
	"github.com/fatih/color"
)

type CliSlotAssignmentHandler struct {
	assignmentService slotassignment.SlotAssignmentMgr
	userService       userservice.UserManager
	slotService       slotservice.SlotMgr
	officeService     officeservice.OfficeMgr
}

func NewCliSlotAssignmentHandler(
	assignmentService slotassignment.SlotAssignmentMgr,
	userService userservice.UserManager,
	slotService slotservice.SlotMgr,
	officeService officeservice.OfficeMgr,
) *CliSlotAssignmentHandler {
	return &CliSlotAssignmentHandler{
		assignmentService: assignmentService,
		userService:       userService,
		slotService:       slotService,
		officeService:     officeService,
	}
}

func (sah *CliSlotAssignmentHandler) ViewVehiclesWithUnassignedSlots(ctx context.Context) {
	unassignedVehicles, err := sah.assignmentService.GetVehiclesWithUnassignedSlots(ctx)
	if err != nil {
		customerrors.DisplayError("error fetching vehicles")
		return
	}

	for i, vehicle := range unassignedVehicles {
		color.Magenta("%d. UserId: %s, NumberPlate: %s, Vehicle Type: %s\n", i, vehicle.UserID.String(), vehicle.NumberPlate, vehicle.VehicleType)
	}
	color.Green(menuconstants.PressEnterToContinue)
	fmt.Scanln()
}

func (sah *CliSlotAssignmentHandler) AssignSlot(ctx context.Context) {
	color.Cyan(menuconstants.SelectVehicleForSlot)
	unassignedVehicles, err := sah.assignmentService.GetVehiclesWithUnassignedSlots(ctx)
	if err != nil {
		customerrors.DisplayError("error fetching vehicles")
		return
	}

	for i, vehicle := range unassignedVehicles {
		color.Magenta("%d. UserId: %s, NumberPlate: %s, Vehicle Type: %s\n", i+1, vehicle.UserID, vehicle.NumberPlate, vehicle.VehicleType)
	}

	var vehicleNumber int
	fmt.Scanf("%d", &vehicleNumber)

	if vehicleNumber < 1 || vehicleNumber > len(unassignedVehicles) {
		customerrors.DisplayError("invalid vehicle number selected")
		return
	}

	vehicle := unassignedVehicles[vehicleNumber-1]

	user, err := sah.userService.GetUserById(ctx, vehicle.UserID.String())
	if err != nil {
		customerrors.DisplayError("error fetching user details")
		return
	}

	office, err := sah.officeService.GetOfficeByName(ctx, user.Office)
	if err != nil {
		customerrors.DisplayError("error fetching office details")
		return
	}

	freeSlots, err := sah.slotService.GetFreeSlotsByBuilding(ctx, office.BuildingID, vehicle.VehicleType)
	if err != nil {
		customerrors.DisplayError("error fetching free slots")
		return
	}

	color.Cyan(menuconstants.SelectSlotForVehicle)

	for i, slot := range freeSlots {
		color.Blue("%d. Floor Number: %d, Slot Number: %d, Slot Type: %s\n", i+1, slot.FloorNumber, slot.SlotNumber, slot.SlotType)
	}

	var slotNumber int
	fmt.Scanf("%d", &slotNumber)

	if slotNumber < 1 || slotNumber > len(freeSlots) {
		customerrors.DisplayError("invalid slot number selected")
		return
	}
	slot := freeSlots[slotNumber-1]
	if err := sah.assignmentService.AssignSlot(ctx, vehicle.VehicleID.String(), slot); err != nil {
		customerrors.DisplayError("error assigning slot to vehicle")
		return
	}
	color.Green(menuconstants.SlotAssignedSuccess, vehicle.NumberPlate)
	color.Green(menuconstants.PressEnterToContinue)
	fmt.Scanln()
}
