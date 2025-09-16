package models

import (
	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	"github.com/google/uuid"
)

type Vehicle struct {
	VehicleID          uuid.UUID                `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	NumberPlate        string                   `gorm:"not null;type:varchar(10);unique"`
	VehicleType        vehicletypes.VehicleType `gorm:"not null"`
	UserID             uuid.UUID
	User               User
	AssignedBuildingID *uuid.UUID `gorm:"type:uuid;default:null"`
	// AssignedBuilding    *Building `gorm:"foreignKey:AssignedBuildingID;references:BuildingID"`
	AssignedFloorNumber *int `gorm:"default:null;type:int"`
	// AssignedFloor       *Floor    `gorm:"foreignKey:AssignedBuildingID,AssignedFloorNumber;references:BuildingID,FloorNumber"`
	AssignedSlotNumber *int  `gorm:"default:null;type:int"`
	AssignedSlot       *Slot `gorm:"foreignKey:AssignedBuildingID,AssignedFloorNumber,AssignedSlotNumber;references:BuildingID,FloorNumber,SlotNumber"`
	IsActive           bool  `gorm:"default:true"`
	// ParkingHistories    []ParkingHistory         `gorm:"foreignKey:VehicleID;references:VehicleID"`
}

// func (v Vehicle) GetID() string {
// 	return v.VehicleID.String()
// }
