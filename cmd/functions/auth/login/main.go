package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Kaushik1766/ParkingManagement/db"
	corsmiddleware "github.com/Kaushik1766/ParkingManagement/internal/middleware/cors_middleware"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	userrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/user_repository"
	authservice "github.com/Kaushik1766/ParkingManagement/internal/service/auth_service"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var authService authservice.AuthenticationManager

func init() {
	gormDb, err := db.InitDB()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize DB: %v", err))
	}

	userRepo := userrepository.NewSQLUserRepository(gormDb)
	authService = authservice.NewAuthService(userRepo)
}

func main() {
	lambda.Start(corsmiddleware.WithCORS(handler))
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req models.LoginRequestDTO

	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		body, _ := json.Marshal(map[string]string{"message": "invalid request body"})
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       string(body),
		}, nil
	}

	token, err := authService.Login(req)
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
