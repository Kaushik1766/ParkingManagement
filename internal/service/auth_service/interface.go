package authservice

import (
	"context"

	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
)

//go:generate mockgen -source=interface.go -destination=../../../mocks/auth_service_mock.go -package=mocks
type AuthenticationManager interface {
	Login(ctx context.Context, loginReq models.LoginRequestDTO) (string, error)
	Signup(ctx context.Context, registerReq models.RegisterRequestDTO, role roles.Role) error
}
