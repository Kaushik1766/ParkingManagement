package parkinghistoryrepository

import (
	"time"

	"github.com/Kaushik1766/ParkingManagement/internal/models"
)

//go:generate mockgen -source=interface.go -destination=../../../mocks/parking_history_storage_mock.go -package=mocks
type ParkingHistoryStorage interface {
	AddParking(vehicle models.Vehicle) (string, error)
	Unpark(id string) error
	GetParkingHistoryByNumberPlate(numberplate string, startTime, endTime time.Time) ([]models.ParkingHistoryDTO, error)
	GetParkingHistoryByUser(userId string, startTime, endTime time.Time) ([]models.ParkingHistoryDTO, error)
	GetActiveUserParkings(userId string) ([]models.ParkingHistoryDTO, error)
	UnparkByNumberPlate(numberplate string) error
}
