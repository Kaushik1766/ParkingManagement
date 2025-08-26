package buildingservice_test

import (
	"context"
	"errors"
	"reflect"
	"slices"
	"testing"

	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	buildingrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/building_repository"
	buildingservice "github.com/Kaushik1766/ParkingManagement/internal/service/building_service"
	"github.com/Kaushik1766/ParkingManagement/mocks"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

var adminCtx = context.WithValue(context.Background(), constants.User, models.UserJwt{
	RegisteredClaims: jwt.RegisteredClaims{
		ID: "afdsfasdfasd",
	},
	Email:  "admin@a.com",
	Role:   roles.Admin,
	Office: constants.AdminOffice,
})

var userCtx = context.WithValue(context.Background(), constants.User, models.UserJwt{
	RegisteredClaims: jwt.RegisteredClaims{
		ID: "afdsfasdfasd",
	},
	Email:  "user@a.com",
	Role:   roles.Customer,
	Office: "asfd",
})

func TestBuildingService_AddBuilding(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockBuildingStorage(ctrl)
	mockStorage.EXPECT().AddBuilding(gomock.Any()).Return(nil)

	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		repo buildingrepository.BuildingStorage
		// Named input parameters for target function.
		buildingName string
		wantErr      bool
		ctx          context.Context
	}{
		{
			name:         "authorized user",
			repo:         mockStorage,
			buildingName: "advant",
			wantErr:      false,
			ctx:          adminCtx,
		},
		{
			name:         "unauthorized user",
			repo:         mockStorage,
			buildingName: "advant",
			wantErr:      true,
			ctx:          userCtx,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := buildingservice.NewBuildingService(tt.repo)
			gotErr := bs.AddBuilding(tt.ctx, tt.buildingName)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("AddBuilding() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("AddBuilding() succeeded unexpectedly")
			}
		})
	}
}

func TestBuildingService_DeleteBuildingByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockBuildingStorage := mocks.NewMockBuildingStorage(ctrl)
	mockBuildingStorage.EXPECT().DeleteBuildingByID("validId").Return(nil)
	mockBuildingStorage.EXPECT().DeleteBuildingByID("invalidId").Return(errors.New("building id not found"))

	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		repo buildingrepository.BuildingStorage
		// Named input parameters for target function.
		buildingID string
		wantErr    bool
		ctx        context.Context
	}{
		{
			name:       "authorized user",
			repo:       mockBuildingStorage,
			buildingID: "validId",
			wantErr:    false,
			ctx:        adminCtx,
		},
		{
			name:       "unauthorized user",
			repo:       mockBuildingStorage,
			buildingID: "validId",
			wantErr:    true,
			ctx:        userCtx,
		},
		{
			name:       "invalid id",
			repo:       mockBuildingStorage,
			buildingID: "invalidId",
			wantErr:    true,
			ctx:        adminCtx,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := buildingservice.NewBuildingService(tt.repo)
			gotErr := bs.DeleteBuildingByID(tt.ctx, tt.buildingID)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("DeleteBuildingByID() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("DeleteBuildingByID() succeeded unexpectedly")
			}
		})
	}
}

