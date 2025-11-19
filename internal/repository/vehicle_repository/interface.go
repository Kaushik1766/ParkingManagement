package vehiclerepository

import (
	"context"

	models "github.com/Kaushik1766/ParkingManagement/internal/models"
	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	"github.com/google/uuid"
)

//go:generate mockgen -source=interface.go -destination=../../../mocks/vehicle_storage_mock.go -package=mocks
type VehicleStorage interface {
	AddVehicle(ctx context.Context, numberplate string, userid uuid.UUID, vehicleType vehicletypes.VehicleType) (models.Vehicle, error)
	RemoveVehicle(ctx context.Context, numberplate string) error
	GetVehicleById(ctx context.Context, vehicleId uuid.UUID) (models.Vehicle, error)
	GetVehiclesByUserId(ctx context.Context, userId uuid.UUID) ([]models.Vehicle, error)
	GetVehicleByNumberPlate(ctx context.Context, numberplate string) (models.Vehicle, error)
	GetVehiclesWithUnassignedSlots(ctx context.Context) (vehicles []models.Vehicle, err error)
	GetParkingStatus(ctx context.Context, numberplate string) (bool, error)
	Save(ctx context.Context, vehicle models.Vehicle) error
}
