package models

import (
	"fmt"
	"time"

	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	"github.com/google/uuid"
)

type ParkingHistory struct {
	ParkingID uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	// NumberPlate string                   `gorm:"not null;type:varchar(10)"`
	// BuildingID uuid.UUID `gorm:"not null;type:uuid"`
	// Building   Building  `gorm:"foreignKey:BuildingID;references:BuildingID"`
	// UserID      uuid.UUID                `gorm:"not null;type:uuid"`
	// User        User                     `gorm:"foreignKey:UserID;references:UserID"`
	VehicleID uuid.UUID `gorm:"not null;type:uuid;"`
	Vehicle   Vehicle   `gorm:"references:VehicleID"`
	// FloorNumber int        `gorm:"not null"`
	// Floor       Floor      `gorm:"foreignKey:BuildingID,FloorNumber;references:BuildingID,FloorNumber"`
	// SlotNumber  int        `gorm:"not null"`
	// Slot        Slot       `gorm:"foreignKey:BuildingID,FloorNumber,SlotNumber;references:BuildingID,FloorNumber,SlotNumber"`
	StartTime time.Time  `gorm:"not null;default:current_timestamp"`
	EndTime   *time.Time `gorm:"default:null"`
	// VehicleType vehicletypes.VehicleType `gorm:"not null"`
}
type ParkingHistoryDTO struct {
	TicketId     string
	NumberPlate  string
	BuildingId   string
	FLoorNumber  int
	SlotNumber   int
	StartTime    time.Time
	EndTime      time.Time
	VechicleType vehicletypes.VehicleType
}

func (phdto *ParkingHistoryDTO) String() string {
	return fmt.Sprintf("TicketId: %s\nNumberPlate: %s\nBuildingId: %s\nFloorNumber: %d\nSlotNumber: %d\nStartTime: %s\nEndTime: %s",
		phdto.TicketId, phdto.NumberPlate, phdto.BuildingId, phdto.FLoorNumber, phdto.SlotNumber, phdto.StartTime, phdto.EndTime)
}
