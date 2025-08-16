package models

import (
	"fmt"
)

type VehicleDTO struct {
	NumberPlate  string `json:"number_plate"`
	VehicleType  string `json:"vehicle_type"`
	AssignedSlot Slot   `json:"assigned_slot"`
}

func (v VehicleDTO) String() string {
	return fmt.Sprintf("%s (%s)", v.NumberPlate, v.VehicleType)
}
