package authservice

import (
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
)

//go:generate mockgen -source=interface.go -destination=../../../mocks/auth_service_mock.go -package=mocks
type AuthenticationManager interface {
	Login(loginReq models.LoginRequestDTO) (string, error)
	Signup(registerReq models.RegisterRequestDTO, role roles.Role) error
}
