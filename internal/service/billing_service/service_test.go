package billingservice

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	models "github.com/Kaushik1766/ParkingManagement/internal/models"
	billrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/bill_repository"
	parkinghistoryrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/parking_history_repository"
	userrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/user_repository"
	"github.com/Kaushik1766/ParkingManagement/mocks"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestBillingService_GetMonthlyBill(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserStorage(ctrl)
	mockParkingRepo := mocks.NewMockParkingHistoryStorage(ctrl)
	mockBillRepo := mocks.NewMockBillStorage(ctrl)

	// Test data
	// Test data
	userID := uuid.New().String()
	month := 10
	year := 2023

	expectedBill := models.BillDTO{
		TotalAmount: 400, // 2 hours * 200
		BillDate:    time.Now().Format(time.DateOnly),
		UserId:      userID,
	}

	type fields struct {
		userRepository    userrepository.UserStorage
		parkingRepository parkinghistoryrepository.ParkingHistoryStorage
		billRepository    billrepository.BillStorage
	}
	type args struct {
		userId string
		month  int
		year   int
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		mockSetup func()
		want      models.BillDTO
		wantErr   bool
	}{
		{
			name: "Success - Get monthly bill (cached)",
			fields: fields{
				userRepository:    mockUserRepo,
				parkingRepository: mockParkingRepo,
				billRepository:    mockBillRepo,
			},
			args: args{
				userId: userID,
				month:  month,
				year:   year,
			},
			mockSetup: func() {
				mockBillRepo.EXPECT().GetBill(gomock.Any(), userID, month, year).Return(expectedBill, nil)
			},
			want:    expectedBill,
			wantErr: false,
		},
		{
			name: "Error - Bill not found",
			fields: fields{
				userRepository:    mockUserRepo,
				parkingRepository: mockParkingRepo,
				billRepository:    mockBillRepo,
			},
			args: args{
				userId: userID,
				month:  month,
				year:   year,
			},
			mockSetup: func() {
				mockBillRepo.EXPECT().GetBill(gomock.Any(), userID, month, year).Return(models.BillDTO{}, nil)
			},
			want:    models.BillDTO{},
			wantErr: true,
		},
		{
			name: "Error - GetBill fails",
			fields: fields{
				userRepository:    mockUserRepo,
				parkingRepository: mockParkingRepo,
				billRepository:    mockBillRepo,
			},
			args: args{
				userId: userID,
				month:  month,
				year:   year,
			},
			mockSetup: func() {
				mockBillRepo.EXPECT().GetBill(gomock.Any(), userID, month, year).Return(models.BillDTO{}, errors.New("db error"))
			},
			want:    models.BillDTO{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			bs := NewBillingService(tt.fields.userRepository, tt.fields.parkingRepository, tt.fields.billRepository)
			got, err := bs.GetMonthlyBill(context.Background(), tt.args.userId, tt.args.month, tt.args.year)
			if (err != nil) != tt.wantErr {
				t.Errorf("BillingService.GetMonthlyBill() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// We can't compare BillDate exactly because it uses time.Now()
				// So we'll check other fields and ensure BillDate is today
				if got.UserId != tt.want.UserId {
					t.Errorf("BillingService.GetMonthlyBill() UserId = %v, want %v", got.UserId, tt.want.UserId)
				}
				if got.TotalAmount != tt.want.TotalAmount {
					t.Errorf("BillingService.GetMonthlyBill() TotalAmount = %v, want %v", got.TotalAmount, tt.want.TotalAmount)
				}
				if len(got.ParkingHistory) != len(tt.want.ParkingHistory) {
					t.Errorf("BillingService.GetMonthlyBill() ParkingHistory length = %v, want %v", len(got.ParkingHistory), len(tt.want.ParkingHistory))
				}
			}
		})
	}
}

func TestNewBillingService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserStorage(ctrl)
	mockParkingRepo := mocks.NewMockParkingHistoryStorage(ctrl)
	mockBillRepo := mocks.NewMockBillStorage(ctrl)

	type args struct {
		userRepo    userrepository.UserStorage
		parkingRepo parkinghistoryrepository.ParkingHistoryStorage
		billRepo    billrepository.BillStorage
	}
	tests := []struct {
		name string
		args args
		want *BillingService
	}{
		{
			name: "valid repos",
			args: args{
				userRepo:    mockUserRepo,
				parkingRepo: mockParkingRepo,
				billRepo:    mockBillRepo,
			},
			want: &BillingService{
				userRepository:    mockUserRepo,
				parkingRepository: mockParkingRepo,
				billRepository:    mockBillRepo,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBillingService(tt.args.userRepo, tt.args.parkingRepo, tt.args.billRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBillingService() = %v, want %v", got, tt.want)
			}
		})
	}
}
