package authservice

import (
	"context"

	"github.com/Kaushik1766/ParkingManagement/internal/models/enums"
)

type AuthenticationManager interface {
	Login(ctx context.Context, email, password string)
	Signup(ctx context.Context, name, email, password string, role enums.Role)
}
