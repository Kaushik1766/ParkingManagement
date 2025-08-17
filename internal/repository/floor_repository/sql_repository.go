package floorrepository

import (
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SQLFloorRepository struct {
	db *gorm.DB
}

func (sqlfr *SQLFloorRepository) AddFloor(buildingId uuid.UUID, floorNumber int) error {
	floor := models.Floor{
		BuildingID:  buildingId,
		FloorNumber: floorNumber,
	}
	return sqlfr.db.Create(&floor).Error
}

func (sqlfr *SQLFloorRepository) DeleteFloor(buildingId uuid.UUID, floorNumber int) error {
	floor := models.Floor{}
	err := sqlfr.db.Where("building_id = ? and floor_number = ?", buildingId, floorNumber).First(&floor).Error
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

func (sqlfr *SQLFloorRepository) GetFloorsByBuildingId(buildingId uuid.UUID) ([]int, error) {
	var floors []models.Floor
	err := sqlfr.db.Where("building_id = ?", buildingId).Find(&floors).Error
	returnedFloors := make([]int, len(floors))
	if err != nil {
		return nil, err
	}
	for i, floor := range floors {
		returnedFloors[i] = floor.FloorNumber
	}
	return returnedFloors, nil
}

func NewSQLFloorRepository(db *gorm.DB) *SQLFloorRepository {
	return &SQLFloorRepository{
		db: db,
	}
}
