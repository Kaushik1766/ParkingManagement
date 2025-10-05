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
	return phs.parkingRepo.GetParkingHistoryByUser(userId, startTime, endTime)
}

// GetParkingHistory retrieves parking history for the authenticated user within the specified time range.
func (phs *ParkingHistoryService) GetParkingHistory(ctx context.Context, startTime, endTime time.Time) ([]models.ParkingHistoryDTO, error) {
	userCtx := ctx.Value(constants.User).(models.UserJwt)

	return phs.parkingRepo.GetParkingHistoryByUser(userCtx.ID, startTime, endTime)
}

func (phs *ParkingHistoryService) GetParkingHistoryByNumberPlate(ctx context.Context, numberplate string, startTime, endTime time.Time) ([]models.ParkingHistoryDTO, error) {
	userCtx := ctx.Value(constants.User).(models.UserJwt)

	vehicle, err := phs.vehicleRepo.GetVehicleByNumberPlate(numberplate)
	if err != nil {
		return nil, err
	}

	if userCtx.Role != roles.Admin && userCtx.ID != vehicle.UserID.String() {
		return nil, customerrors.Unauthorized{}
	}

	return phs.parkingRepo.GetParkingHistoryByNumberPlate(numberplate, startTime, endTime)
}

func (phs *ParkingHistoryService) GetActiveUserParkings(ctx context.Context) ([]models.ParkingHistoryDTO, error) {
	userCtx := ctx.Value(constants.User).(models.UserJwt)

	return phs.parkingRepo.GetActiveUserParkings(userCtx.ID)
}

func NewParkingHistoryService(parkingRepo parkinghistoryrepository.ParkingHistoryStorage, vehicleRepo vehiclerepository.VehicleStorage) *ParkingHistoryService {
	return &ParkingHistoryService{
		parkingRepo: parkingRepo,
		vehicleRepo: vehicleRepo,
	}
}
