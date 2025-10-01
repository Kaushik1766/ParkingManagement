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

	uid := uuid.MustParse(ctxUser.ID)
	userVehicles, err := sas.vehicleRepo.GetVehiclesByUserId(uid)
	if err != nil {
		return err
	}

	vehicleUuid, err := uuid.Parse(vehicleId)
	if err != nil {
		return err
	}
	vehicle, err := sas.vehicleRepo.GetVehicleById(vehicleUuid)
	if err != nil {
		return err
	}

	for _, val := range userVehicles {
		// if found update and exit
		if val.VehicleType == vehicle.VehicleType && val.AssignedBuildingID != uuid.Nil {
			// vehicle.AssignedSlot = val.AssignedSlot
			log.Println("slot repeating")
			vehicle.AssignedBuildingID = val.AssignedBuildingID
			vehicle.AssignedFloorNumber = val.AssignedFloorNumber
			vehicle.AssignedSlotNumber = val.AssignedSlotNumber
			err = sas.vehicleRepo.Save(vehicle)
			if err != nil {
				return err
			}
			return nil
		}
	}

	// if first vehicle of type
	userOffice, err := sas.officeRepo.GetOfficeByName(ctxUser.Office)
	if err != nil {
		return err
	}

	freeSlots, err := sas.slotRepo.GetFreeSlotsByFloor(userOffice.BuildingID, userOffice.FloorNumber)
	if err != nil {
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
			err := sas.vehicleRepo.Save(vehicle)
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

	vehicle, err := sas.vehicleRepo.GetVehicleById(uuid.MustParse(vehicleId))
	if err != nil {
		return err
	}

	userVehicles, err := sas.vehicleRepo.GetVehiclesByUserId(vehicle.UserID)
	if err != nil {
		return err
	}

	for i, val := range userVehicles {
		if val.VehicleType == vehicle.VehicleType {
			userVehicles[i].AssignedSlot = slot
			sas.vehicleRepo.Save(userVehicles[i])
		}
	}
	return sas.slotRepo.Save(slot)
}
