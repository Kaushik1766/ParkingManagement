package models

import (
	"github.com/google/uuid"
)

type Floor struct {
	BuildingID  uuid.UUID `gorm:"primaryKey;type:uuid"`
	Building    Building
	FloorNumber int     `gorm:"primaryKey"`
	Slots       []Slot  `gorm:"foreignKey:BuildingID,FloorNumber;references:BuildingID,FloorNumber"`
	Office      *Office `gorm:"foreignKey:BuildingID,FloorNumber"`
	// OfficeID    *uuid.UUID `gorm:"type:uuid;default:null"`
}

// func (f Floor) GetID() string {
// 	return fmt.Sprintf("%v%v", f.BuildingID, f.FloorNumber)
// }

type FloorDTO struct {
	BuildingID     string `json:"buildingId"`
	FloorNumber    int    `json:"floorNumber"`
	TotalSlots     int    `json:"totalSlots"`
	AvailableSlots int    `json:"availableSlots"`
	AssignedOffice string `json:"assignedOffice,omitempty"`
	//Slots       []Slot `json:"Slots"`
}

type FloorSummary struct {
	BuildingID     uuid.UUID
	FloorNumber    int
	TotalSlots     int
	AvailableSlots int
	AssignedOffice string
}
