package slotassignment

import (
	"context"

	"github.com/Kaushik1766/ParkingManagement/internal/models/slot"
)

type SlotAssignmentMgr interface {
	AutoAssignSlot(ctx context.Context, vehicleId string) error
	UnassignSlot(ctx context.Context, vehicleId string) error
	AssignSlot(ctx context.Context, vehicleId string, slot slot.Slot) error
}
