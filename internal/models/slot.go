package models

import (
	"fmt"

	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	"github.com/google/uuid"
)

// Slot is slot, isoccupied is for is assigned
type Slot struct {
	BuildingID  uuid.UUID                `gorm:"primaryKey;type:uuid"`
	FloorNumber int                      `gorm:"primaryKey;type:int"`
	SlotNumber  int                      `gorm:"primaryKey;type:int"`
	SlotType    vehicletypes.VehicleType `gorm:"not null"`
	Vehicles    []Vehicle                `gorm:"foreignKey:AssignedBuildingID,AssignedFloorNumber,AssignedSlotNumber"`
}

type SlotDTO struct {
	BuildingID  string `json:"building_id"`
	FloorNumber int    `json:"floor_number"`
	SlotNumber  int    `json:"slot_number"`
	SlotType    string `json:"slot_type"`
	IsOccupied  bool   `json:"is_occupied,omitempty"`
}

func (s Slot) String() string {
	if s.BuildingID == uuid.Nil {
		return "unassigned"
	}
	return fmt.Sprintf("%v_%v_%v", s.BuildingID, s.FloorNumber, s.SlotNumber)
}

func (s Slot) ToDTO() *SlotDTO {
	return &SlotDTO{
		BuildingID:  s.BuildingID.String(),
		FloorNumber: s.FloorNumber,
		SlotNumber:  s.SlotNumber,
		SlotType:    s.SlotType.String(),
		IsOccupied:  len(s.Vehicles) > 0,
	}
}
