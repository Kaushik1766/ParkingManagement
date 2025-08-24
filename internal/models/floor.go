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
	BuildingID  string `json:"building_id"`
	FloorNumber int    `json:"floor_number"`
}
