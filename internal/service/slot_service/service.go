package slotservice

import (
	"context"
	"errors"

	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	"github.com/Kaushik1766/ParkingManagement/internal/models/slot"
	userjwt "github.com/Kaushik1766/ParkingManagement/internal/models/user_jwt"
	buildingrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/building_repository"
	floorrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/floor_repository"
	slotrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/slot_repository"
)

type SlotService struct {
	slotRepo     slotrepository.SlotStorage
	buildingRepo buildingrepository.BuildingStorage
	floorRepo    floorrepository.FloorStorage
}

func (ss *SlotService) GetSlotsByFloor(ctx context.Context, buildingName string, floorNumber int) ([]slot.Slot, error) {
	ctxUser := ctx.Value(constants.User).(userjwt.UserJwt)
	if ctxUser.Role != roles.Admin {
		return nil, errors.New("unauthorized: only admin or user can view slots")
	}

	building, err := ss.buildingRepo.GetBuildingByName(buildingName)
	if err != nil {
		return nil, err
	}

	_, err = ss.floorRepo.GetFloor(building.BuildingId, floorNumber)
	if err != nil {
		return nil, err
	}

	slots, err := ss.slotRepo.GetSlotsByFloor(building.BuildingId, floorNumber)
	if err != nil {
		return nil, err
	}

	return slots, nil
}

func NewSlotService(slotRepo slotrepository.SlotStorage, buildingRepo buildingrepository.BuildingStorage, floorRepo floorrepository.FloorStorage) *SlotService {
	return &SlotService{
		slotRepo:     slotRepo,
		buildingRepo: buildingRepo,
		floorRepo:    floorRepo,
	}
}

func (ss *SlotService) AddSlots(ctx context.Context, buildingName string, floorNumber int, slotNumbers []int, slotType vehicletypes.VehicleType) error {
	ctxUser := ctx.Value(constants.User).(userjwt.UserJwt)
	if ctxUser.Role != roles.Admin {
		return errors.New("unauthorized: only admin can add slots")
	}

	building, err := ss.buildingRepo.GetBuildingByName(buildingName)
	if err != nil {
		return err
	}

	_, err = ss.floorRepo.GetFloor(building.BuildingId, floorNumber)
	if err != nil {
		return err
	}

	for _, slotNumber := range slotNumbers {
		err = ss.slotRepo.AddSlot(building.BuildingId, floorNumber, slotNumber, slotType)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ss *SlotService) DeleteSlots(ctx context.Context, buildingName string, floorNumber int, slotNumbers []int) error {
	ctxUser := ctx.Value(constants.User).(userjwt.UserJwt)
	if ctxUser.Role != roles.Admin {
		return errors.New("unauthorized: only admin can delete slots")
	}

	building, err := ss.buildingRepo.GetBuildingByName(buildingName)
	if err != nil {
		return err
	}

	_, err = ss.floorRepo.GetFloor(building.BuildingId, floorNumber)
	if err != nil {
		return err
	}

	for _, slotNumber := range slotNumbers {
		err = ss.slotRepo.DeleteSlot(building.BuildingId, floorNumber, slotNumber)
		if err != nil {
			return err
		}
	}
	return nil
}
