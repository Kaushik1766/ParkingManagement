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

func (sqlbr *SQLBuildingRepository) DeleteBuildingByName(buildingName string) error {
	if buildingName == constants.AdminBuilding {
		return errors.New("buildingrepo: cannot delete admin building")
	}
	building := models.Building{}
	err := sqlbr.db.Where("building_name = ?", buildingName).First(&building).Error
	if err != nil {
		return err
	}
	return sqlbr.db.Delete(&building).Error
}

func (sqlbr *SQLBuildingRepository) GetBuildingByName(buildingName string) (models.Building, error) {
	if buildingName == constants.AdminBuilding {
		return models.Building{}, errors.New("buildingrepo: cannot get admin building")
	}
	building := models.Building{}
	err := sqlbr.db.Where("building_name = ?", buildingName).First(&building).Error
	return building, err
}

func (sqlbr *SQLBuildingRepository) GetAllBuildings() ([]models.Building, error) {
	buildings := []models.Building{}
	err := sqlbr.db.Where("building_name <> ?", constants.AdminBuilding).Find(&buildings).Error
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
	return sqlbr.db.Create(&building).Error
}
