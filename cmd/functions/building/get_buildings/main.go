package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Kaushik1766/ParkingManagement/db"
	authenticationmiddleware "github.com/Kaushik1766/ParkingManagement/internal/middleware/authentication_middleware"
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

	buildings, err := buildingService.GetAllBuildings(ctx)
	if err != nil {
		return customerrors.LambdaError(401, err.Error()), nil
	}

	res, _ := json.Marshal(buildings)
	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Body:       string(res),
	}, nil
}
