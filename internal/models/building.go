package models

import (
	"github.com/google/uuid"
)

type Building struct {
	BuildingID   uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	BuildingName string    `gorm:"type:varchar(255);not null;unique"`
	Floors       []Floor   `gorm:"foreignKey:BuildingID;references:BuildingID"`
}

type BuildingDTO struct {
	BuildingID     string `json:"buildingId"`
	Name           string `json:"name"`
	AvailableSlots int    `json:"availableSlots"`
	TotalSlots     int    `json:"totalSlots"`
	TotalFloors    int    `json:"totalFloors"`
}

type BuildingSummary struct {
	BuildingId     uuid.UUID
	BuildingName   string
	AvailableSlots int
	TotalSlots     int
	TotalFloors    int
}
