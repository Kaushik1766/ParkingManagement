package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/Kaushik1766/ParkingManagement/db"
	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	authenticationmiddleware "github.com/Kaushik1766/ParkingManagement/internal/middleware/authentication_middleware"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	slotrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/slot_repository"
	slotservice "github.com/Kaushik1766/ParkingManagement/internal/service/slot_service"
	customerrors "github.com/Kaushik1766/ParkingManagement/pkg/customErrors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var slotService slotservice.SlotMgr

func init() {
	gormDb, err := db.InitDB()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize DB: %v", err))
	}

	slotRepo := slotrepository.NewSQLSlotRepository(gormDb)
	slotService = slotservice.NewSlotService(slotRepo)
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

	floorIDStr := strings.TrimSpace(event.PathParameters["floorId"])
	if floorIDStr == "" {
		return customerrors.LambdaError(400, "missing floorId"), nil
	}

	floorNumber, err := strconv.Atoi(floorIDStr)
	if err != nil {
		return customerrors.LambdaError(400, "invalid floorId"), nil
	}

	slots, err := slotService.GetSlotsByFloor(ctx, buildingID, floorNumber)
	if err != nil {
		return lambdaError(err), nil
	}

	body, _ := json.Marshal(slots)
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
