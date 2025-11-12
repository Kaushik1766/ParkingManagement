package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Kaushik1766/ParkingManagement/db"
	authenticationmiddleware "github.com/Kaushik1766/ParkingManagement/internal/middleware/authentication_middleware"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
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
)

var userService userservice.UserManager

func init() {
	gormDb, err := db.InitDB()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize DB: %v", err))
	}

	userRepo := userrepository.NewSQLUserRepository(gormDb)
	vehicleRepo := vehiclerepository.NewSQLVehicleRepository(gormDb)
	officeRepo := officerepository.NewSQLOfficeRepository(gormDb)
	buildingRepo := buildingrepository.NewSQLBuildingRepository(gormDb)
	floorRepo := floorrepository.NewSQLFloorRepository(gormDb)
	slotRepo := slotrepository.NewSQLSlotRepository(gormDb)
	assignmentService := slotassignment.NewSlotAssignmentService(vehicleRepo, floorRepo, buildingRepo, slotRepo, officeRepo)

	userService = userservice.NewUserService(userRepo, vehicleRepo, officeRepo, assignmentService, buildingRepo)
}

func main() {
	lambda.Start(authenticationmiddleware.AuthorizedInvoke(handler))
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userID := event.PathParameters["userId"]
	if userID == "" {
		return customerrors.LambdaError(400, "missing userId"), nil
	}

	var req models.UpdateUserDTO
	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		return customerrors.LambdaError(400, "invalid request body"), nil
	}

	if err := userService.UpdateProfile(ctx, userID, req); err != nil {
		return lambdaError(err), nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
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
