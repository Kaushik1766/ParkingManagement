package authservice

import (
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums"
	"github.com/golang-jwt/jwt/v5"
)

type AuthenticationManager interface {
	Login(email, password string) (*jwt.Token, error)
	Signup(name, email, password string, role enums.Role) error
}
