package slotrepository

import (
	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	"github.com/Kaushik1766/ParkingManagement/internal/models/slot"
	"github.com/google/uuid"
)

type SlotStorage interface {
	AddSlot(buildingId uuid.UUID, floorNumber, slotNumber int, slotType vehicletypes.VehicleType) error
	DeleteSlot(buildingId uuid.UUID, floorNumber, slotNumber int) error
	GetSlotsByFloor(buildingId uuid.UUID, floorNumber int) ([]slot.Slot, error)
	GetFreeSlotsByFloor(buildingId uuid.UUID, floorNumber int) ([]slot.Slot, error)
	SetSlotOccupied(buildingId uuid.UUID, floorNumber, slotNumber int, isOccupied bool) error
	GetFreeSlotsByBuilding(buildingId uuid.UUID) ([]slot.Slot, error)
	Save(slot slot.Slot) error
}
