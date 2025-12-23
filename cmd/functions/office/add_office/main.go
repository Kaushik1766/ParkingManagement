package main

import (
	"context"
	"encoding/json"
	"strings"

	authenticationmiddleware "github.com/Kaushik1766/ParkingManagement/internal/middleware/authentication_middleware"
	corsmiddleware "github.com/Kaushik1766/ParkingManagement/internal/middleware/cors_middleware"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	officerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/office_repository"
	officeservice "github.com/Kaushik1766/ParkingManagement/internal/service/office_service"
	customerrors "github.com/Kaushik1766/ParkingManagement/pkg/customErrors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var officeService officeservice.OfficeMgr

func init() {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic("aws config not found")
	}

	client := dynamodb.NewFromConfig(cfg)
	officeRepo := officerepository.NewNOSQLOfficeRepository(client)
	officeService = officeservice.NewOfficeService(officeRepo)
}

func main() {
	lambda.Start(corsmiddleware.WithCORS(authenticationmiddleware.AuthorizedInvoke(handler)))
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	buildingID := strings.TrimSpace(event.PathParameters["buildingId"])
	if buildingID == "" {
		return customerrors.LambdaError(400, "missing buildingId"), nil
	}

	var officeReq models.OfficeDTO
	if err := json.Unmarshal([]byte(event.Body), &officeReq); err != nil {
		return customerrors.LambdaError(400, "invalid request body"), nil
	}

	if strings.TrimSpace(officeReq.OfficeName) == "" {
		return customerrors.LambdaError(400, "officeName is required"), nil
	}

	if err := officeService.AddOffice(ctx, officeReq.OfficeName, buildingID, officeReq.FloorNumber); err != nil {
		return lambdaError(err), nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
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
