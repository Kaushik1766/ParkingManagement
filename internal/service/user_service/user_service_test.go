package userservice_test

import (
	"context"
	"testing"

	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	models "github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	officerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/office_repository"
	userrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/user_repository"
	vehiclerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/vehicle_repository"
	slotassignment "github.com/Kaushik1766/ParkingManagement/internal/service/slot_assignment"
	userservice "github.com/Kaushik1766/ParkingManagement/internal/service/user_service"
	"github.com/Kaushik1766/ParkingManagement/utils"
	"github.com/golang-jwt/jwt/v5"
)

var (
	db, _       = utils.GetDB()
	vehicleRepo = vehiclerepository.NewSQLVehicleRepository(db)
	officeRepo  = officerepository.NewSQLOfficeRepository(db)
	userRepo    = userrepository.NewSQLUserRepository(db)
)

func TestUserService_GetRegisteredVehicles(t *testing.T) {
	baseCtx := context.Background()
	userCtx := context.WithValue(baseCtx, constants.User, models.UserJwt{
		Email:  "kaushik@a.com",
		Role:   roles.Customer,
		Office: "samsung",
		RegisteredClaims: jwt.RegisteredClaims{
			ID: "9aa170ec-105d-408c-bd8b-b2e66053474c",
		},
	})
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		repo              userrepository.UserStorage
		vehicRepo         vehiclerepository.VehicleStorage
		officeRepo        officerepository.OfficeStorage
		assignmentService slotassignment.SlotAssignmentMgr
		want              []models.VehicleDTO
	}{
		{
			repo:              userRepo,
			vehicRepo:         vehicleRepo,
			officeRepo:        officeRepo,
			assignmentService: nil,
			name:              "Get Registered Vehicles",
			want:              []models.VehicleDTO{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := userservice.NewUserService(tt.repo, tt.vehicRepo, tt.officeRepo, tt.assignmentService)
			got := us.GetRegisteredVehicles(userCtx)

			if true {
				t.Errorf("GetRegisteredVehicles() = %v, want %v", got, tt.want)
			}
		})
	}
}
