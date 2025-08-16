package floor

import (
	"github.com/Kaushik1766/ParkingManagement/internal/models/office"
	"github.com/Kaushik1766/ParkingManagement/internal/models/slot"
	"github.com/google/uuid"
)

type Floor struct {
	BuildingID  uuid.UUID      `gorm:"primaryKey;type:uuid;not null"`
	FloorNumber int            `gorm:"primaryKey;type:int;not null"`
	Slots       []slot.Slot    `gorm:"foreignKey:BuildingID,FloorNumber;references:BuildingID,FloorNumber"`
	Office      *office.Office `gorm:"foreignKey:BuildingID,FloorNumber;references:BuildingID,FloorNumber"`
}

// func (f Floor) GetID() string {
// 	return fmt.Sprintf("%v%v", f.BuildingID, f.FloorNumber)
// }
