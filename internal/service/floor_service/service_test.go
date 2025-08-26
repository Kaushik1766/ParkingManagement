package floorservice

import (
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	floorrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/floor_repository"
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

var customerCtx = context.WithValue(context.Background(), constants.User, models.UserJwt{
	RegisteredClaims: jwt.RegisteredClaims{
		ID: "afdsfasdfasd",
	},
	Email:  "user@a.com",
	Role:   roles.Customer,
	Office: "asfd",
})

func TestFloorService_AddFloor(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	floorRepo := mocks.NewMockFloorStorage(ctrl)
	floorRepo.EXPECT().AddFloor(gomock.Any(), gomock.Any()).Return(nil)

	type fields struct {
		floorRepo floorrepository.FloorStorage
	}
	type args struct {
		ctx         context.Context
		buildingId  string
		floorNumber int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "adminctx",
			fields: fields{
				floorRepo: floorRepo,
			},
			args: args{
				ctx:         adminCtx,
				floorNumber: 1,
				buildingId:  "afdsfasdfasd",
			},
			wantErr: false,
		},
		{
			name: "customer ctx",
			fields: fields{
				floorRepo: floorRepo,
			},
			args: args{
				ctx:         customerCtx,
				floorNumber: 1,
				buildingId:  "afdsfasdfasd",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &FloorService{
				floorRepo: tt.fields.floorRepo,
			}
			if err := fs.AddFloor(tt.args.ctx, tt.args.buildingId, tt.args.floorNumber); (err != nil) != tt.wantErr {
				t.Errorf("AddFloor() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFloorService_AddFloorByBuildingId(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	floorRepo := mocks.NewMockFloorStorage(ctrl)
	floorRepo.EXPECT().AddFloor(gomock.Any(), gomock.Any()).Return(nil)

	type fields struct {
		floorRepo floorrepository.FloorStorage
	}
	type args struct {
		ctx         context.Context
		buildingId  string
		floorNumber int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "adminctx",
			fields: fields{
				floorRepo: floorRepo,
			},
			args: args{
				ctx:         adminCtx,
				buildingId:  "afdsfasdfasd",
				floorNumber: 1,
			},
			wantErr: false,
		},
		{
			name: "customer ctx",
			fields: fields{
				floorRepo: floorRepo,
			},
			args: args{
				ctx:         customerCtx,
				buildingId:  "afdsfasdfasd",
				floorNumber: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &FloorService{
				floorRepo: tt.fields.floorRepo,
			}
			if err := fs.AddFloorByBuildingId(tt.args.ctx, tt.args.buildingId, tt.args.floorNumber); (err != nil) != tt.wantErr {
				t.Errorf("AddFloorByBuildingId() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFloorService_DeleteFloor(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	floorRepo := mocks.NewMockFloorStorage(ctrl)
	floorRepo.EXPECT().DeleteFloor(gomock.Any(), gomock.Any()).Return(nil)

	type fields struct {
		floorRepo floorrepository.FloorStorage
	}
	type args struct {
		ctx         context.Context
		buildingId  string
		floorNumber int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "adminctx",
			fields: fields{
				floorRepo: floorRepo,
			},
			args: args{
				ctx:         adminCtx,
				buildingId:  "afdsfasdfasd",
				floorNumber: 1,
			},
			wantErr: false,
		},
		{
			name: "customer ctx",
			fields: fields{
				floorRepo: floorRepo,
			},
			args: args{
				ctx:         customerCtx,
				buildingId:  "afdsfasdfasd",
				floorNumber: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &FloorService{
				floorRepo: tt.fields.floorRepo,
			}
			if err := fs.DeleteFloor(tt.args.ctx, tt.args.buildingId, tt.args.floorNumber); (err != nil) != tt.wantErr {
				t.Errorf("DeleteFloor() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFloorService_GetFloorsByBuildingId(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	floorRepo := mocks.NewMockFloorStorage(ctrl)
	floorRepo.EXPECT().GetFloorsByBuildingId(gomock.Any()).Return([]models.Floor{
		{
			BuildingID:  uuid.Nil,
			FloorNumber: 1,
		},
	}, nil)
	floorRepo.EXPECT().GetFloorsByBuildingId(gomock.Any()).Return([]models.Floor{}, errors.New("no floors in building"))

	type fields struct {
		floorRepo floorrepository.FloorStorage
	}
	type args struct {
		ctx        context.Context
		buildingId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.FloorDTO
		wantErr bool
	}{
		{
			name: "adminctx",
			fields: fields{
				floorRepo: floorRepo,
			},
			args: args{
				ctx:        adminCtx,
				buildingId: "afdsfasdfasd",
			},
			want: []models.FloorDTO{
				{
					FloorNumber: 1,
					BuildingID:  uuid.Nil.String(),
				},
			},
			wantErr: false,
		},
		{
			name: "customer ctx",
			fields: fields{
				floorRepo: floorRepo,
			},
			args: args{
				ctx:        customerCtx,
				buildingId: "afdsfasdfasd",
			},
			want:    []models.FloorDTO{},
			wantErr: true,
		},
		{
			name: "no floors",
			fields: fields{
				floorRepo: floorRepo,
			},
			args: args{
				ctx:        adminCtx,
				buildingId: "afdsfasdfasd",
			},
			want:    []models.FloorDTO{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &FloorService{
				floorRepo: tt.fields.floorRepo,
			}
			got, err := fs.GetFloorsByBuildingId(tt.args.ctx, tt.args.buildingId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFloorsByBuildingId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !slices.Equal(got, tt.want) {
				t.Errorf("GetFloorsByBuildingId() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewFloorService(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	floorRepo := mocks.NewMockFloorStorage(ctrl)

	type args struct {
		floorRepo floorrepository.FloorStorage
	}
	tests := []struct {
		name string
		args args
		want *FloorService
	}{
		{
			name: "valid floor repo",
			args: args{
				floorRepo: floorRepo,
			},
			want: &FloorService{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFloorService(tt.args.floorRepo); got == nil {
				t.Errorf("NewFloorService() = %v, want %v", got, tt.want)
			}
		})
	}
}
