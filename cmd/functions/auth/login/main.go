package main

import (
	"context"
	"encoding/json"
	"log"

	corsmiddleware "github.com/Kaushik1766/ParkingManagement/internal/middleware/cors_middleware"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	userrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/user_repository"
	authservice "github.com/Kaushik1766/ParkingManagement/internal/service/auth_service"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var authService authservice.AuthenticationManager

func init() {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic("aws config not found")
	}

	client := dynamodb.NewFromConfig(cfg)

	userRepo := userrepository.NewNOSQLUserRepository(client)
	authService = authservice.NewAuthService(userRepo)

	log.Println("initialized login function")
}

func main() {
	lambda.Start(corsmiddleware.WithCORS(handler))
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("%+v\n", ctx)
	log.Printf("%+v\n", event)

	var req models.LoginRequestDTO

	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		body, _ := json.Marshal(map[string]string{"message": "invalid request body"})
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       string(body),
		}, nil
	}

	token, err := authService.Login(ctx, req)
	if err != nil {
		body, _ := json.Marshal(map[string]string{"message": "Invalid credentials"})
		return events.APIGatewayProxyResponse{
			StatusCode: 401,
			Body:       string(body),
		}, nil
	}

	body, _ := json.Marshal(map[string]string{"jwt": token})
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(body),
	}, nil
}
