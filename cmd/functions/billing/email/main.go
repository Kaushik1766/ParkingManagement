package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Kaushik1766/ParkingManagement/internal/config"
	billrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/bill_repository"
	parkinghistoryrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/parking_history_repository"
	userrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/user_repository"
	billingservice "github.com/Kaushik1766/ParkingManagement/internal/service/billing_service"
	emailservice "github.com/Kaushik1766/ParkingManagement/internal/service/email_service"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var billingService billingservice.BillingMgr
var emailService emailservice.EmailManager

func init() {
	cfg, err := awsconfig.LoadDefaultConfig(context.TODO(), awsconfig.WithRegion(config.AwsRegion))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := dynamodb.NewFromConfig(cfg)
	userRepo := userrepository.NewNOSQLUserRepository(client)
	parkingRepo := parkinghistoryrepository.NewNOSQLParkingRepository(client)
	billRepo := billrepository.NewNOSQLBillRepository(client)

	billingService = billingservice.NewBillingService(userRepo, parkingRepo, billRepo)
	emailService = emailservice.NewEmailService()
}

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, event events.SQSMessage) error {
	fmt.Println(event.Body)
	return nil
}
