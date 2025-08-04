package floorservice

import "context"

type FloorMgr interface {
	AddFloor(ctx context.Context, buildingName string, floorNumber int) error
	DeleteFloor(ctx context.Context, buildingName string, floorNumber int) error
	DeleteFloors(ctx context.Context, buildingName string, floorNumbers []int) error
	AddFloors(ctx context.Context, buildingName string, floorNumbers []int) error
	GetFloorsByBuildingId(ctx context.Context, buildingName string) ([]int, error)
}
