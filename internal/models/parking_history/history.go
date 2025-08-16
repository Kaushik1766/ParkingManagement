package parkinghistory

import (
	"os/user"
	"time"

	"github.com/Kaushik1766/ParkingManagement/internal/models/building"
	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	"github.com/Kaushik1766/ParkingManagement/internal/models/floor"
	"github.com/Kaushik1766/ParkingManagement/internal/models/slot"
	"github.com/google/uuid"
)

type ParkingHistory struct {
	ParkingID   uuid.UUID                `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	NumberPlate string                   `gorm:"not null;type:varchar(10)"`
	BuildingID  uuid.UUID                `gorm:"not null;type:uuid"`
	Building    building.Building        `gorm:"foreignKey:BuildingID;references:BuildingID"`
	UserID      uuid.UUID                `gorm:"not null;type:uuid"`
	User        user.User                `gorm:"foreignKey:UserID;references:UserID"`
	FloorNumber int                      `gorm:"not null"`
	Floor       floor.Floor              `gorm:"foreignKey:BuildingID,FloorNumber;references:BuildingID,FloorNumber"`
	SlotNumber  int                      `gorm:"not null"`
	Slot        slot.Slot                `gorm:"foreignKey:BuildingID,FloorNumber,SlotNumber;references:BuildingID,FloorNumber,SlotNumber"`
	StartTime   time.Time                `gorm:"not null;default:current_timestamp"`
	EndTime     *time.Time               `gorm:"default:null"`
	VehicleType vehicletypes.VehicleType `gorm:"not null"`
}

func (p ParkingHistory) GetID() string {
	return p.ParkingID.String()
}
