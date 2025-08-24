package models

import (
	"fmt"
)

type VehicleDTO struct {
	NumberPlate  string `json:"number_plate"`
	VehicleType  string `json:"vehicle_type"`
	AssignedSlot Slot   `json:"assigned_slot"`
}

type AddVehicleDTO struct {
	NumberPlate string `json:"number_plate" binding:"required"`
	VehicleType int    `json:"vehicle_type"`
}

func (v VehicleDTO) String() string {
	return fmt.Sprintf("%s (%s)", v.NumberPlate, v.VehicleType)
}
