package vehiclerepository

import (
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums"
	"github.com/Kaushik1766/ParkingManagement/internal/models/vehicle"
	"github.com/google/uuid"
)

type VehicleStorage interface {
	AddVehicle(numberplate string, userid uuid.UUID, vehicleType enums.VehicleType) error
	RemoveVehicle(numberplate string) error
	GetVehicleById(vehicleId uuid.UUID) (vehicle.Vehicle, error)
	GetVehiclesByUserId(userId uuid.UUID) ([]vehicle.Vehicle, error)
}
