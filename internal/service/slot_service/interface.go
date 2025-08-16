package slotservice

import (
	"context"

	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	"github.com/Kaushik1766/ParkingManagement/internal/models/slot"
	"github.com/google/uuid"
)

type SlotMgr interface {
	AddSlots(ctx context.Context, buildingName string, floorNumber int, slotNumbers []int, slotType vehicletypes.VehicleType) error
	DeleteSlots(ctx context.Context, buildingName string, floorNumber int, slotNumbers []int) error
	GetSlotsByFloor(ctx context.Context, buildingName string, floorNumber int) ([]slot.Slot, error)
	GetFreeSlotsByBuilding(ctx context.Context, buildingID uuid.UUID, vehicleType vehicletypes.VehicleType) ([]slot.Slot, error)
}
