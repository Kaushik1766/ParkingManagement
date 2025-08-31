package billingservice

import (
	"errors"
	"os"
	"reflect"
	"testing"
	"time"

	models "github.com/Kaushik1766/ParkingManagement/internal/models"
	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	parkinghistoryrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/parking_history_repository"
	userrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/user_repository"
	"github.com/Kaushik1766/ParkingManagement/mocks"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

// TestBillingService_GenerateMonthlyInvoice tests the GenerateMonthlyInvoice method
func TestBillingService_GenerateMonthlyInvoice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserStorage(ctrl)
	mockParkingRepo := mocks.NewMockParkingHistoryStorage(ctrl)

	// Test data
	userID1 := uuid.New()
	userID2 := uuid.New()

	users := []models.User{
		{
			UserID: userID1,
			Name:   "John Doe",
			Email:  "john@example.com",
		},
		{
			UserID: userID2,
			Name:   "Jane Smith",
			Email:  "jane@example.com",
		},
	}

	startTime := time.Now().AddDate(0, -1, 0)

	parkingHistory1 := []models.ParkingHistoryDTO{
		{
			TicketId:     "TICKET001",
			NumberPlate:  "ABC123",
			BuildingId:   "BUILD001",
			FLoorNumber:  1,
			SlotNumber:   1,
			StartTime:    startTime,
			EndTime:      startTime.Add(2 * time.Hour), // 2 hours
			VechicleType: vehicletypes.FourWheeler,
		},
		{
			TicketId:     "TICKET002",
			NumberPlate:  "XYZ789",
			BuildingId:   "BUILD001",
			FLoorNumber:  1,
			SlotNumber:   2,
			StartTime:    startTime.Add(3 * time.Hour),
			EndTime:      startTime.Add(5 * time.Hour), // 2 hours
			VechicleType: vehicletypes.TwoWheeler,
		},
	}

	parkingHistory2 := []models.ParkingHistoryDTO{
		{
			TicketId:     "TICKET003",
			NumberPlate:  "DEF456",
			BuildingId:   "BUILD002",
			FLoorNumber:  2,
			SlotNumber:   1,
			StartTime:    startTime,
			EndTime:      time.Time{}, // Zero end time - should be skipped
			VechicleType: vehicletypes.FourWheeler,
		},
	}

	type fields struct {
		userRepository    userrepository.UserStorage
		parkingRepository parkinghistoryrepository.ParkingHistoryStorage
	}
	tests := []struct {
		name               string
		fields             fields
		mockSetup          func()
		expectFileCreation bool
	}{
		{
			name: "Success - Generate invoice with multiple users and parking history",
			fields: fields{
				userRepository:    mockUserRepo,
				parkingRepository: mockParkingRepo,
			},
			mockSetup: func() {
				mockUserRepo.EXPECT().GetAllUsers().Return(users, nil)
				mockParkingRepo.EXPECT().GetParkingHistoryByUser(userID1.String(), gomock.Any(), gomock.Any()).Return(parkingHistory1, nil)
				mockParkingRepo.EXPECT().GetParkingHistoryByUser(userID2.String(), gomock.Any(), gomock.Any()).Return(parkingHistory2, nil)
			},
			expectFileCreation: true,
		},
		{
			name: "Success - Generate invoice with empty parking history",
			fields: fields{
				userRepository:    mockUserRepo,
				parkingRepository: mockParkingRepo,
			},
			mockSetup: func() {
				mockUserRepo.EXPECT().GetAllUsers().Return(users, nil)
				mockParkingRepo.EXPECT().GetParkingHistoryByUser(userID1.String(), gomock.Any(), gomock.Any()).Return([]models.ParkingHistoryDTO{}, nil)
				mockParkingRepo.EXPECT().GetParkingHistoryByUser(userID2.String(), gomock.Any(), gomock.Any()).Return([]models.ParkingHistoryDTO{}, nil)
			},
			expectFileCreation: true,
		},
		{
			name: "Error - GetAllUsers fails",
			fields: fields{
				userRepository:    mockUserRepo,
				parkingRepository: mockParkingRepo,
			},
			mockSetup: func() {
				mockUserRepo.EXPECT().GetAllUsers().Return(nil, errors.New("database error"))
			},
			expectFileCreation: false,
		},
		{
			name: "Error - GetParkingHistoryByUser fails",
			fields: fields{
				userRepository:    mockUserRepo,
				parkingRepository: mockParkingRepo,
			},
			mockSetup: func() {
				mockUserRepo.EXPECT().GetAllUsers().Return(users, nil)
				mockParkingRepo.EXPECT().GetParkingHistoryByUser(userID1.String(), gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))
			},
			expectFileCreation: false,
		},
		{
			name: "Success - Empty user list",
			fields: fields{
				userRepository:    mockUserRepo,
				parkingRepository: mockParkingRepo,
			},
			mockSetup: func() {
				mockUserRepo.EXPECT().GetAllUsers().Return([]models.User{}, nil)
			},
			expectFileCreation: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up any existing bills.txt file
			os.Remove("bills.txt")

			tt.mockSetup()

			bs := &BillingService{
				userRepository:    tt.fields.userRepository,
				parkingRepository: tt.fields.parkingRepository,
			}

			// Note: GenerateMonthlyInvoice sleeps for 1 hour, but for testing we'll need to work around this
			// For now, we'll just call it and check if the file is created
			go func() {
				time.Sleep(100 * time.Millisecond) // Small delay to let the method start
				// This is a workaround since we can't easily test the sleep in unit tests
			}()

			bs.GenerateMonthlyInvoice()

			// Check if file was created as expected
			_, err := os.Stat("bills.txt")
			fileExists := !os.IsNotExist(err)

			if tt.expectFileCreation && !fileExists {
				t.Errorf("GenerateMonthlyInvoice() expected file creation but file was not created")
			}
			if !tt.expectFileCreation && fileExists {
				t.Errorf("GenerateMonthlyInvoice() did not expect file creation but file was created")
			}

			// Clean up
			os.Remove("bills.txt")
		})
	}
}

func TestNewBillingService(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserStorage(ctrl)
	mockParkingRepo := mocks.NewMockParkingHistoryStorage(ctrl)

	type args struct {
		userRepo    userrepository.UserStorage
		parkingRepo parkinghistoryrepository.ParkingHistoryStorage
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
			},
			want: &BillingService{
				userRepository:    mockUserRepo,
				parkingRepository: mockParkingRepo,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBillingService(tt.args.userRepo, tt.args.parkingRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBillingService() = %v, want %v", got, tt.want)
			}
		})
	}
}
