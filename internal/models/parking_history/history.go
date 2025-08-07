package parkinghistory

import (
	"time"

	"github.com/google/uuid"
)

type ParkingHistory struct {
	ParkingId   uuid.UUID
	NumberPlate string
	BuildingId  string
	UserId      uuid.UUID
	FLoorNumber int
	SlotNumber  int
	StartTime   time.Time
	EndTime     time.Time
}

func (p ParkingHistory) GetID() string {
	return p.ParkingId.String()
}
