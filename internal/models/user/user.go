package user

import (
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums"
	"github.com/google/uuid"
)

type User struct {
	UserId   uuid.UUID
	Name     string
	Email    string
	Password string
	Role     enums.Role
	IsActive bool
}

func (u User) GetID() string {
	return u.UserId.String()
}
