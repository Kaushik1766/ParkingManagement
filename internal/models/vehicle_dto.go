package models

import (
	"fmt"
)

type VehicleDTO struct {
	NumberPlate          string `json:"number_plate"`
	VehicleType          string `json:"vehicle_type"`
	IsParked             bool   `json:"is_parked"`
	AssignedBuildingName string `json:"assigned_building_name"`
	AssignedBuildingID   string `json:"assigned_building_id"`
	AssignedFloorNumber  int    `json:"assigned_floor_number"`
	AssignedSlotNumber   int    `json:"assigned_slot_number"`
}

type AddVehicleDTO struct {
	NumberPlate string `json:"numberplate" binding:"required"`
	VehicleType int    `json:"type"`
}

func (v VehicleDTO) String() string {
	return fmt.Sprintf("%s (%s)", v.NumberPlate, v.VehicleType)
}
