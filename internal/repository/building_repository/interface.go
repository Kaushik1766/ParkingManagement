package buildingrepository

import (
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/google/uuid"
)

type BuildingStorage interface {
	AddBuilding(name string) error
	DeleteBuilding(name string) error
	GetBuildingByName(name string) (models.Building, error)
	GetAllBuildings() ([]models.Building, error)
	GetBuildingByID(buildingID uuid.UUID) (models.Building, error)
}
