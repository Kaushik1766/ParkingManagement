package slotservice

import (
	"context"

	models "github.com/Kaushik1766/ParkingManagement/internal/models"
	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	"github.com/google/uuid"
)

//go:generate mockgen -source=interface.go -destination=../../../mocks/slot_service_mock.go -package=mocks
type SlotMgr interface {
	// AddSlots(ctx context.Context, buildingName string, floorNumber int, slotNumbers []int, slotType vehicletypes.VehicleType) error
	// DeleteSlots(ctx context.Context, buildingName string, floorNumber int, slotNumbers []int) error
	GetSlotsByFloor(ctx context.Context, buildingId string, floorNumber int) ([]models.SlotDTO, error)
	GetFreeSlotsByBuilding(ctx context.Context, buildingID uuid.UUID, vehicleType vehicletypes.VehicleType) ([]models.Slot, error)
}
