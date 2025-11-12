package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Kaushik1766/ParkingManagement/db"
	corsmiddleware "github.com/Kaushik1766/ParkingManagement/internal/middleware/cors_middleware"
	officerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/office_repository"
	officeservice "github.com/Kaushik1766/ParkingManagement/internal/service/office_service"
	customerrors "github.com/Kaushik1766/ParkingManagement/pkg/customErrors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var officeService officeservice.OfficeMgr

func init() {
	gormDb, err := db.InitDB()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize DB: %v", err))
	}

	officeRepo := officerepository.NewSQLOfficeRepository(gormDb)
	officeService = officeservice.NewOfficeService(officeRepo)
}

func main() {
	lambda.Start(corsmiddleware.WithCORS(handler))
}

func handler(ctx context.Context, _ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	offices, err := officeService.GetAllOfficeNames(ctx)
	if err != nil {
		return lambdaError(err), nil
	}

	body, _ := json.Marshal(offices)
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
