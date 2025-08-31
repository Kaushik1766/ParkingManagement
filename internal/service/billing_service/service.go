package billingservice

import (
	"fmt"
	"log"
	"os"
	"time"

	billingrates "github.com/Kaushik1766/ParkingManagement/internal/constants/billing_rates"
	models "github.com/Kaushik1766/ParkingManagement/internal/models"
	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	parkinghistoryrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/parking_history_repository"
	userrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/user_repository"
)

type BillingService struct {
	// userService           userservice.UserManager
	// parkingHistoryService parkinghistoryservice.ParkingHistoryMgr
	userRepository    userrepository.UserStorage
	parkingRepository parkinghistoryrepository.ParkingHistoryStorage
}

// func NewBillingService(userService userservice.UserManager, parkingHistoryService parkinghistoryservice.ParkingHistoryMgr) *BillingService {
// 	return &BillingService{
// 		userService:           userService,
// 		parkingHistoryService: parkingHistoryService,
// 	}
// }

func NewBillingService(userRepo userrepository.UserStorage, parkingRepo parkinghistoryrepository.ParkingHistoryStorage) *BillingService {
	return &BillingService{
		userRepository:    userRepo,
		parkingRepository: parkingRepo,
	}
}

func (bs *BillingService) GenerateMonthlyInvoice() {
	// time.Sleep(config.BillingDuration)
	// log.Println("billingservice: Generating monthly invoice...")
	users, err := bs.userRepository.GetAllUsers()
	if err != nil {
		log.Println("billingservice: Error fetching users:", err)
		return
	}

	billsString := ""
	startTime := time.Now().AddDate(0, -1, 0)
	endTime := time.Now()

	for _, user := range users {
		parkingHistory, err := bs.parkingRepository.GetParkingHistoryByUser(user.UserID.String(), startTime, endTime)
		if err != nil {
			log.Printf("billingservice: Error fetching parking history for user %s: %v\n", user.UserID, err)
			return
		}

		var totalAmount float64 = 0
		for _, ph := range parkingHistory {
			if ph.EndTime.IsZero() {
				log.Printf("billingservice: Parking end time is zero for user %s, skipping...\n", user.UserID)
				continue
			}

			totalTime := ph.EndTime.Sub(ph.StartTime).Hours()
			if ph.VechicleType == vehicletypes.TwoWheeler {
				totalAmount += totalTime * billingrates.TwoWheeler
			} else {
				totalAmount += totalTime * billingrates.FourWheeler
			}
		}
		curBill := models.BillDTO{
			ParkingHistory: parkingHistory,
			TotalAmount:    totalAmount,
			BillDate:       time.Now().Format(time.DateOnly),
			UserId:         user.UserID.String(),
		}
		billsString += curBill.String()
	}

	os.WriteFile("bills.txt", []byte(billsString), 0666)
	fmt.Println("for demo: bill generated")
}
