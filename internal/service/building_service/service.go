package buildingservice

import (
	"context"
	"errors"

	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	userjwt "github.com/Kaushik1766/ParkingManagement/internal/models/user_jwt"
	buildingrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/building_repository"
)

type BuildingService struct {
	buildingRepo buildingrepository.BuildingStorage
}

func (bs *BuildingService) GetAllBuildings(ctx context.Context) ([]string, error) {
	ctxUser := ctx.Value(constants.User).(userjwt.UserJwt)
	if ctxUser.Role != roles.Admin {
		return nil, errors.New("unauthorized: only admin can view buildings")
	}
	buildings, err := bs.buildingRepo.GetAllBuildings()
	if err != nil {
		return nil, err
	}
	var buildingNames []string
	for _, building := range buildings {
		buildingNames = append(buildingNames, building.BuildingName)
	}
	return buildingNames, nil
}

func NewBuildingService(repo buildingrepository.BuildingStorage) *BuildingService {
	return &BuildingService{
		buildingRepo: repo,
	}
}

func (bs *BuildingService) AddBuilding(ctx context.Context, name string) error {
	ctxUser := ctx.Value(constants.User).(userjwt.UserJwt)
	if ctxUser.Role != roles.Admin {
		return errors.New("unauthorized: only admin can add buildings")
	}
	err := bs.buildingRepo.AddBuilding(name)
	return err
}

func (bs *BuildingService) DeleteBuilding(ctx context.Context, name string) error {
	ctxUser := ctx.Value(constants.User).(userjwt.UserJwt)
	if ctxUser.Role != roles.Admin {
		return errors.New("unauthorized: only admin can delete buildings")
	}
	err := bs.buildingRepo.DeleteBuilding(name)
	return err
}
