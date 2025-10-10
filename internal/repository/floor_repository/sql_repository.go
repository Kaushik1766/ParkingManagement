package floorrepository

import (
	"errors"

	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SQLFloorRepository struct {
	db *gorm.DB
}

func (sqlfr *SQLFloorRepository) AddFloor(buildingId string, floorNumber int) error {
	var slots []models.Slot

	buildingUUID, err := uuid.Parse(buildingId)
	if err != nil {
		return err
	}

	for i, s := range constants.SlotLayout {
		if s == '0' {
			slots = append(slots, models.Slot{
				BuildingID:  buildingUUID,
				FloorNumber: floorNumber,
				SlotNumber:  i,
				SlotType:    vehicletypes.TwoWheeler,
			})
		} else {
			slots = append(slots, models.Slot{
				BuildingID:  buildingUUID,
				FloorNumber: floorNumber,
				SlotNumber:  i,
				SlotType:    vehicletypes.FourWheeler,
			})
		}
	}

	floor := models.Floor{
		BuildingID:  buildingUUID,
		FloorNumber: floorNumber,
		Slots:       slots,
	}
	err = sqlfr.db.Create(&floor).Error
	if err != nil {
		return errors.New("duplicate floor number not allowed")
	}
	return nil
}

func (sqlfr *SQLFloorRepository) DeleteFloor(buildingId string, floorNumber int) error {
	buildingUUID, err := uuid.Parse(buildingId)
	if err != nil {
		return err
	}
	floor := models.Floor{}
	err = sqlfr.db.Where("building_id = ? and floor_number = ?", buildingUUID, floorNumber).First(&floor).Error
	if err != nil {
		return err
	}

	return sqlfr.db.Delete(&floor).Error
}

func (sqlfr *SQLFloorRepository) GetFloor(buildingId uuid.UUID, floorNumber int) (int, error) {
	var floor models.Floor
	err := sqlfr.db.Where("building_id = ? and floor_number = ?", buildingId, floorNumber).First(&floor).Error
	if err != nil {
		return 0, err
	}
	return floor.FloorNumber, nil
}

func (sqlfr *SQLFloorRepository) GetFloorsByBuildingId(buildingId string) ([]models.FloorSummary, error) {
	buildingUUID, err := uuid.Parse(buildingId)
	if err != nil {
		return nil, err
	}
	var floors []models.FloorSummary
	//err = sqlfr.db.Where("building_id = ?", buildingUUID).Find(&floors).Error
	err = sqlfr.db.
		Table("floors").
		Select(`
			floors.building_id as building_id,
			floors.floor_number as floor_number,
			offices.office_name as assigned_office,
			count(distinct case when slots.slot_number is not null then
				(slots.building_id, slots.floor_number, slots.slot_number) end) as total_slots,
			count(distinct case when vehicles.assigned_slot_number is null then
				(slots.building_id,slots.floor_number, slots.slot_number) end) as available_slots
		`).
		Joins("left join slots on floors.building_id = slots.building_id and floors.floor_number = slots.floor_number").
		Joins("left join offices on offices.building_id = floors.building_id and offices.floor_number = floors.floor_number").
		Joins("left join vehicles on slots.building_id = vehicles.assigned_building_id and slots.floor_number = vehicles.assigned_floor_number and slots.slot_number = vehicles.assigned_slot_number").
		Where("floors.building_id = ?", buildingUUID).
		Group("floors.building_id, floors.floor_number, offices.office_id").
		Find(&floors).Error
	return floors, err
}

func NewSQLFloorRepository(db *gorm.DB) *SQLFloorRepository {
	return &SQLFloorRepository{
		db: db,
	}
}
