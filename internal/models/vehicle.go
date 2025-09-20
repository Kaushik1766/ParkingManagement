package models

import (
	"fmt"

	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	"github.com/google/uuid"
)

type Vehicle struct {
	VehicleID           uuid.UUID                `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	NumberPlate         string                   `gorm:"not null;type:varchar(10);unique"`
	VehicleType         vehicletypes.VehicleType `gorm:"not null"`
	UserID              uuid.UUID
	User                User
	AssignedBuildingID  *uuid.UUID `gorm:"type:uuid;default:null"`
	AssignedFloorNumber *int       `gorm:"default:null;type:int"`
	AssignedSlotNumber  *int       `gorm:"default:null;type:int"`
	AssignedSlot        *Slot      `gorm:"foreignKey:AssignedBuildingID,AssignedFloorNumber,AssignedSlotNumber;references:BuildingID,FloorNumber,SlotNumber"`
	IsActive            bool       `gorm:"default:true"`
}

type VehicleDTO struct {
	NumberPlate  string   `json:"number_plate"`
	VehicleType  string   `json:"vehicle_type"`
	AssignedSlot *SlotDTO `json:"assigned_slot"`
}

type AddVehicleDTO struct {
	NumberPlate string `json:"numberplate" binding:"required"`
	VehicleType int    `json:"type"`
}

func (v VehicleDTO) String() string {
	return fmt.Sprintf("%s (%s)", v.NumberPlate, v.VehicleType)
}
