package vehicle

import (
	"fmt"

	"github.com/Kaushik1766/ParkingManagement/internal/models/slot"
)

type VehicleDTO struct {
	NumberPlate  string    `json:"number_plate"`
	VehicleType  string    `json:"vehicle_type"`
	AssignedSlot slot.Slot `json:"assigned_slot"`
}

func (v VehicleDTO) String() string {
	return fmt.Sprintf("%s (%s)", v.NumberPlate, v.VehicleType)
}
