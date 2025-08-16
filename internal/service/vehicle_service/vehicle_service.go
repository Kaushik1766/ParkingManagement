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

func (vs *VehicleService) Park(ctx context.Context, numberplate string) (string, error) {
	userCtx := ctx.Value(constants.User).(models.UserJwt)

	vehicle, err := vs.vehicleRepo.GetVehicleByNumberPlate(numberplate)
	if err != nil {
		return "", err
	}

	if vehicle.UserID.String() != userCtx.ID {
		return "", customerrors.Unathorized{}
	}

	ticketId, err := vs.parkingRepo.AddParking(vehicle)
	if err != nil {
		return "", err
	}

	return ticketId, nil
}

func (vs *VehicleService) Unpark(ctx context.Context, ticketId string) error {
	// userCtx := ctx.Value(constants.User).(models.UserJwt)

	// vehicle, err := vs.vehicleRepo.GetVehicleByNumberPlate(ticketId)
	// if err != nil {
	// 	return err
	// }
	//
	// if vehicle.UserId.String() != userCtx.ID {
	// 	return customerrors.Unathorized{}
	// }

	err := vs.parkingRepo.Unpark(ticketId)
	if err != nil {
		return err
	}

	return nil
}

func NewVehicleService(vehicleRepo vehiclerepository.VehicleStorage, parkingRepo parkinghistoryrepository.ParkingHistoryStorage) *VehicleService {
	return &VehicleService{
		vehicleRepo: vehicleRepo,
		parkingRepo: parkingRepo,
	}
}
