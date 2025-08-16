package userrepository

import (
	models "github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
)

type UserStorage interface {
	GetUserByEmail(email string) (models.User, error)
	GetUserById(id string) (models.User, error)
	GetAllUsers() ([]models.User, error)
	Save(user models.User) error
	CreateUser(name, email, password, office string, role roles.Role) error
}
