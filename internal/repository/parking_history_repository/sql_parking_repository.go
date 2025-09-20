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

func (sqlpr *SQLParkingRepository) UnparkByNumberPlate(numberplate string) error {
	// Using raw SQL for better performance with JOIN in UPDATE
	err := sqlpr.db.Exec(`
		UPDATE parking_histories 
		SET end_time = ? 
		FROM vehicles 
		WHERE vehicles.vehicle_id = parking_histories.vehicle_id 
		AND vehicles.number_plate = ? 
		AND parking_histories.end_time IS NULL`,
		time.Now(), numberplate).Error
	return err
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
	var historyDTO []models.ParkingHistoryDTO

	rows, err := sqlpr.db.Raw(`
		SELECT 
			ph.parking_id,
			v.number_plate,
			v.assigned_building_id,
			v.assigned_floor_number,
			v.assigned_slot_number,
			ph.start_time,
			ph.end_time,
			v.vehicle_type
		FROM parking_histories ph
		JOIN vehicles v ON v.vehicle_id = ph.vehicle_id
		WHERE v.number_plate = ? 
		AND ph.start_time >= ? 
		AND ph.end_time <= ? 
		AND ph.end_time IS NOT NULL
		ORDER BY ph.start_time DESC`,
		numberplate, startTime, endTime).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var dto models.ParkingHistoryDTO
		var buildingID uuid.UUID
		var floorNumber, slotNumber int

		err := rows.Scan(
			&dto.TicketId,
			&dto.NumberPlate,
			&buildingID,
			&floorNumber,
			&slotNumber,
			&dto.StartTime,
			&dto.EndTime,
			&dto.VechicleType,
		)
		if err != nil {
			return nil, err
		}

		dto.BuildingId = buildingID.String()
		dto.FLoorNumber = floorNumber
		dto.SlotNumber = slotNumber
		dto.StartTime = dto.StartTime.Local()
		dto.EndTime = dto.EndTime.Local()

		historyDTO = append(historyDTO, dto)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return historyDTO, nil
}

func (sqlpr *SQLParkingRepository) GetParkingHistoryByUser(userId string, startTime time.Time, endTime time.Time) ([]models.ParkingHistoryDTO, error) {
	var historyDTO []models.ParkingHistoryDTO

	rows, err := sqlpr.db.Raw(`
		SELECT 
			ph.parking_id,
			v.number_plate,
			v.assigned_building_id,
			v.assigned_floor_number,
			v.assigned_slot_number,
			ph.start_time,
			ph.end_time,
			v.vehicle_type
		FROM parking_histories ph
		JOIN vehicles v ON v.vehicle_id = ph.vehicle_id
		WHERE v.user_id = ? 
		AND ph.start_time >= ? 
		AND ph.end_time <= ?
		AND ph.end_time IS NOT NULL
		ORDER BY ph.start_time DESC`,
		uuid.MustParse(userId), startTime, endTime).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var dto models.ParkingHistoryDTO
		var buildingID uuid.UUID
		var floorNumber, slotNumber int

		err := rows.Scan(
			&dto.TicketId,
			&dto.NumberPlate,
			&buildingID,
			&floorNumber,
			&slotNumber,
			&dto.StartTime,
			&dto.EndTime,
			&dto.VechicleType,
		)
		if err != nil {
			return nil, err
		}

		dto.BuildingId = buildingID.String()
		dto.FLoorNumber = floorNumber
		dto.SlotNumber = slotNumber
		dto.StartTime = dto.StartTime.Local()
		dto.EndTime = dto.EndTime.Local()

		historyDTO = append(historyDTO, dto)
	}

	if err = rows.Err(); err != nil {
		return nil, err
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
