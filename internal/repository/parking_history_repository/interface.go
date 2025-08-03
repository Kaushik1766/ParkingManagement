package parkinghistoryrepository

import parkinghistory "github.com/Kaushik1766/ParkingManagement/internal/models/parking_history"

type ParkingHistoryStorage interface {
	AddParking(numberplate string) (string, error)
	GetParkingHistoryByNumberPlate(numberplate string, startTime, endTime string) ([]parkinghistory.ParkingHistory, error)
}
