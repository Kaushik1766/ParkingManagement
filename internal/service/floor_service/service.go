package floorservice

import (
	"context"
	"errors"

	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	models "github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	floorrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/floor_repository"
)

type FloorService struct {
	floorRepo floorrepository.FloorStorage
}

func NewFloorService(floorRepo floorrepository.FloorStorage) *FloorService {
	return &FloorService{
		floorRepo: floorRepo,
	}
}

func (fs *FloorService) AddFloorByBuildingId(ctx context.Context, buildingId string, floorNumber int) error {
	ctxUser := ctx.Value(constants.User).(models.UserJwt)
	if ctxUser.Role != roles.Admin {
		return errors.New("unauthorized: only admin can add floors")
	}

	return fs.floorRepo.AddFloor(ctx, buildingId, floorNumber)
}

func (fs *FloorService) GetFloorsByBuildingId(ctx context.Context, buildingId string) ([]models.FloorDTO, error) {
	ctxUser := ctx.Value(constants.User).(models.UserJwt)
	if ctxUser.Role != roles.Admin {
		return []models.FloorDTO{}, errors.New("unauthorized: only admin can view floors")
	}

	floors, err := fs.floorRepo.GetFloorsByBuildingId(ctx, buildingId)
	if err != nil {
		return []models.FloorDTO{}, err
	}

	floorsDTO := []models.FloorDTO{}
	for _, floor := range floors {
		floorsDTO = append(floorsDTO, models.FloorDTO{
			BuildingID:     floor.BuildingID.String(),
			FloorNumber:    floor.FloorNumber,
			TotalSlots:     floor.TotalSlots,
			AvailableSlots: floor.AvailableSlots,
			AssignedOffice: floor.AssignedOffice,
		})
	}

	return floorsDTO, nil
}

func (fs *FloorService) AddFloor(ctx context.Context, buildingId string, floorNumber int) error {
	ctxUser := ctx.Value(constants.User).(models.UserJwt)
	if ctxUser.Role != roles.Admin {
		return errors.New("unauthorized: only admin can add floors")
	}
	return fs.floorRepo.AddFloor(ctx, buildingId, floorNumber)
}

func (fs *FloorService) DeleteFloor(ctx context.Context, buildingId string, floorNumber int) error {
	ctxUser := ctx.Value(constants.User).(models.UserJwt)
	if ctxUser.Role != roles.Admin {
		return errors.New("unauthorized: only admin can delete floors")
	}

	return fs.floorRepo.DeleteFloor(ctx, buildingId, floorNumber)
}
