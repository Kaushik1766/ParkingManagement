package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/Kaushik1766/ParkingManagement/db"
	authenticationmiddleware "github.com/Kaushik1766/ParkingManagement/internal/middleware/authentication_middleware"
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
	lambda.Start(authenticationmiddleware.AuthorizedInvoke(handler))
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	officeID := strings.TrimSpace(event.PathParameters["officeId"])
	if officeID == "" {
		return customerrors.LambdaError(400, "missing officeId"), nil
	}

	if err := officeService.RemoveOffice(ctx, officeID); err != nil {
		return lambdaError(err), nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 204,
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
