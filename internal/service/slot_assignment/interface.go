package slotassignment

import (
	"context"
)

//go:generate mockgen -source=interface.go -destination=../../../mocks/slot_assignment_service_mock.go -package=mocks
type SlotAssignmentMgr interface {
	AutoAssignSlot(ctx context.Context, numberplate string) error
}
