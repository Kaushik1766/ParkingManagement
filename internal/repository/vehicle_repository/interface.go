package vehiclerepository

import (
	models "github.com/Kaushik1766/ParkingManagement/internal/models"
	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	"github.com/google/uuid"
)

type VehicleStorage interface {
	AddVehicle(numberplate string, userid uuid.UUID, vehicleType vehicletypes.VehicleType) (models.Vehicle, error)
	RemoveVehicle(numberplate string) error
	GetVehicleById(vehicleId uuid.UUID) (models.Vehicle, error)
	GetVehiclesByUserId(userId uuid.UUID) ([]models.Vehicle, error)
	GetVehicleByNumberPlate(numberplate string) (models.Vehicle, error)
	GetVehiclesWithUnassignedSlots() (vehicles []models.Vehicle, err error)
	Save(vehicle models.Vehicle) error
}
