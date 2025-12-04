package billingservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	billingrates "github.com/Kaushik1766/ParkingManagement/internal/constants/billing_rates"
	models "github.com/Kaushik1766/ParkingManagement/internal/models"
	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	billrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/bill_repository"
	parkinghistoryrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/parking_history_repository"
	userrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/user_repository"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type BillingService struct {
	// userService           userservice.UserManager
	// parkingHistoryService parkinghistoryservice.ParkingHistoryMgr
	userRepository    userrepository.UserStorage
	parkingRepository parkinghistoryrepository.ParkingHistoryStorage
	billRepository    billrepository.BillStorage
	sqsClient         *sqs.Client
}

// func NewBillingService(userService userservice.UserManager, parkingHistoryService parkinghistoryservice.ParkingHistoryMgr) *BillingService {
// 	return &BillingService{
// 		userService:           userService,
// 		parkingHistoryService: parkingHistoryService,
// 	}
// }

func NewBillingService(userRepo userrepository.UserStorage, parkingRepo parkinghistoryrepository.ParkingHistoryStorage, billRepo billrepository.BillStorage, sqsClient *sqs.Client) *BillingService {
	return &BillingService{
		userRepository:    userRepo,
		parkingRepository: parkingRepo,
		billRepository:    billRepo,
		sqsClient:         sqsClient,
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

	// gen for current month
	now := time.Now()
	startTime := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	endTime := now
	month := int(startTime.Month())
	year := startTime.Year()

	log.Printf("billingservice: Generating bills for %d-%d", month, year)

	for _, user := range users {
		userEmail := user.Email
		userId := user.UserID.String()

		// check if already generated
		// existingBill, err := bs.billRepository.GetBill(ctx, userId, month, year)
		// if err == nil && existingBill.UserId != "" {
		// 	log.Printf("billingservice: Bill already exists for user %s, skipping...\n", userEmail)
		// 	continue
		// }

		parkingHistory, err := bs.parkingRepository.GetParkingHistoryByUser(ctx, userId, startTime, endTime)
		if err != nil {
			log.Printf("billingservice: Error fetching parking history for user %s: %v\n", userId, err)
			continue
		}

		log.Printf("billingservice: Found %d parking records for user %s (email: %s) in period %v to %v", len(parkingHistory), userId, userEmail, startTime, endTime)

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
			UserEmail:      userEmail,
			UserId:         userId,
			BillingMonth:   month,
			BillingYear:    year,
		}

		err = bs.sendToSQS(ctx, bill)
		if err != nil {
			log.Println(err)
			continue
		}

		err = bs.billRepository.SaveBill(ctx, bill)
		if err != nil {
			log.Printf("billingservice: Error saving bill for user %s: %v\n", userEmail, err)
		} else {
			log.Printf("billingservice: Generated and saved bill for user %s\n", userEmail)
		}
	}
}

func (bs *BillingService) sendToSQS(ctx context.Context, bill models.BillDTO) error {
	queueUrl := os.Getenv("EMAIL_SQS")
	emailMessage := models.SQSEmailMessage{
		To:     bill.UserEmail,
		Header: fmt.Sprintf(constants.BillEmailHeader, bill.BillDate),
		Body:   formatBillBody(bill),
	}

	messageBytes, err := json.Marshal(emailMessage)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = bs.sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody: aws.String(string(messageBytes)),
		QueueUrl:    aws.String(queueUrl),
	})
	return err
}

func formatBillBody(bill models.BillDTO) string {
	rows := strings.Builder{}

	for _, parking := range bill.ParkingHistory {
		if parking.EndTime.IsZero() {
			continue
		}
		duration := parking.EndTime.Sub(parking.StartTime).Hours()
		var cost float64
		if parking.VechicleType == vehicletypes.TwoWheeler.String() {
			cost = duration * billingrates.TwoWheeler
		} else {
			cost = duration * billingrates.FourWheeler
		}

		rows.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%.2f</td><td>%.2f</td></tr>",
			parking.StartTime.Format("2006-01-02"),
			parking.VechicleType,
			parking.NumberPlate,
			duration,
			cost,
		))
	}

	return fmt.Sprintf(constants.BillEmailTemplate, rows.String(), bill.TotalAmount)
}
