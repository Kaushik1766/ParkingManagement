package buildingrepository

import (
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/google/uuid"
)

//go:generate mockgen -source=interface.go -destination=../../../mocks/building_storage_mock.go -package=mocks
type BuildingStorage interface {
	AddBuilding(buildingName string) error
	DeleteBuildingByName(buildingName string) error
	DeleteBuildingByID(buildingID string) error
	GetBuildingByName(buildingName string) (models.Building, error)
	GetAllBuildings() ([]models.Building, error)
	GetBuildingByID(buildingID uuid.UUID) (models.Building, error)
}
