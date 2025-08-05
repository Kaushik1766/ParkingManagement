package floorservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	userjwt "github.com/Kaushik1766/ParkingManagement/internal/models/user_jwt"
	buildingrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/building_repository"
	floorrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/floor_repository"
)

type FloorService struct {
	floorRepo    floorrepository.FloorStorage
	buildingRepo buildingrepository.BuildingStorage
}

func (fs *FloorService) DeleteFloors(ctx context.Context, buildingName string, floorNumbers []int) error {
	ctxUser := ctx.Value(constants.User).(userjwt.UserJwt)
	if ctxUser.Role != roles.Admin {
		return errors.New("unauthorized: only admin can delete floors")
	}
	building, err := fs.buildingRepo.GetBuildingByName(buildingName)
	if err != nil {
		return err
	}
	for _, floorNumber := range floorNumbers {
		err = fs.floorRepo.DeleteFloor(building.BuildingId, floorNumber)
		if err != nil {
			return err
		}
	}
	return nil
}

func (fs *FloorService) GetFloorsByBuildingId(ctx context.Context, buildingName string) ([]int, error) {
	ctxUser := ctx.Value(constants.User).(userjwt.UserJwt)
	if ctxUser.Role != roles.Admin {
		return nil, errors.New("unauthorized: only admin can view floors")
	}
	building, err := fs.buildingRepo.GetBuildingByName(buildingName)
	if err != nil {
		return nil, err
	}
	floors, err := fs.floorRepo.GetFloorsByBuildingId(building.BuildingId)
	if err != nil {
		return nil, err
	}
	return floors, nil
}

func (fs *FloorService) AddFloors(ctx context.Context, buildingName string, floorNumbers []int) error {
	ctxUser := ctx.Value(constants.User).(userjwt.UserJwt)
	if ctxUser.Role != roles.Admin {
		return errors.New("unauthorized: only admin can add floors")
	}
	building, err := fs.buildingRepo.GetBuildingByName(buildingName)
	if err != nil {
		return err
	}
	for _, floorNumber := range floorNumbers {
		err = fs.floorRepo.AddFloor(building.BuildingId, floorNumber)
		if err != nil {
			return fmt.Errorf("error adding floor %d to building %s: %w", floorNumber, building.BuildingName, err)
		}
	}
	return nil
}

func (fs *FloorService) AddFloor(ctx context.Context, buildingName string, floorNumber int) error {
	ctxUser := ctx.Value(constants.User).(userjwt.UserJwt)
	if ctxUser.Role != roles.Admin {
		return errors.New("unauthorized: only admin can add floors")
	}
	building, err := fs.buildingRepo.GetBuildingByName(buildingName)
	if err != nil {
		return err
	}
	return fs.floorRepo.AddFloor(building.BuildingId, floorNumber)
}

func (fs *FloorService) DeleteFloor(ctx context.Context, buildingName string, floorNumber int) error {
	ctxUser := ctx.Value(constants.User).(userjwt.UserJwt)
	if ctxUser.Role != roles.Admin {
		return errors.New("unauthorized: only admin can delete floors")
	}
	building, err := fs.buildingRepo.GetBuildingByName(buildingName)
	if err != nil {
		return err
	}
	return fs.floorRepo.DeleteFloor(building.BuildingId, floorNumber)
}

func NewFloorService(floorRepo floorrepository.FloorStorage, buildingRepo buildingrepository.BuildingStorage) *FloorService {
	return &FloorService{
		floorRepo:    floorRepo,
		buildingRepo: buildingRepo,
	}
}
