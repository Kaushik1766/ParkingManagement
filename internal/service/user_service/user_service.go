package userservice

import (
	"context"
	"errors"

	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums"
	user "github.com/Kaushik1766/ParkingManagement/internal/models/user"
	"github.com/Kaushik1766/ParkingManagement/internal/models/vehicle"
	userrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/user_repository"
	vehiclerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/vehicle_repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo    userrepository.UserStorage
	vehicleRepo vehiclerepository.VehicleStorage
}

func (us *UserService) RegisterVehicle(ctx context.Context, numberplate string, vehicleType enums.VehicleType) error {
	currentUser := ctx.Value(constants.User).(user.User)
	err := us.vehicleRepo.AddVehicle(numberplate, currentUser.UserId, vehicleType)
	return err
}

func (us *UserService) UnregisterVehicle(ctx context.Context, numberplate string) error {
	currentUser := ctx.Value(constants.User).(user.User)
	userVehicles, err := us.vehicleRepo.GetVehiclesByUserId(currentUser.UserId)
	if err != nil {
		return err
	}
	for _, v := range userVehicles {
		if v.NumberPlate == numberplate {
			if v.IsActive {
				return us.vehicleRepo.RemoveVehicle(numberplate)
			} else {
				return nil
			}
		}
	}
	return errors.New("vehicle not found for the user")
}

func (us *UserService) GetRegisteredVehicles(ctx context.Context) []vehicle.VehicleDTO {
	currentUser := ctx.Value(constants.User).(user.User)
	userVehicles, err := us.vehicleRepo.GetVehiclesByUserId(currentUser.UserId)
	if err != nil {
		return []vehicle.VehicleDTO{}
	}
	var userVehicleDTO []vehicle.VehicleDTO
	for _, v := range userVehicles {
		userVehicleDTO = append(userVehicleDTO, vehicle.VehicleDTO{
			NumberPlate:  v.NumberPlate,
			VehicleType:  v.VehicleType.String(),
			AssignedSlot: v.AssignedSlot,
		})
	}
	return userVehicleDTO
}

func NewUserService(repo userrepository.UserStorage, vehicRepo vehiclerepository.VehicleStorage) *UserService {
	return &UserService{
		userRepo:    repo,
		vehicleRepo: vehicRepo,
	}
}

func (us *UserService) UpdateProfile(ctx context.Context, name, email, password string) error {
	currentUser := ctx.Value(constants.User).(user.User)
	currentUser.Name = name
	currentUser.Email = email
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	currentUser.Password = string(hashedPassword)
	err = us.userRepo.Save(currentUser)
	return err
}

func (us *UserService) DeleteProfile(ctx context.Context) error {
	currentUser := ctx.Value("user").(user.User)
	currentUser.IsActive = false
	err := us.userRepo.Save(currentUser)
	return err
}
