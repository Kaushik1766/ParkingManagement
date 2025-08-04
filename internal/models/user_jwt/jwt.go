package userjwt

import (
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	"github.com/golang-jwt/jwt/v5"
)

type UserJwt struct {
	jwt.RegisteredClaims
	ID    string     `json:"id"`
	Email string     `json:"email"`
	Role  roles.Role `json:"role"`
}
