package billingservice_test

import (
	"testing"

	billingservice "github.com/Kaushik1766/ParkingManagement/internal/service/billing_service"
	parkinghistoryservice "github.com/Kaushik1766/ParkingManagement/internal/service/parking_history_service"
	userservice "github.com/Kaushik1766/ParkingManagement/internal/service/user_service"
)

func TestBillingService_GenerateMonthlyInvoice(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		userService           userservice.UserManager
		parkingHistoryService parkinghistoryservice.ParkingHistoryMgr
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := billingservice.NewBillingService(tt.userService, tt.parkingHistoryService)
			bs.GenerateMonthlyInvoice()
		})
	}
}
