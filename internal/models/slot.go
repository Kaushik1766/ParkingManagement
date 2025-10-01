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
	Vehicles    []Vehicle                `json:"-"gorm:"foreignKey:AssignedBuildingID,AssignedFloorNumber,AssignedSlotNumber"`
}

type SlotDTO struct {
	BuildingID  string `json:"building_id"`
	FloorNumber int    `json:"floor_number"`
	SlotNumber  int    `json:"slot_number"`
	SlotType    string `json:"slot_type"`
	IsOccupied  bool   `json:"is_occupied"`
}

// func (s Slot) GetID() string {
// 	return fmt.Sprintf("%v%v%v", s.BuildingID, s.FloorNumber, s.SlotNumber)
// }

func (s Slot) String() string {
	if s.BuildingID == uuid.Nil {
		return "unassigned"
	}
	return fmt.Sprintf("%v_%v_%v", s.BuildingID, s.FloorNumber, s.SlotNumber)
}

// func (s Slot) ToIdentifiableSlot(slotString string) (*Slot, error) {
// 	parts := strings.Split(slotString, "_")
// 	if len(parts) != 3 {
// 		return nil, fmt.Errorf("invalid slot string format: %s", slotString)
// 	}
// 	buildingId, err := uuid.Parse(parts[0])
// 	if err != nil {
// 		return nil, fmt.Errorf("invalid building ID: %s", parts[0])
// 	}
// 	floorNumber, err := strconv.Atoi(parts[1])
// 	if err != nil {
// 		return nil, fmt.Errorf("invalid floor number: %s", parts[1])
// 	}
// 	slotNumber, err := strconv.Atoi(parts[2])
// 	if err != nil {
// 		return nil, fmt.Errorf("invalid slot number: %s", parts[2])
// 	}
// 	return &Slot{
// 		BuildingID:  buildingId,
// 		FloorNumber: floorNumber,
// 		SlotNumber:  slotNumber,
// 	}, nil
// }
