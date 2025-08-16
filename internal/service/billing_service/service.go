package billingservice

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Kaushik1766/ParkingManagement/internal/config"
	billingrates "github.com/Kaushik1766/ParkingManagement/internal/constants/billing_rates"
	"github.com/Kaushik1766/ParkingManagement/internal/models/billing"
	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	parkinghistoryservice "github.com/Kaushik1766/ParkingManagement/internal/service/parking_history_service"
	userservice "github.com/Kaushik1766/ParkingManagement/internal/service/user_service"
)

type BillingService struct {
	userService           userservice.UserManager
	parkingHistoryService parkinghistoryservice.ParkingHistoryMgr
}

func NewBillingService(userService userservice.UserManager, parkingHistoryService parkinghistoryservice.ParkingHistoryMgr) *BillingService {
	return &BillingService{
		userService:           userService,
		parkingHistoryService: parkingHistoryService,
	}
}

func (bs *BillingService) GenerateMonthlyInvoice() {
	time.Sleep(config.BillingDuration)
	log.Println("billingservice: Generating monthly invoice...")
	users, err := bs.userService.GetAllUsers(context.Background())
	if err != nil {
		log.Println("billingservice: Error fetching users:", err)
		return
	}

	billsString := ""
	startTime := time.Now().AddDate(0, -1, 0)
	endTime := time.Now()

	for _, user := range users {
		parkingHistory, err := bs.parkingHistoryService.GetParkingHistoryById(user.UserID.String(), startTime, endTime)
		if err != nil {
			log.Printf("billingservice: Error fetching parking history for user %s: %v\n", user.UserID, err)
			return
		}

		var totalAmount float64 = 0
		for _, ph := range parkingHistory {
			// layout := "2006-01-02 15:04:05.999999999 -0700 MST"
			// startTime, err := time.Parse(layout, ph.StartTime)
			// if err != nil {
			// 	log.Printf("billingservice: Error parsing start time for user %s: %v\n", user.UserId, err)
			// 	continue
			// }
			//
			// endTime, err := time.Parse(layout, ph.EndTime)
			// if err != nil {
			// 	log.Printf("billingservice: Error parsing end time for user %s: %v\n", user.UserId, err)
			// 	continue
			// }
			if ph.EndTime.IsZero() {
				log.Printf("billingservice: Parking end time is zero for user %s, skipping...\n", user.UserID)
				continue
			}

			// log.Printf("billingservice: (%v) (%v)", endTime.UTC(), endTime.UTC().IsZero())
			totalTime := ph.EndTime.Sub(ph.StartTime).Hours()
			if ph.VechicleType == vehicletypes.TwoWheeler {
				totalAmount += totalTime * billingrates.TwoWheeler
			} else {
				totalAmount += totalTime * billingrates.FourWheeler
			}
		}
		curBill := billing.BillDTO{
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
