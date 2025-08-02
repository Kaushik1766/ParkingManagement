package authenticationmiddleware

import (
	"context"
	"errors"
	"time"

	"github.com/Kaushik1766/ParkingManagement/internal/config"
	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	"github.com/Kaushik1766/ParkingManagement/internal/models/user"
	userjwt "github.com/Kaushik1766/ParkingManagement/internal/models/user_jwt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func CliAuthenticate(ctx context.Context, token string) (context.Context, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return []byte(config.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := parsedToken.Claims.(userjwt.UserJwt); ok {
		if claims.ExpiresAt.Compare(time.Now()) == -1 {
			return nil, errors.New("token expired")
		}

		id, err := uuid.FromBytes([]byte(claims.ID))
		if err != nil {
			return nil, err
		}

		userCtx := context.WithValue(ctx, constants.User, user.UserContext{
			Id:    id,
			Email: claims.Email,
			Role:  claims.Role,
		})

		return userCtx, nil
	} else {
		return nil, nil
	}
}
