package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	authenticationmiddleware "github.com/Kaushik1766/ParkingManagement/internal/middleware/authentication_middleware"
	corsmiddleware "github.com/Kaushik1766/ParkingManagement/internal/middleware/cors_middleware"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	billrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/bill_repository"
	parkinghistoryrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/parking_history_repository"
	userrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/user_repository"
	billingservice "github.com/Kaushik1766/ParkingManagement/internal/service/billing_service"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

var billingService billingservice.BillingMgr

func init() {
	cfg, err := awsconfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := dynamodb.NewFromConfig(cfg)
	sqsClient := sqs.NewFromConfig(cfg)
	userRepo := userrepository.NewNOSQLUserRepository(client)
	parkingRepo := parkinghistoryrepository.NewNOSQLParkingRepository(client)
	billRepo := billrepository.NewNOSQLBillRepository(client)

	billingService = billingservice.NewBillingService(userRepo, parkingRepo, billRepo, sqsClient)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var userId string
	if ctxUser, ok := ctx.Value(constants.User).(models.UserJwt); ok {
		userId = ctxUser.Email
	}

	if userId == "" {
		log.Println("Unauthorized: User ID not found in context")
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       `{"message": "Unauthorized: User ID not found in context"}`,
		}, nil
	}

	monthStr := request.QueryStringParameters["month"]
	yearStr := request.QueryStringParameters["year"]

	if monthStr == "" || yearStr == "" {
		log.Println("Missing required query parameters: month, year")
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"message": "Missing required query parameters: month, year"}`,
		}, nil
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		log.Println("Invalid month value")
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"message": "Invalid month"}`,
		}, nil
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		log.Println("Invalid year value")
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"message": "Invalid year"}`,
		}, nil
	}

	bill, err := billingService.GetMonthlyBill(ctx, userId, month, year)
	if err != nil {
		log.Printf("Error getting monthly bill: %v\n", err)
		if err.Error() == "bill not found" {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusNotFound,
				Body:       `{"message": "Bill not generated for this period"}`,
			}, nil
		}
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"message": "Internal server error"}`,
		}, nil
	}

	responseBody, err := json.Marshal(bill)
	if err != nil {
		log.Printf("Error marshalling response: %v\n", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"message": "Internal server error"}`,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(responseBody),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(corsmiddleware.WithCORS(authenticationmiddleware.AuthorizedInvoke(handler)))
}
