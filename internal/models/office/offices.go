package office

import "github.com/google/uuid"

type Office struct {
	OfficeID    uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	OfficeName  string    `gorm:"not null;unique"`
	BuildingID  uuid.UUID `gorm:"not null;type:uuid"`
	FloorNumber int       `gorm:"not null"`
}
