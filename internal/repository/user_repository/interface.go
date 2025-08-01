package userrepository

import (
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums"
	user "github.com/Kaushik1766/ParkingManagement/internal/models/user"
)

type UserStorage interface {
	GetUserByEmail(email string) (user.User, error)
	GetUserById(id string) (user.User, error)
	Save(user user.User) error
	CreateUser(name, email, password string, role enums.Role) error
}
