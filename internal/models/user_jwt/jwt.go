package userjwt

import (
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums"
	"github.com/golang-jwt/jwt/v5"
)

type UserJwt struct {
	jwt.RegisteredClaims
	ID    string     `json:"id"`
	Email string     `json:"email"`
	Role  enums.Role `json:"role"`
}
