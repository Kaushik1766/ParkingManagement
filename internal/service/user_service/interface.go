package userservice

import (
	"context"

	models "github.com/Kaushik1766/ParkingManagement/internal/models"
	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
)

type UserManager interface {
	UpdateProfile(ctx context.Context, name, email, password, office string) error
	DeleteProfile(ctx context.Context) error
	RegisterVehicle(ctx context.Context, numberplate string, vehicleType vehicletypes.VehicleType) error
	UnregisterVehicle(ctx context.Context, numberplate string) error
	GetRegisteredVehicles(ctx context.Context) []models.VehicleDTO
	GetUserProfile(ctx context.Context) (models.UserDTO, error)
	GetUserById(ctx context.Context, userId string) (models.UserDTO, error)
	GetAllUsers(ctx context.Context) ([]models.User, error)
}
