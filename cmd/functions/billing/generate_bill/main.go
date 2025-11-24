package main

import (
	"context"
	"log"

	"github.com/Kaushik1766/ParkingManagement/internal/config"
	billrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/bill_repository"
	parkinghistoryrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/parking_history_repository"
	userrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/user_repository"
	billingservice "github.com/Kaushik1766/ParkingManagement/internal/service/billing_service"
	"github.com/aws/aws-lambda-go/lambda"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var billingService billingservice.BillingMgr

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
}

func handler(ctx context.Context) error {
	log.Println("Starting monthly bill generation...")
	billingService.GenerateMonthlyBills(ctx)
	log.Println("Monthly bill generation completed.")
	return nil
}

func main() {
	lambda.Start(handler)
}
