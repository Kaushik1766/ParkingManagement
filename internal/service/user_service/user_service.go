package userservice

import (
	"context"
	"errors"
	"log"

	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	user "github.com/Kaushik1766/ParkingManagement/internal/models/user"
	userjwt "github.com/Kaushik1766/ParkingManagement/internal/models/user_jwt"
	"github.com/Kaushik1766/ParkingManagement/internal/models/vehicle"
	userrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/user_repository"
	vehiclerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/vehicle_repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo    userrepository.UserStorage
	vehicleRepo vehiclerepository.VehicleStorage
}

func (us *UserService) GetUserProfile(ctx context.Context) (user.UserDTO, error) {
	ctxUser := ctx.Value(constants.User).(userjwt.UserJwt)
	currentUser, err := us.userRepo.GetUserById(ctxUser.ID)
	if err != nil {
		return user.UserDTO{}, err
	}
	userDto := user.UserDTO{
		UserId: currentUser.UserId.String(),
		Name:   currentUser.Name,
		Email:  currentUser.Email,
		Role:   currentUser.Role.String(),
	}
	return userDto, nil
}

func (us *UserService) RegisterVehicle(ctx context.Context, numberplate string, vehicleType vehicletypes.VehicleType) error {
	ctxUser := ctx.Value(constants.User).(userjwt.UserJwt)
	currentUser, err := us.userRepo.GetUserById(ctxUser.ID)
	if err != nil {
		return err
	}
	err = us.vehicleRepo.AddVehicle(numberplate, currentUser.UserId, vehicleType)
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
	currentUser := ctx.Value(constants.User).(userjwt.UserJwt)
	// fmt.Println(currentUser.ID)
	uid, err := uuid.Parse(currentUser.ID)
	if err != nil {
		log.Println(err)
		return []vehicle.VehicleDTO{}
	}
	userVehicles, err := us.vehicleRepo.GetVehiclesByUserId(uid)
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
	ctxVal := ctx.Value(constants.User)
	if ctxVal == nil {
		return errors.New("invalid context")
	}
	currentUser := ctxVal.(userjwt.UserJwt)

	updatedUser, err := us.userRepo.GetUserById(currentUser.ID)
	if err != nil {
		return err
	}

	if name != "" {
		updatedUser.Name = name
	}
	if email != "" {
		updatedUser.Email = email
	}
	if password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), 12)
		if err != nil {
			return err
		}
		updatedUser.Password = string(hashedPassword)
	}
	err = us.userRepo.Save(updatedUser)
	return err
}

func (us *UserService) DeleteProfile(ctx context.Context) error {
	currentUser := ctx.Value(constants.User).(user.User)
	currentUser.IsActive = false
	err := us.userRepo.Save(currentUser)
	return err
}
