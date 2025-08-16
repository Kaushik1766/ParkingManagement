package vehicle

import (
	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	"github.com/Kaushik1766/ParkingManagement/internal/models/slot"
	"github.com/google/uuid"
)

type Vehicle struct {
	VehicleID    uuid.UUID                `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	NumberPlate  string                   `gorm:"not null;type:varchar(10);unique"`
	VehicleType  vehicletypes.VehicleType `gorm:"not null"`
	UserID       uuid.UUID                `gorm:"type:uuid;not null"`
	BuildingID   uuid.UUID                `gorm:"type:uuid;not null"`
	FloorNumber  int                      `gorm:"not null"`
	SlotNumber   int                      `gorm:"not null"`
	AssignedSlot slot.Slot                `gorm:"foreignKey:BuildingID,FloorNumber,SlotNumber;references:BuildingID,FloorNumber,SlotNumber"`
	IsActive     bool                     `gorm:"default:true"`
}

// func (v Vehicle) GetID() string {
// 	return v.VehicleID.String()
// }
