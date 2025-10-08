package floorrepository

import (
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/google/uuid"
)

//go:generate mockgen -source=interface.go -destination=../../../mocks/floor_storage_mock.go -package=mocks
type FloorStorage interface {
	AddFloor(buildingId string, floorNumber int) error
	DeleteFloor(buildingId string, floorNumber int) error
	GetFloor(buildingId uuid.UUID, floorNumber int) (int, error)
	GetFloorsByBuildingId(buildingId string) ([]models.FloorSummary, error)
}
