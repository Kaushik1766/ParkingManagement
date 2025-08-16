package buildingrepository

import (
	"github.com/Kaushik1766/ParkingManagement/internal/models/building"
	"github.com/google/uuid"
)

type BuildingStorage interface {
	AddBuilding(name string) error
	DeleteBuilding(name string) error
	GetBuildingByName(name string) (building.Building, error)
	GetAllBuildings() ([]building.Building, error)
	GetBuildingByID(buildingID uuid.UUID) (building.Building, error)
}
