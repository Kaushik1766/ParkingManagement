package userservice

import (
	"context"

	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	user "github.com/Kaushik1766/ParkingManagement/internal/models/user"
	"github.com/Kaushik1766/ParkingManagement/internal/models/vehicle"
)

type UserManager interface {
	UpdateProfile(ctx context.Context, name, email, password, office string) error
	DeleteProfile(ctx context.Context) error
	RegisterVehicle(ctx context.Context, numberplate string, vehicleType vehicletypes.VehicleType) error
	UnregisterVehicle(ctx context.Context, numberplate string) error
	GetRegisteredVehicles(ctx context.Context) []vehicle.VehicleDTO
	GetUserProfile(ctx context.Context) (user.UserDTO, error)
	GetUserById(ctx context.Context, userId string) (user.UserDTO, error)
	GetAllUsers(ctx context.Context) ([]user.User, error)
}
