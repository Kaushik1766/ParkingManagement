package vehicle

import (
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums"
	"github.com/google/uuid"
)

type Vehicle struct {
	VehicleId    uuid.UUID
	NumberPlate  string
	VehicleType  enums.VehicleType
	UserId       uuid.UUID
	AssignedSlot string
	IsActive     bool
}

func (v Vehicle) GetID() string {
	return v.VehicleId.String()
}
