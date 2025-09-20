package vehiclerepository

import (
	"database/sql"

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
		NumberPlate:         numberplate,
		UserID:              userid,
		VehicleType:         vehicleType,
		AssignedBuildingID:  nil,
		AssignedFloorNumber: nil,
		AssignedSlotNumber:  nil,
	}

	err := sqlvr.db.Create(&vehicle).Error
	if err != nil {
		return models.Vehicle{}, err
	}

	return vehicle, nil
}

func (sqlvr *SQLVehicleRepository) RemoveVehicle(numberplate string) error {
	err := sqlvr.db.Model(&models.Vehicle{}).Where("number_plate = ?", numberplate).Update("is_active", false).Error
	if err != nil {
		return err
	}
	return nil
}

func (sqlvr *SQLVehicleRepository) GetVehicleById(vehicleId uuid.UUID) (models.Vehicle, error) {
	var vehicle models.Vehicle
	err := sqlvr.db.Where("vehicle_id = ? AND is_active = ?", vehicleId, true).First(&vehicle).Error
	if err != nil {
		return models.Vehicle{}, err
	}
	return vehicle, nil
}

func (sqlvr *SQLVehicleRepository) GetVehiclesByUserId(userId uuid.UUID) ([]models.Vehicle, error) {
	var vehicles []models.Vehicle

	rows, err := sqlvr.db.Raw(`
	select v.vehicle_id, v.number_plate, v.user_id, v.vehicle_type, 
	       v.assigned_building_id, v.assigned_floor_number, v.assigned_slot_number, v.is_active,
	       s.building_id, s.floor_number, s.slot_number, s.slot_type
	from vehicles as v
	left join slots as s on s.slot_number = v.assigned_slot_number
	and s.floor_number = v.assigned_floor_number
	and s.building_id = v.assigned_building_id
	where v.user_id = ? AND v.is_active = true`, userId).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var vehicle models.Vehicle
		var slot models.Slot
		var slotBuildingID sql.Null[uuid.UUID]
		var slotFloorNumber sql.Null[int]
		var slotNumber sql.Null[int]
		var slotType sql.Null[vehicletypes.VehicleType]

		err := rows.Scan(
			&vehicle.VehicleID,
			&vehicle.NumberPlate,
			&vehicle.UserID,
			&vehicle.VehicleType,
			&vehicle.AssignedBuildingID,
			&vehicle.AssignedFloorNumber,
			&vehicle.AssignedSlotNumber,
			&vehicle.IsActive,
			&slotBuildingID,
			&slotFloorNumber,
			&slotNumber,
			&slotType,
		)
		if err != nil {
			return nil, err
		}

		if slotBuildingID.Valid && slotFloorNumber.Valid && slotNumber.Valid && slotType.Valid {
			slot.BuildingID = slotBuildingID.V
			slot.FloorNumber = slotFloorNumber.V
			slot.SlotNumber = slotNumber.V
			slot.SlotType = slotType.V
			vehicle.AssignedSlot = &slot
		} else {
			vehicle.AssignedSlot = nil
		}

		vehicles = append(vehicles, vehicle)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return vehicles, nil
}

func (sqlvr *SQLVehicleRepository) GetVehicleByNumberPlate(numberplate string) (models.Vehicle, error) {
	var vehicle models.Vehicle
	var slot models.Slot
	var slotBuildingID sql.Null[uuid.UUID]
	var slotFloorNumber sql.Null[int]
	var slotNumber sql.Null[int]
	var slotType sql.Null[vehicletypes.VehicleType]

	err := sqlvr.db.Raw(`
	select v.vehicle_id, v.number_plate, v.user_id, v.vehicle_type, 
	       v.assigned_building_id, v.assigned_floor_number, v.assigned_slot_number, v.is_active,
	       s.building_id, s.floor_number, s.slot_number, s.slot_type
	from vehicles as v
	left join slots as s on s.slot_number = v.assigned_slot_number
	and s.floor_number = v.assigned_floor_number
	and s.building_id = v.assigned_building_id
	where v.number_plate = ? AND v.is_active = true`, numberplate).Row().Scan(
		&vehicle.VehicleID,
		&vehicle.NumberPlate,
		&vehicle.UserID,
		&vehicle.VehicleType,
		&vehicle.AssignedBuildingID,
		&vehicle.AssignedFloorNumber,
		&vehicle.AssignedSlotNumber,
		&vehicle.IsActive,
		&slotBuildingID,
		&slotFloorNumber,
		&slotNumber,
		&slotType,
	)

	if err != nil {
		return models.Vehicle{}, err
	}

	// Only populate slot if it exists (not NULL)
	if slotBuildingID.Valid && slotFloorNumber.Valid && slotNumber.Valid && slotType.Valid {
		slot.BuildingID = slotBuildingID.V
		slot.FloorNumber = slotFloorNumber.V
		slot.SlotNumber = slotNumber.V
		slot.SlotType = slotType.V
		vehicle.AssignedSlot = &slot
	} else {
		vehicle.AssignedSlot = nil
	}

	return vehicle, nil
}

func (sqlvr *SQLVehicleRepository) GetVehiclesWithUnassignedSlots() (vehicles []models.Vehicle, err error) {
	err = sqlvr.db.
		Where("assigned_building_id IS NULL AND assigned_floor_number IS NULL AND assigned_slot_number IS NULL AND is_active = ?", true).
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
