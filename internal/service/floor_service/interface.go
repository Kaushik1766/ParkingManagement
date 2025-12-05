package floorservice

import (
	"context"

	"github.com/Kaushik1766/ParkingManagement/internal/models"
)

//go:generate mockgen -source=interface.go -destination=../../../mocks/floor_service_mock.go -package=mocks
type FloorMgr interface {
	AddFloor(ctx context.Context, buildingName string, floorNumber int) error
	DeleteFloor(ctx context.Context, buildignId string, floorNumber int) error
	AddFloorByBuildingId(ctx context.Context, buildingId string, floorNumber int) error
	GetFloorsByBuildingId(ctx context.Context, buildingId string) ([]models.FloorDTO, error)
}
