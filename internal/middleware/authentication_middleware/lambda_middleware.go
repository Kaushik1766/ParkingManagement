package authenticationmiddleware

import (
	"context"

	customerrors "github.com/Kaushik1766/ParkingManagement/pkg/customErrors"
	"github.com/aws/aws-lambda-go/events"
)

func AuthorizedInvoke(fn func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		authHeader := req.Headers["Authorization"]
		if authHeader == "" {
			return customerrors.LambdaError(401, "Unauthorized"), nil
		}

		token := authHeader[len("Bearer "):]
		authCtx, err := CliAuthenticate(ctx, token)
		if err != nil {
			return customerrors.LambdaError(401, "Unauthorized: "+err.Error()), nil
		}

		return fn(authCtx, req)
	}
}
