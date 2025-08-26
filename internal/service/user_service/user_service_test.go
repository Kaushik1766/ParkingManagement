package userservice

import (
	"context"
	"reflect"
	"testing"

	"github.com/Kaushik1766/ParkingManagement/internal/models"
	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	officerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/office_repository"
	userrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/user_repository"
	vehiclerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/vehicle_repository"
	slotassignment "github.com/Kaushik1766/ParkingManagement/internal/service/slot_assignment"
)

func TestNewUserService(t *testing.T) {
	type args struct {
		repo              userrepository.UserStorage
		vehicRepo         vehiclerepository.VehicleStorage
		officeRepo        officerepository.OfficeStorage
		assignmentService slotassignment.SlotAssignmentMgr
	}
	tests := []struct {
		name string
		args args
		want *UserService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserService(tt.args.repo, tt.args.vehicRepo, tt.args.officeRepo, tt.args.assignmentService); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_DeleteProfile(t *testing.T) {
	type fields struct {
		userRepo          userrepository.UserStorage
		vehicleRepo       vehiclerepository.VehicleStorage
		officeRepo        officerepository.OfficeStorage
		assignmentService slotassignment.SlotAssignmentMgr
	}
	type args struct {
		ctx    context.Context
		userId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := UserService{
				userRepo:          tt.fields.userRepo,
				vehicleRepo:       tt.fields.vehicleRepo,
				officeRepo:        tt.fields.officeRepo,
				assignmentService: tt.fields.assignmentService,
			}
			if err := us.DeleteProfile(tt.args.ctx, tt.args.userId); (err != nil) != tt.wantErr {
				t.Errorf("DeleteProfile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserService_GetAllUsers(t *testing.T) {
	type fields struct {
		userRepo          userrepository.UserStorage
		vehicleRepo       vehiclerepository.VehicleStorage
		officeRepo        officerepository.OfficeStorage
		assignmentService slotassignment.SlotAssignmentMgr
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.UserDTO
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &UserService{
				userRepo:          tt.fields.userRepo,
				vehicleRepo:       tt.fields.vehicleRepo,
				officeRepo:        tt.fields.officeRepo,
				assignmentService: tt.fields.assignmentService,
			}
			got, err := us.GetAllUsers(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllUsers() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_GetRegisteredVehicles(t *testing.T) {
	type fields struct {
		userRepo          userrepository.UserStorage
		vehicleRepo       vehiclerepository.VehicleStorage
		officeRepo        officerepository.OfficeStorage
		assignmentService slotassignment.SlotAssignmentMgr
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.VehicleDTO
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &UserService{
				userRepo:          tt.fields.userRepo,
				vehicleRepo:       tt.fields.vehicleRepo,
				officeRepo:        tt.fields.officeRepo,
				assignmentService: tt.fields.assignmentService,
			}
			got, err := us.GetRegisteredVehicles(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRegisteredVehicles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRegisteredVehicles() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_GetUserById(t *testing.T) {
	type fields struct {
		userRepo          userrepository.UserStorage
		vehicleRepo       vehiclerepository.VehicleStorage
		officeRepo        officerepository.OfficeStorage
		assignmentService slotassignment.SlotAssignmentMgr
	}
	type args struct {
		ctx    context.Context
		userId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.UserDTO
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &UserService{
				userRepo:          tt.fields.userRepo,
				vehicleRepo:       tt.fields.vehicleRepo,
				officeRepo:        tt.fields.officeRepo,
				assignmentService: tt.fields.assignmentService,
			}
			got, err := us.GetUserById(tt.args.ctx, tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserById() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_GetUserProfile(t *testing.T) {
	type fields struct {
		userRepo          userrepository.UserStorage
		vehicleRepo       vehiclerepository.VehicleStorage
		officeRepo        officerepository.OfficeStorage
		assignmentService slotassignment.SlotAssignmentMgr
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.UserDTO
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &UserService{
				userRepo:          tt.fields.userRepo,
				vehicleRepo:       tt.fields.vehicleRepo,
				officeRepo:        tt.fields.officeRepo,
				assignmentService: tt.fields.assignmentService,
			}
			got, err := us.GetUserProfile(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserProfile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_RegisterVehicle(t *testing.T) {
	type fields struct {
		userRepo          userrepository.UserStorage
		vehicleRepo       vehiclerepository.VehicleStorage
		officeRepo        officerepository.OfficeStorage
		assignmentService slotassignment.SlotAssignmentMgr
	}
	type args struct {
		ctx         context.Context
		numberplate string
		vehicleType vehicletypes.VehicleType
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &UserService{
				userRepo:          tt.fields.userRepo,
				vehicleRepo:       tt.fields.vehicleRepo,
				officeRepo:        tt.fields.officeRepo,
				assignmentService: tt.fields.assignmentService,
			}
			if err := us.RegisterVehicle(tt.args.ctx, tt.args.numberplate, tt.args.vehicleType); (err != nil) != tt.wantErr {
				t.Errorf("RegisterVehicle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserService_UnregisterVehicle(t *testing.T) {
	type fields struct {
		userRepo          userrepository.UserStorage
		vehicleRepo       vehiclerepository.VehicleStorage
		officeRepo        officerepository.OfficeStorage
		assignmentService slotassignment.SlotAssignmentMgr
	}
	type args struct {
		ctx         context.Context
		numberplate string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &UserService{
				userRepo:          tt.fields.userRepo,
				vehicleRepo:       tt.fields.vehicleRepo,
				officeRepo:        tt.fields.officeRepo,
				assignmentService: tt.fields.assignmentService,
			}
			if err := us.UnregisterVehicle(tt.args.ctx, tt.args.numberplate); (err != nil) != tt.wantErr {
				t.Errorf("UnregisterVehicle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserService_UpdateProfile(t *testing.T) {
	type fields struct {
		userRepo          userrepository.UserStorage
		vehicleRepo       vehiclerepository.VehicleStorage
		officeRepo        officerepository.OfficeStorage
		assignmentService slotassignment.SlotAssignmentMgr
	}
	type args struct {
		ctx       context.Context
		userId    string
		updateReq models.UpdateUserDTO
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &UserService{
				userRepo:          tt.fields.userRepo,
				vehicleRepo:       tt.fields.vehicleRepo,
				officeRepo:        tt.fields.officeRepo,
				assignmentService: tt.fields.assignmentService,
			}
			if err := us.UpdateProfile(tt.args.ctx, tt.args.userId, tt.args.updateReq); (err != nil) != tt.wantErr {
				t.Errorf("UpdateProfile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
