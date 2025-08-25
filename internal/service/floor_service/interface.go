package floorservice

import (
	"context"

	"github.com/Kaushik1766/ParkingManagement/internal/models"
)

type FloorMgr interface {
	AddFloor(ctx context.Context, buildingName string, floorNumber int) error
	DeleteFloor(ctx context.Context, buildignId string, floorNumber int) error
	AddFloors(ctx context.Context, buildingName string, floorNumbers []int) error
	AddFloorByBuildingId(ctx context.Context, buildingId string, floorNumber int) error
	GetFloorsByBuildingId(ctx context.Context, buildingId string) ([]models.FloorDTO, error)
}
