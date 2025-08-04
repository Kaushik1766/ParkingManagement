package buildingrepository

import "github.com/Kaushik1766/ParkingManagement/internal/models/building"

type BuildingStorage interface {
	AddBuilding(name string) error
	DeleteBuilding(name string) error
	GetBuildingByName(name string) (building.Building, error)
}
