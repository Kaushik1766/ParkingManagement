package vehicleservice

import (
	"context"

	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	models "github.com/Kaushik1766/ParkingManagement/internal/models"
	parkinghistoryrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/parking_history_repository"
	vehiclerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/vehicle_repository"
	customerrors "github.com/Kaushik1766/ParkingManagement/pkg/customErrors"
)

type VehicleService struct {
	vehicleRepo vehiclerepository.VehicleStorage
	parkingRepo parkinghistoryrepository.ParkingHistoryStorage
}

func (vs *VehicleService) UnparkByNumberPlate(ctx context.Context, numberplate string) error {
	userCtx := ctx.Value(constants.User).(models.UserJwt)

	vehicle, err := vs.vehicleRepo.GetVehicleByNumberPlate(numberplate)
	if err != nil {
		return err
	}

	if vehicle.UserID.String() != userCtx.ID {
		return customerrors.Unathorized{}
	}

	err = vs.parkingRepo.UnparkByNumberPlate(numberplate)
	if err != nil {
		return err
	}

	return nil
}

func (vs *VehicleService) Park(ctx context.Context, numberplate string) (string, error) {
	userCtx := ctx.Value(constants.User).(models.UserJwt)

	vehicle, err := vs.vehicleRepo.GetVehicleByNumberPlate(numberplate)
	if err != nil {
		return "", err
	}

	if vehicle.UserID.String() != userCtx.ID {
		return "", customerrors.Unathorized{}
	}

	return vs.parkingRepo.AddParking(vehicle)

}

func (vs *VehicleService) Unpark(ctx context.Context, ticketId string) error {
	// TODO: check if auth is correct in usecase
	// userCtx := ctx.Value(constants.User).(models.UserJwt)

	// vehicle, err := vs.vehicleRepo.GetVehicleByNumberPlate(ticketId)
	// if err != nil {
	// 	return err
	// }
	//
	// if vehicle.UserId.String() != userCtx.ID {
	// 	return customerrors.Unathorized{}
	// }

	return vs.parkingRepo.Unpark(ticketId)
}

func NewVehicleService(vehicleRepo vehiclerepository.VehicleStorage, parkingRepo parkinghistoryrepository.ParkingHistoryStorage) *VehicleService {
	return &VehicleService{
		vehicleRepo: vehicleRepo,
		parkingRepo: parkingRepo,
	}
}
