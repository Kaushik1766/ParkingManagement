package userrepository

import (
	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	models "github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SQLUserRepository struct {
	db *gorm.DB
}

func (sqlur *SQLUserRepository) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := sqlur.db.Where("email = ? AND is_active = true", email).Preload("Office").First(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (sqlur *SQLUserRepository) GetUserById(id string) (models.User, error) {
	var user models.User
	uid, err := uuid.Parse(id)
	if err != nil {
		return models.User{}, err
	}
	err = sqlur.db.Where("user_id = ? AND is_active = true", uid).Preload("Office").First(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (sqlur *SQLUserRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := sqlur.db.Where("is_active = true").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (sqlur *SQLUserRepository) Save(user models.User) error {
	err := sqlur.db.Save(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func (sqlur *SQLUserRepository) CreateAdminOffice() error {
	var office models.Office
	err := sqlur.db.Where("office_name = ?", constants.AdminOffice).First(&office).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			building := models.Building{
				BuildingName: constants.AdminBuilding,
				Floors: []models.Floor{
					{
						FloorNumber: 0,
						Office: &models.Office{
							OfficeName: constants.AdminOffice,
						},
					},
				},
			}
			err = sqlur.db.Create(&building).Error
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func (sqlur *SQLUserRepository) SeedBuildingAndOfice() error {
	err := sqlur.db.Create(&models.Building{
		BuildingName: constants.TestBuilding,
		Floors: []models.Floor{
			{
				FloorNumber: 1,
				Office: &models.Office{
					OfficeName: constants.TestOffice,
				},
			},
		},
	}).Error
	return err
}

func (sqlur *SQLUserRepository) CreateUser(name string, email string, password string, office string, role roles.Role) error {
	var officeStruct models.Office

	err := sqlur.db.Where("office_name = ?", office).First(&officeStruct).Error
	if err != nil {
		return err
	}

	user := models.User{
		Name:     name,
		Email:    email,
		Password: password,
		Role:     role,
		OfficeID: officeStruct.OfficeID,
	}

	err = sqlur.db.Create(&user).Error
	return err
}

func NewSQLUserRepository(db *gorm.DB) *SQLUserRepository {
	return &SQLUserRepository{
		db: db,
	}
}
