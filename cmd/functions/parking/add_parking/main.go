package main

import (
	"context"
	"encoding/json"
	"strings"

	authenticationmiddleware "github.com/Kaushik1766/ParkingManagement/internal/middleware/authentication_middleware"
	corsmiddleware "github.com/Kaushik1766/ParkingManagement/internal/middleware/cors_middleware"
	parkinghistoryrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/parking_history_repository"
	vehiclerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/vehicle_repository"
	vehicleservice "github.com/Kaushik1766/ParkingManagement/internal/service/vehicle_service"
	customerrors "github.com/Kaushik1766/ParkingManagement/pkg/customErrors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var vehicleService vehicleservice.VehicleMgr

func init() {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic("aws config not found")
	}

	client := dynamodb.NewFromConfig(cfg)

	parkingRepo := parkinghistoryrepository.NewNOSQLParkingRepository(client)
	vehicleRepo := vehiclerepository.NewNOSQLVehicleRepository(client)

	vehicleService = vehicleservice.NewVehicleService(vehicleRepo, parkingRepo)
}

func main() {
	lambda.Start(corsmiddleware.WithCORS(authenticationmiddleware.AuthorizedInvoke(handler)))
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req struct {
		NumberPlate string `json:"numberplate"`
	}

	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		return customerrors.LambdaError(400, "invalid request body"), nil
	}

	if strings.TrimSpace(req.NumberPlate) == "" {
		return customerrors.LambdaError(400, "numberplate is required"), nil
	}

	ticketID, err := vehicleService.Park(ctx, req.NumberPlate)
	if err != nil {
		return lambdaError(err), nil
	}

	body, _ := json.Marshal(map[string]string{"ticketId": ticketID})
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
