package parkinghistory

import (
	"time"

	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
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
	VehicleType vehicletypes.VehicleType
}

func (p ParkingHistory) GetID() string {
	return p.ParkingId.String()
}
