package parkinghistoryrepository

import parkinghistory "github.com/Kaushik1766/ParkingManagement/internal/models/parking_history"

type FileParkingRepository struct{}

func (fpr *FileParkingRepository) AddParking(numberplate string) (string, error) {
	panic("not implemented") // TODO: Implement
}

func (fpr *FileParkingRepository) GetParkingHistoryByNumberPlate(numberplate string, startTime string, endTime string) ([]parkinghistory.ParkingHistory, error) {
	panic("not implemented") // TODO: Implement
}
