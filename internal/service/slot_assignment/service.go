package slotassignment

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	models "github.com/Kaushik1766/ParkingManagement/internal/models"
	buildingrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/building_repository"
	floorrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/floor_repository"
	officerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/office_repository"
	slotrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/slot_repository"
	vehiclerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/vehicle_repository"
	"github.com/google/uuid"
)

type SlotAssignmentService struct {
	vehicleRepo  vehiclerepository.VehicleStorage
	floorRepo    floorrepository.FloorStorage
	buildingRepo buildingrepository.BuildingStorage
	slotRepo     slotrepository.SlotStorage
	officeRepo   officerepository.OfficeStorage
}

func NewSlotAssignmentService(
	vehicleRepo vehiclerepository.VehicleStorage,
	floorRepo floorrepository.FloorStorage,
	buildingRepo buildingrepository.BuildingStorage,
	slotRepo slotrepository.SlotStorage,
	officeRepo officerepository.OfficeStorage,
) *SlotAssignmentService {
	return &SlotAssignmentService{
		vehicleRepo:  vehicleRepo,
		floorRepo:    floorRepo,
		buildingRepo: buildingRepo,
		slotRepo:     slotRepo,
		officeRepo:   officeRepo,
	}
}

func (sas *SlotAssignmentService) AutoAssignSlot(ctx context.Context, vehicleId string) error {
	ctxUser := ctx.Value(constants.User).(models.UserJwt)

	vehicleUuid, err := uuid.Parse(vehicleId)
	if err != nil {
		return err
	}
	vehicle, err := sas.vehicleRepo.GetVehicleById(ctx, vehicleUuid)
	if err != nil {
		return err
	}

	// Check if vehicle already has a slot assigned (from reusing another vehicle's slot)
	if vehicle.AssignedBuildingID != uuid.Nil {
		log.Println("Vehicle already has slot assigned:", vehicle.AssignedBuildingID, vehicle.AssignedFloorNumber, vehicle.AssignedSlotNumber)
		return nil
	}

	// Vehicle doesn't have a slot, need to assign a new one
	// This means it's the first vehicle of this type for the user
	userOffice, err := sas.officeRepo.GetOfficeByName(ctx, ctxUser.Office)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	freeSlots, err := sas.slotRepo.GetFreeSlotsByFloor(ctx, userOffice.BuildingID, userOffice.FloorNumber)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	if len(freeSlots) == 0 {
		return errors.New("no free slots available please contact admin")
	}

	fmt.Println(freeSlots)

	for _, val := range freeSlots {
		if val.SlotType == vehicle.VehicleType {
			vehicle.AssignedBuildingID = val.BuildingID
			vehicle.AssignedFloorNumber = val.FloorNumber
			vehicle.AssignedSlotNumber = val.SlotNumber
			log.Println("free slot found, assigning it")
			err := sas.vehicleRepo.Save(ctx, vehicle)
			if err != nil {
				return err
			}

			return nil
		}
	}

	return errors.New("no free slot available please contact the admin")
}

func (sas *SlotAssignmentService) AssignSlot(ctx context.Context, vehicleId string, slot models.Slot) error {
	// ctxUser := ctx.Value(constants.User).(models.UserJwt)

	vehicle, err := sas.vehicleRepo.GetVehicleById(ctx, uuid.MustParse(vehicleId))
	if err != nil {
		return err
	}

	userVehicles, err := sas.vehicleRepo.GetVehiclesByUserId(ctx, vehicle.UserID)
	if err != nil {
		return err
	}

	for i, val := range userVehicles {
		if val.VehicleType == vehicle.VehicleType {
			userVehicles[i].AssignedSlot = slot
			sas.vehicleRepo.Save(ctx, userVehicles[i])
		}
	}
	return sas.slotRepo.Save(ctx, slot)
}
