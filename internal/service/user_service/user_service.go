package userservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	errorcodes "github.com/Kaushik1766/ParkingManagement/internal/constants/error_codes"
	models "github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	buildingrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/building_repository"
	officerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/office_repository"
	userrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/user_repository"
	vehiclerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/vehicle_repository"
	slotassignment "github.com/Kaushik1766/ParkingManagement/internal/service/slot_assignment"
	customerrors "github.com/Kaushik1766/ParkingManagement/pkg/customErrors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo          userrepository.UserStorage
	vehicleRepo       vehiclerepository.VehicleStorage
	officeRepo        officerepository.OfficeStorage
	assignmentService slotassignment.SlotAssignmentMgr
	buildingRepo      buildingrepository.BuildingStorage
}

func (us *UserService) GetUserProfile(ctx context.Context) (models.UserDTO, error) {
	ctxUser := ctx.Value(constants.User).(models.UserJwt)
	currentUser, err := us.userRepo.GetUserByEmail(ctx, ctxUser.Email)
	if err != nil {
		return models.UserDTO{}, err
	}
	userDto := models.UserDTO{
		UserId: currentUser.UserID.String(),
		Name:   currentUser.Name,
		Email:  currentUser.Email,
		Role:   currentUser.Role.String(),
		Office: currentUser.Office.OfficeName,
		// Office: currentUser.Office,
	}
	return userDto, nil
}

func (us *UserService) GetUserById(ctx context.Context, userId string) (models.UserDTO, error) {
	userStruct, err := us.userRepo.GetUserById(ctx, userId)
	if err != nil {
		return models.UserDTO{}, err
	}

	return models.UserDTO{
		UserId: userStruct.UserID.String(),
		Name:   userStruct.Name,
		Email:  userStruct.Email,
		Role:   userStruct.Role.String(),
		Office: userStruct.Office.OfficeName,
	}, nil
}

func (us *UserService) RegisterVehicle(ctx context.Context, numberplate string, vehicleType vehicletypes.VehicleType) error {
	if len(numberplate) == 0 {
		return errors.New("numberplate cannot be empty")
	}

	if len(numberplate) != 10 {
		return errors.New("numberplate must be 10 characters long")
	}
	ctxUser := ctx.Value(constants.User).(models.UserJwt)
	newVehicle, err := us.vehicleRepo.AddVehicle(ctx, numberplate, uuid.MustParse(ctxUser.ID), vehicleType)
	if err != nil {
		return err
	}

	err = us.assignmentService.AutoAssignSlot(ctx, newVehicle.VehicleID.String())
	if err != nil {
		return fmt.Errorf("failed to assign slot: %w", err)
	}

	return err
}

func (us *UserService) GetAllUsers(ctx context.Context) ([]models.UserDTO, error) {
	allUsers, err := us.userRepo.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	var activeUsers []models.User
	for _, u := range allUsers {
		if u.IsActive {
			activeUsers = append(activeUsers, u)
		}
	}

	activerUsersDTO := make([]models.UserDTO, len(activeUsers))
	for i, u := range activeUsers {
		activerUsersDTO[i] = models.UserDTO{
			UserId: u.UserID.String(),
			Name:   u.Name,
			Email:  u.Email,
			Role:   u.Role.String(),
			Office: u.Office.OfficeName,
		}
	}
	return activerUsersDTO, nil
}

func (us *UserService) UnregisterVehicle(ctx context.Context, numberplate string) error {
	currentUser := ctx.Value(constants.User).(models.UserJwt)

	vehicle, err := us.vehicleRepo.GetVehicleByNumberPlate(ctx, numberplate)
	if err != nil {
		return errors.New("error fetching vehicles")
	}

	if vehicle.UserID.String() != currentUser.ID {
		return errors.New("vehicle not owned by user")
	}

	if vehicle.IsActive {
		return us.vehicleRepo.RemoveVehicle(ctx, numberplate)
	} else {
		return nil
	}
	//userVehicles, err := us.vehicleRepo.GetVehiclesByUserId(uuid.MustParse(currentUser.ID))
	//if err != nil {
	//	return err
	//}
	//for _, v := range userVehicles {
	//	if v.NumberPlate == numberplate {
	//		if v.IsActive {
	//			return us.vehicleRepo.RemoveVehicle(numberplate)
	//		} else {
	//			return nil
	//		}
	//	}
	//}
	//return errors.New("vehicle not found for the user")
}

