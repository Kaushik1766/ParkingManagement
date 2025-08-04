package floorrepository

import "github.com/google/uuid"

type FloorStorage interface {
	AddFloor(buildingId uuid.UUID, floorNumber int) error
	DeleteFloor(buildingId uuid.UUID, floorNumber int) error
}
