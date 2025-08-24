package buildingservice

import (
	"context"

	"github.com/Kaushik1766/ParkingManagement/internal/models"
)

type BuildingMgr interface {
	AddBuilding(ctx context.Context, name string) error
	DeleteBuilding(ctx context.Context, name string) error
	DeleteBuildingByID(ctx context.Context, buildingID string) error
	GetAllBuildings(ctx context.Context) ([]models.BuildingDTO, error)
	GetBuildingByID(ctx context.Context, buildingID string) (models.BuildingDTO, error)
}
