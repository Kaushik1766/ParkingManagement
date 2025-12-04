package models

import (
	"fmt"

	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	"github.com/google/uuid"
)

type Slot struct {
	BuildingID  uuid.UUID                `gorm:"primaryKey;type:uuid"`
	FloorNumber int                      `gorm:"primaryKey;type:int"`
	Floor       Floor                    `gorm:"-"`
	SlotNumber  int                      `gorm:"primaryKey;type:int"`
	SlotType    vehicletypes.VehicleType `gorm:"not null"`
	Vehicles    []Vehicle                `json:"-" gorm:"foreignKey:AssignedBuildingID,AssignedFloorNumber,AssignedSlotNumber"`
}

type SlotDTO struct {
	BuildingID    string            `json:"buildingId"`
	FloorNumber   int               `json:"floorNumber"`
	SlotNumber    int               `json:"slotNumber"`
	SlotType      string            `json:"slotType"`
	IsAssigned    bool              `json:"isAssigned"`
	ParkingStatus *ParkingStatusDTO `json:"parkingStatus,omitempty"`
}

type ParkingStatusDTO struct {
	NumberPlate string `json:"numberPlate,omitempty"`
	ParkedAt    string `json:"parkedAt,omitempty"`
	UserName    string `json:"userName,omitempty"`
	UserEmail   string `json:"userEmail,omitempty"`
}

func (s Slot) String() string {
	if s.BuildingID == uuid.Nil {
		return "unassigned"
	}
	return fmt.Sprintf("%v_%v_%v", s.BuildingID, s.FloorNumber, s.SlotNumber)
}
