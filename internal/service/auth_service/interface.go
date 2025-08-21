package authservice

import (
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
)

type AuthenticationManager interface {
	Login(email, password string) (string, error)
	Signup(registerReq models.RegisterRequestDTO, role roles.Role) error
}
