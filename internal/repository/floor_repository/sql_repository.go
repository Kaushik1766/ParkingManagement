package floorrepository

import (
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
	return sqlfr.db.Create(&floor).Error
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

func (sqlfr *SQLFloorRepository) GetFloorsByBuildingId(buildingId string) ([]models.Floor, error) {
	buildingUUID, err := uuid.Parse(buildingId)
	if err != nil {
		return nil, err
	}
	var floors []models.Floor
	err = sqlfr.db.Where("building_id = ?", buildingUUID).Find(&floors).Error
	if err != nil {
		return nil, err
	}
	return floors, nil
}

func NewSQLFloorRepository(db *gorm.DB) *SQLFloorRepository {
	return &SQLFloorRepository{
		db: db,
	}
}
