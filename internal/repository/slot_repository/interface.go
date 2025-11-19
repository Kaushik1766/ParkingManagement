package slotrepository

import (
	"context"

	slot "github.com/Kaushik1766/ParkingManagement/internal/models"
	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	"github.com/google/uuid"
)

//go:generate mockgen -source=interface.go -destination=../../../mocks/slot_storage_mock.go -package=mocks
type SlotStorage interface {
	AddSlot(ctx context.Context, buildingId uuid.UUID, floorNumber, slotNumber int, slotType vehicletypes.VehicleType) error
	DeleteSlot(ctx context.Context, buildingId uuid.UUID, floorNumber, slotNumber int) error
	GetSlotsByFloor(ctx context.Context, buildingId uuid.UUID, floorNumber int) ([]slot.Slot, error)
	GetFreeSlotsByFloor(ctx context.Context, buildingId uuid.UUID, floorNumber int) ([]slot.Slot, error)
	GetFreeSlotsByBuilding(ctx context.Context, buildingId uuid.UUID) ([]slot.Slot, error)
	Save(ctx context.Context, slot slot.Slot) error
}
