package slotassignment

import (
	"context"

	"github.com/Kaushik1766/ParkingManagement/internal/models"
)

//go:generate mockgen -source=interface.go -destination=../../../mocks/slot_assignment_service_mock.go -package=mocks
type SlotAssignmentMgr interface {
	AutoAssignSlot(ctx context.Context, vehicleId string) error
	UnassignSlot(ctx context.Context, vehicleId string) error
	AssignSlot(ctx context.Context, vehicleId string, slot models.Slot) error
	GetVehiclesWithUnassignedSlots(ctx context.Context) ([]models.Vehicle, error)
}
