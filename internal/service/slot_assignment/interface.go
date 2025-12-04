package slotassignment

import (
	"context"

	slot "github.com/Kaushik1766/ParkingManagement/internal/models"
)

//go:generate mockgen -source=interface.go -destination=../../../mocks/slot_assignment_service_mock.go -package=mocks
type SlotAssignmentMgr interface {
	AutoAssignSlot(ctx context.Context, numberplate string) error
	AssignSlot(ctx context.Context, numberplate string, slot slot.Slot) error
}
