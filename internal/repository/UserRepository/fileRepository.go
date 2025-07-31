package userrepository

import (
	"sync"

	user "github.com/Kaushik1766/ParkingManagement/internal/models/User"
	customerrors "github.com/Kaushik1766/ParkingManagement/pkg/customErrors"
)

type FileUserRepository struct {
	*sync.Mutex
	users []user.User
}

func (db *FileUserRepository) GetUserById(id string) (user.User, error) {
	db.Lock()
	defer db.Unlock()
	for _, val := range db.users {
		if val.UserId.String() == id {
			return val, nil
		}
	}
	return user.User{}, customerrors.UserNotFound{}
}

func (db *FileUserRepository) GetUserByEmail(email string) (user.User, error) {
	db.Lock()
	defer db.Unlock()
	for _, val := range db.users {
		if val.Email == email {
			return val, nil
		}
	}
	return user.User{}, customerrors.UserNotFound{}
}

func (db *FileUserRepository) Save(user user.User) error {
	db.Lock()
	defer db.Unlock()
	for i, val := range db.users {
		if user.UserId == val.UserId {
			db.users[i] = user
			return nil
		}
	}
	db.users = append(db.users, user)
	return nil
}
