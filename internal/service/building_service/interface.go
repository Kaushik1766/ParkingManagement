package buildingservice

import "context"

type BuildingMgr interface {
	AddBuilding(ctx context.Context, name string) error
	DeleteBuilding(ctx context.Context, name string) error
}
