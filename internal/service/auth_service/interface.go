package authservice

import (
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums"
)

type AuthenticationManager interface {
	Login(email, password string) (string, error)
	Signup(name, email, password string, role enums.Role) error
}
