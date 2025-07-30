package userservice

import (
	"context"

	gofiledb "github.com/Kaushik1766/GoFileDB"
	userModel "github.com/Kaushik1766/ParkingManagement/internal/models/User"
)

type UserService struct {
	db gofiledb.Repository[userModel.User]
}

func (s UserService) UpdateProfile(ctx context.Context, name, email, password string) error {
	prev, err := s.db.GetByParameter()
}
