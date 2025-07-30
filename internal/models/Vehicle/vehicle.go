package vehicle

import (
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums"
	"github.com/google/uuid"
)

type Vehicle struct {
	VehicleId   uuid.UUID
	NumberPlate string
	UserId      uuid.UUID
	VehicleType enums.VehicleType
}

func (v Vehicle) GetID() string {
	return v.VehicleId.String()
}
