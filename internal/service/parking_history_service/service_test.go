package parkinghistoryservice

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	parkinghistoryrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/parking_history_repository"
	vehiclerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/vehicle_repository"
	"github.com/Kaushik1766/ParkingManagement/mocks"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

var user = models.UserJwt{
	RegisteredClaims: jwt.RegisteredClaims{
		ID: uuid.NewString(),
	},
	Email:  "kaushik@.acom",
	Role:   roles.Customer,
	Office: uuid.NewString(),
}

var userCtx = context.WithValue(context.Background(), constants.User, user)

func TestNewParkingHistoryService(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParkingRepo := mocks.NewMockParkingHistoryStorage(ctrl)
	mockVehicleRepo := mocks.NewMockVehicleStorage(ctrl)

	type args struct {
		parkingRepo parkinghistoryrepository.ParkingHistoryStorage
		vehicleRepo vehiclerepository.VehicleStorage
	}
	tests := []struct {
		name string
		args args
		want *ParkingHistoryService
	}{
		{
			name: "valid repos",
			args: args{
				parkingRepo: mockParkingRepo,
				vehicleRepo: mockVehicleRepo,
			},
			want: &ParkingHistoryService{
				parkingRepo: mockParkingRepo,
				vehicleRepo: mockVehicleRepo,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewParkingHistoryService(tt.args.parkingRepo, tt.args.vehicleRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewParkingHistoryService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParkingHistoryService_GetActiveUserParkings(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParkingRepo := mocks.NewMockParkingHistoryStorage(ctrl)
	mockVehicleRepo := mocks.NewMockVehicleStorage(ctrl)
	mockParkingRepo.EXPECT().GetActiveUserParkings(gomock.Any(), gomock.Any()).Return([]models.ParkingHistoryDTO{
		{
			TicketId:     "asdf",
			NumberPlate:  "asdf",
			BuildingId:   "asdf",
			FLoorNumber:  0,
			SlotNumber:   0,
			StartTime:    time.Time{},
			EndTime:      time.Time{},
			VechicleType: vehicletypes.FourWheeler.String(),
		},
	}, nil)

	type fields struct {
		parkingRepo parkinghistoryrepository.ParkingHistoryStorage
		vehicleRepo vehiclerepository.VehicleStorage
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.ParkingHistoryDTO
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				parkingRepo: mockParkingRepo,
				vehicleRepo: mockVehicleRepo,
			},
			args: args{
				ctx: userCtx,
			},
			want: []models.ParkingHistoryDTO{
				{
					TicketId:     "asdf",
					NumberPlate:  "asdf",
					BuildingId:   "asdf",
					FLoorNumber:  0,
					SlotNumber:   0,
					StartTime:    time.Time{},
					EndTime:      time.Time{},
					VechicleType: vehicletypes.FourWheeler.String(),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			phs := &ParkingHistoryService{
				parkingRepo: tt.fields.parkingRepo,
				vehicleRepo: tt.fields.vehicleRepo,
			}
			got, err := phs.GetActiveUserParkings(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetActiveUserParkings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetActiveUserParkings() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParkingHistoryService_GetParkingHistory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParkingRepo := mocks.NewMockParkingHistoryStorage(ctrl)
	mockVehicleRepo := mocks.NewMockVehicleStorage(ctrl)

	mockParkingRepo.EXPECT().GetParkingHistoryByUser(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]models.ParkingHistoryDTO{}, nil)

	type fields struct {
		parkingRepo parkinghistoryrepository.ParkingHistoryStorage
		vehicleRepo vehiclerepository.VehicleStorage
	}
	type args struct {
		ctx       context.Context
		startTime time.Time
		endTime   time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.ParkingHistoryDTO
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				parkingRepo: mockParkingRepo,
				vehicleRepo: mockVehicleRepo,
			},
			args: args{
				ctx:       userCtx,
				startTime: time.Now(),
				endTime:   time.Now(),
			},
			want:    []models.ParkingHistoryDTO{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			phs := &ParkingHistoryService{
				parkingRepo: tt.fields.parkingRepo,
				vehicleRepo: tt.fields.vehicleRepo,
			}
			got, err := phs.GetParkingHistory(tt.args.ctx, tt.args.startTime, tt.args.endTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetParkingHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetParkingHistory() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParkingHistoryService_GetParkingHistoryByNumberPlate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParkingRepo := mocks.NewMockParkingHistoryStorage(ctrl)
	mockVehicleRepo := mocks.NewMockVehicleStorage(ctrl)

	mockVehicleRepo.EXPECT().GetVehicleByNumberPlate(gomock.Any(), "asdf").Return(models.Vehicle{
		UserID: uuid.MustParse(user.ID),
	}, nil)
	mockVehicleRepo.EXPECT().GetVehicleByNumberPlate(gomock.Any(), "dasdf").Return(models.Vehicle{
		UserID: uuid.New(),
	}, nil)
	mockVehicleRepo.EXPECT().GetVehicleByNumberPlate(gomock.Any(), "invalidNumberplate").Return(models.Vehicle{}, errors.New("invalid numberplate"))
	mockParkingRepo.EXPECT().GetParkingHistoryByNumberPlate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]models.ParkingHistoryDTO{}, nil).AnyTimes()

	type fields struct {
		parkingRepo parkinghistoryrepository.ParkingHistoryStorage
		vehicleRepo vehiclerepository.VehicleStorage
	}
	type args struct {
		ctx         context.Context
		numberplate string
		startTime   time.Time
		endTime     time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.ParkingHistoryDTO
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				parkingRepo: mockParkingRepo,
				vehicleRepo: mockVehicleRepo,
			},
			args: args{
				ctx:         userCtx,
				numberplate: "asdf",
				startTime:   time.Now(),
				endTime:     time.Now(),
			},
			want:    []models.ParkingHistoryDTO{},
			wantErr: false,
		},
		{
			name: "unauthorized",
			fields: fields{
				parkingRepo: mockParkingRepo,
				vehicleRepo: mockVehicleRepo,
			},
			args: args{
				ctx:         userCtx,
				numberplate: "dasdf",
				startTime:   time.Now(),
				endTime:     time.Now(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid",
			fields: fields{
				parkingRepo: mockParkingRepo,
				vehicleRepo: mockVehicleRepo,
			},
			args: args{
				ctx:         userCtx,
				numberplate: "invalidNumberplate",
				startTime:   time.Now(),
				endTime:     time.Now(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			phs := &ParkingHistoryService{
				parkingRepo: tt.fields.parkingRepo,
				vehicleRepo: tt.fields.vehicleRepo,
			}
			got, err := phs.GetParkingHistoryByNumberPlate(tt.args.ctx, tt.args.numberplate, tt.args.startTime, tt.args.endTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetParkingHistoryByNumberPlate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetParkingHistoryByNumberPlate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParkingHistoryService_GetParkingHistoryByUserId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParkingRepo := mocks.NewMockParkingHistoryStorage(ctrl)
	mockVehicleRepo := mocks.NewMockVehicleStorage(ctrl)
	mockParkingRepo.EXPECT().GetParkingHistoryByUser(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]models.ParkingHistoryDTO{}, nil)
	type fields struct {
		parkingRepo parkinghistoryrepository.ParkingHistoryStorage
		vehicleRepo vehiclerepository.VehicleStorage
	}
	type args struct {
		ctx       context.Context
		userId    string
		startTime time.Time
		endTime   time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.ParkingHistoryDTO
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				parkingRepo: mockParkingRepo,
				vehicleRepo: mockVehicleRepo,
			},
			args: args{
				ctx:       userCtx,
				userId:    userCtx.Value(constants.User).(models.UserJwt).ID,
				startTime: time.Now(),
				endTime:   time.Now(),
			},
			want:    []models.ParkingHistoryDTO{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			phs := &ParkingHistoryService{
				parkingRepo: tt.fields.parkingRepo,
				vehicleRepo: tt.fields.vehicleRepo,
			}
			got, err := phs.GetParkingHistoryByUserId(tt.args.ctx, tt.args.userId, tt.args.startTime, tt.args.endTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetParkingHistoryByUserId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetParkingHistoryByUserId() got = %v, want %v", got, tt.want)
			}
		})
	}
}
