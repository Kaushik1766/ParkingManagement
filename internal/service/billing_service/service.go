package billingservice

import (
	"context"
	"errors"
	"log"
	"time"

	billingrates "github.com/Kaushik1766/ParkingManagement/internal/constants/billing_rates"
	models "github.com/Kaushik1766/ParkingManagement/internal/models"
	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	billrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/bill_repository"
	parkinghistoryrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/parking_history_repository"
	userrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/user_repository"
)

type BillingService struct {
	// userService           userservice.UserManager
	// parkingHistoryService parkinghistoryservice.ParkingHistoryMgr
	userRepository    userrepository.UserStorage
	parkingRepository parkinghistoryrepository.ParkingHistoryStorage
	billRepository    billrepository.BillStorage
}

// func NewBillingService(userService userservice.UserManager, parkingHistoryService parkinghistoryservice.ParkingHistoryMgr) *BillingService {
// 	return &BillingService{
// 		userService:           userService,
// 		parkingHistoryService: parkingHistoryService,
// 	}
// }

func NewBillingService(userRepo userrepository.UserStorage, parkingRepo parkinghistoryrepository.ParkingHistoryStorage, billRepo billrepository.BillStorage) *BillingService {
	return &BillingService{
		userRepository:    userRepo,
		parkingRepository: parkingRepo,
		billRepository:    billRepo,
	}
}

func (bs *BillingService) GetMonthlyBill(ctx context.Context, userId string, month, year int) (models.BillDTO, error) {
	// Check if bill already exists
	existingBill, err := bs.billRepository.GetBill(ctx, userId, month, year)
	if err != nil {
		log.Printf("billingservice: Error fetching bill for user %s: %v\n", userId, err)
		return models.BillDTO{}, err
	}
	if existingBill.UserId == "" {
		return models.BillDTO{}, errors.New("bill not found")
	}

	return existingBill, nil
}

func (bs *BillingService) GenerateMonthlyBills(ctx context.Context) {
	users, err := bs.userRepository.GetAllUsers(ctx)
	if err != nil {
		log.Println("billingservice: Error fetching users:", err)
		return
	}

	// Generate for previous month
	now := time.Now()
	startTime := time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, time.Local)
	endTime := startTime.AddDate(0, 1, 0).Add(-time.Nanosecond)
	month := int(startTime.Month())
	year := startTime.Year()

	log.Printf("billingservice: Generating bills for %d-%d", month, year)

	for _, user := range users {
		userId := user.UserID.String()

		// Check if bill already exists
		existingBill, err := bs.billRepository.GetBill(ctx, userId, month, year)
		if err == nil && existingBill.UserId != "" {
			log.Printf("billingservice: Bill already exists for user %s, skipping...\n", userId)
			continue
		}

		parkingHistory, err := bs.parkingRepository.GetParkingHistoryByUser(ctx, userId, startTime, endTime)
		if err != nil {
			log.Printf("billingservice: Error fetching parking history for user %s: %v\n", userId, err)
			continue
		}

		var totalAmount float64 = 0
		for _, ph := range parkingHistory {
			if ph.EndTime.IsZero() {
				continue
			}

			totalTime := ph.EndTime.Sub(ph.StartTime).Hours()
			if ph.VechicleType == vehicletypes.TwoWheeler.String() {
				totalAmount += totalTime * billingrates.TwoWheeler
			} else {
				totalAmount += totalTime * billingrates.FourWheeler
			}
		}

		bill := models.BillDTO{
			ParkingHistory: parkingHistory,
			TotalAmount:    totalAmount,
			BillDate:       time.Now().Format(time.DateOnly),
			UserId:         userId,
		}

		// Store the bill
		err = bs.billRepository.SaveBill(ctx, bill)
		if err != nil {
			log.Printf("billingservice: Error saving bill for user %s: %v\n", userId, err)
		} else {
			log.Printf("billingservice: Generated and saved bill for user %s\n", userId)
		}
	}
}
