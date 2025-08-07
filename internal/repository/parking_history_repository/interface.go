package parkinghistoryrepository

import (
	"time"

	parkinghistory "github.com/Kaushik1766/ParkingManagement/internal/models/parking_history"
	"github.com/Kaushik1766/ParkingManagement/internal/models/vehicle"
)

type ParkingHistoryStorage interface {
	AddParking(vehicle vehicle.Vehicle) (string, error)
	Unpark(id string) error
	GetParkingHistoryByNumberPlate(numberplate string, startTime, endTime time.Time) ([]parkinghistory.ParkingHistoryDTO, error)
	GetParkingHistoryByUser(userId string, startTime, endTime time.Time) ([]parkinghistory.ParkingHistoryDTO, error)
	GetActiveUserParkings(userId string) ([]parkinghistory.ParkingHistoryDTO, error)
}
