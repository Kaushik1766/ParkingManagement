package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Kaushik1766/ParkingManagement/db"
	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	authenticationmiddleware "github.com/Kaushik1766/ParkingManagement/internal/middleware/authentication_middleware"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	floorrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/floor_repository"
	floorservice "github.com/Kaushik1766/ParkingManagement/internal/service/floor_service"
	customerrors "github.com/Kaushik1766/ParkingManagement/pkg/customErrors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var floorService floorservice.FloorMgr

func init() {
	gormDb, err := db.InitDB()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize DB: %v", err))
	}

	floorRepo := floorrepository.NewSQLFloorRepository(gormDb)
	floorService = floorservice.NewFloorService(floorRepo)
}

func main() {
	lambda.Start(authenticationmiddleware.AuthorizedInvoke(handler))
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userCtx := ctx.Value(constants.User).(models.UserJwt)
	if userCtx.Role != roles.Admin {
		return customerrors.LambdaError(401, "Unauthorized access"), nil
	}

	buildingID := strings.TrimSpace(event.PathParameters["buildingId"])
	if buildingID == "" {
		return customerrors.LambdaError(400, "missing buildingId"), nil
	}

	var req struct {
		FloorNumber int `json:"floor_number"`
	}

	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		return customerrors.LambdaError(400, "invalid request body"), nil
	}

	if req.FloorNumber <= 0 {
		return customerrors.LambdaError(400, "invalid floor number"), nil
	}

	if err := floorService.AddFloorByBuildingId(ctx, buildingID, req.FloorNumber); err != nil {
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
