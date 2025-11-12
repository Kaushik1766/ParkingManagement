package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Kaushik1766/ParkingManagement/db"
	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	authenticationmiddleware "github.com/Kaushik1766/ParkingManagement/internal/middleware/authentication_middleware"
	corsmiddleware "github.com/Kaushik1766/ParkingManagement/internal/middleware/cors_middleware"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
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
	lambda.Start(corsmiddleware.WithCORS(authenticationmiddleware.AuthorizedInvoke(handler)))
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userCtx := ctx.Value(constants.User).(models.UserJwt)

	var (
		vehicles []models.VehicleDTO
		err      error
	)

	if userCtx.Role == roles.Admin {
		if id := strings.TrimSpace(event.QueryStringParameters["userId"]); id != "" {
			vehicles, err = userService.GetVehiclesByUserId(ctx, id)
		}
	} else {
		vehicles, err = userService.GetRegisteredVehicles(ctx)
	}

	if err != nil {
		return lambdaError(err), nil
	}

	if vehicles == nil {
		vehicles = []models.VehicleDTO{}
	}

	body, _ := json.Marshal(vehicles)
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
