package userrepository

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/Kaushik1766/ParkingManagement/internal/config"
	models "github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	customerrors "github.com/Kaushik1766/ParkingManagement/pkg/customErrors"
	"github.com/google/uuid"
)

type FileUserRepository struct {
	*sync.Mutex
	users []models.User
}

func (db *FileUserRepository) GetAllUsers() ([]models.User, error) {
	db.Lock()
	defer db.Unlock()
	if len(db.users) == 0 {
		return []models.User{}, customerrors.UserNotFound{}
	}
	return db.users, nil
}

func NewFileUserRepository() *FileUserRepository {
	data, err := os.ReadFile(config.UsersPath)
	if err != nil {
		os.WriteFile(config.UsersPath, []byte("[]"), 0666)
		data, err = json.Marshal([]models.User{})
		if err != nil {
			fmt.Println("unable to marshal")
		}
	}

	var userData []models.User
	err = json.Unmarshal(data, &userData)
	if err != nil {
		fmt.Println(err)
		panic("corrupted data")
	}
	return &FileUserRepository{
		Mutex: &sync.Mutex{},
		users: userData,
	}
}

func (db *FileUserRepository) GetUserById(id string) (models.User, error) {
	db.Lock()
	defer db.Unlock()
	for _, val := range db.users {
		if val.UserID.String() == id {
			return val, nil
		}
	}
	return models.User{}, customerrors.UserNotFound{}
}

func (db *FileUserRepository) GetUserByEmail(email string) (models.User, error) {
	db.Lock()
	defer db.Unlock()
	for _, val := range db.users {
		if val.Email == email {
			return val, nil
		}
	}
	return models.User{}, customerrors.UserNotFound{}
}

func (db *FileUserRepository) Save(user models.User) error {
	db.Lock()
	defer db.Unlock()
	for i, val := range db.users {
		if user.UserID == val.UserID {
			db.users[i] = user
			return nil
		}
	}
	db.users = append(db.users, user)
	return nil
}

func (db *FileUserRepository) CreateUser(name, email, password, office string, role roles.Role) error {
	for _, val := range db.users {
		if val.Email == email {
			return errors.New("email already used")
		}
	}

	db.users = append(db.users, models.User{
		UserID:   uuid.New(),
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
	err = os.WriteFile(config.UsersPath, data, 0666)
	return err
}
