package buildingrepository

import (
	"errors"

	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SQLBuildingRepository struct {
	db *gorm.DB
}

func (sqlbr *SQLBuildingRepository) DeleteBuildingByID(buildingID string) error {
	buildingUUID, err := uuid.Parse(buildingID)
	if err != nil {
		return err
	}

	err = sqlbr.db.
		Delete(&models.Building{
			BuildingID: buildingUUID,
		}).
		Error
	return err
}

func (sqlbr *SQLBuildingRepository) GetBuildingByName(buildingName string) (models.Building, error) {
	if buildingName == constants.AdminBuilding {
		return models.Building{}, errors.New("buildingrepo: cannot get admin building")
	}
	building := models.Building{}
	err := sqlbr.db.Where("building_name = ?", buildingName).First(&building).Error
	return building, err
}

func (sqlbr *SQLBuildingRepository) GetAllBuildingSummary() ([]models.BuildingSummary, error) {
	var buildings []models.BuildingSummary
	err := sqlbr.db.
		Table("buildings").
		Select(`
			buildings.building_id as building_id,
			buildings.building_name as building_name,
			count(distinct case when floors.floor_number is not null then
				(floors.building_id, floors.floor_number) end) as total_floors,
			count(distinct case when slots.slot_number is not null then
				(slots.building_id,slots.floor_number, slots.slot_number) end) as total_slots,
			count(distinct case when vehicles.assigned_slot_number is null then
				(slots.building_id,slots.floor_number, slots.slot_number) end) as available_slots
		`).
		Joins("left join floors on buildings.building_id = floors.building_id").
		Joins("left join slots on floors.building_id = slots.building_id and floors.floor_number = slots.floor_number").
		Joins("left join vehicles on slots.building_id = vehicles.assigned_building_id and slots.floor_number = vehicles.assigned_floor_number and slots.slot_number = vehicles.assigned_slot_number").
		Group("buildings.building_id").
		Where("buildings.building_name <> 'ADMIN_BUILDING'").
		Find(&buildings).Error

	return buildings, err
}

func (sqlbr *SQLBuildingRepository) GetAllBuildings() ([]models.Building, error) {
	var buildings []models.Building
	err := sqlbr.db.
		Where("building_name <> ?", constants.AdminBuilding).
		Preload("Floors.Slots.Vehicles").
		Find(&buildings).Error
	return buildings, err
}

func (sqlbr *SQLBuildingRepository) GetBuildingByID(buildingID uuid.UUID) (models.Building, error) {
	building := models.Building{}
	err := sqlbr.db.
		Where("building_id = ? AND building_name <> ?", buildingID, constants.AdminBuilding).
		First(&building).Error
	return building, err
}

func NewSQLBuildingRepository(db *gorm.DB) *SQLBuildingRepository {
	return &SQLBuildingRepository{
		db: db,
	}
}

func (sqlbr *SQLBuildingRepository) AddBuilding(buildingName string) error {
	if buildingName == constants.AdminBuilding {
		return errors.New("buildingrepo: cannot add admin building")
	}
	building := models.Building{
		BuildingName: buildingName,
		Floors:       nil,
	}
	err := sqlbr.db.Create(&building).Error
	if err != nil {
		return errors.New("building name should be unique")
	}
	return nil
}
