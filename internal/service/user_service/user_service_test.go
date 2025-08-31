package userservice

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	officerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/office_repository"
	userrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/user_repository"
	vehiclerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/vehicle_repository"
	slotassignment "github.com/Kaushik1766/ParkingManagement/internal/service/slot_assignment"
	"github.com/Kaushik1766/ParkingManagement/mocks"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

var userId = uuid.New()
var adminId = uuid.New()

var adminCtx = context.WithValue(context.Background(), constants.User, models.UserJwt{
	RegisteredClaims: jwt.RegisteredClaims{
		ID: adminId.String(),
	},
	Email:  "admin@a.com",
	Role:   roles.Admin,
	Office: "wg",
})
var userCtx = context.WithValue(context.Background(), constants.User, models.UserJwt{
	RegisteredClaims: jwt.RegisteredClaims{
		ID: userId.String(),
	},
	Email:  "kaushik@a.com",
	Role:   roles.Customer,
	Office: "wg",
})
var unauthorizedCtx = context.WithValue(context.Background(), constants.User, models.UserJwt{
	RegisteredClaims: jwt.RegisteredClaims{
		ID: uuid.New().String(),
	},
	Email:  "unauthorized@a.com",
	Role:   roles.Customer,
	Office: "wg",
})

func TestNewUserService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserStorage(ctrl)
	mockVehicleRepo := mocks.NewMockVehicleStorage(ctrl)
	mockOfficeRepo := mocks.NewMockOfficeStorage(ctrl)
	mockAssignmentService := mocks.NewMockSlotAssignmentMgr(ctrl)

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
		{
			name: "valid repos",
			args: args{
				repo:              mockUserRepo,
				vehicRepo:         mockVehicleRepo,
				officeRepo:        mockOfficeRepo,
				assignmentService: mockAssignmentService,
			},
			want: &UserService{
				userRepo:          mockUserRepo,
				vehicleRepo:       mockVehicleRepo,
				officeRepo:        mockOfficeRepo,
				assignmentService: mockAssignmentService,
			},
		},
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserStorage(ctrl)

	type fields struct {
		userRepo userrepository.UserStorage
	}
	type args struct {
		ctx    context.Context
		userId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func()
		wantErr bool
	}{
		{
			name:   "adminctx",
			fields: fields{userRepo: mockUserRepo},
			args: args{
				ctx:    adminCtx,
				userId: userId.String(),
			},
			mock: func() {
				mockUserRepo.EXPECT().GetUserById(gomock.Any()).Return(models.User{
					UserID: userId,
				}, nil)
				mockUserRepo.EXPECT().Save(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "userctx with own profile",
			fields: fields{
				userRepo: mockUserRepo,
			},
			args: args{
				ctx:    userCtx,
				userId: userId.String(),
			},
			mock: func() {
				mockUserRepo.EXPECT().GetUserById(userId.String()).Return(models.User{
					UserID: userId,
				}, nil)
				mockUserRepo.EXPECT().Save(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "userctx with other users id",
			fields: fields{
				userRepo: mockUserRepo,
			},
			args: args{
				ctx:    userCtx,
				userId: uuid.NewString(),
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "user not found",
			fields: fields{
				userRepo: mockUserRepo,
			},
			args: args{
				ctx:    adminCtx,
				userId: userId.String(),
			},
			mock: func() {
				mockUserRepo.EXPECT().GetUserById(gomock.Any()).Return(models.User{}, errors.New("user not found"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			us := &UserService{
				userRepo: tt.fields.userRepo,
			}
			if err := us.DeleteProfile(tt.args.ctx, tt.args.userId); (err != nil) != tt.wantErr {
				t.Errorf("DeleteProfile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserService_GetAllUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserStorage(ctrl)
	users := []models.User{
		{UserID: uuid.New(), Name: "User 1", Email: "user1@example.com", Role: roles.Customer, Office: models.Office{OfficeName: "Office1"}, IsActive: true},
		{UserID: uuid.New(), Name: "User 2", Email: "user2@example.com", Role: roles.Admin, Office: models.Office{OfficeName: "Office2"}, IsActive: true},
		{UserID: uuid.New(), Name: "User 3", Email: "user3@example.com", Role: roles.Customer, Office: models.Office{OfficeName: "Office1"}, IsActive: false},
	}

	type fields struct {
		userRepo userrepository.UserStorage
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func()
		want    []models.UserDTO
		wantErr bool
	}{
		{
			name:   "Success",
			fields: fields{userRepo: mockUserRepo},
			args:   args{ctx: context.Background()},
			mock: func() {
				mockUserRepo.EXPECT().GetAllUsers().Return(users, nil)
			},
			want: []models.UserDTO{
				{UserId: users[0].UserID.String(), Name: "User 1", Email: "user1@example.com", Role: "Customer", Office: "Office1"},
				{UserId: users[1].UserID.String(), Name: "User 2", Email: "user2@example.com", Role: "Admin", Office: "Office2"},
			},
			wantErr: false,
		},
		{
			name:   "repo error",
			fields: fields{userRepo: mockUserRepo},
			args:   args{ctx: context.Background()},
			mock: func() {
				mockUserRepo.EXPECT().GetAllUsers().Return(nil, errors.New("db error"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			us := &UserService{
				userRepo: tt.fields.userRepo,
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVehicleRepo := mocks.NewMockVehicleStorage(ctrl)

	vehicles := []models.Vehicle{
		{
			NumberPlate:  "asdf",
			VehicleType:  vehicletypes.TwoWheeler,
			AssignedSlot: nil,
		},
		{
			NumberPlate:  "asde",
			VehicleType:  vehicletypes.TwoWheeler,
			AssignedSlot: &models.Slot{},
		},
	}

	type fields struct {
		vehicleRepo vehiclerepository.VehicleStorage
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func()
		want    []models.VehicleDTO
		wantErr bool
	}{
		{
			name:   "success",
			fields: fields{vehicleRepo: mockVehicleRepo},
			args: args{
				ctx: userCtx,
			},
			mock: func() {
				mockVehicleRepo.EXPECT().GetVehiclesByUserId(gomock.Any()).Return(vehicles, nil)
			},
			want: []models.VehicleDTO{
				{
					AssignedSlot: models.Slot{},
					VehicleType:  vehicletypes.TwoWheeler.String(),
					NumberPlate:  "asdf",
				},
				{
					AssignedSlot: models.Slot{},
					VehicleType:  vehicletypes.TwoWheeler.String(),
					NumberPlate:  "asde",
				},
			},
			wantErr: false,
		},
		{
			name:   "repo error",
			fields: fields{vehicleRepo: mockVehicleRepo},
			args:   args{ctx: userCtx},
			mock: func() {
				mockVehicleRepo.EXPECT().GetVehiclesByUserId(gomock.Any()).Return(nil, errors.New("db error"))
			},
			want:    []models.VehicleDTO{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			us := &UserService{
				vehicleRepo: tt.fields.vehicleRepo,
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserStorage(ctrl)
	userID := uuid.New()
	user := models.User{UserID: userID, Name: "kaushik", Email: "kaushik@a.com", Role: roles.Customer, Office: models.Office{OfficeName: "wg"}}

	type fields struct {
		userRepo userrepository.UserStorage
	}
	type args struct {
		ctx    context.Context
		userId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func()
		want    models.UserDTO
		wantErr bool
	}{
		{
			name:   "Success",
			fields: fields{userRepo: mockUserRepo},
			args:   args{ctx: context.Background(), userId: userID.String()},
			mock: func() {
				mockUserRepo.EXPECT().GetUserById(userID.String()).Return(user, nil)
			},
			want: models.UserDTO{
				UserId: userID.String(),
				Name:   "kaushik",
				Email:  "kaushik@a.com",
				Role:   roles.Customer.String(),
				Office: "wg",
			},
			wantErr: false,
		},
		{
			name:   "user not found",
			fields: fields{userRepo: mockUserRepo},
			args:   args{ctx: context.Background(), userId: userID.String()},
			mock: func() {
				mockUserRepo.EXPECT().GetUserById(userID.String()).Return(models.User{}, errors.New("not found"))
			},
			want:    models.UserDTO{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			us := &UserService{
				userRepo: tt.fields.userRepo,
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserStorage(ctrl)
	user := models.User{UserID: userId, Name: "kaushik", Email: "kaushik@a.com", Role: roles.Customer, Office: models.Office{OfficeName: "wg"}}

	type fields struct {
		userRepo userrepository.UserStorage
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func()
		want    models.UserDTO
		wantErr bool
	}{
		{
			name:   "Success",
			fields: fields{userRepo: mockUserRepo},
			args:   args{ctx: userCtx},
			mock: func() {
				mockUserRepo.EXPECT().GetUserById(userId.String()).Return(user, nil)
			},
			want: models.UserDTO{
				UserId: userId.String(),
				Name:   "kaushik",
				Email:  "kaushik@a.com",
				Role:   "Customer",
				Office: "wg",
			},
			wantErr: false,
		},
		{
			name:   "Failure - User not found",
			fields: fields{userRepo: mockUserRepo},
			args:   args{ctx: userCtx},
			mock: func() {
				mockUserRepo.EXPECT().GetUserById(userId.String()).Return(models.User{}, errors.New("not found"))
			},
			want:    models.UserDTO{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			us := &UserService{
				userRepo: tt.fields.userRepo,
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVehicleRepo := mocks.NewMockVehicleStorage(ctrl)
	mockAssignmentService := mocks.NewMockSlotAssignmentMgr(ctrl)
	vehicleID := uuid.New()

	type fields struct {
		vehicleRepo       vehiclerepository.VehicleStorage
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
		mock    func()
		wantErr bool
	}{
		{
			name:   "Success",
			fields: fields{vehicleRepo: mockVehicleRepo, assignmentService: mockAssignmentService},
			args:   args{ctx: userCtx, numberplate: "VALID12345", vehicleType: vehicletypes.FourWheeler},
			mock: func() {
				mockVehicleRepo.EXPECT().AddVehicle("VALID12345", userId, vehicletypes.FourWheeler).Return(models.Vehicle{VehicleID: vehicleID}, nil)
				mockAssignmentService.EXPECT().AutoAssignSlot(userCtx, vehicleID.String()).Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "Failure - Empty numberplate",
			fields:  fields{},
			args:    args{ctx: userCtx, numberplate: "", vehicleType: vehicletypes.FourWheeler},
			mock:    func() {},
			wantErr: true,
		},
		{
			name:    "Failure - Invalid numberplate length",
			fields:  fields{},
			args:    args{ctx: userCtx, numberplate: "SHORT", vehicleType: vehicletypes.FourWheeler},
			mock:    func() {},
			wantErr: true,
		},
		{
			name:   "Failure - AddVehicle error",
			fields: fields{vehicleRepo: mockVehicleRepo},
			args:   args{ctx: userCtx, numberplate: "VALID12345", vehicleType: vehicletypes.FourWheeler},
			mock: func() {
				mockVehicleRepo.EXPECT().AddVehicle("VALID12345", userId, vehicletypes.FourWheeler).Return(models.Vehicle{}, errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name:   "Failure - AutoAssignSlot error",
			fields: fields{vehicleRepo: mockVehicleRepo, assignmentService: mockAssignmentService},
			args:   args{ctx: userCtx, numberplate: "VALID12345", vehicleType: vehicletypes.FourWheeler},
			mock: func() {
				mockVehicleRepo.EXPECT().AddVehicle("VALID12345", userId, vehicletypes.FourWheeler).Return(models.Vehicle{VehicleID: vehicleID}, nil)
				mockAssignmentService.EXPECT().AutoAssignSlot(userCtx, vehicleID.String()).Return(errors.New("assignment error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			us := &UserService{
				vehicleRepo:       tt.fields.vehicleRepo,
				assignmentService: tt.fields.assignmentService,
			}
			if err := us.RegisterVehicle(tt.args.ctx, tt.args.numberplate, tt.args.vehicleType); (err != nil) != tt.wantErr {
				t.Errorf("RegisterVehicle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserService_UnregisterVehicle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVehicleRepo := mocks.NewMockVehicleStorage(ctrl)
	vehicles := []models.Vehicle{
		{NumberPlate: "PLATE12345", IsActive: true},
		{NumberPlate: "PLATE67890", IsActive: false},
	}

	type fields struct {
		vehicleRepo vehiclerepository.VehicleStorage
	}
	type args struct {
		ctx         context.Context
		numberplate string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func()
		wantErr bool
	}{
		{
			name:   "Success - Unregister active vehicle",
			fields: fields{vehicleRepo: mockVehicleRepo},
			args:   args{ctx: userCtx, numberplate: "PLATE12345"},
			mock: func() {
				mockVehicleRepo.EXPECT().GetVehiclesByUserId(userId).Return(vehicles, nil)
				mockVehicleRepo.EXPECT().RemoveVehicle("PLATE12345").Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "Success - Vehicle already inactive",
			fields: fields{vehicleRepo: mockVehicleRepo},
			args:   args{ctx: userCtx, numberplate: "PLATE67890"},
			mock: func() {
				mockVehicleRepo.EXPECT().GetVehiclesByUserId(userId).Return(vehicles, nil)
			},
			wantErr: false,
		},
		{
			name:   "Failure - GetVehiclesByUserId error",
			fields: fields{vehicleRepo: mockVehicleRepo},
			args:   args{ctx: userCtx, numberplate: "PLATE12345"},
			mock: func() {
				mockVehicleRepo.EXPECT().GetVehiclesByUserId(userId).Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name:   "Failure - Vehicle not found for user",
			fields: fields{vehicleRepo: mockVehicleRepo},
			args:   args{ctx: userCtx, numberplate: "NOTFOUND12"},
			mock: func() {
				mockVehicleRepo.EXPECT().GetVehiclesByUserId(userId).Return(vehicles, nil)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			us := &UserService{
				vehicleRepo: tt.fields.vehicleRepo,
			}
			if err := us.UnregisterVehicle(tt.args.ctx, tt.args.numberplate); (err != nil) != tt.wantErr {
				t.Errorf("UnregisterVehicle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserService_UpdateProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserStorage(ctrl)
	mockOfficeRepo := mocks.NewMockOfficeStorage(ctrl)

	user := models.User{UserID: userId, Name: "kaushik", Email: "kaushik@a.com"}
	office := models.Office{OfficeID: uuid.New(), OfficeName: "wg"}

	type fields struct {
		userRepo   userrepository.UserStorage
		officeRepo officerepository.OfficeStorage
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
		mock    func()
		wantErr bool
	}{
		{
			name:   "admin updates user",
			fields: fields{userRepo: mockUserRepo, officeRepo: mockOfficeRepo},
			args:   args{ctx: adminCtx, userId: userId.String(), updateReq: models.UpdateUserDTO{Name: "kaushik", Office: "wg", Email: "kaushik@a.com", Password: "asdf"}},
			mock: func() {
				mockUserRepo.EXPECT().GetUserById(userId.String()).Return(user, nil)
				mockOfficeRepo.EXPECT().GetOfficeByName("wg").Return(office, nil)
				mockUserRepo.EXPECT().Save(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "password too long",
			fields: fields{userRepo: mockUserRepo, officeRepo: mockOfficeRepo},
			args: args{
				ctx:    adminCtx,
				userId: userId.String(),
				updateReq: models.UpdateUserDTO{
					Name:   "kaushik",
					Office: "wg",
					Email:  "kaushik@a.com",
					Password: func() string {
						p := ""
						for range 100 {
							p += "a"
						}
						return p
					}(),
				}},
			mock: func() {
				mockUserRepo.EXPECT().GetUserById(userId.String()).Return(user, nil)
				mockOfficeRepo.EXPECT().GetOfficeByName("wg").Return(office, nil)
				// No Save expectation since bcrypt should fail with a long password
			},
			wantErr: true, // bcrypt will fail with very long passwords (>72 bytes)
		},
		{
			name:    "unauthorized",
			fields:  fields{},
			args:    args{ctx: unauthorizedCtx, userId: userId.String(), updateReq: models.UpdateUserDTO{Name: "kaushik"}},
			mock:    func() {},
			wantErr: true,
		},
		{
			name:   "user not found",
			fields: fields{userRepo: mockUserRepo},
			args:   args{ctx: userCtx, userId: userId.String(), updateReq: models.UpdateUserDTO{Name: "kaushik"}},
			mock: func() {
				mockUserRepo.EXPECT().GetUserById(userId.String()).Return(models.User{}, errors.New("not found"))
			},
			wantErr: true,
		},
		{
			name:   "office does not exist",
			fields: fields{userRepo: mockUserRepo, officeRepo: mockOfficeRepo},
			args:   args{ctx: userCtx, userId: userId.String(), updateReq: models.UpdateUserDTO{Office: "NonExistent Office"}},
			mock: func() {
				mockUserRepo.EXPECT().GetUserById(userId.String()).Return(user, nil)
				mockOfficeRepo.EXPECT().GetOfficeByName("NonExistent Office").Return(models.Office{}, errors.New("not found"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			us := &UserService{
				userRepo:   tt.fields.userRepo,
				officeRepo: tt.fields.officeRepo,
			}
			if err := us.UpdateProfile(tt.args.ctx, tt.args.userId, tt.args.updateReq); (err != nil) != tt.wantErr {
				t.Errorf("UpdateProfile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserService_GetVehiclesByUserId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVehicleRepo := mocks.NewMockVehicleStorage(ctrl)

	vehicles := []models.Vehicle{
		{
			VehicleID:   uuid.New(),
			UserID:      userId,
			NumberPlate: "TEST123456",
			VehicleType: vehicletypes.FourWheeler,
			IsActive:    true,
			AssignedSlot: &models.Slot{
				BuildingID:  uuid.New(),
				FloorNumber: 1,
				SlotNumber:  1,
				SlotType:    vehicletypes.FourWheeler,
			},
		},
		{
			VehicleID:    uuid.New(),
			UserID:       userId,
			NumberPlate:  "TEST789012",
			VehicleType:  vehicletypes.TwoWheeler,
			IsActive:     true,
			AssignedSlot: nil,
		},
	}

	type fields struct {
		vehicleRepo vehiclerepository.VehicleStorage
	}
	type args struct {
		ctx    context.Context
		userId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func()
		want    []models.VehicleDTO
		wantErr bool
	}{
		{
			name:   "Success - Get vehicles with assigned slots",
			fields: fields{vehicleRepo: mockVehicleRepo},
			args:   args{ctx: context.Background(), userId: userId.String()},
			mock: func() {
				mockVehicleRepo.EXPECT().GetVehiclesByUserId(userId).Return(vehicles, nil)
			},
			want: []models.VehicleDTO{
				{
					NumberPlate:  "TEST123456",
					VehicleType:  vehicletypes.FourWheeler.String(),
					AssignedSlot: *vehicles[0].AssignedSlot,
				},
				{
					NumberPlate:  "TEST789012",
					VehicleType:  vehicletypes.TwoWheeler.String(),
					AssignedSlot: models.Slot{}, // Empty slot when no slot is assigned
				},
			},
			wantErr: false,
		},
		{
			name:    "Failure - Invalid user ID",
			fields:  fields{vehicleRepo: mockVehicleRepo},
			args:    args{ctx: context.Background(), userId: "invalid-uuid"},
			mock:    func() {},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "Failure - Repository error",
			fields: fields{vehicleRepo: mockVehicleRepo},
			args:   args{ctx: context.Background(), userId: userId.String()},
			mock: func() {
				mockVehicleRepo.EXPECT().GetVehiclesByUserId(userId).Return(nil, errors.New("db error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "Success - Empty vehicle list",
			fields: fields{vehicleRepo: mockVehicleRepo},
			args:   args{ctx: context.Background(), userId: userId.String()},
			mock: func() {
				mockVehicleRepo.EXPECT().GetVehiclesByUserId(userId).Return([]models.Vehicle{}, nil)
			},
			want:    []models.VehicleDTO{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			us := &UserService{
				vehicleRepo: tt.fields.vehicleRepo,
			}
			got, err := us.GetVehiclesByUserId(tt.args.ctx, tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVehiclesByUserId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetVehiclesByUserId() got = %v, want %v", got, tt.want)
			}
		})
	}
}
