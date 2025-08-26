package slotservice

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	slotrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/slot_repository"
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

func TestNewSlotService(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSlotRepo := mocks.NewMockSlotStorage(ctrl)

	type args struct {
		slotRepo slotrepository.SlotStorage
	}
	tests := []struct {
		name string
		args args
		want *SlotService
	}{
		{
			name: "valid slot repo",
			args: args{
				slotRepo: mockSlotRepo,
			},
			want: &SlotService{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSlotService(tt.args.slotRepo); got == nil {
				t.Errorf("NewSlotService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlotService_GetFreeSlotsByBuilding(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSlotRepo := mocks.NewMockSlotStorage(ctrl)
	mockSlotRepo.EXPECT().GetFreeSlotsByBuilding(gomock.Any()).Return([]models.Slot{
		{
			BuildingID:  uuid.Nil,
			FloorNumber: 1,
			SlotNumber:  1,
			SlotType:    vehicletypes.TwoWheeler,
		},
	}, nil)
	mockSlotRepo.EXPECT().GetFreeSlotsByBuilding(gomock.Any()).Return(nil, errors.New("error"))

	type fields struct {
		slotRepo slotrepository.SlotStorage
	}
	type args struct {
		ctx         context.Context
		buildingID  uuid.UUID
		vehicleType vehicletypes.VehicleType
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.Slot
		wantErr bool
	}{
		{
			name: "adminctx",
			fields: fields{
				slotRepo: mockSlotRepo,
			},
			args: args{
				ctx:         adminCtx,
				buildingID:  uuid.New(),
				vehicleType: vehicletypes.TwoWheeler,
			},
			want: []models.Slot{
				{
					BuildingID:  uuid.Nil,
					FloorNumber: 1,
					SlotNumber:  1,
					SlotType:    vehicletypes.TwoWheeler,
				},
			},
			wantErr: false,
		},
		{
			name: "customer ctx",
			fields: fields{
				slotRepo: mockSlotRepo,
			},
			args: args{
				ctx:         customerCtx,
				buildingID:  uuid.New(),
				vehicleType: vehicletypes.TwoWheeler,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "no free slots",
			fields: fields{
				slotRepo: mockSlotRepo,
			},
			args: args{
				ctx:         adminCtx,
				buildingID:  uuid.New(),
				vehicleType: vehicletypes.TwoWheeler,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := &SlotService{
				slotRepo: tt.fields.slotRepo,
			}
			got, err := ss.GetFreeSlotsByBuilding(tt.args.ctx, tt.args.buildingID, tt.args.vehicleType)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFreeSlotsByBuilding() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFreeSlotsByBuilding() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlotService_GetSlotsByFloor(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSlotRepo := mocks.NewMockSlotStorage(ctrl)
	mockSlotRepo.EXPECT().GetSlotsByFloor(gomock.Any(), gomock.Any()).Return([]models.Slot{
		{
			BuildingID:  uuid.Nil,
			FloorNumber: 1,
			SlotNumber:  1,
			SlotType:    vehicletypes.TwoWheeler,
			Vehicles:    []models.Vehicle{},
		},
	}, nil)
	mockSlotRepo.EXPECT().GetSlotsByFloor(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))

	type fields struct {
		slotRepo slotrepository.SlotStorage
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
		want    []models.SlotDTO
		wantErr bool
	}{
		{
			name: "adminctx",
			fields: fields{
				slotRepo: mockSlotRepo,
			},
			args: args{
				ctx:         adminCtx,
				floorNumber: 1,
				buildingId:  uuid.New().String(),
			},
			want: []models.SlotDTO{
				{
					BuildingID:  uuid.Nil.String(),
					FloorNumber: 1,
					SlotNumber:  1,
					SlotType:    vehicletypes.TwoWheeler.String(),
					IsOccupied:  false,
				},
			},
			wantErr: false,
		},
		{
			name: "customer ctx",
			fields: fields{
				slotRepo: mockSlotRepo,
			},
			args: args{
				ctx:         customerCtx,
				floorNumber: 1,
				buildingId:  uuid.New().String(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid uuid",
			fields: fields{
				slotRepo: mockSlotRepo,
			},
			args: args{
				ctx:         adminCtx,
				floorNumber: 1,
				buildingId:  "fdsafsdf",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "repo error",
			fields: fields{
				slotRepo: mockSlotRepo,
			},
			args: args{
				ctx:         adminCtx,
				floorNumber: 1,
				buildingId:  uuid.New().String(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := &SlotService{
				slotRepo: tt.fields.slotRepo,
			}
			got, err := ss.GetSlotsByFloor(tt.args.ctx, tt.args.buildingId, tt.args.floorNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSlotsByFloor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSlotsByFloor() got = %v, want %v", got, tt.want)
			}
		})
	}
}