func TestBuildingService_TestNewBuildingService(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockBuildingStorage := mocks.NewMockBuildingStorage(ctrl)

	type args struct {
		repo buildingrepository.BuildingStorage
	}
	tests := []struct {
		name string
		args args
		want *buildingservice.BuildingService
	}{
		{
			name: "valid repo",
			args: args{
				mockBuildingStorage,
			},
			want: buildingservice.NewBuildingService(mockBuildingStorage),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildingservice.NewBuildingService(tt.args.repo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBuildingService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildingService_GetBuildingByID(t *testing.T) {

	ctrl := gomock.NewController(t)

	invalidUUID := uuid.New()

	mockBuildingStorage := mocks.NewMockBuildingStorage(ctrl)
	mockBuildingStorage.EXPECT().GetBuildingByID(uuid.Nil).Return(models.Building{
		BuildingID:   uuid.Nil,
		BuildingName: "advant",
	}, nil)
	mockBuildingStorage.EXPECT().GetBuildingByID(invalidUUID).Return(models.Building{},
		errors.New("building id not found"))

	type fields struct {
		buildingRepo buildingrepository.BuildingStorage
	}
	type args struct {
		ctx        context.Context
		buildingID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.BuildingDTO
		wantErr bool
	}{
		{
			name: "validId",
			fields: fields{
				buildingRepo: mockBuildingStorage,
			},
			args: args{
				ctx:        adminCtx,
				buildingID: uuid.Nil.String(),
			},
			want: models.BuildingDTO{
				Name:       "advant",
				BuildingID: uuid.Nil.String(),
			},
			wantErr: false,
		},
		{
			name: "invalidUUID",
			fields: fields{
				buildingRepo: mockBuildingStorage,
			},
			args: args{
				ctx:        adminCtx,
				buildingID: "adfasdf",
			},
			want:    models.BuildingDTO{},
			wantErr: true,
		},
		{
			name: "invalid building id",
			fields: fields{
				buildingRepo: mockBuildingStorage,
			},
			args: args{
				ctx:        adminCtx,
				buildingID: invalidUUID.String(),
			},
			want:    models.BuildingDTO{},
			wantErr: true,
		},
		{
			name: "unauthorized user",
			fields: fields{
				buildingRepo: mockBuildingStorage,
			},
			args: args{
				ctx:        userCtx,
				buildingID: uuid.Nil.String(),
			},
			want:    models.BuildingDTO{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := buildingservice.NewBuildingService(tt.fields.buildingRepo)
			got, err := bs.GetBuildingByID(tt.args.ctx, tt.args.buildingID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBuildingByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBuildingByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildingService_GetAllBuildings(t *testing.T) {

	ctrl := gomock.NewController(t)
	mockBuildingStorage := mocks.NewMockBuildingStorage(ctrl)
	mockBuildingStorage.EXPECT().GetAllBuildings().Return([]models.Building{
		{
			BuildingID:   uuid.Nil,
			BuildingName: "advant",
			Floors:       nil,
		},
	}, nil)

	type fields struct {
		buildingRepo buildingrepository.BuildingStorage
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.BuildingDTO
		wantErr bool
	}{
		{
			name: "authorised user",
			fields: fields{
				buildingRepo: mockBuildingStorage,
			},
			args: args{
				ctx: adminCtx,
			},
			want: []models.BuildingDTO{
				{
					BuildingID: uuid.Nil.String(),
					Name:       "advant",
				},
			},
			wantErr: false,
		},
		{
			name: "unauthorised user",
			fields: fields{
				buildingRepo: mockBuildingStorage,
			},
			args: args{
				ctx: userCtx,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := buildingservice.NewBuildingService(tt.fields.buildingRepo)
			got, err := bs.GetAllBuildings(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllBuildings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !slices.Equal(got, tt.want) {
				t.Errorf("GetAllBuildings() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildingService_DeleteBuilding(t *testing.T) {

	ctrl := gomock.NewController(t)
	mockBuildingStorage := mocks.NewMockBuildingStorage(ctrl)
	mockBuildingStorage.EXPECT().DeleteBuildingByName(gomock.Any()).Return(nil)

	type fields struct {
		buildingRepo buildingrepository.BuildingStorage
	}
	type args struct {
		ctx          context.Context
		buildingName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "authorised user",
			fields: fields{
				buildingRepo: mockBuildingStorage,
			},
			args: args{
				ctx:          adminCtx,
				buildingName: "advant",
			},
			wantErr: false,
		},
		{
			name: "unauthorised user",
			fields: fields{
				buildingRepo: mockBuildingStorage,
			},
			args: args{
				ctx:          userCtx,
				buildingName: "advant",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := buildingservice.NewBuildingService(tt.fields.buildingRepo)
			if err := bs.DeleteBuilding(tt.args.ctx, tt.args.buildingName); (err != nil) != tt.wantErr {
				t.Errorf("DeleteBuilding() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
