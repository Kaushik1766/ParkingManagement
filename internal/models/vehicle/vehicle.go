package vehicle

import (
	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	"github.com/google/uuid"
)

type Vehicle struct {
	VehicleId    uuid.UUID
	NumberPlate  string
	VehicleType  vehicletypes.VehicleType
	UserId       uuid.UUID
	AssignedSlot string
	IsActive     bool
}

func (v Vehicle) GetID() string {
	return v.VehicleId.String()
}
