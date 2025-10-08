package models

import (
	"github.com/google/uuid"
)

type Office struct {
	OfficeID    uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	OfficeName  string    `gorm:"not null;unique"`
	BuildingID  uuid.UUID `gorm:"uniqueIndex:idx_building_floor_office"`
	FloorNumber int       `gorm:"uniqueIndex:idx_building_floor_office"`
}

type OfficeDTO struct {
	BuildingID  string `json:"building_id"`
	FloorNumber int    `json:"floor_number"`
	OfficeName  string `json:"office_name"`
	OfficeID    string `json:"office_id"`
}
