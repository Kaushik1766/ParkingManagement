package authenticationmiddleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Kaushik1766/ParkingManagement/internal/config"
	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	customerrors "github.com/Kaushik1766/ParkingManagement/pkg/customErrors"
	"github.com/golang-jwt/jwt/v5"
)

func CliAuthenticate(ctx context.Context, token string) (context.Context, error) {
	var tokenClaims models.UserJwt

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
	fmt.Println(tokenClaims)
	userCtx := context.WithValue(ctx, constants.User, tokenClaims)
	return userCtx, nil
}

func AuthenticatedRoute(fn func(ctx context.Context, w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			customerrors.UnauthorizedError(w, customerrors.WebError{
				Message: "Missing Authorization header",
				Code:    http.StatusUnauthorized,
			})
			return
		}

		token := authHeader[len("Bearer "):]
		ctx, err := CliAuthenticate(r.Context(), token)
		if err != nil {
			customerrors.UnauthorizedError(w, customerrors.WebError{
				Message: "Unauthorized: " + err.Error(),
				Code:    http.StatusUnauthorized,
			})
			return
		}

		fn(ctx, w, r)
	}
}
