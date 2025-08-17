package models

import (
	"time"

	"github.com/google/uuid"
)

type ParkingHistory struct {
	ParkingID uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	// NumberPlate string                   `gorm:"not null;type:varchar(10)"`
	// BuildingID uuid.UUID `gorm:"not null;type:uuid"`
	// Building   Building  `gorm:"foreignKey:BuildingID;references:BuildingID"`
	// UserID      uuid.UUID                `gorm:"not null;type:uuid"`
	// User        User                     `gorm:"foreignKey:UserID;references:UserID"`
	VehicleID uuid.UUID `gorm:"not null;type:uuid"`
	Vehicle   Vehicle   `gorm:"foreignKey:VehicleID;references:VehicleID"`
	// FloorNumber int        `gorm:"not null"`
	// Floor       Floor      `gorm:"foreignKey:BuildingID,FloorNumber;references:BuildingID,FloorNumber"`
	// SlotNumber  int        `gorm:"not null"`
	// Slot        Slot       `gorm:"foreignKey:BuildingID,FloorNumber,SlotNumber;references:BuildingID,FloorNumber,SlotNumber"`
	StartTime time.Time  `gorm:"not null;default:current_timestamp"`
	EndTime   *time.Time `gorm:"default:null"`
	// VehicleType vehicletypes.VehicleType `gorm:"not null"`
}

func (p ParkingHistory) GetID() string {
	return p.ParkingID.String()
}
