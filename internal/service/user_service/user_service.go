package userservice

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	user "github.com/Kaushik1766/ParkingManagement/internal/models/user"
	userjwt "github.com/Kaushik1766/ParkingManagement/internal/models/user_jwt"
	"github.com/Kaushik1766/ParkingManagement/internal/models/vehicle"
	officerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/office_repository"
	userrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/user_repository"
	vehiclerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/vehicle_repository"
	slotassignment "github.com/Kaushik1766/ParkingManagement/internal/service/slot_assignment"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo          userrepository.UserStorage
	vehicleRepo       vehiclerepository.VehicleStorage
	officeRepo        officerepository.OfficeStorage
	assignmentService slotassignment.SlotAssignmentMgr
}

func (us *UserService) GetUserProfile(ctx context.Context) (user.UserDTO, error) {
	ctxUser := ctx.Value(constants.User).(userjwt.UserJwt)
	currentUser, err := us.userRepo.GetUserById(ctxUser.ID)
	if err != nil {
		return user.UserDTO{}, err
	}
	userDto := user.UserDTO{
		UserId: currentUser.UserID.String(),
		Name:   currentUser.Name,
		Email:  currentUser.Email,
		Role:   currentUser.Role.String(),
		// Office: currentUser.Office,
	}
	return userDto, nil
}

func (us *UserService) GetUserById(ctx context.Context, userId string) (user.UserDTO, error) {
	userStruct, err := us.userRepo.GetUserById(userId)
	if err != nil {
		return user.UserDTO{}, err
	}

	return user.UserDTO{
		UserId: userStruct.UserID.String(),
		Name:   userStruct.Name,
		Email:  userStruct.Email,
		Role:   userStruct.Role.String(),
		// Office: userStruct.Office,
	}, nil
}

func (us *UserService) RegisterVehicle(ctx context.Context, numberplate string, vehicleType vehicletypes.VehicleType) error {
	if len(numberplate) == 0 {
		return errors.New("numberplate cannot be empty")
	}

	if len(numberplate) != 10 {
		return errors.New("numberplate must be 10 characters long")
	}
	ctxUser := ctx.Value(constants.User).(userjwt.UserJwt)
	currentUser, err := us.userRepo.GetUserById(ctxUser.ID)
	if err != nil {
		return err
	}
	newVehicle, err := us.vehicleRepo.AddVehicle(numberplate, currentUser.UserID, vehicleType)
	if err != nil {
		return err
	}

	err = us.assignmentService.AutoAssignSlot(ctx, newVehicle.VehicleID.String())
	if err != nil {
		return fmt.Errorf("failed to assign slot: %w", err)
	}

	return err
}

func (us *UserService) GetAllUsers(ctx context.Context) ([]user.User, error) {
	allUsers, err := us.userRepo.GetAllUsers()
	if err != nil {
		return nil, err
	}
	var activeUsers []user.User
	for _, u := range allUsers {
		if u.IsActive {
			activeUsers = append(activeUsers, u)
		}
	}
	return activeUsers, nil
}

func (us *UserService) UnregisterVehicle(ctx context.Context, numberplate string) error {
	currentUser := ctx.Value(constants.User).(userjwt.UserJwt)
	userVehicles, err := us.vehicleRepo.GetVehiclesByUserId(uuid.MustParse(currentUser.ID))
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

func NewUserService(
	repo userrepository.UserStorage,
	vehicRepo vehiclerepository.VehicleStorage,
	officeRepo officerepository.OfficeStorage,
	assignmentService slotassignment.SlotAssignmentMgr,
) *UserService {
	return &UserService{
		userRepo:          repo,
		vehicleRepo:       vehicRepo,
		officeRepo:        officeRepo,
		assignmentService: assignmentService,
	}
}

func (us *UserService) UpdateProfile(ctx context.Context, name, email, password, office string) error {
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
	if office != "" {
		_, err = us.officeRepo.GetOfficeByName(office)
		if err != nil {
			return errors.New("office does not exist")
		}
		// updatedUser.Office = office
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
