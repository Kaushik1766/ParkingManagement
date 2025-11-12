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
	buildingrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/building_repository"
	buildingservice "github.com/Kaushik1766/ParkingManagement/internal/service/building_service"
	customerrors "github.com/Kaushik1766/ParkingManagement/pkg/customErrors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var buildingService buildingservice.BuildingMgr

func init() {
	gormDb, err := db.InitDB()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize DB: %v", err))
	}

	buildingRepo := buildingrepository.NewSQLBuildingRepository(gormDb)
	buildingService = buildingservice.NewBuildingService(buildingRepo)
}

func main() {
	lambda.Start(authenticationmiddleware.AuthorizedInvoke(handler))
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userCtx := ctx.Value(constants.User).(models.UserJwt)
	if userCtx.Role != roles.Admin {
		return customerrors.LambdaError(401, "Unauthorized access"), nil
	}

	var req struct {
		Name string `json:"buildingName"`
	}

	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		return customerrors.LambdaError(400, "invalid request body"), nil
	}

	if strings.TrimSpace(req.Name) == "" {
		return customerrors.LambdaError(400, "building name required"), nil
	}

	if err := buildingService.AddBuilding(ctx, req.Name); err != nil {
		return lambdaError(err), nil
	}

	body, _ := json.Marshal(map[string]string{"message": "Building added successfully"})
	return events.APIGatewayProxyResponse{
		StatusCode: 201,
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
