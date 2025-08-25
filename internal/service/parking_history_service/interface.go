package parkinghistoryservice

import (
	"context"
	"time"

	models "github.com/Kaushik1766/ParkingManagement/internal/models"
)

type ParkingHistoryMgr interface {
	GetParkingHistoryByNumberPlate(ctx context.Context, numberplate string, startTime, endTime time.Time) ([]models.ParkingHistoryDTO, error)
	GetActiveUserParkings(ctx context.Context) ([]models.ParkingHistoryDTO, error)
	GetParkingHistory(ctx context.Context, startTime, endTime time.Time) ([]models.ParkingHistoryDTO, error)
	GetParkingHistoryByUserId(userId string, startTime, endTime time.Time) ([]models.ParkingHistoryDTO, error)
}
