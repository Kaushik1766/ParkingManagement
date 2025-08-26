package billingservice

import (
	"reflect"
	"testing"

	parkinghistoryrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/parking_history_repository"
	userrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/user_repository"
	"github.com/Kaushik1766/ParkingManagement/mocks"
	"go.uber.org/mock/gomock"
)

// TODO: discuss invoice generation strategy with mentors for rest
//func TestBillingService_GenerateMonthlyInvoice(t *testing.T) {
//
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockUserRepo := mocks.NewMockUserStorage(ctrl)
//	mockParkingRepo := mocks.NewMockParkingHistoryStorage(ctrl)
//	mockUserRepo.EXPECT().GetAllUsers().AnyTimes().Return([]models.User{}, nil)
//	mockParkingRepo.EXPECT().GetParkingHistoryByUser(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return([]models.ParkingHistoryDTO{}, nil)
//	type fields struct {
//		userRepository    userrepository.UserStorage
//		parkingRepository parkinghistoryrepository.ParkingHistoryStorage
//	}
//	tests := []struct {
//		name   string
//		fields fields
//	}{
//		{
//			name: "valid",
//			fields: fields{
//				userRepository:    mockUserRepo,
//				parkingRepository: mockParkingRepo,
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			bs := &BillingService{
//				userRepository:    tt.fields.userRepository,
//				parkingRepository: tt.fields.parkingRepository,
//			}
//			bs.GenerateMonthlyInvoice()
//		})
//	}
//}

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
