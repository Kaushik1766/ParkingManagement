package parkinghistoryrepository

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	parkinghistory "github.com/Kaushik1766/ParkingManagement/internal/models/parking_history"
	"github.com/Kaushik1766/ParkingManagement/internal/models/slot"
	"github.com/Kaushik1766/ParkingManagement/internal/models/vehicle"
	"github.com/google/uuid"
)

type FileParkingRepository struct {
	*sync.Mutex
	parkings []parkinghistory.ParkingHistory
}

func (fpr *FileParkingRepository) Unpark(id string) error {
	fpr.Lock()
	defer fpr.Unlock()

	for i, parking := range fpr.parkings {
		if parking.ParkingId.String() == id && parking.EndTime.IsZero() {
			fpr.parkings[i].EndTime = time.Now()
			return nil
		}
	}

	return errors.New("parking not found or already ended")
}

func (fpr *FileParkingRepository) GetParkingHistoryByUser(userId string, startTime, endTime time.Time) ([]parkinghistory.ParkingHistoryDTO, error) {
	fpr.Lock()
	defer fpr.Unlock()

	var history []parkinghistory.ParkingHistoryDTO

	for _, parking := range fpr.parkings {
		if parking.UserId.String() == userId && parking.StartTime.After(startTime) && parking.EndTime.Before(endTime) {
			history = append(history, parkinghistory.ParkingHistoryDTO{
				TicketId:    parking.ParkingId.String(),
				NumberPlate: parking.NumberPlate,
				BuildingId:  parking.BuildingId,
				FLoorNumber: parking.FLoorNumber,
				SlotNumber:  parking.SlotNumber,
				StartTime:   parking.StartTime.Local().String(),
				EndTime:     parking.EndTime.Local().String(),
			})
		}
	}
	return history, nil
}

func (fpr *FileParkingRepository) AddParking(vehicle vehicle.Vehicle) (string, error) {
	fpr.Lock()
	defer fpr.Unlock()

	for _, parking := range fpr.parkings {
		if parking.BuildingId == vehicle.AssignedSlot.BuildingId.String() &&
			parking.FLoorNumber == vehicle.AssignedSlot.FloorNumber && parking.SlotNumber == vehicle.AssignedSlot.SlotNumber &&
			parking.EndTime.IsZero() {
			return "", errors.New("vehicle already parked in this slot")
		}
	}

	newParking := parkinghistory.ParkingHistory{
		ParkingId:   uuid.New(),
		NumberPlate: vehicle.NumberPlate,
		UserId:      vehicle.UserId,
		BuildingId:  vehicle.AssignedSlot.BuildingId.String(),
		FLoorNumber: vehicle.AssignedSlot.FloorNumber,
		SlotNumber:  vehicle.AssignedSlot.SlotNumber,
		StartTime:   time.Now(),
		EndTime:     time.Time{},
	}
	log.Printf("Adding new parking: %+v", newParking)
	fpr.parkings = append(fpr.parkings, newParking)

	return newParking.ParkingId.String(), nil
}

func (fpr *FileParkingRepository) GetParkingHistoryByNumberPlate(numberplate string, startTime, endTime time.Time) ([]parkinghistory.ParkingHistoryDTO, error) {
	fpr.Lock()
	defer fpr.Unlock()

	var history []parkinghistory.ParkingHistoryDTO

	for _, parking := range fpr.parkings {
		if parking.NumberPlate == numberplate && parking.StartTime.After(startTime) && parking.EndTime.Before(endTime) {
			history = append(history, parkinghistory.ParkingHistoryDTO{
				TicketId:    parking.ParkingId.String(),
				NumberPlate: parking.NumberPlate,
				BuildingId:  parking.BuildingId,
				FLoorNumber: parking.FLoorNumber,
				SlotNumber:  parking.SlotNumber,
				StartTime:   parking.StartTime.Local().String(),
				EndTime:     parking.EndTime.Local().String(),
			})
		}
	}
	return history, nil
}

func (fpr *FileParkingRepository) GetActiveUserParkings(userId string) ([]parkinghistory.ParkingHistoryDTO, error) {
	fpr.Lock()
	defer fpr.Unlock()

	var activeParkings []parkinghistory.ParkingHistoryDTO

	for _, parking := range fpr.parkings {
		if parking.UserId.String() == userId && parking.EndTime.IsZero() {
			activeParkings = append(activeParkings, parkinghistory.ParkingHistoryDTO{
				TicketId:    parking.ParkingId.String(),
				NumberPlate: parking.NumberPlate,
				BuildingId:  parking.BuildingId,
				FLoorNumber: parking.FLoorNumber,
				SlotNumber:  parking.SlotNumber,
				StartTime:   parking.StartTime.Local().String(),
				EndTime:     parking.EndTime.Local().String(),
			})
		}
	}
	return activeParkings, nil
}

func NewFileParkingHistoryRepository() *FileParkingRepository {
	data, err := os.ReadFile("parkinghistory.json")
	if err != nil {
		os.WriteFile("parkinghistory.json", []byte("[]"), 0666)
		data, err = json.Marshal([]slot.Slot{})
		if err != nil {
			fmt.Println("unable to marshal")
		}
	}

	var parkingData []parkinghistory.ParkingHistory
	err = json.Unmarshal(data, &parkingData)
	if err != nil {
		fmt.Println(err)
		panic("corrupted data")
	}
	return &FileParkingRepository{
		Mutex:    &sync.Mutex{},
		parkings: parkingData,
	}
}

func (fpr *FileParkingRepository) SerializeData() {
	data, err := json.Marshal(fpr.parkings)
	if err != nil {
		fmt.Println("unable to marshal parking data")
		return
	}
	err = os.WriteFile("parkinghistory.json", data, 0666)
	if err != nil {
		fmt.Println("unable to write parking data to file")
	}
}
