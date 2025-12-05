package buildingservice

import (
	"context"

	"github.com/Kaushik1766/ParkingManagement/internal/models"
)

//go:generate mockgen -source=interface.go -destination=../../../mocks/building_service_mock.go -package=mocks
type BuildingMgr interface {
	AddBuilding(ctx context.Context, name string) error
	DeleteBuildingByID(ctx context.Context, buildingID string) error
	GetAllBuildings(ctx context.Context) ([]models.BuildingDTO, error)
	GetBuildingByID(ctx context.Context, buildingID string) (models.BuildingDTO, error)
}
