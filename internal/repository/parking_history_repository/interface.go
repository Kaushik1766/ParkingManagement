package parkinghistoryrepository

import (
	"context"
	"time"

	"github.com/Kaushik1766/ParkingManagement/internal/models"
)

//go:generate mockgen -source=interface.go -destination=../../../mocks/parking_history_storage_mock.go -package=mocks
type ParkingHistoryStorage interface {
	AddParking(ctx context.Context, vehicle models.Vehicle) (string, error)
	Unpark(ctx context.Context, id string) error
	GetParkingHistoryByNumberPlate(ctx context.Context, numberplate string, startTime, endTime time.Time) ([]models.ParkingHistoryDTO, error)
	GetParkingHistoryByUser(ctx context.Context, userId string, startTime, endTime time.Time) ([]models.ParkingHistoryDTO, error)
	UnparkByNumberPlate(ctx context.Context, numberplate string) error
}
