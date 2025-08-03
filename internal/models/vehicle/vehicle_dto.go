package vehicle

type VehicleDTO struct {
	NumberPlate  string `json:"number_plate"`
	VehicleType  string `json:"vehicle_type"`
	AssignedSlot string `json:"assigned_slot"`
}
