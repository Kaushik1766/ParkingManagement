package userservice

import (
	"context"

	models "github.com/Kaushik1766/ParkingManagement/internal/models"
	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
)

//go:generate mockgen -source=interface.go -destination=../../../mocks/user_service_mock.go -package=mocks
type UserManager interface {
	UpdateProfile(ctx context.Context, userId string, updateReq models.UpdateUserDTO) error
	DeleteProfile(ctx context.Context, userId string) error
	RegisterVehicle(ctx context.Context, numberplate string, vehicleType vehicletypes.VehicleType) error
	UnregisterVehicle(ctx context.Context, numberplate string) error
	GetRegisteredVehicles(ctx context.Context) ([]models.VehicleDTO, error)
	GetUserProfile(ctx context.Context) (models.UserDTO, error)
	GetUserById(ctx context.Context, userId string) (models.UserDTO, error)
	GetAllUsers(ctx context.Context) ([]models.UserDTO, error)
}
