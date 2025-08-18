package parkinghistoryrepository

import (
	"errors"
	"time"

	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SQLParkingRepository struct {
	db *gorm.DB
}

func (sqlpr *SQLParkingRepository) AddParking(vehicle models.Vehicle) (string, error) {
	if vehicle.AssignedSlot == nil {
		return "", errors.New("parkingrepo: vehicle does not have an assigned slot")
	}

	var pastParking models.ParkingHistory
	err := sqlpr.db.Joins("left join vehicles on parking_histories.vehicle_id = vehicles.vehicle_id").
		Where("vehicles.assigned_building_id = ? AND vehicles.assigned_floor_number = ? AND vehicles.assigned_slot_number = ? AND end_time IS NULL", vehicle.AssignedBuildingID, vehicle.AssignedFloorNumber, vehicle.AssignedSlot.SlotNumber).
		First(&pastParking).Error
	// fmt.Println(pastParking)

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// fmt.Println(err)
		return "", errors.New("parkingrepo: vehicle is already parked")
	}

	parking := models.ParkingHistory{
		Vehicle: vehicle,
	}
	err = sqlpr.db.Create(&parking).Error
	return parking.ParkingID.String(), err
}

func (sqlpr *SQLParkingRepository) Unpark(id string) error {
	err := sqlpr.db.Model(&models.ParkingHistory{}).
		Where("parking_id = ? AND end_time IS NULL", id).
		Update("end_time", time.Now()).Error
	return err
}

func (sqlpr *SQLParkingRepository) GetParkingHistoryByNumberPlate(numberplate string, startTime time.Time, endTime time.Time) ([]models.ParkingHistoryDTO, error) {
	var history []models.ParkingHistory
	err := sqlpr.db.Where("number_plate = ? AND start_time >= ? AND end_time <= ? AND end_time IS NOT NULL", numberplate, startTime, endTime).
		Preload("Vehicle").
		Find(&history).Error
	if err != nil {
		return nil, err
	}
	var historyDTO []models.ParkingHistoryDTO
	for _, parking := range history {
		historyDTO = append(historyDTO, models.ParkingHistoryDTO{
			TicketId:     parking.ParkingID.String(),
			NumberPlate:  parking.Vehicle.NumberPlate,
			BuildingId:   parking.Vehicle.AssignedBuildingID.String(),
			FLoorNumber:  *parking.Vehicle.AssignedFloorNumber,
			SlotNumber:   *parking.Vehicle.AssignedSlotNumber,
			StartTime:    parking.StartTime.Local(),
			EndTime:      parking.EndTime.Local(),
			VechicleType: parking.Vehicle.VehicleType,
		})
	}
	return historyDTO, nil
}

func (sqlpr *SQLParkingRepository) GetParkingHistoryByUser(userId string, startTime time.Time, endTime time.Time) ([]models.ParkingHistoryDTO, error) {
	var history []models.ParkingHistory
	err := sqlpr.db.Joins("left join vehicles on vehicles.vehicle_id = parking_histories.vehicle_id").Where("vehicles.user_id = ? AND start_time >= ? AND end_time <= ? AND end_time IS NOT NULL", uuid.MustParse(userId), startTime, endTime).Preload("Vehicle.AssignedSlot").Find(&history).Error
	if err != nil {
		return nil, err
	}
	var historyDTO []models.ParkingHistoryDTO
	for _, parking := range history {
		historyDTO = append(historyDTO, models.ParkingHistoryDTO{
			TicketId:     parking.ParkingID.String(),
			NumberPlate:  parking.Vehicle.NumberPlate,
			BuildingId:   parking.Vehicle.AssignedSlot.BuildingID.String(),
			FLoorNumber:  parking.Vehicle.AssignedSlot.FloorNumber,
			SlotNumber:   parking.Vehicle.AssignedSlot.SlotNumber,
			StartTime:    parking.StartTime.Local(),
			EndTime:      parking.EndTime.Local(),
			VechicleType: parking.Vehicle.VehicleType,
		})
	}
	return historyDTO, nil
}

func (sqlpr *SQLParkingRepository) GetActiveUserParkings(userId string) ([]models.ParkingHistoryDTO, error) {
	var activeParkings []models.ParkingHistory
	err := sqlpr.db.Joins("left join vehicles on vehicles.vehicle_id = parking_histories.vehicle_id").Where("vehicles.user_id = ? AND end_time IS NULL", uuid.MustParse(userId)).Preload("Vehicle.AssignedSlot").Find(&activeParkings).Error
	if err != nil {
		return nil, err
	}
	var activeParkingsDTO []models.ParkingHistoryDTO
	for _, parking := range activeParkings {
		if parking.Vehicle.AssignedSlot != nil {
			activeParkingsDTO = append(activeParkingsDTO, models.ParkingHistoryDTO{
				TicketId:     parking.ParkingID.String(),
				NumberPlate:  parking.Vehicle.NumberPlate,
				BuildingId:   parking.Vehicle.AssignedSlot.String(),
				FLoorNumber:  parking.Vehicle.AssignedSlot.FloorNumber,
				SlotNumber:   parking.Vehicle.AssignedSlot.SlotNumber,
				StartTime:    parking.StartTime.Local(),
				VechicleType: parking.Vehicle.VehicleType,
			})
		}
	}
	return activeParkingsDTO, nil
}

func NewSQLParkingRepository(db *gorm.DB) *SQLParkingRepository {
	return &SQLParkingRepository{
		db: db,
	}
}
