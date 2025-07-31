package userrepository

import user "github.com/Kaushik1766/ParkingManagement/internal/models/User"

type UserStorage interface {
	GetUserByEmail(email string) (user.User, error)
	GetUserById(id string) (user.User, error)
	Save(user user.User) error
}
