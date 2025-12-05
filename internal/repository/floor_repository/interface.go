package floorrepository

import (
	"context"

	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/google/uuid"
)

//go:generate mockgen -source=interface.go -destination=../../../mocks/floor_storage_mock.go -package=mocks
type FloorStorage interface {
	AddFloor(ctx context.Context, buildingId string, floorNumber int) error
	DeleteFloor(ctx context.Context, buildingId string, floorNumber int) error
	GetFloor(ctx context.Context, buildingId uuid.UUID, floorNumber int) (int, error)
	GetFloorsByBuildingId(ctx context.Context, buildingId string) ([]models.FloorSummary, error)
}