func (us *UserService) GetRegisteredVehicles(ctx context.Context) ([]models.VehicleDTO, error) {
	currentUser := ctx.Value(constants.User).(models.UserJwt)
	uid, _ := uuid.Parse(currentUser.ID)
	userVehicles, err := us.vehicleRepo.GetVehiclesByUserId(ctx, uid)
	if err != nil {
		return []models.VehicleDTO{}, err
	}

	var userVehicleDTO []models.VehicleDTO
	for _, v := range userVehicles {
		building, err := us.buildingRepo.GetBuildingByID(ctx, v.AssignedBuildingID)
		if err != nil {
			return []models.VehicleDTO{}, err
		}

		isParked, err := us.vehicleRepo.GetParkingStatus(ctx, v.NumberPlate)
		if err != nil {
			return []models.VehicleDTO{}, err
		}
		userVehicleDTO = append(userVehicleDTO, models.VehicleDTO{
			NumberPlate:          v.NumberPlate,
			VehicleType:          v.VehicleType.String(),
			AssignedBuildingID:   v.AssignedBuildingID.String(),
			IsParked:             isParked,
			AssignedFloorNumber:  v.AssignedFloorNumber,
			AssignedSlotNumber:   v.AssignedSlotNumber,
			AssignedBuildingName: building.BuildingName,
		})
	}
	return userVehicleDTO, nil
}

func NewUserService(
	repo userrepository.UserStorage,
	vehicRepo vehiclerepository.VehicleStorage,
	officeRepo officerepository.OfficeStorage,
	assignmentService slotassignment.SlotAssignmentMgr,
	buildingRepo buildingrepository.BuildingStorage,
) *UserService {
	return &UserService{
		userRepo:          repo,
		vehicleRepo:       vehicRepo,
		officeRepo:        officeRepo,
		assignmentService: assignmentService,
		buildingRepo:      buildingRepo,
	}
}

func (us *UserService) UpdateProfile(ctx context.Context, userId string, updateReq models.UpdateUserDTO) error {
	currentUser := ctx.Value(constants.User).(models.UserJwt)

	if currentUser.Role != roles.Admin && currentUser.ID != userId {
		return errors.New("unauthorized to update other user's profile")
	}

	updatedUser, err := us.userRepo.GetUserByEmail(ctx, currentUser.Email)
	if err != nil {
		return err
	}

	if updateReq.Name != "" {
		updatedUser.Name = updateReq.Name
	}
	if updateReq.Email != "" {
		updatedUser.Email = updateReq.Email
	}
	if updateReq.OfficeId != "" {
		_, err = us.officeRepo.GetOfficeById(ctx, updateReq.OfficeId)
		if err != nil {
			return errors.New("office does not exist")
		}
		updatedUser.OfficeID = uuid.MustParse(updateReq.OfficeId)
	}
	if updateReq.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateReq.Password), 12)
		if err != nil {
			return err
		}
		updatedUser.Password = string(hashedPassword)
	}
	return us.userRepo.Save(ctx, updatedUser)
}

func (us *UserService) DeleteProfile(ctx context.Context, userId string) error {
	ctxUser := ctx.Value(constants.User).(models.UserJwt)
	if ctxUser.Role != roles.Admin && ctxUser.ID != userId {
		return errors.New("unauthorized to delete other user's profile")
	}

	user, err := us.userRepo.GetUserByEmail(ctx, ctxUser.Email)
	if err != nil {
		return err
	}

	user.IsActive = false
	err = us.userRepo.Save(ctx, user)
	return err
}

// TODO:FIX
func (us *UserService) GetVehiclesByUserId(ctx context.Context, userId string) ([]models.VehicleDTO, error) {
	uid, err := uuid.Parse(userId)
	if err != nil {
		return nil, customerrors.NewWebError(err, errorcodes.InvalidInput)
	}
	userVehicles, err := us.vehicleRepo.GetVehiclesByUserId(ctx, uid)
	if err != nil {
		return nil, err
	}

	var vehicleDTOs []models.VehicleDTO
	for _, v := range userVehicles {
		building, err := us.buildingRepo.GetBuildingByID(ctx, v.AssignedBuildingID)
		if err != nil {
			return []models.VehicleDTO{}, err
		}

		isParked, err := us.vehicleRepo.GetParkingStatus(ctx, v.NumberPlate)
		if err != nil {
			return []models.VehicleDTO{}, err
		}
		vehicleDTO := models.VehicleDTO{
			NumberPlate:          v.NumberPlate,
			VehicleType:          v.VehicleType.String(),
			AssignedBuildingName: building.BuildingName,
			IsParked:             isParked,
			AssignedFloorNumber:  v.AssignedFloorNumber,
			AssignedSlotNumber:   v.AssignedSlotNumber,
			AssignedBuildingID:   v.AssignedBuildingID.String(),
		}
		vehicleDTOs = append(vehicleDTOs, vehicleDTO)
	}
	if vehicleDTOs == nil {
		vehicleDTOs = []models.VehicleDTO{}
	}
	return vehicleDTOs, nil
}
