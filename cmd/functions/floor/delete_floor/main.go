package main

import (
	"context"
	"strconv"
	"strings"

	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	authenticationmiddleware "github.com/Kaushik1766/ParkingManagement/internal/middleware/authentication_middleware"
	corsmiddleware "github.com/Kaushik1766/ParkingManagement/internal/middleware/cors_middleware"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	floorrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/floor_repository"
	floorservice "github.com/Kaushik1766/ParkingManagement/internal/service/floor_service"
	customerrors "github.com/Kaushik1766/ParkingManagement/pkg/customErrors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var floorService floorservice.FloorMgr

func init() {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic("aws config not found")
	}

	client := dynamodb.NewFromConfig(cfg)

	floorRepo := floorrepository.NewNOSQLFloorRepository(client)
	floorService = floorservice.NewFloorService(floorRepo)
}

func main() {
	lambda.Start(corsmiddleware.WithCORS(authenticationmiddleware.AuthorizedInvoke(handler)))
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

	floorIDStr := strings.TrimSpace(event.PathParameters["floorId"])
	if floorIDStr == "" {
		return customerrors.LambdaError(400, "missing floorId"), nil
	}

	floorNumber, err := strconv.Atoi(floorIDStr)
	if err != nil {
		return customerrors.LambdaError(400, "invalid floorId"), nil
	}

	if err := floorService.DeleteFloor(ctx, buildingID, floorNumber); err != nil {
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
