package buildingrepository

import (
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/google/uuid"
)

type BuildingStorage interface {
	AddBuilding(buildingName string) error
	DeleteBuildingByName(buildingName string) error
	DeleteBuildingByID(buildingID string) error
	GetBuildingByName(buildingName string) (models.Building, error)
	GetAllBuildings() ([]models.Building, error)
	GetBuildingByID(buildingID uuid.UUID) (models.Building, error)
}
