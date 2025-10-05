package parkinghistoryrepository

import (
	"errors"
	"time"

	"github.com/Kaushik1766/ParkingManagement/internal/models"
	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SQLParkingRepository struct {
	db *gorm.DB
}

func (sqlpr *SQLParkingRepository) UnparkByNumberPlate(numberplate string) error {
	var ph models.ParkingHistory
	err := sqlpr.db.Joins("JOIN vehicles ON vehicles.vehicle_id = parking_histories.vehicle_id").
		Where("vehicles.number_plate = ? AND parking_histories.end_time IS NULL", numberplate).
		First(&ph).Error
	if err != nil {
		return err
	}

	now := time.Now()
	ph.EndTime = &now

	return sqlpr.db.Save(&ph).Error
}

func (sqlpr *SQLParkingRepository) AddParking(vehicle models.Vehicle) (string, error) {
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
			FLoorNumber:  parking.Vehicle.AssignedFloorNumber,
			SlotNumber:   parking.Vehicle.AssignedSlotNumber,
			StartTime:    parking.StartTime.Local(),
			EndTime:      parking.EndTime.Local(),
			VechicleType: parking.Vehicle.VehicleType.String(),
		})
	}
	return historyDTO, nil
}

func (sqlpr *SQLParkingRepository) GetParkingHistoryByUser(userId string, startTime time.Time, endTime time.Time) ([]models.ParkingHistoryDTO, error) {
	type Result struct {
		ParkingID    uuid.UUID `gorm:"column:parking_id"`
		StartTime    time.Time `gorm:"column:start_time"`
		EndTime      time.Time `gorm:"column:end_time"`
		NumberPlate  string    `gorm:"column:number_plate"`
		VehicleType  int       `gorm:"column:vehicle_type"`
		SlotNumber   int       `gorm:"column:slot_number"`
		FloorNumber  int       `gorm:"column:floor_number"`
		BuildingID   uuid.UUID `gorm:"column:building_id"`
		BuildingName string    `gorm:"column:building_name"`
	}

	var results []Result

	rawSQL := `
		SELECT
			ph.parking_id,
			ph.start_time,
			ph.end_time,
			v.number_plate,
			v.vehicle_type,
			s.slot_number,
			s.floor_number,
			b.building_id,
			b.building_name
		FROM parking_histories ph
		JOIN vehicles v ON v.vehicle_id = ph.vehicle_id
		JOIN slots s ON s.building_id = v.assigned_building_id
		             AND s.floor_number = v.assigned_floor_number
		             AND s.slot_number = v.assigned_slot_number
		JOIN floors f ON f.building_id = s.building_id
		             AND f.floor_number = s.floor_number
		JOIN buildings b ON b.building_id = f.building_id
		WHERE v.user_id = ? 
		  AND ph.start_time >= ? 
		  AND ph.end_time <= ? 
		  AND ph.end_time IS NOT NULL
		ORDER BY ph.start_time DESC
	`

	err := sqlpr.db.Raw(rawSQL, uuid.MustParse(userId), startTime, endTime).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	// Map results into DTO
	historyDTO := make([]models.ParkingHistoryDTO, 0, len(results))
	for _, r := range results {
		historyDTO = append(historyDTO, models.ParkingHistoryDTO{
			TicketId:     r.ParkingID.String(),
			NumberPlate:  r.NumberPlate,
			BuildingId:   r.BuildingID.String(),
			BuildingName: r.BuildingName,
			FLoorNumber:  r.FloorNumber,
			SlotNumber:   r.SlotNumber,
			StartTime:    r.StartTime.Local(),
			EndTime:      r.EndTime.Local(),
			VechicleType: vehicletypes.VehicleType(r.VehicleType).String(),
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
		activeParkingsDTO = append(activeParkingsDTO, models.ParkingHistoryDTO{
			TicketId:     parking.ParkingID.String(),
			NumberPlate:  parking.Vehicle.NumberPlate,
			BuildingId:   parking.Vehicle.AssignedSlot.String(),
			FLoorNumber:  parking.Vehicle.AssignedSlot.FloorNumber,
			SlotNumber:   parking.Vehicle.AssignedSlot.SlotNumber,
			StartTime:    parking.StartTime.Local(),
			VechicleType: parking.Vehicle.VehicleType.String(),
		})
	}
	return activeParkingsDTO, nil
}

func NewSQLParkingRepository(db *gorm.DB) *SQLParkingRepository {
	return &SQLParkingRepository{
		db: db,
	}
}
