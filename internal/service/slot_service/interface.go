package slotservice

import (
	"context"

	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
)

type SlotMgr interface {
	AddSlots(ctx context.Context, buildingName string, floorNumber int, slotNumbers []int, slotType vehicletypes.VehicleType) error
	DeleteSlots(ctx context.Context, buildingName string, floorNumber int, slotNumbers []int) error
}
