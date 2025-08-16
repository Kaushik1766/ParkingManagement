package models

import (
	"github.com/google/uuid"
)

type Floor struct {
	BuildingID  uuid.UUID `gorm:"primaryKey;type:uuid;not null"`
	FloorNumber int       `gorm:"primaryKey;type:int;not null"`
	Slots       []Slot    `gorm:"foreignKey:BuildingID,FloorNumber;references:BuildingID,FloorNumber"`
	Office      *Office   `gorm:"foreignKey:BuildingID,FloorNumber;references:BuildingID,FloorNumber"`
}

// func (f Floor) GetID() string {
// 	return fmt.Sprintf("%v%v", f.BuildingID, f.FloorNumber)
// }
