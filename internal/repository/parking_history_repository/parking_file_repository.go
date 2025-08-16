package parkinghistoryrepository

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/Kaushik1766/ParkingManagement/internal/config"
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
		if parking.ParkingID.String() == id && parking.EndTime.IsZero() {
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
		if parking.UserID.String() == userId && (parking.StartTime.After(startTime) || parking.StartTime.Equal(startTime)) && (parking.EndTime.Before(endTime) || parking.EndTime.Equal(endTime)) {
			history = append(history, parkinghistory.ParkingHistoryDTO{
				TicketId:     parking.ParkingID.String(),
				NumberPlate:  parking.NumberPlate,
				BuildingId:   parking.BuildingID,
				FLoorNumber:  parking.FloorNumber,
				SlotNumber:   parking.SlotNumber,
				StartTime:    parking.StartTime.Local(),
				EndTime:      parking.EndTime.Local(),
				VechicleType: parking.VehicleType,
			})
		}
	}
	return history, nil
}

func (fpr *FileParkingRepository) AddParking(vehicle vehicle.Vehicle) (string, error) {
	fpr.Lock()
	defer fpr.Unlock()

	if vehicle.AssignedSlot.BuildingID == uuid.Nil {
		return "", errors.New("no slot assigned contact the admin")
	}

	for _, parking := range fpr.parkings {
		if parking.BuildingID == vehicle.AssignedSlot.BuildingID.String() &&
			parking.FloorNumber == vehicle.AssignedSlot.FloorNumber && parking.SlotNumber == vehicle.AssignedSlot.SlotNumber &&
			parking.EndTime.IsZero() {
			return "", errors.New("vehicle already parked in this slot")
		}
	}

	newParking := parkinghistory.ParkingHistory{
		ParkingID:   uuid.New(),
		NumberPlate: vehicle.NumberPlate,
		UserID:      vehicle.UserID,
		BuildingID:  vehicle.AssignedSlot.BuildingID.String(),
		FloorNumber: vehicle.AssignedSlot.FloorNumber,
		SlotNumber:  vehicle.AssignedSlot.SlotNumber,
		StartTime:   time.Now(),
		EndTime:     time.Time{},
		VehicleType: vehicle.VehicleType,
	}
	log.Printf("Adding new parking: %+v", newParking)
	fpr.parkings = append(fpr.parkings, newParking)

	return newParking.ParkingID.String(), nil
}

func (fpr *FileParkingRepository) GetParkingHistoryByNumberPlate(numberplate string, startTime, endTime time.Time) ([]parkinghistory.ParkingHistoryDTO, error) {
	fpr.Lock()
	defer fpr.Unlock()

	var history []parkinghistory.ParkingHistoryDTO

	for _, parking := range fpr.parkings {
		if parking.NumberPlate == numberplate && parking.StartTime.After(startTime) && parking.EndTime.Before(endTime) {
			history = append(history, parkinghistory.ParkingHistoryDTO{
				TicketId:     parking.ParkingID.String(),
				NumberPlate:  parking.NumberPlate,
				BuildingId:   parking.BuildingID,
				FLoorNumber:  parking.FloorNumber,
				SlotNumber:   parking.SlotNumber,
				StartTime:    parking.StartTime.Local(),
				EndTime:      parking.EndTime.Local(),
				VechicleType: parking.VehicleType,
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
		if parking.UserID.String() == userId && parking.EndTime.IsZero() {
			activeParkings = append(activeParkings, parkinghistory.ParkingHistoryDTO{
				TicketId:     parking.ParkingID.String(),
				NumberPlate:  parking.NumberPlate,
				BuildingId:   parking.BuildingID,
				FLoorNumber:  parking.FloorNumber,
				SlotNumber:   parking.SlotNumber,
				StartTime:    parking.StartTime.Local(),
				EndTime:      parking.EndTime.Local(),
				VechicleType: parking.VehicleType,
			})
		}
	}
	return activeParkings, nil
}

func NewFileParkingHistoryRepository() *FileParkingRepository {
	data, err := os.ReadFile(config.ParkingHistoryPath)
	if err != nil {
		os.WriteFile(config.ParkingHistoryPath, []byte("[]"), 0666)
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
	err = os.WriteFile(config.ParkingHistoryPath, data, 0666)
	if err != nil {
		fmt.Println("unable to write parking data to file")
	}
}
