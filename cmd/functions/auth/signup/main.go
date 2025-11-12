package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Kaushik1766/ParkingManagement/db"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
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
	lambda.Start(handler)
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req models.RegisterRequestDTO

	if err := json.Unmarshal([]byte(event.Body), &req); err != nil {
		body, _ := json.Marshal(map[string]string{"message": "invalid request body"})
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       string(body),
		}, nil
	}

	err := authService.Signup(req, roles.Customer)
	if err != nil {
		body, _ := json.Marshal(map[string]string{"message": "Invalid credentials"})
		return events.APIGatewayProxyResponse{
			StatusCode: 409,
			Body:       string(body),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
	}, nil
}
