package parkinghistoryservice

import (
	"context"
	"time"

	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	parkinghistory "github.com/Kaushik1766/ParkingManagement/internal/models/parking_history"
	userjwt "github.com/Kaushik1766/ParkingManagement/internal/models/user_jwt"
	parkinghistoryrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/parking_history_repository"
	vehiclerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/vehicle_repository"
	customerrors "github.com/Kaushik1766/ParkingManagement/pkg/customErrors"
)

type ParkingHistoryService struct {
	parkingRepo parkinghistoryrepository.ParkingHistoryStorage
	vehicleRepo vehiclerepository.VehicleStorage
}

func (phs *ParkingHistoryService) GetParkingHistoryById(userId string, startTime, endTime time.Time) ([]parkinghistory.ParkingHistoryDTO, error) {
	// startDate, err := time.Parse(time.DateOnly, startTime)
	// if err != nil {
	// 	return []parkinghistory.ParkingHistoryDTO{}, err
	// }
	//
	// endDate, err := time.Parse(time.DateOnly, endTime)
	// if err != nil {
	// 	return []parkinghistory.ParkingHistoryDTO{}, err
	// }

	parkingHistory, err := phs.parkingRepo.GetParkingHistoryByUser(userId, startTime, endTime)
	if err != nil {
		return []parkinghistory.ParkingHistoryDTO{}, err
	}

	return parkingHistory, nil
}

func (phs *ParkingHistoryService) GetParkingHistory(ctx context.Context, startTime, endTime time.Time) ([]parkinghistory.ParkingHistoryDTO, error) {
	userCtx := ctx.Value(constants.User).(userjwt.UserJwt)

	// startTimeParsed, err := time.Parse(time.DateOnly, startTime)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// endTimeParsed, err := time.Parse(time.DateOnly, endTime)
	// if err != nil {
	// 	return nil, err
	// }

	parkingHistory, err := phs.parkingRepo.GetParkingHistoryByUser(userCtx.ID, startTime, endTime)
	if err != nil {
		return nil, err
	}

	return parkingHistory, nil
}

func (phs *ParkingHistoryService) GetParkingHistoryByNumberPlate(ctx context.Context, numberplate string, startTime string, endTime string) ([]parkinghistory.ParkingHistoryDTO, error) {
	userCtx := ctx.Value(constants.User).(userjwt.UserJwt)

	vehicle, err := phs.vehicleRepo.GetVehicleByNumberPlate(numberplate)
	if err != nil {
		return nil, err
	}

	if userCtx.Role != roles.Admin && userCtx.ID != vehicle.UserId.String() {
		return nil, customerrors.Unathorized{}
	}

	startTimeParsed, err := time.Parse(time.DateOnly, startTime)
	if err != nil {
		return nil, err
	}

	endTimeParsed, err := time.Parse(time.DateOnly, endTime)
	if err != nil {
		return nil, err
	}

	parkingHistory, err := phs.parkingRepo.GetParkingHistoryByNumberPlate(numberplate, startTimeParsed, endTimeParsed)
	if err != nil {
		return nil, err
	}

	return parkingHistory, err
}

func (phs *ParkingHistoryService) GetParkingHistoryByUser(ctx context.Context, userId string, startTime string, endTime string) ([]parkinghistory.ParkingHistoryDTO, error) {
	userCtx := ctx.Value(constants.User).(userjwt.UserJwt)

	if userCtx.Role != roles.Admin && userCtx.ID != userId {
		return nil, customerrors.Unathorized{}
	}

	startTimeParsed, err := time.Parse(time.DateOnly, startTime)
	if err != nil {
		return nil, err
	}

	endTimeParsed, err := time.Parse(time.DateOnly, endTime)
	if err != nil {
		return nil, err
	}

	parkingHistory, err := phs.parkingRepo.GetParkingHistoryByUser(userId, startTimeParsed, endTimeParsed)
	if err != nil {
		return nil, err
	}

	return parkingHistory, nil
}

func (phs *ParkingHistoryService) GetActiveUserParkings(ctx context.Context) ([]parkinghistory.ParkingHistoryDTO, error) {
	userCtx := ctx.Value(constants.User).(userjwt.UserJwt)

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
