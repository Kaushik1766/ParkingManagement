package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Kaushik1766/ParkingManagement/db"
	authenticationmiddleware "github.com/Kaushik1766/ParkingManagement/internal/middleware/authentication_middleware"
	corsmiddleware "github.com/Kaushik1766/ParkingManagement/internal/middleware/cors_middleware"
	parkinghistoryrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/parking_history_repository"
	vehiclerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/vehicle_repository"
	parkinghistoryservice "github.com/Kaushik1766/ParkingManagement/internal/service/parking_history_service"
	customerrors "github.com/Kaushik1766/ParkingManagement/pkg/customErrors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var parkingService parkinghistoryservice.ParkingHistoryMgr

func init() {
	gormDb, err := db.InitDB()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize DB: %v", err))
	}

	parkingRepo := parkinghistoryrepository.NewSQLParkingRepository(gormDb)
	vehicleRepo := vehiclerepository.NewSQLVehicleRepository(gormDb)

	parkingService = parkinghistoryservice.NewParkingHistoryService(parkingRepo, vehicleRepo)
}

func main() {
	lambda.Start(corsmiddleware.WithCORS(authenticationmiddleware.AuthorizedInvoke(handler)))
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	start := time.Now().AddDate(0, -1, 0)
	end := time.Now()

	if startParam := strings.TrimSpace(event.QueryStringParameters["startTime"]); startParam != "" {
		parsed, err := time.Parse(time.RFC3339, startParam)
		if err != nil {
			return customerrors.LambdaError(400, "invalid startTime format"), nil
		}
		start = parsed
	}

	if endParam := strings.TrimSpace(event.QueryStringParameters["endTime"]); endParam != "" {
		parsed, err := time.Parse(time.RFC3339, endParam)
		if err != nil {
			return customerrors.LambdaError(400, "invalid endTime format"), nil
		}
		end = parsed
	}

	parkings, err := parkingService.GetParkingHistory(ctx, start, end)
	if err != nil {
		return lambdaError(err), nil
	}

	body, _ := json.Marshal(parkings)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(body),
	}, nil
}

func lambdaError(err error) events.APIGatewayProxyResponse {
	msg := err.Error()
	status := 500
	lower := strings.ToLower(msg)
	switch {
	case strings.Contains(lower, "unauthorized"):
		status = 401
	case strings.Contains(lower, "invalid"),
		strings.Contains(lower, "not found"),
		strings.Contains(lower, "duplicate"),
		strings.Contains(lower, "cannot"),
		strings.Contains(lower, "already"),
		strings.Contains(lower, "empty"):
		status = 400
	}

	return customerrors.LambdaError(status, msg)
}
