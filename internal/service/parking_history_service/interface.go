package parkinghistoryservice

import (
	"context"

	parkinghistory "github.com/Kaushik1766/ParkingManagement/internal/models/parking_history"
)

type ParkingHistoryMgr interface {
	GetParkingHistoryByNumberPlate(ctx context.Context, numberplate string, startTime, endTime string) ([]parkinghistory.ParkingHistoryDTO, error)
	GetParkingHistoryByUser(ctx context.Context, userId string, startTime, endTime string) ([]parkinghistory.ParkingHistoryDTO, error)
}
