package userservice

import (
	"context"

	"github.com/Kaushik1766/ParkingManagement/internal/models/enums"
	"github.com/Kaushik1766/ParkingManagement/internal/models/vehicle"
)

// will pass primaryKey in the context

type UserManager interface {
	UpdateProfile(ctx context.Context, name, email, password string) error
	DeleteProfile(ctx context.Context) error
	RegisterVehicle(ctx context.Context, numberplate string, vehicleType enums.VehicleType) error
	UnregisterVehicle(ctx context.Context, numberplate string) error
	GetRegisteredVehicles(ctx context.Context) []vehicle.VehicleDTO
}
