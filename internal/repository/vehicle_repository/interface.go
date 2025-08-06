package vehiclerepository

import (
	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	"github.com/Kaushik1766/ParkingManagement/internal/models/vehicle"
	"github.com/google/uuid"
)

type VehicleStorage interface {
	AddVehicle(numberplate string, userid uuid.UUID, vehicleType vehicletypes.VehicleType) error
	RemoveVehicle(numberplate string) error
	GetVehicleById(vehicleId uuid.UUID) (vehicle.Vehicle, error)
	GetVehiclesByUserId(userId uuid.UUID) ([]vehicle.Vehicle, error)
	Save(vehicle vehicle.Vehicle) error
}
