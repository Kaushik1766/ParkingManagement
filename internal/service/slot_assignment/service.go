package slotassignment

import (
	"context"
	"errors"

	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	"github.com/Kaushik1766/ParkingManagement/internal/models/slot"
	userjwt "github.com/Kaushik1766/ParkingManagement/internal/models/user_jwt"
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
	ctxUser := ctx.Value(constants.User).(userjwt.UserJwt)

	uid, err := uuid.Parse(ctxUser.ID)
	if err != nil {
		return err
	}
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
		if val.VehicleType == vehicle.VehicleType && val.AssignedSlot.BuildingId != uuid.Nil {
			vehicle.AssignedSlot = val.AssignedSlot
			err = sas.vehicleRepo.Save(vehicle)
			if err != nil {
				return err
			}
			return nil
		}
	}

	userOffice, err := sas.officeRepo.GetOfficeByName(ctxUser.Office)
	if err != nil {
		return err
	}

	officeBuilding, err := sas.buildingRepo.GetBuildingByName(userOffice.BuildingName)
	if err != nil {
		return err
	}

	freeSlots, err := sas.slotRepo.GetFreeSlotsByFloor(officeBuilding.BuildingId, userOffice.FloorNumber)
	if err != nil {
		return err
	}

	if len(freeSlots) == 0 {
		return errors.New("no free slots available please contact admin")
	}

	vehicle.AssignedSlot = freeSlots[0]
	freeSlots[0].IsOccupied = true
	err = sas.vehicleRepo.Save(vehicle)
	if err != nil {
		return err
	}
	err = sas.slotRepo.Save(freeSlots[0])
	if err != nil {
		return err
	}
	return nil
}

func (sas *SlotAssignmentService) UnassignSlot(ctx context.Context, vehicleId string) error {
	ctxUser := ctx.Value(constants.User).(userjwt.UserJwt)

	uid, err := uuid.Parse(ctxUser.ID)
	if err != nil {
		return err
	}
	userVehicles, err := sas.vehicleRepo.GetVehiclesByUserId(uid)
	if err != nil {
		return err
	}

	vehicleUuid, err := uuid.Parse(vehicleId)
	if err != nil {
		return err
	}
	newVehicle, err := sas.vehicleRepo.GetVehicleById(vehicleUuid)
	if err != nil {
		return err
	}

	cnt := 0
	for _, val := range userVehicles {
		if val.AssignedSlot.SlotNumber == newVehicle.AssignedSlot.SlotNumber {
			cnt++
		}
	}

	if cnt == 1 {
		err = sas.slotRepo.Save(slot.Slot{
			BuildingId:  newVehicle.AssignedSlot.BuildingId,
			FloorNumber: newVehicle.AssignedSlot.FloorNumber,
			SlotNumber:  newVehicle.AssignedSlot.SlotNumber,
			IsOccupied:  false,
			SlotType:    newVehicle.AssignedSlot.SlotType,
		})
		if err != nil {
			return err
		}

	}

	newVehicle.AssignedSlot = slot.Slot{}
	err = sas.vehicleRepo.Save(newVehicle)
	if err != nil {
		return err
	}
	return nil
}

func (sas *SlotAssignmentService) AssignSlot(ctx context.Context, vehicleId string, slot slot.Slot) error {
	// ctxUser := ctx.Value(constants.User).(userjwt.UserJwt)

	vehicle, err := sas.vehicleRepo.GetVehicleById(uuid.MustParse(vehicleId))
	if err != nil {
		return err
	}

	userVehicles, err := sas.vehicleRepo.GetVehiclesByUserId(vehicle.UserId)
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
