package parkinghistoryservice

import (
	"context"
	"time"

	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	models "github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	parkinghistoryrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/parking_history_repository"
	vehiclerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/vehicle_repository"
	customerrors "github.com/Kaushik1766/ParkingManagement/pkg/customErrors"
)

type ParkingHistoryService struct {
	parkingRepo parkinghistoryrepository.ParkingHistoryStorage
	vehicleRepo vehiclerepository.VehicleStorage
}

func (phs *ParkingHistoryService) GetParkingHistoryByUserId(userId string, startTime, endTime time.Time) ([]models.ParkingHistoryDTO, error) {
	parkingHistory, err := phs.parkingRepo.GetParkingHistoryByUser(userId, startTime, endTime)
	if err != nil {
		return []models.ParkingHistoryDTO{}, err
	}

	return parkingHistory, nil
}

func (phs *ParkingHistoryService) GetParkingHistory(ctx context.Context, startTime, endTime time.Time) ([]models.ParkingHistoryDTO, error) {
	userCtx := ctx.Value(constants.User).(models.UserJwt)

	parkingHistory, err := phs.parkingRepo.GetParkingHistoryByUser(userCtx.ID, startTime, endTime)
	if err != nil {
		return nil, err
	}

	return parkingHistory, nil
}

func (phs *ParkingHistoryService) GetParkingHistoryByNumberPlate(ctx context.Context, numberplate string, startTime, endTime time.Time) ([]models.ParkingHistoryDTO, error) {
	userCtx := ctx.Value(constants.User).(models.UserJwt)

	vehicle, err := phs.vehicleRepo.GetVehicleByNumberPlate(numberplate)
	if err != nil {
		return nil, err
	}

	if userCtx.Role != roles.Admin && userCtx.ID != vehicle.UserID.String() {
		return nil, customerrors.Unathorized{}
	}

	parkingHistory, err := phs.parkingRepo.GetParkingHistoryByNumberPlate(numberplate, startTime, endTime)
	if err != nil {
		return nil, err
	}

	return parkingHistory, err
}

func (phs *ParkingHistoryService) GetActiveUserParkings(ctx context.Context) ([]models.ParkingHistoryDTO, error) {
	userCtx := ctx.Value(constants.User).(models.UserJwt)

	activeParkings, err := phs.parkingRepo.GetActiveUserParkings(userCtx.ID)
	if err != nil {
		return nil, err
	}

	return activeParkings, nil
}

func NewParkingHistoryService(parkingRepo parkinghistoryrepository.ParkingHistoryStorage, vehicleRepo vehiclerepository.VehicleStorage) *ParkingHistoryService {
	return &ParkingHistoryService{
		parkingRepo: parkingRepo,
		vehicleRepo: vehicleRepo,
	}
}
