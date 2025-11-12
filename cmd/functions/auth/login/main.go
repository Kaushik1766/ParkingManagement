package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Kaushik1766/ParkingManagement/db"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	userrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/user_repository"
	authservice "github.com/Kaushik1766/ParkingManagement/internal/service/auth_service"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// var lambdaApp *lambda_app.LambdaApp

var authService authservice.AuthenticationManager

func init() {
	gormDb, _ := db.InitDB()

	userRepo := userrepository.NewSQLUserRepository(gormDb)

	authService = authservice.NewAuthService(userRepo)
}

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req models.LoginRequestDTO

	err := json.Unmarshal([]byte(event.Body), &req)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body: fmt.Sprintf("%v", map[string]string{
				"message": "invalid request body",
			}),
		}, errors.New("bad input")
	}

	token, err := authService.Login(req)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 401,
			Body: fmt.Sprintf("%v", map[string]string{
				"message": "Invalid credentials",
			}),
		}, errors.New("invalid credentials")
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       `{"jwt":"` + token + `"}`,
	}, nil
}
