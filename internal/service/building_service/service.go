package buildingservice

import (
	"context"
	"errors"

	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	models "github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	buildingrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/building_repository"
	"github.com/google/uuid"
)

type BuildingService struct {
	buildingRepo buildingrepository.BuildingStorage
}

func (bs *BuildingService) DeleteBuildingByID(ctx context.Context, buildingID string) error {
	ctxUser := ctx.Value(constants.User).(models.UserJwt)
	if ctxUser.Role != roles.Admin {
		return errors.New("unauthorized: only admin can delete buildings")
	}

	err := bs.buildingRepo.DeleteBuildingByID(ctx, buildingID)
	if err != nil {
		return err
	}
	return nil
}

func (bs *BuildingService) GetBuildingByID(ctx context.Context, buildingID string) (models.BuildingDTO, error) {
	ctxUser := ctx.Value(constants.User).(models.UserJwt)
	if ctxUser.Role != roles.Admin {
		return models.BuildingDTO{}, errors.New("unauthorized: only admin can view buildings")
	}

	buildingUUID, err := uuid.Parse(buildingID)
	if err != nil {
		return models.BuildingDTO{}, err
	}

	building, err := bs.buildingRepo.GetBuildingByID(ctx, buildingUUID)
	if err != nil {
		return models.BuildingDTO{}, err
	}
	return models.BuildingDTO{
		BuildingID: building.BuildingID.String(),
		Name:       building.BuildingName,
	}, nil
}

func (bs *BuildingService) GetAllBuildings(ctx context.Context) ([]models.BuildingDTO, error) {
	ctxUser := ctx.Value(constants.User).(models.UserJwt)
	if ctxUser.Role != roles.Admin {
		return nil, errors.New("unauthorized: only admin can view buildings")
	}
	buildings, err := bs.buildingRepo.GetAllBuildingSummary(ctx)
	if err != nil {
		return nil, err
	}
	var res []models.BuildingDTO
	for _, building := range buildings {
		res = append(res, models.BuildingDTO{
			BuildingID:     building.BuildingId.String(),
			Name:           building.BuildingName,
			AvailableSlots: building.AvailableSlots,
			TotalSlots:     building.TotalSlots,
			TotalFloors:    building.TotalFloors,
		})
	}
	return res, nil
}

func NewBuildingService(repo buildingrepository.BuildingStorage) *BuildingService {
	return &BuildingService{
		buildingRepo: repo,
	}
}

func (bs *BuildingService) AddBuilding(ctx context.Context, buildingName string) error {
	ctxUser := ctx.Value(constants.User).(models.UserJwt)
	if ctxUser.Role != roles.Admin {
		return errors.New("unauthorized: only admin can add buildings")
	}
	return bs.buildingRepo.AddBuilding(ctx, buildingName)
}
