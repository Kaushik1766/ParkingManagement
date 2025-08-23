package officerepository

import (
	"errors"

	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SQLOfficeRepository struct {
	db *gorm.DB
}

func NewSQLOfficeRepository(db *gorm.DB) *SQLOfficeRepository {
	return &SQLOfficeRepository{
		db: db,
	}
}

func (sqlor *SQLOfficeRepository) AddOffice(officeName string, buildingID uuid.UUID, floorNumber int) error {
	office := models.Office{
		BuildingID:  buildingID,
		FloorNumber: floorNumber,
		OfficeName:  officeName,
	}

	return sqlor.db.Create(&office).Error
}

func (sqlor *SQLOfficeRepository) DeleteOffice(officeName string) error {
	if officeName == constants.AdminOffice {
		return errors.New("officerepo: cannot delete admin office")
	}
	return sqlor.db.Where("office_name = ?", officeName).Delete(&models.Office{}).Error
}

func (sqlor *SQLOfficeRepository) GetBuildingAndFloorByOffice(officeName string) (uuid.UUID, int, error) {
	var office models.Office
	err := sqlor.db.Where("office_name = ?", officeName).First(&office).Error
	if err != nil {
		return uuid.UUID{}, 0, err
	}

	return office.BuildingID, office.FloorNumber, nil
}

func (sqlor *SQLOfficeRepository) GetOfficesByBuilding(buildingID uuid.UUID) ([]models.Office, error) {
	var offices []models.Office

	err := sqlor.db.Where("building_id = ? AND office_name <> ?", buildingID, constants.AdminOffice).Find(&offices).Error
	if err != nil {
		return nil, err
	}
	return offices, nil
}

func (sqlor *SQLOfficeRepository) GetAllOffices() ([]models.Office, error) {
	var offices []models.Office
	err := sqlor.db.Where("office_name <> ?", constants.AdminOffice).Find(&offices).Error
	if err != nil {
		return nil, err
	}
	return offices, nil
}

func (sqlor *SQLOfficeRepository) GetOfficeByName(officeName string) (models.Office, error) {
	// if officeName == constants.AdminOffice {
	// 	return models.Office{}, errors.New("officerepo: cannot get admin office by name")
	// }
	var office models.Office
	err := sqlor.db.Where("office_name = ?", officeName).First(&office).Error
	return office, err
}
