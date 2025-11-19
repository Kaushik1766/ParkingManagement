package main

import (
	"context"
	"encoding/json"
	"strings"

	authenticationmiddleware "github.com/Kaushik1766/ParkingManagement/internal/middleware/authentication_middleware"
	corsmiddleware "github.com/Kaushik1766/ParkingManagement/internal/middleware/cors_middleware"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	buildingrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/building_repository"
	floorrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/floor_repository"
	officerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/office_repository"
	slotrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/slot_repository"
	userrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/user_repository"
	vehiclerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/vehicle_repository"
	slotassignment "github.com/Kaushik1766/ParkingManagement/internal/service/slot_assignment"
	userservice "github.com/Kaushik1766/ParkingManagement/internal/service/user_service"
	customerrors "github.com/Kaushik1766/ParkingManagement/pkg/customErrors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var userService userservice.UserManager

func init() {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic("aws config not found")
	}

	client := dynamodb.NewFromConfig(cfg)

	userRepo := userrepository.NewNOSQLUserRepository(client)
	vehicleRepo := vehiclerepository.NewNOSQLVehicleRepository(client)
	officeRepo := officerepository.NewNOSQLOfficeRepository(client)
	buildingRepo := buildingrepository.NewNOSQLBuidlingRepository(client)
	floorRepo := floorrepository.NewNOSQLFloorRepository(client)
	slotRepo := slotrepository.NewNOSQLSlotRepository(client)
	assignmentService := slotassignment.NewSlotAssignmentService(vehicleRepo, floorRepo, buildingRepo, slotRepo, officeRepo)

	userService = userservice.NewUserService(userRepo, vehicleRepo, officeRepo, assignmentService, buildingRepo)
}

func main() {
	lambda.Start(corsmiddleware.WithCORS(authenticationmiddleware.AuthorizedInvoke(handler)))
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req models.AddVehicleDTO
	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		return customerrors.LambdaError(400, "invalid request body"), nil
	}

	if strings.TrimSpace(req.NumberPlate) == "" {
		return customerrors.LambdaError(400, "numberplate is required"), nil
	}

	if req.VehicleType != int(vehicletypes.TwoWheeler) && req.VehicleType != int(vehicletypes.FourWheeler) {
		return customerrors.LambdaError(400, "invalid vehicle type"), nil
	}

	if err := userService.RegisterVehicle(ctx, req.NumberPlate, vehicletypes.VehicleType(req.VehicleType)); err != nil {
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
