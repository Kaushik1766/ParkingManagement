package userrepository

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/Kaushik1766/ParkingManagement/internal/models/enums"
	user "github.com/Kaushik1766/ParkingManagement/internal/models/user"
	customerrors "github.com/Kaushik1766/ParkingManagement/pkg/customErrors"
	"github.com/google/uuid"
)

type FileUserRepository struct {
	*sync.Mutex
	users []user.User
}

func NewFileUserRepository() *FileUserRepository {
	data, err := os.ReadFile("users.json")
	if err != nil {
		os.WriteFile("users.json", []byte("[]"), 0666)
		data, err = json.Marshal([]user.User{})
		if err != nil {
			fmt.Println("unable to marshal")
		}
	}

	var userData []user.User
	err = json.Unmarshal(data, &userData)
	if err != nil {
		fmt.Println(err)
		panic("corrupted data")
	}
	return &FileUserRepository{
		users: userData,
	}
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

func (db *FileUserRepository) CreateUser(name, email, password string, role enums.Role) error {
	for _, val := range db.users {
		if val.Email == email {
			return errors.New("email already used")
		}
	}

	db.users = append(db.users, user.User{
		UserId:   uuid.New(),
		Name:     name,
		Email:    email,
		Password: password,
		Role:     role,
		IsActive: true,
	})
	return nil
}

func (db *FileUserRepository) SerializeData() error {
	db.Lock()
	defer db.Unlock()
	data, err := json.Marshal(db.users)
	if err != nil {
		return err
	}
	err = os.WriteFile("users.json", data, 0666)
	return err
}
