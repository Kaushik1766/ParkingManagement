package billingservice

import (
	parkinghistoryservice "github.com/Kaushik1766/ParkingManagement/internal/service/parking_history_service"
	userservice "github.com/Kaushik1766/ParkingManagement/internal/service/user_service"
)

type BillingService struct {
	userService           userservice.UserManager
	parkingHistoryService parkinghistoryservice.ParkingHistoryMgr
}

// func (bs *BillingService) GenerateMonthlyInvoice(customerID string) (string, error) {
// 	users, err := bs.userService.
// }
