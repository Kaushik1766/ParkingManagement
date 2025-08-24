package models

import (
	"github.com/google/uuid"
)

type Building struct {
	BuildingID   uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	BuildingName string    `gorm:"type:varchar(255);not null;unique"`
	Floors       []Floor   `gorm:"foreignKey:BuildingID;references:BuildingID"`
}

// func (b Building) GetID() string {
// 	return b.BuildingID.String()
// }
//

type BuildingDTO struct {
	BuildingID string `json:"building_id"`
	Name       string `json:"name"`
}
