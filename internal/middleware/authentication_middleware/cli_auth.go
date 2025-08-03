package authenticationmiddleware

import (
	"context"
	"errors"
	"time"

	"github.com/Kaushik1766/ParkingManagement/internal/config"
	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	userjwt "github.com/Kaushik1766/ParkingManagement/internal/models/user_jwt"
	"github.com/golang-jwt/jwt/v5"
)

func CliAuthenticate(ctx context.Context, token string) (context.Context, error) {
	var tokenClaims userjwt.UserJwt

	parsedToken, err := jwt.ParseWithClaims(token, &tokenClaims, func(t *jwt.Token) (any, error) {
		return []byte(config.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, errors.New("invalid jwt")
	}
	if tokenClaims.ExpiresAt.Compare(time.Now()) == -1 {
		return nil, errors.New("token expired")
	}
	userCtx := context.WithValue(ctx, constants.User, tokenClaims)
	return userCtx, nil
}
