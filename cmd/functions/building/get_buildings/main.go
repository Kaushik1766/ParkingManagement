package main

import (
	"context"
	"encoding/json"

	authenticationmiddleware "github.com/Kaushik1766/ParkingManagement/internal/middleware/authentication_middleware"
	corsmiddleware "github.com/Kaushik1766/ParkingManagement/internal/middleware/cors_middleware"
	buildingrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/building_repository"
	buildingservice "github.com/Kaushik1766/ParkingManagement/internal/service/building_service"
	customerrors "github.com/Kaushik1766/ParkingManagement/pkg/customErrors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var buildingService buildingservice.BuildingMgr

func init() {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic("aws config not found")
	}

	client := dynamodb.NewFromConfig(cfg)

	buildingRepo := buildingrepository.NewNOSQLBuidlingRepository(client)
	buildingService = buildingservice.NewBuildingService(buildingRepo)
}

func main() {
	lambda.Start(corsmiddleware.WithCORS(authenticationmiddleware.AuthorizedInvoke(handler)))
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
