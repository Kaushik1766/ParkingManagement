package vehiclerepository

import (
	models "github.com/Kaushik1766/ParkingManagement/internal/models"
	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SQLVehicleRepository struct {
	db *gorm.DB
}

func (sqlvr *SQLVehicleRepository) AddVehicle(numberplate string, userid uuid.UUID, vehicleType vehicletypes.VehicleType) (models.Vehicle, error) {
	vehicle := models.Vehicle{
		NumberPlate: numberplate,
		UserID:      userid,
		VehicleType: vehicleType,
	}

	err := sqlvr.db.Create(&vehicle).Error
	if err != nil {
		return models.Vehicle{}, err
	}

	return vehicle, nil
}

func (sqlvr *SQLVehicleRepository) RemoveVehicle(numberplate string) error {
	err := sqlvr.db.Where("number_plate = ?", numberplate).Delete(&models.Vehicle{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (sqlvr *SQLVehicleRepository) GetVehicleById(vehicleId uuid.UUID) (models.Vehicle, error) {
	var vehicle models.Vehicle
	err := sqlvr.db.Where("vehicle_id = ?", vehicleId).First(&vehicle).Error
	if err != nil {
		return models.Vehicle{}, err
	}
	return vehicle, nil
}

func (sqlvr *SQLVehicleRepository) GetVehiclesByUserId(userId uuid.UUID) ([]models.Vehicle, error) {
	var vehicles []models.Vehicle
	err := sqlvr.db.Where("user_id = ?", userId).Find(&vehicles).Error
	if err != nil {
		return nil, err
	}
	return vehicles, nil
}

func (sqlvr *SQLVehicleRepository) GetVehicleByNumberPlate(numberplate string) (models.Vehicle, error) {
	var vehicle models.Vehicle
	err := sqlvr.db.Where("number_plate = ?", numberplate).First(&vehicle).Error
	if err != nil {
		return models.Vehicle{}, err
	}
	return vehicle, nil
}

func (sqlvr *SQLVehicleRepository) GetVehiclesWithUnassignedSlots() (vehicles []models.Vehicle, err error) {
	err = sqlvr.db.
		Where("assigned_building_id IS NULL AND assigned_floor_number IS NULL AND assigned_slot_number IS NULL").
		Find(&vehicles).Error
	if err != nil {
		return nil, err
	}
	return vehicles, nil
}

func (sqlvr *SQLVehicleRepository) Save(vehicle models.Vehicle) error {
	err := sqlvr.db.Save(&vehicle).Error
	if err != nil {
		return err
	}
	return nil
}

func NewSQLVehicleRepository(db *gorm.DB) *SQLVehicleRepository {
	return &SQLVehicleRepository{
		db: db,
	}
}
