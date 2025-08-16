package buildingrepository

import (
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SQLBuildingRepository struct {
	db *gorm.DB
}

func (sqlbr *SQLBuildingRepository) DeleteBuilding(name string) error {
	building := models.Building{}
	err := sqlbr.db.Where("building_name = ?", name).First(&building).Error
	if err != nil {
		return err
	}
	return sqlbr.db.Delete(&building).Error
}

func (sqlbr *SQLBuildingRepository) GetBuildingByName(name string) (models.Building, error) {
	building := models.Building{}
	err := sqlbr.db.Where("building_name = ?", name).First(&building).Error
	return building, err
}

func (sqlbr *SQLBuildingRepository) GetAllBuildings() ([]models.Building, error) {
	buildings := []models.Building{}
	err := sqlbr.db.Find(&buildings).Error
	return buildings, err
}

func (sqlbr *SQLBuildingRepository) GetBuildingByID(buildingID uuid.UUID) (models.Building, error) {
	building := models.Building{}
	err := sqlbr.db.Where("building_id = ?", buildingID).First(&building).Error
	return building, err
}

func NewSQLBuildingRepository(db *gorm.DB) *SQLBuildingRepository {
	return &SQLBuildingRepository{
		db: db,
	}
}

func (sqlbr *SQLBuildingRepository) AddBuilding(buildingName string) error {
	building := models.Building{
		BuildingName: buildingName,
		Floors:       nil,
	}
	return sqlbr.db.Create(&building).Error
}
