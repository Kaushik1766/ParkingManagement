package buildingrepository

import (
	"context"

	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/google/uuid"
)

//go:generate mockgen -source=interface.go -destination=../../../mocks/building_storage_mock.go -package=mocks
type BuildingStorage interface {
	AddBuilding(ctx context.Context, buildingName string) error
	DeleteBuildingByID(ctx context.Context, buildingID string) error
	// GetBuildingByName(buildingName string) (models.Building, error)
	GetAllBuildings(ctx context.Context) ([]models.Building, error)
	GetBuildingByID(ctx context.Context, buildingID uuid.UUID) (models.Building, error)
	GetAllBuildingSummary(ctx context.Context) ([]models.BuildingSummary, error)
}
