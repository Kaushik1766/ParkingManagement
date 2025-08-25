package floorrepository

import (
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/google/uuid"
)

type FloorStorage interface {
	AddFloor(buildingId string, floorNumber int) error
	DeleteFloor(buildingId string, floorNumber int) error
	GetFloor(buildingId uuid.UUID, floorNumber int) (int, error)
	GetFloorsByBuildingId(buildingId string) ([]models.Floor, error)
}
